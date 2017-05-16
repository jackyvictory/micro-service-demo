package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	// "sourcegraph.com/sourcegraph/appdash"

	"github.com/CardInfoLink/log"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	pb "github.com/jackyvictory/micro-service-demo/service/pb"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	grpclb "github.com/jackyvictory/micro-service-demo/facility/grpclb"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var (
	port          = flag.Int("port", 9001, "transaction service listening port")
	transServ     = flag.String("trans service", "trans", "transaction service name")
	etcdReg       = flag.String("reg", "http://192.168.99.40:2379,http://192.168.99.50:2379,http://192.168.99.60:2379", "register etcd address")
	endpoint      = flag.String("zipkin endpoint", "http://10.30.1.20:9411/api/v1/spans", "Endpoint to send Zipkin spans to")
	debug         = flag.Bool("debug mode", false, "zipkin debug mode")
	sameSpan      = flag.Bool("same span", true, "same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)")
	traceID128Bit = flag.Bool("trace id 128 bit", true, "make Tracer generate 128 bit traceID's for root spans.")
)

type transServer struct{}

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
	// appDashCollector := appdash.NewRemoteCollector(appDashServAddr)
	// tracer := appdashot.NewTracer(appDashCollector)
	// opentracing.InitGlobalTracer(tracer)

	// Init zipkin tracing service
	zipkinCollector, err := zipkin.NewHTTPCollector(*endpoint)
	if err != nil {
		log.Error(err)
		return
	}
	defer zipkinCollector.Close()
	zipkinRecoder := zipkin.NewRecorder(zipkinCollector, *debug, fmt.Sprintf(":%d", *port), *transServ)
	tracer, err := zipkin.NewTracer(zipkinRecoder, zipkin.ClientServerSameSpan(*sameSpan), zipkin.TraceID128Bit(*traceID128Bit))
	if err != nil {
		log.Error(err)
		return
	}
	opentracing.InitGlobalTracer(tracer)

	// Listen gRPC trans service port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Errorf("failed to listen, error is %v", err)
		return
	}

	// Regist grpc load balancer
	err = grpclb.Register(*transServ, localIP(), *port, *etcdReg, time.Second*10, 15)
	if err != nil {
		log.Errorf("failed to register grpclb, error is %v", err)
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	// Init gRPC trans service
	// All future RPC activity involving `s` will be automatically traced.
	s := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterTransactionServer(s, &transServer{})
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		log.Errorf("failed to serve, error is %v", err)
		return
	}
}

func (s *transServer) Add(ctx context.Context, in *pb.Trans) (*pb.Resp, error) {
	log.Infof("Adding trans %#+v\n", in)
	return &pb.Resp{Ok: true}, nil
}

func (s *transServer) Update(ctx context.Context, in *pb.Trans) (*pb.Resp, error) {
	log.Infof("Updating trans %#+v\n", in)
	return &pb.Resp{Ok: true}, nil
}

func (s *transServer) Find(ctx context.Context, in *pb.QueryCond) (*pb.TransList, error) {
	// create new span using span found in context as parent (if none is found, our span becomes the trace root).
	// span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("RPC %s.Find", *transServ))
	// defer span.Finish()

	log.Infof("Finding trans %#+v\n", in)
	time.Sleep(30 * time.Millisecond)
	list, err := dbMockFind(ctx)
	return list, err
}

func dbMockFind(ctx context.Context) (*pb.TransList, error) {
	// create new span using span found in context as parent (if none is found, our span becomes the trace root).
	resourceSpan, _ := opentracing.StartSpanFromContext(
		ctx,
		fmt.Sprintf("DB %s.dbMockFind", *transServ),
		// opentracing.StartTime(time.Now()),
	)
	defer func() {
		resourceSpan.Finish()
		// log.Debugf("%#+v", resourceSpan)
	}()
	// mark span as resource type
	ext.SpanKind.Set(resourceSpan, "resource")
	// name of the resource we try to reach
	ext.PeerService.Set(resourceSpan, "MySQL")
	// hostname of the resource
	ext.PeerHostname.Set(resourceSpan, "localhost")
	// port of the resource
	ext.PeerPort.Set(resourceSpan, 8806)
	// let's binary annotate the query we run
	resourceSpan.SetTag(
		"query", "SELECT * FROM test",
	)

	// Let's assume the query is going to take some time. Finding the right
	// world domination recipes is like searching for a needle in a haystack.
	time.Sleep(40 * time.Millisecond)

	return &pb.TransList{
		T: []*pb.Trans{
			&pb.Trans{Id: "id1", OrderNum: "order1", TransAmt: 1},
			&pb.Trans{Id: "id2", OrderNum: "order2", TransAmt: 2},
			&pb.Trans{Id: "id3", OrderNum: "order3", TransAmt: 3},
		},
	}, nil
}

// localIP 本机 IP
func localIP() (localIP string) {
	inter := "eth1"

	ifi, err := net.InterfaceByName(inter)
	if err != nil {
		log.Error(err)
		return
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		log.Error(err)
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				break
			}
		}
	}
	// log.Debugf("local ip is %v", localIP)
	return
}
