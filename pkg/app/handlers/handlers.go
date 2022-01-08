package handlers

import (
	"echo-http/pkg/app/processor"
	"encoding/json"
	"net/http"
)

type Handler struct {
	proc *processor.Processor
}

func NewHandler(proc *processor.Processor) *Handler {
	return &Handler{
		proc: proc,
	}
}

func (h *Handler) JsonData(w http.ResponseWriter, r *http.Request) {
	Req, err := h.proc.GetInfo(r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}
	JsonReq, err := json.Marshal(Req)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(JsonReq)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}
}
