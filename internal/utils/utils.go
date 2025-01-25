package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func Help() {
	i := `./frappuccino --help
Coffee Shop Management System

Usage:
  frappuccino [--port <N>] [--dir <S>] 
  frappuccino --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`
	fmt.Println(i)
	os.Exit(0)
}

func Send_Request(mystruct any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mystruct); err != nil {
		slog.Error("JSON encode error:", "Log_Err_Handler function ", err)
		return err
	}
	return nil
}

func CheckPort() error {
	num, err := strconv.Atoi(*models.Port)
	if err != nil {
		return err
	}
	if num < 1024 || num > 49151 {
		return errors.New("port number range:1024-49151")
	}
	return nil
}
