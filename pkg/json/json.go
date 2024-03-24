package json

import (
	"encoding/json"
	"net/http"
)

func HttpError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	JsonResponse(w, code, ErrorResponse{
		Error: msg,
	})
}

func JsonResponse(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Invalid JSON"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func DecodeJSONRequest(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(target)
}
