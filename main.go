package main

import (
	"github.com/VolticFroogo/rush-hour-api/v1"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		ServerName: os.Getenv("SERVER_NAME"),
		Debug:      true,
	})
	if err != nil {
		log.Fatalf("Error initialising Sentry: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(10 * time.Second)
	defer sentry.Recover()

	r := gin.Default()
	r.Use(recoveryMiddleware)

	v1.Init(r)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}

func recoveryMiddleware(c *gin.Context) {
	defer sentry.Recover()
	c.Next()
}
