package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "golang_grpc_gin_jaeger_B/hello"
	"golang_grpc_gin_jaeger_B/tracing"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)

func initJaegerLog() {

	fmt.Println("jaeger init")
	tracer, closer = tracing.Init("Client A")
	// defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
}

func Run() {

	initJaegerLog()

	addr := "172.17.0.2:9997"
	// conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
	)

	if err != nil {
		log.Fatalf("Can not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Message: "Moto"})
	if err != nil {
		log.Fatalf("Could not get nonce: %v", err)
	}

	fmt.Println("Response:", r.GetMessage())
}
