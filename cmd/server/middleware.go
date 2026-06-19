package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/keyring"
)

func injectCacheMiddleware(cache *sync.Map, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), datautils.MakeKeyTo[*sync.Map](keyring.Cache), cache)
		nr := r.WithContext(ctx)
		next.ServeHTTP(w, nr)
	})
}
