package app

import (
	"log"
	"net/http"
)

type App struct {
	server *http.Server
}

func New(port string) *App {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return &App{
		server: server,
	}
}

func (a *App) Run() error {
	log.Println("server started on", a.server.Addr)
	return a.server.ListenAndServe()
}
