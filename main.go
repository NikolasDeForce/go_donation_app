package main

import (
	"donation/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var port = ":8010"

var mux = http.NewServeMux()

func main() {
	arguments := os.Args
	if len(arguments) >= 2 {
		port = ":" + arguments[1]
	}

	mux.Handle("/error", http.HandlerFunc(handlers.MethodNotAllowedHandler))

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.Handle("/", http.HandlerFunc(handlers.MainHandler))

	mux.Handle("/donation/static/", http.StripPrefix("/static", fs))
	mux.Handle("/donation", http.HandlerFunc(handlers.DonationHanler))

	s := &http.Server{
		Addr:         port,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	go func() {
		log.Println("Listening to", port)
		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	log.Println("Quitting after signal:", sig)
	time.Sleep(5 * time.Second)
	s.Shutdown(nil)
}
