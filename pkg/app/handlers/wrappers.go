package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WrapErrorWithStatus(w http.ResponseWriter, err error, httpStatus int) {
	var m = map[string]string{
		"result": "error",
		"data":   err.Error(),
	}

	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)
	fmt.Fprintln(w, string(res))
}

func WrapOK(w http.ResponseWriter, m interface{}) {
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(res))
}
