package main

import (
	"io"
	"log"
	"time"

	"github.com/gorilla/handlers"
)

func customResponseLogger(w io.Writer, params handlers.LogFormatterParams) {
	log.Printf("%s %s (Size %d) took %v", params.Request.Method, params.Request.RequestURI, params.Size, time.Since(params.TimeStamp))
}
