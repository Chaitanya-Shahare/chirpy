package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chrips map[int]Chirp `json:"chirps"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

var (
	ErrInvalidChirp = errors.New("invalid chirp")
)

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.Stat(db.path)

	if os.IsNotExist(err) {
		dbFile, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer dbFile.Close()
		initialDB := DBStructure{Chrips: make(map[int]Chirp)}
		return db.writeDB(initialDB)
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStruct DBStructure
	if err := json.Unmarshal(data, &dbStruct); err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.MarshalIndent(dbStruct, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	if body == "" {
		return Chirp{}, ErrInvalidChirp
	}

	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	nextID := 1

	if len(dbStruct.Chrips) > 0 {
		for id := range dbStruct.Chrips {
			if id >= nextID {
				nextID = id + 1
			}
		}
	}

	chirp := Chirp{
		ID:   nextID,
		Body: body,
	}

	dbStruct.Chrips[nextID] = chirp

	if err := db.writeDB(dbStruct); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chrips))

	for _, chirp := range dbStruct.Chrips {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	return chirps, nil
}
