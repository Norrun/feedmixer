package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/Norrun/feedmixer/internal/wire"
)

func RenderError(renderer func(sw *wire.SnitchResponseWriter, r *http.Request), next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sw := wire.NewSnitchResponceWriter(w)
		defer func() {
			rec := recover()
			if rec != nil {
				log.Printf("panic occured in some handler: %v", rec)
				debug.PrintStack()
				if sw.Status() == 0 {
					w.WriteHeader(500)
				}
				if renderer != nil && sw.IsWritten() {
					renderer(sw, r)
				}
			}
		}()

		next.ServeHTTP(w, r)

	}
}
