package db

import "database/sql"

// Repository defines data access.
type Repository interface {
	GetUser(id string) (string, error)
	SaveUser(id, name string) error
}

// PostgresRepo implements Repository.
type PostgresRepo struct {
	db *sql.DB
}

func (p *PostgresRepo) GetUser(id string) (string, error) {
	// raw SQL query for testing
	row := p.db.QueryRow("SELECT name FROM users WHERE id = $1", id)
	var name string
	err := row.Scan(&name)
	return name, err
}

func (p *PostgresRepo) SaveUser(id, name string) error {
	_, err := p.db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", id, name)
	return err
}

// MemoryRepo is a mock repository.
type MemoryRepo struct {
	data map[string]string
}

func (m *MemoryRepo) GetUser(id string) (string, error) {
	return m.data[id], nil
}

func (m *MemoryRepo) SaveUser(id, name string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[id] = name
	return nil
}
