package handlers

import (
	"net/http"
	"time"

	"github.com/vedantwankhade/gh-readme-cards/internal/config"
	"github.com/vedantwankhade/gh-readme-cards/internal/services"
)

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	app := config.GetApp()
	q := r.URL.Query()
	cardName := q.Get("card")
	img, err := services.GetCard(cardName)
	if err != nil {
		app.Error(err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	http.ServeContent(w, r, "card.png", time.Now(), img)
}
