package bridges

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vinodhalaharvi/purekernels/pkg/result"
)

// DecodeJSON decodes JSON from reader into type T
func DecodeJSON[T any](r io.Reader) result.Result[T] {
	var value T
	if err := json.NewDecoder(r).Decode(&value); err != nil {
		return result.Err[T](err)
	}
	return result.Ok(value)
}

// DecodeJSONRequest decodes JSON from HTTP request body
func DecodeJSONRequest[T any](r *http.Request) result.Result[T] {
	return DecodeJSON[T](r.Body)
}

// EncodeJSON encodes value to JSON
func EncodeJSON[T any](w io.Writer, value T) result.Result[struct{}] {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return result.Err[struct{}](err)
	}
	return result.Ok(struct{}{})
}

// RespondJSON sends JSON response
func RespondJSON[T any](w http.ResponseWriter, status int, value T) result.Result[struct{}] {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return EncodeJSON(w, value)
}

// RespondError sends JSON error response
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}
