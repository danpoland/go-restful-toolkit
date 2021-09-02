package render

import (
	"encoding/json"
	"net/http"
)

// OK renders a 200 with the provided body.
func OK(w http.ResponseWriter, body interface{}) error {
	return Response(w, http.StatusOK, body)
}

// Accepted renders a 202 response with an empty body.
func Accepted(w http.ResponseWriter) error {
	return Response(w, http.StatusAccepted, nil)
}

// NoContent renders a 204 with an empty body.
func NoContent(w http.ResponseWriter) error {
	return Response(w, http.StatusNoContent, nil)
}

// BadRequest renders a 400 response with the provided body.
func BadRequest(w http.ResponseWriter, body interface{}) error {
	return Response(w, http.StatusBadRequest, body)
}

// UnprocessableEntity renders a 422 response with the provided body.
func UnprocessableEntity(w http.ResponseWriter, body interface{}) error {
	return Response(w, http.StatusUnprocessableEntity, body)
}

// Unauthorized renders a 401 with the provided body.
func Unauthorized(w http.ResponseWriter) error {
	return Response(w, http.StatusUnauthorized, nil)
}

// NotFound renders a 404 response.
func NotFound(w http.ResponseWriter) error {
	return Response(w, http.StatusNotFound, nil)
}

// ServerError renders a 500 response.
func ServerError(w http.ResponseWriter) error {
	return Response(w, http.StatusInternalServerError, nil)
}

// Response writes a http response using the value passed in body as JSON.
// If it cannot convert the value to JSON, it returns an error.
func Response(w http.ResponseWriter, statusCode int, body interface{}) error {
	w.WriteHeader(statusCode)

	if body == nil {
		return nil
	}

	js, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}