package handler

import (
	"net/http"
)

func (h *DataBaseHandler) ClearDBHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.usecase.CleanUp()
		w.WriteHeader(200)
	}
	w.WriteHeader(400)
}
