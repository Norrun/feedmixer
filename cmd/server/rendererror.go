package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func RenderError(renderer func(w http.ResponseWriter, r *http.Request), next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//sw := wire.SnitchResponseWriter{}
		defer func() {
			rec := recover()
			if rec != nil {
				log.Printf("panic occured in some handler: %v", rec)
				debug.PrintStack()

				w.WriteHeader(500)
				if renderer != nil {
					renderer(w, r)
				}
			}
		}()

		next.ServeHTTP(w, r)

	}
}
