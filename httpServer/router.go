package httpServer

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	olog "github.com/opentracing/opentracing-go/log"
)

func Ping(c *gin.Context) {

	// tracer, closer := tracing.Init("microservice")
	// defer closer.Close()
	// opentracing.SetGlobalTracer(tracer)

	span := tracer.StartSpan("Ping_Y")
	span.SetTag("Ping_Y", "pong")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	ping2(ctx)
	ping3(ctx)

	c.JSON(200, gin.H{ // 返回一个JSON，状态码是200，gin.H是map[string]interface{}的简写
		"message": "pong",
	})
}

func ping2(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Ping_Y2")
	defer span.Finish()

	span.LogFields(
		olog.String("event", "Ping_Y2 event"),
		olog.String("value", "Ping_Y2 value"),
	)
}

func ping3(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Ping_Y3")
	defer span.Finish()

	span.LogFields(
		olog.String("event", "Ping_Y3 event"),
		olog.String("value", "Ping_Y3 value"),
	)
	{
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		ping4(ctx)
	}
}

func ping4(ctx context.Context) {

	span, _ := opentracing.StartSpanFromContext(ctx, "Ping_Y4")
	defer span.Finish()

	url := "http://172.17.0.2:8080/users?name=barbara"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("err: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err: ", err)
	}
	span.LogFields(
		olog.String("event", string(body)),
	)
}

//=====================================================

func GetXByName(c *gin.Context) {

	// tracer, closer := tracing.Init("microservice")
	// defer closer.Close()
	// opentracing.SetGlobalTracer(tracer)

	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := tracer.StartSpan("GetXByName", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	span.SetTag(":name", c.Param("name"))

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	ping4(ctx)

	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func GetXByName1(ctx context.Context) {

	span, _ := opentracing.StartSpanFromContext(ctx, "GetXByName1")
	defer span.Finish()

	span.SetTag("GetXByName1", "adam")

	span.LogFields(
		olog.String("event", "GetXByName1 event"),
		olog.String("value", "GetXByName1 value"),
	)

	{
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		GetXByName2(ctx)
	}
}

func GetXByName2(ctx context.Context) {

	span, _ := opentracing.StartSpanFromContext(ctx, "GetXByName2")
	defer span.Finish()

	span.LogFields(
		olog.String("event", "GetXByName2 event"),
		olog.String("value", "GetXByName2 value"),
	)
	{
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		GetXByName3(ctx)
	}
}

func GetXByName3(ctx context.Context) {

	span, _ := opentracing.StartSpanFromContext(ctx, "GetXByName3")
	defer span.Finish()

	span.LogFields(
		olog.String("event", "GetXByName3 event"),
		olog.String("value", "GetXByName3 value"),
	)
}

//=====================================

func GetUsers(c *gin.Context) {

	name := c.Query("name")
	role := c.DefaultQuery("role", "teacher")

	c.String(http.StatusOK, "%s is a %s", name, role)
}

func PostForm(c *gin.Context) {

	username := c.PostForm("username")
	password := c.DefaultPostForm("password", "000000") // 可设置默认值

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"password": password,
	})
}

func PostMap(c *gin.Context) {

	ids := c.QueryMap("ids")
	names := c.PostFormMap("names")

	fmt.Printf("ids: %v", ids)
	fmt.Printf("names: %v", names)

	c.JSON(http.StatusOK, gin.H{
		"ids":   ids,
		"names": names,
	})
}

func GetDirect(c *gin.Context) {

	c.Redirect(http.StatusMovedPermanently, "/destination")

	// OR
	// c.Request.URL.Path = "/"
	// r.HandleContext(c)
}

func GetDestination(c *gin.Context) {

	c.String(http.StatusOK, "Direct Destination")
}

func UploadSingleFile(c *gin.Context) {

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusCreated, "upload successful \n")
}

func UploadMultipleFiles(c *gin.Context) {

	err := c.Request.ParseMultipartForm(200000)
	if err != nil {
		log.Fatal(err)
	}

	// 获取表单
	form := c.Request.MultipartForm

	// 获取参数upload后面的多个文件名，存放到数组files里面，
	files := form.File["upload"]

	// 遍历数组，每取出一个file就拷贝一次
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		fileName := files[i].Filename
		fmt.Println(fileName)

		out, err := os.Create(fileName)
		defer out.Close()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}

		c.String(http.StatusCreated, "uploadFiles success! \n")
	}
}

func PostRawData(c *gin.Context) {

	// 获取原始字节
	d, err := c.GetRawData()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(d))

	c.String(200, "ok")
}
