package middleware

import "net/http"

// responseWithError(w, 403, fmt.Spintf(Auth error: %v, err))

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				err, ok := e.(error)
				if !ok {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal server error"))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})

}
