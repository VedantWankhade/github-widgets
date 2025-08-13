package main

import (
	"fmt"
	"net/http"

	"github.com/vedantwankhade/gh-readme-cards/internal/config"
	"github.com/vedantwankhade/gh-readme-cards/internal/handlers"
)

func main() {
	app := config.InitApp()
	app.Info("App initialized")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", handlers.GenericHandler)

	server := http.Server{
		Addr:     ":" + app.Port,
		Handler:  mux,
		ErrorLog: app.GetErrLogger(),
	}
	app.Info(fmt.Sprintf("Server started on port %s", app.Port))
	app.Fatal(server.ListenAndServe())
}
