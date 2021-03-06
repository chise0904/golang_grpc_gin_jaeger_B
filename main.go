package main

import (
	// "golang_grpc_gin_jaeger_B/grpcServer"
	// "golang_grpc_gin_jaeger_B/httpServer"
	"golang_grpc_gin_jaeger_B/client"
)

// func main() {

// 	ch := make(chan struct{})

// 	go grpcServer.Run()
// 	go httpServer.Run()

// 	<-ch
// }

func main() {

	client.Run()
}
