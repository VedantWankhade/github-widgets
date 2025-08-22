package widgets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/vedantwankhade/github-widgets/internal/config"
)

func Health(params url.Values) (io.ReadSeeker, error) {
	app := config.GetApp()

	service := params.Get("service")
	if service == "" {
		app.Error("parameter 'service' not found in url")
		return nil, fmt.Errorf("please provide 'service' parameter for 'health' widget")
	}

	// TODO)) validate service url

	var msg string

	res, err := http.Get(service)
	if err != nil {
		app.Error(err.Error())
		msg = fmt.Sprintf("something went wrong with get url: %s", err.Error())
	} else {
		msg = fmt.Sprintf("%v", res.Status)
	}

	app.Info(msg)

	svg := fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="110" height="20" role="img" aria-label="build: passing">
	  <title>build: passing</title>
	  <rect width="50" height="20" fill="#555"/>   <!-- Left (label background) -->
	  <rect x="50" width="60" height="20" fill="#4c1"/> <!-- Right (status background) -->
	  
	  <!-- Divider -->
	  <rect x="50" width="1" height="20" fill="#000" fill-opacity="0.2"/>
	  
	  <!-- Left text -->
	  <text x="25" y="14"
			fill="#fff"
			font-family="Verdana,DejaVu Sans,sans-serif"
			font-size="11"
			text-anchor="middle">
		status
	  </text>
	  
	  <!-- Right text -->
	  <text x="80" y="14"
			fill="#fff"
			font-family="Verdana,DejaVu Sans,sans-serif"
			font-size="11"
			text-anchor="middle">
	   %s 
	  </text>
	</svg>
	`, msg)

	return bytes.NewReader([]byte(svg)), nil
}
