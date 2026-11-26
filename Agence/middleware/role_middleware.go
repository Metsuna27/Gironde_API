package middleware

import (
	"net/http"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			claims, ok := GetClaims(r.Context())

			if !ok || claims == nil {

				writeError(
					w,
					http.StatusUnauthorized,
					"Utilisateur non authentifié",
				)

				return
			}

			for _, role := range allowedRoles {

				if claims.Role == role {

					next.ServeHTTP(w, r)
					return
				}
			}

			writeError(
				w,
				http.StatusForbidden,
				"Accès interdit : permissions insuffisantes",
			)
		})
	}
}
