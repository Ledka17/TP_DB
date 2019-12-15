package handler

import (
	"encoding/json"
	"net/http"
)

func (h *DataBaseHandler) ClearDBHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.usecase.CleanUp()
		w.WriteHeader(200)
	}
	w.WriteHeader(400)
}

func (h *DataBaseHandler) GetServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		body, err := json.Marshal(h.usecase.GetStatusDB())
		checkErr(err)

		w.WriteHeader(200)
		w.Write(body)
		return
	}
	w.WriteHeader(400)
}
