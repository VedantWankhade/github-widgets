package handlers

import (
	"net/http"
	"time"

	"github.com/vedantwankhade/github-widgets/internal/services"
)

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	img, err := services.GetCard(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeContent(w, r, "card.png", time.Now(), img)
}
