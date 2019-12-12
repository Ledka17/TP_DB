package handler

import (
	"encoding/json"
	"net/http"
)

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
