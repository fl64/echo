package handlers

import (
	"echo/internal/app-http/processor"
	"net/http"
	"strconv"
	"strings"
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

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	l := r.FormValue("len")
	if l == "" {
		l = "1"
	}
	n, err := strconv.Atoi(l)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(strings.Repeat(strings.Repeat("#", 1), n)))
	w.Header().Set("Content-Type", "application/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
