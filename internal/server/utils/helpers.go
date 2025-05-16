package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type JSON map[string]any

func ParseJSONRequestBody(body io.ReadCloser, model interface{}) error {
	// Read the request body
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	// Validate the JSON
	if !json.Valid(data) {
		return fmt.Errorf("invalid json data")
	}
	// Decode into the model
	if err := json.Unmarshal(data, model); err != nil {
		return err
	}
	// Validate the model using the validator
	// TODO: improve validation messages
	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		return err
	}
	return nil
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, message JSON) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		jsonData = []byte(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(jsonData))
}
