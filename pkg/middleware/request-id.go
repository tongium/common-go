package middleware

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/tongium/common-go/pkg/constant"
)

const defaultHeaderKey string = "X-Request-ID"

func randomString(size int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, size)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

func generator() string {
	return randomString(32)
}

func RequestIDMiddleware(headerKey string) func(next http.Handler) http.Handler {
	if headerKey == "" {
		headerKey = defaultHeaderKey
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(headerKey)
			if rid == "" {
				rid = generator()
			}

			ctx := context.WithValue(r.Context(), constant.RequestIDContextKey, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
