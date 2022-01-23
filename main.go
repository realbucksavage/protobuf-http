package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"protobuf-http/greeter"
	"sync"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	var (
		svc  = greeter.NewService()
		port = flag.Int("port", 8080, "specify port to start the http server on")
	)
	flag.Parse()

	var wg sync.WaitGroup
	shutdown := make(chan bool, 1)

	var handler http.Handler
	{
		router := mux.NewRouter()
		greeter.AddRoutes(router, svc)

		c := cors.New(cors.Options{AllowedOrigins: []string{"http://localhost:3000"}, AllowCredentials: true})
		handler = c.Handler(router)
		handler = handlers.CustomLoggingHandler(os.Stdout, handler, customResponseLogger)
	}

	httpServer := http.Server{Handler: handler}
	{
		go func() {
			lis, err := net.Listen("tcp", fmt.Sprintf(`:%d`, *port))
			if err != nil {
				panic(err)
			}

			log.Printf("listening on %s", lis.Addr().String())
			wg.Add(1)
			if err := httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()

		go func() {
			<-shutdown
			log.Printf("server shutdown result: %v", httpServer.Shutdown(context.Background()))
			wg.Done()
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("received shutdown signal %v", <-sig)
	close(shutdown)
	wg.Wait()
}
