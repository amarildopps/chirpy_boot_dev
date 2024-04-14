package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amarildopps/chirpy_boot_dev/internal/database"
)

func (db *databaseConfig) addChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if params.Body == "" {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 500, "Chirp is too long")
		return
	}

	chirpyResult := database.Chirp{}
	fmt.Println("Body:", params.Body)
	chirpyResult, err = db.DB.CreateChirp(params.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 201, chirpyResult)

}

func (db *databaseConfig) getChirpys(w http.ResponseWriter, r *http.Request) {

	chirps, err := db.DB.GetChirps()
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, chirps)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}
