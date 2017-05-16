# micro-service-demo

Just a demo

# Steps

* [Install gRPC and protobuf](https://github.com/jackyvictory/micro-service-demo/tree/master/service/pb)

* [Install etcd cluster on vagrant virtual boxex](https://github.com/jackyvictory/micro-service-demo/tree/master/vagrant)

* [Install and startup zipkin for tracing](https://github.com/jackyvictory/micro-service-demo/tree/master/facility/tracing/zipkin)

* [Run trans service demo for gRPC server](https://github.com/jackyvictory/micro-service-demo/tree/master/service/trans)

* [Run platform demo for gRPC client](https://github.com/jackyvictory/micro-service-demo/tree/master/product/cilPlatform)

* [Launch a gRPC request](http://localhost:9000/trans/)

* [See the traces of that request](http://localhost:9411/)

* [Monitor etcd cluster status](http://192.168.99.40:3000/)

* See gRPC requests' load balance

    >* Login vagrant virtual boxex app1 to app3 with user vagrant
    >* See the logs in /opt/services/transService.log while launching gRPC requests
    >* Also you can terminate some(or all) trans services in every box, and launch request again
