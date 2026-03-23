package exception

import (
	"net/http"
	"task-management/helper"
	"task-management/model/web"
	"github.com/go-playground/validator/v10"
)

func notFoundError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Add("content-type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code: http.StatusNotFound,
			Status: "Not Found",
			Data: exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func validationError(writer http.ResponseWriter, _ *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Add("content-type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code: http.StatusBadRequest,
			Status: "Bad Request",
			Data: exception.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, _ *http.Request, err interface{}) {
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Code: http.StatusInternalServerError,
		Status: "Internal Server Error",
		Data: err,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if notFoundError(writer, request, err) {
		return
	}

	if validationError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}