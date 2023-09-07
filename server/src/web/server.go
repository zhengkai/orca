package web

import (
	"net/http"
	"project/config"
	"project/core"
	"project/vertexai"
	"project/zj"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server ...
func Server() {

	mux := http.NewServeMux()

	mux.Handle(`/_metrics`, promhttp.Handler())
	mux.Handle(`/v1/moderations`, core.NewCore())
	mux.Handle(`/v1/completions`, core.NewCore())
	mux.Handle(`/v1/chat/completions`, core.NewCore())
	mux.Handle(`/va/chat`, vertexai.WebChat)
	mux.Handle(`/va/text`, vertexai.WebText)

	s := &http.Server{
		Addr:         config.WebAddr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	zj.J(`start web server`, config.WebAddr)

	s.ListenAndServe()
}
