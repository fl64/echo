package handlers

import (
	"echo/internal/app-http/processor"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type Handler struct {
	proc       *processor.Processor
	respStatus *atomic.Int32
}

func NewHandler(proc *processor.Processor, respStatus *atomic.Int32) *Handler {
	return &Handler{
		proc:       proc,
		respStatus: respStatus,
	}
}

func (h *Handler) JsonAllInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.proc.Do(r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusInternalServerError)
		return
	}
	WrapOK(w, info, int(h.respStatus.Load()))
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
