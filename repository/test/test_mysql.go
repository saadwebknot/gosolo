package testrepo

import (
	"context"
	"database/sql"

	"github.com/s1s1ty/go-mysql-crud/models"
	repository "github.com/s1s1ty/go-mysql-crud/repository"
)

// NewSQLTestRepo returns a mysql-backed TestRepo implementation.
func NewSQLTestRepo(conn *sql.DB) repository.TestRepo {
	return &mysqlTestRepo{Conn: conn}
}

type mysqlTestRepo struct {
	Conn *sql.DB
}

func (m *mysqlTestRepo) FetchAll(ctx context.Context) ([]*models.Test, error) {
	query := "SELECT id, name, age FROM test"

	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*models.Test, 0)
	for rows.Next() {
		record := new(models.Test)
		if err := rows.Scan(
			&record.ID,
			&record.Name,
			&record.Age,
		); err != nil {
			return nil, err
		}
		payload = append(payload, record)
	}

	return payload, rows.Err()
}
