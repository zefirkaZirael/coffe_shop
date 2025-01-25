package utils

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type Err struct {
	Err_Message string `json:"Error"`
	Status      int    `json:"Status code"`
}

func Err_Handler(w http.ResponseWriter, r *http.Request) {
	Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
}

func Log_Err_Handler(err error, status_Code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_Code)
	Err := Err{
		Err_Message: err.Error(),
		Status:      status_Code,
	}
	if err := json.NewEncoder(w).Encode(Err); err != nil {
		slog.Error("JSON encode error:", "Log_Err_Handler function", err)
	}
}
