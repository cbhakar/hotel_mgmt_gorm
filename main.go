package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_task/api"
	"test_task/service"
)

func main() {
	r := mux.NewRouter()

	log.Println("server started at port : 8080")
	r.HandleFunc("/offer", api.AddDataHandler).Methods("POST")

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic("error starting server")
	}

	//handler process signal
	go HandleOSSignals(func() {
		err := service.Stop()
		if err != nil {
			panic(err)
		}
	})
	fmt.Println("closing app")
}

func HandleOSSignals(fn func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)

	for sig := range signals {
		switch sig {
		case syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM:
			fn()
		}
	}
}
