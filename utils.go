package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	resp, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error while writing response :%w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	type errResponse struct {
		Error string `json:"error"`
	}

	resp := errResponse{msg}
	err := respondWithJSON(w, code, resp)
	if err != nil {
		return err
	}
	return nil
}
