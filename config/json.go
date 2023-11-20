package config

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ClientResponse struct {
	Rsp any `json:"rsp"`
}

type ErrorRsp struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Rsp     any    `json:"rsp,omitempty"`
}

func ReadBody(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one mg

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func ReadQuery(w http.ResponseWriter, r *http.Request, data any) error {
	return nil
}

func WriteResponse(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil

}

func ErrorResponse(w http.ResponseWriter, msg string, data any, status int) error {

	response := ErrorRsp{
		Error:   true,
		Message: msg,
		Rsp:     data,
	}

	return WriteResponse(w, status, response)
}
