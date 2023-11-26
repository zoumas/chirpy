// Package database uses a JSON file as a database.
// You would always use a DBMS here.
// This is only for the purposes of the course.
package database

import (
	"encoding/json"
	"os"
	"sync"
)

// DB serves as a database using an underlying JSON file.
type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp       `json:"chirps"`
	Users         map[int]User        `json:"users"`
	RevokedTokens map[string]struct{} `json:"revoked_tokens"`
}

func NewDBStructure() DBStructure {
	return DBStructure{
		Chirps:        make(map[int]Chirp),
		Users:         make(map[int]User),
		RevokedTokens: make(map[string]struct{}),
	}
}

// New creates a new database file on the filesystem or truncates it if it already exists.
// This implementation always starts with an empty database.
func New(path string) (*DB, error) {
	_, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err = db.Persist(NewDBStructure())
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Load read the file from DB.path, unmarshalls it into JSON and returns a DBStructure
func (db *DB) Load() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbs := DBStructure{}
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return DBStructure{}, err
	}
	return dbs, nil
}

// Persist writes the JSON encoding of a given DBStructure to the file from DB.path
func (db *DB) Persist(dbs DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.MarshalIndent(dbs, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}
