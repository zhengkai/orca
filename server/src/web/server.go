package web

import (
	"net/http"
	"project/config"
	"project/core"
	"project/zj"
	"time"
)

// Server ...
func Server() {

	mux := http.NewServeMux()

	mux.HandleFunc(`/`, core.NewCore().WebHandle)

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
