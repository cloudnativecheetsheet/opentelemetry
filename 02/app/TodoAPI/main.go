// main.go

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todoapi/app/controllers"
)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Process Shutdown...")
		os.Exit(1)
	}()
}

func main() {
	controllers.StartMainServer()
}
