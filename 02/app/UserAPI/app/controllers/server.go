package controllers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"userapi/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var serverPort = config.Config.Port

func StartMainServer() {
	log.Println("info: Start Server" + "port: " + serverPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Otel Collecotor への接続設定
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// router 設定
	r := gin.New()

	r.Use(otelgin.Middleware("UserAPI-server"))

	//--- handler 設定
	r.POST("/createUser", createUser)
	r.POST("/getUserByEmail", getUserByEmail)

	r.POST("/encrypt", Encrypt)

	r.Run(":" + serverPort)
}
