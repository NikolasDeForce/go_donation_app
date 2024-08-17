package main

import (
	"donation/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

var port = ":8010"

var rMux = mux.NewRouter()

func main() {
	arguments := os.Args
	if len(arguments) >= 2 {
		port = ":" + arguments[1]
	}

	NotAllowed := handlers.NotAllowedHandler{}
	rMux.MethodNotAllowedHandler = NotAllowed

	rMux.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	getMux := rMux.Methods(http.MethodGet).Subrouter()

	getMux.HandleFunc("/", handlers.MainHandler)
	getMux.HandleFunc("/donation", handlers.DonationHanler)

	//API
	getMux.HandleFunc("/register/{login}/{password}", handlers.GetTokenHandler)
	getMux.HandleFunc("/api/{token}/donates", handlers.GetDonatesHandler)

	s := &http.Server{
		Addr:         port,
		Handler:      rMux,
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
