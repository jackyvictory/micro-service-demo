package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/CardInfoLink/log"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	pb "github.com/jackyvictory/micro-service-demo/service/pb"
	opentracing "github.com/opentracing/opentracing-go"
	grpc "google.golang.org/grpc"
	// "sourcegraph.com/sourcegraph/appdash"
	// appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	grpclb "github.com/jackyvictory/micro-service-demo/facility/grpclb"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	port          = flag.Int("port", 9000, "cil platform listening port")
	platServ      = flag.String("cil platform", "cilPlatform", "cil platform name")
	transServ     = flag.String("trans service", "trans", "transaction service name")
	etcdReg       = flag.String("reg", "http://192.168.99.40:2379,http://192.168.99.50:2379,http://192.168.99.60:2379", "register etcd address")
	endpoint      = flag.String("zipkin endpoint", "http://10.30.1.20:9411/api/v1/spans", "Endpoint to send Zipkin spans to")
	debug         = flag.Bool("debug mode", false, "zipkin debug mode")
	sameSpan      = flag.Bool("same span", true, "same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)")
	traceID128Bit = flag.Bool("trace id 128 bit", true, "make Tracer generate 128 bit traceID's for root spans.")
)

var transClient pb.TransactionClient

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// go func() {
	// 	ticker := time.NewTicker(time.Second)
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			log.Debug(time.Now().UnixNano())
	// 		}
	// 	}
	// }()

	// Init appdash tracing service
	// appDashCollector := appdash.NewRemoteCollector(addDashServAddr)
	// tracer := appdashot.NewTracer(appDashCollector)
	// opentracing.InitGlobalTracer(tracer)

	// Init zipkin tracing service
	zipkinCollector, err := zipkin.NewHTTPCollector(*endpoint)
	if err != nil {
		log.Error(err)
		return
	}
	defer zipkinCollector.Close()
	zipkinRecoder := zipkin.NewRecorder(zipkinCollector, *debug, fmt.Sprintf(":%d", *port), *platServ)
	tracer, err := zipkin.NewTracer(zipkinRecoder, zipkin.ClientServerSameSpan(*sameSpan), zipkin.TraceID128Bit(*traceID128Bit))
	if err != nil {
		log.Error(err)
		return
	}
	opentracing.InitGlobalTracer(tracer)

	// Init grpc load balancer
	r := grpclb.NewResolver(*transServ)
	b := grpc.RoundRobin(r)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Init gRPC trans service
	// All future RPC activity involving `conn` will be automatically traced.
	conn, err := grpc.DialContext(ctx, *etcdReg, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Errorf("failed to connect transaction service, error is %v\n", err)
		return
	}
	defer conn.Close()

	transClient = pb.NewTransactionClient(conn)

	// Init cilPlatform handler and serve on port platformServAddr
	http.Handle("/", httpRouter())
	log.Infof("serving cil platform on port %v\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Errorf("failed to serve http, error is %v\n", err)
		return
	}
}

func httpRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", login)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/trans", trans)
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	retStr := "hello, cil platform requires login"
	retBytes := []byte(retStr)
	w.Write(retBytes)
}

func trans(w http.ResponseWriter, r *http.Request) {
	var retStr string

	// Create Root Span for duration of the interaction
	span := opentracing.StartSpan(
		fmt.Sprintf("GET %s", r.URL.Path),
	)
	defer func() {
		span.Finish()
		// log.Debugf("%#+v", span)
	}()

	// Put root span in context so it will be used in our calls to the client.
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	time.Sleep(10 * time.Millisecond)
	transList, err := transClient.Find(ctx, &pb.QueryCond{})
	if err != nil {
		retStr = fmt.Sprintf("failed to call find on trans server, error is %v", err)
		log.Errorf(retStr)
		return
	}
	time.Sleep(20 * time.Millisecond)

	transSlice := transList.GetT()
	if transSlice == nil {
		retStr = fmt.Sprintf("nil slice returned from trans server")
		log.Errorf(retStr)
		return
	}

	for _, t := range transSlice {
		retStr += fmt.Sprintf("id: %s, orderNum: %s, transAmt: %d\n", t.GetId(), t.GetOrderNum(), t.GetTransAmt())
	}
	log.Infof(retStr)
	w.Write([]byte(retStr))
}
