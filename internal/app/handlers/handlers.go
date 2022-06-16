package handlers

import (
	"echo-http/internal/app/processor"
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

func (h *Handler) JsonAllInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.proc.Do(r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}
	WrapOK(w, info)
}
