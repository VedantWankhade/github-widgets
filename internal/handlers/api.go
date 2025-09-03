package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vedantwankhade/github-widgets/internal/config"
	"github.com/vedantwankhade/github-widgets/internal/widgets"
)

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	app := config.GetApp()
	params := r.URL.Query()
	var f io.ReadSeeker
	var err error

	widgetName := params.Get("widget")
	if widgetName == "" {
		app.Error("parameter 'widget' not found in url")
		http.Error(w, "please provide 'widget' parameter in url", http.StatusNotFound)
		return
	}

	start := time.Now()

	switch widgetName {
	case "health":
		f, err = widgets.Health(params)
	case "commitgraph":
		f, err = widgets.CommitGraph(params)
	default:
		app.Error("invalid widget name")
		f, err = nil, fmt.Errorf("invalid widget name")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	app.Info(fmt.Sprintf("[Total time to serve request: %dms]", time.Since(start).Milliseconds()))
	http.ServeContent(w, r, "widget.svg", time.Now(), f)
}
