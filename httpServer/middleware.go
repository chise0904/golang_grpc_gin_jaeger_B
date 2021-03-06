package httpServer

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func MyBenchLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		t := time.Now()
		// 给Context实例设置一个值
		c.Set("geektutu", "1111")
		// 请求前
		c.Next()
		// 请求后
		latency := time.Since(t)
		log.Println(latency)
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		t := time.Now()
		// 给Context实例设置一个值
		c.Set("geektutu", "1111")
		// 请求前
		c.Next()
		// 请求后
		latency := time.Since(t)
		log.Println(latency)
	}
}
