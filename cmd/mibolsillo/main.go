package main

import (
	"context"
	"flag"
	"log"
	"mibolsillo/api"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/mux"
)


func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*5, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", api.CreateInvoice)
// Agregar mas handlings


	srv := &http.Server{
		Addr: "0.0.0.0:5000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
