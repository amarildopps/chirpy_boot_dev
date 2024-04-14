package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type DB struct {
	Path string
	Mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

	db := &DB{
		Path: path,
		Mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return db, nil

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File doens't exists")
		file, err := os.Create(db.Path)
		if err != nil {
			return err
		}
		file.Close()
	} else if err != nil {
		return err
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

	// Initialize the DBStructure
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
	}

	// Read the current database file
	file, err := os.Open(db.Path)
	if err != nil {
		return DBStructure{}, err
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&dbStructure)
		if err != nil {
			if err == io.EOF {
				// If EOF, it means the file is empty; initialize empty structure
				dbStructure.Chirps = make(map[int]Chirp)
			} else {
				return DBStructure{}, err
			}
		}
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

	file, err := os.Create(db.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&dbStructure)
	if err != nil {
		return err
	}
	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {

	// Lock the mutex to ensure thread safety
	db.Mux.Lock()
	defer db.Mux.Unlock()

	// Initialize the DBStructure
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := 1
	for id := range dbStructure.Chirps {
		if id >= newId {
			newId = id + 1
		}
	}

	newChirp := Chirp{
		Id:   newId,
		Body: body,
	}

	dbStructure.Chirps[newId] = newChirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	var chirpSlice []Chirp
	for _, chirp := range dbStructure.Chirps {
		chirpSlice = append(chirpSlice, chirp)
	}

	return chirpSlice, nil
}
