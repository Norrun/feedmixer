package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/Norrun/feedmixer/internal/wire"
)

func RenderError(renderer func(aw *wire.ApproveResponseWriter, r *http.Request), next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aw := wire.NewApprovalWriter(w, func(arw *wire.ApproveResponseWriter) bool {
			return arw.Status() < 400 && arw.Status() > 0
		})
		defer func() {
			rec := recover()
			if rec != nil {
				log.Printf("panic occured in some handler: %v", rec)
				debug.PrintStack()
				if aw.Approved() {
					return
				}
				aw.Approve()
				aw.WriteHeader(500)
				renderer(aw, r)
			}
		}()
		next.ServeHTTP(aw, r)
		if aw.Approved() {
			return
		}
		if renderer != nil {
			renderer(aw, r)
		}

	}
}
