package main

import (
	"fmt"
	"net/http"

	"github.com/vedantwankhade/github-widgets/internal/config"
	"github.com/vedantwankhade/github-widgets/internal/handlers"
)

func main() {
	app := config.GetApp()
	app.Info("App initialized")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /get", handlers.GenericHandler)

	server := http.Server{
		Addr:     ":" + app.Port,
		Handler:  mux,
		ErrorLog: app.GetErrLogger(),
	}
	app.Info(fmt.Sprintf("Server started on port %s", app.Port))
	app.Fatal(server.ListenAndServe())
}
