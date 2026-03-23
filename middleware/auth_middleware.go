package middleware

import (
	"net/http"
	"task-management/helper"
	"task-management/model/web"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{
		Handler: handler,
	}
}

func (middleware AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if "SecretKey" == request.Header.Get("X-API-Key") {
		middleware.Handler.ServeHTTP(writer, request)
	} else {
		writer.Header().Add("content-type", "application.json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code: http.StatusUnauthorized,
			Status: "Unauthorized",
		}

		helper.WriteToResponseBody(writer, webResponse)
	}
}