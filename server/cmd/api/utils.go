package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

/** Methodes de manipulation des fichiers JSON pour eviter la duplication recurrente de code. */
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // pour ne pas inclure si pas de data
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction variable avec possiblement plusieurs headers tous de type http.Header permettant d'ecrire
un 'data' dans un JSON. */
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
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
		log.Printf("writeJSON not ok")
		return err
	}

	return nil
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant de lire un JSON dans un 'data'. */
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {

	maxBytes := 1024 * 1024 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return fmt.Errorf("decode error: %w", syntaxErr)
		}
		return fmt.Errorf("decode error: %w", err)
	}

	// to check if more than a one JSON file in the request body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only containt a single json value")
	}

	return nil
}

/*--------------------------------------------------------------------------------------------------------*/
/** Fonction permettant de donner une erreur 'err'. */
func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest // default iff status not specified

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
/*--------------------------------------------------------------------------------------------------------*/
