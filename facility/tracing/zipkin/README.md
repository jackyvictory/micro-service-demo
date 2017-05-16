# Reference
    http://zipkin.io/
    https://github.com/openzipkin/zipkin/
    https://github.com/openzipkin/zipkin-go-opentracing

# Environment
  localhost: OS X 10.10.3

# Requirements
  Elasticsearch on 192.168.1.198:9200

# Installation
tracing server and dashboard, for many languages

> using Java with JDK 1.8 or higher to install and run

      $ wget -O zipkin.jar 'https://search.maven.org/remote_content?g=io.zipkin.java&a=zipkin-server&v=LATEST&c=exec'
      $ STORAGE_TYPE=elasticsearch ES_HOSTS=http://192.168.1.198:9200 java -jar zipkin.jar

# Check zipkin status

      http://localhost:9411/
