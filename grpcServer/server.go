package grpcServer

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang_grpc_gin_jaeger_B/EchoServer"
	"golang_grpc_gin_jaeger_B/client"

	pb "golang_grpc_gin_jaeger_B/hello"

	"google.golang.org/grpc"

	"golang_grpc_gin_jaeger_B/tracing"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)

func initJaegerLog() {

	fmt.Println("jaeger init")
	tracer, closer = tracing.Init("gRPC service X")
	// defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
}

func Run() {

	apiListener, err := net.Listen("tcp", ":9997")
	if err != nil {
		log.Println(err)
		return
	}

	// 註冊 grpc
	es := &EchoServer.EchoServer{}

	// grpc := grpc.NewServer()
	grpc := grpc.NewServer(grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		// add opentracing stream interceptor to chain
		grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
	)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// add opentracing unary interceptor to chain
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))

	pb.RegisterGreeterServer(grpc, es)

	go runClient()

	if err := grpc.Serve(apiListener); err != nil {
		log.Fatal(" grpc.Serve Error: ", err)
		return
	}

}

func runClient() {

	time.Sleep(10 * time.Second)
	fmt.Println("After 3 seconds")

	client.Run()
}
