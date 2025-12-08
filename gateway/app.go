package main

import (
	"RedRockMidAssessment/core/fruit"
	"RedRockMidAssessment/core/utils/banner"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	banner.Banner()

	fruit.GenesisFruit()

	<-quit
	log.Println("Server Starting to Shutdown...")

	fruit.WorldEndingFruit()

	log.Println("Server Shutdown, Bye!")
}
