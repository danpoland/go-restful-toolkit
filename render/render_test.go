package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	body := map[string]string{"key": "val"}
	_ = OK(w, body)

	r := w.Result()

	if r.StatusCode != 200 {
		t.Errorf("Expected 200, got '%v'", r.StatusCode)
	}
}

func TestAccepted(t *testing.T) {
	w := httptest.NewRecorder()
	_ = Accepted(w)

	r := w.Result()

	if r.StatusCode != 202 {
		t.Errorf("Expected 202, got '%v'", r.StatusCode)
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	_ = NoContent(w)

	r := w.Result()

	if r.StatusCode != 204 {
		t.Errorf("Expected 204, got '%v'", r.StatusCode)
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	body := map[string]string{"key": "val"}
	_ = BadRequest(w, body)

	r := w.Result()

	if r.StatusCode != 400 {
		t.Errorf("Expected 400, got '%v'", r.StatusCode)
	}
}

func TestUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	body := map[string]string{"key": "val"}
	_ = UnprocessableEntity(w, body)

	r := w.Result()

	if r.StatusCode != 422 {
		t.Errorf("Expected 422, got '%v'", r.StatusCode)
	}
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	_ = Unauthorized(w)

	r := w.Result()

	if r.StatusCode != 401 {
		t.Errorf("Expected 401, got '%v'", r.StatusCode)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	_ = NotFound(w)

	r := w.Result()

	if r.StatusCode != 404 {
		t.Errorf("Expected 404, got '%v'", r.StatusCode)
	}
}

func TestServerError(t *testing.T) {
	w := httptest.NewRecorder()
	_ = ServerError(w)

	r := w.Result()

	if r.StatusCode != 500 {
		t.Errorf("Expected 500, got '%v'", r.StatusCode)
	}
}

func TestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	code := http.StatusOK
	body := map[string]string{"key": "val"}

	_ = Response(w, code, body)
	r := w.Result()

	if w.Code != code {
		t.Errorf("Expected 200, got '%v'", r.StatusCode)
	}

	expectedBody := "{\"key\":\"val\"}"
	rBody := w.Body.String()
	if rBody != expectedBody {
		t.Errorf("Expected body %v, got %v", expectedBody, rBody)
	}
}
