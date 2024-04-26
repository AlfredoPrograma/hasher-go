package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func handleError(err error, statusCode int, w http.ResponseWriter) {
	log.Printf("%4s %s", " ", err.Error())
	w.WriteHeader(statusCode)
}

func computeHash(content string) (string, int, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(content))

	if err != nil {
		return "", 0, err
	}

	hash := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash), hasher.Size(), nil
}

func main() {
	port := os.Args[1]

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Println("(HANDLER): GET /")
		w.Write([]byte("Hello world"))
	})

	http.HandleFunc("POST /encode", func(w http.ResponseWriter, r *http.Request) {
		var encodingPayload struct {
			Content string
		}

		log.Println("(HANDLER): POST /encode")

		rawBody, err := io.ReadAll(r.Body)

		if err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		if err = json.Unmarshal(rawBody, &encodingPayload); err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		hash, size, err := computeHash(encodingPayload.Content)

		if err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		encodingResponse := struct {
			Size int
			Hash string
		}{
			Size: size,
			Hash: hash,
		}

		rawResponse, err := json.Marshal(&encodingResponse)

		if err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		w.Write(rawResponse)
	})

	http.HandleFunc("POST /verify-hash", func(w http.ResponseWriter, r *http.Request) {
		log.Println("(HANDLER): POST /verify-hash")

		var verifyHashPayload struct {
			Content string
			Hash    string
		}

		rawBody, err := io.ReadAll(r.Body)

		if err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		if err = json.Unmarshal(rawBody, &verifyHashPayload); err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		hash, _, err := computeHash(verifyHashPayload.Content)

		if err != nil {
			handleError(err, http.StatusInternalServerError, w)
			return
		}

		if hash != verifyHashPayload.Hash {
			handleError(errors.New("hashes dont match"), http.StatusUnauthorized, w)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Verified"))
	})

	http.ListenAndServe(port, nil)
}
