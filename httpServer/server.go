package httpServer

import (
	"fmt"
	"golang_grpc_gin_jaeger_B/tracing"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	opentracing "github.com/opentracing/opentracing-go"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
)

func initJaegerLog() {

	fmt.Println("jaeger init")
	tracer, closer = tracing.Init("microservice Y")
	// defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
}

func setupServer() *gin.Engine {

	initJaegerLog()

	r := gin.Default() // 使用默认中间件（logger和recovery）

	// GET
	r.GET("/ping", Ping)
	// named parameter
	r.GET("/user/:name", GetXByName)
	// query string
	r.GET("/users", GetUsers)

	// POST FORM
	r.POST("/postForm", PostForm)
	// POST MAP
	r.POST("/postMap", PostMap)
	// POST rawData
	r.POST("/postRawData", PostRawData)

	// redirect
	r.GET("/redirect", GetDirect)
	r.GET("/destination", GetDestination)

	// ===== group routes 分组路由 =====

	defaultHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": c.FullPath(),
		})
	}
	// group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/posts", defaultHandler)
		v1.GET("/series", defaultHandler)
	}
	// group: v2
	v2 := r.Group("/v2")
	{
		v2.GET("/posts", defaultHandler)
		v2.GET("/series", defaultHandler)
	}

	// ===== 文件上傳 =====

	// 上传文件(单个文件)
	r.POST("/uploadSingleFile", UploadSingleFile)
	// 上传文件(多个文件)
	r.POST("/uploadMultipleFiles", UploadMultipleFiles)

	// ===== 中间件(Middleware) =====

	exampleHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": c.FullPath(),
		})
	}

	// 作用于全局的Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 作用于单个路由的Middleware
	r.GET("/benchmark", MyBenchLogger(), exampleHandler)

	// 作用于某个组的Middleware
	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", exampleHandler)
		authorized.POST("/submit", exampleHandler)
	}

	return r
}

// func setupLog(){

// 	Tracer, Closer = tracing.Init("microservice")
// 	defer closer.Close()
// 	opentracing.SetGlobalTracer(tracer)
// }

func Run() {

	r := setupServer()
	r.Run(":8081")

	closer.Close()
}
