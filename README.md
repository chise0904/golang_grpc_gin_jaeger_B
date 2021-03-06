example of     
    
gin http server +   
grpc server +   
jaeger opentracing in grpc server  
    
jaeger opentracing across microservices & between functions call  

# jaeger opentracing request chain

golang_grpc_gin_jaeger_A is a grcp server    
golang_grpc_gin_jaeger_B is a grcp client    

# Deps

  go get "github.com/opentracing/opentracing-go"        
  go get "github.com/uber/jaeger-client-go"        
  go get "github.com/grpc-ecosystem/go-grpc-middleware"      
  go get "github.com/gin-gonic/gin"    
  go get "google.golang.org/grpc"     
  go get "github.com/uber/jaeger-lib/metrics/prometheus"      

# jaeger

  $ docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.21
