package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"
	"time"
	"todobff/app/SessionInfo"
	"todobff/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var serverPort = config.Config.Port
var pathStatic = config.Config.Static
var EpUserApi = config.Config.EpUserApi
var EpTodoAPI = config.Config.EpTodoApi

var LoginInfo SessionInfo.Session

func parseURL(fn func(*gin.Context, int)) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := tracer.Start(c.Request.Context(), "parseURL")
		defer span.End()

		fmt.Println(c.Request.URL.Path)
		q := validPath.FindStringSubmatch(c.Request.URL.Path)
		if q == nil {
			http.NotFound(c.Writer, c.Request)
			return
		}

		id, _ := strconv.Atoi(q[2])
		fmt.Println(id)
		fn(c, id)
	}
}

var validPath = regexp.MustCompile("^/menu/todos/(edit|save|update|delete)/([0-9]+)$")

func generateHTML(c *gin.Context, data interface{}, procname string, filenames ...string) {
	_, span := tracer.Start(c.Request.Context(), "generateHTML : "+procname)
	defer span.End()

	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(c.Writer, "layout", data)
}

func LoggerAndCreateSpan(c *gin.Context, msg string) trace.Span {
	_, span := tracer.Start(c.Request.Context(), msg)

	SpanId := span.SpanContext().SpanID().String()
	TraceId := span.SpanContext().TraceID().String()

	span.SetAttributes(
		attribute.Int("status", c.Writer.Status()),
		attribute.String("method", c.Request.Method),
		attribute.String("client_ip", c.ClientIP()),
		attribute.String("message", msg),
		attribute.String("span_id", SpanId),
		attribute.String("trace_id", TraceId),
	)

	start := time.Now()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()
	logger.Info("Logger",
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.Duration("elapsed", time.Since(start)),
		zap.String("message", msg),
		zap.String("span_id", SpanId),
		zap.String("trace_id", TraceId),
	)

	return span
}
