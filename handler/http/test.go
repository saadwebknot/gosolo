package handler

import (
	"net/http"

	"github.com/s1s1ty/go-mysql-crud/driver"
	repository "github.com/s1s1ty/go-mysql-crud/repository"
	testrepo "github.com/s1s1ty/go-mysql-crud/repository/test"
)

// NewTestHandler wires dependencies for test routes.
func NewTestHandler(db *driver.DB) *Test {
	return &Test{
		repo: testrepo.NewSQLTestRepo(db.SQL),
	}
}

// Test exposes read-only handlers for the test table.
type Test struct {
	repo repository.TestRepo
}

// Fetch responds with all records stored in the test table.
func (t *Test) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := t.repo.FetchAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusOK, payload)
}
