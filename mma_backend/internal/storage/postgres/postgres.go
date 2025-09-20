package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"mma_api/internal/config"
	"mma_api/internal/types"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	connStr := cfg.Conn_Str
	log.Println("Connecting to PostgreSQL with:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %w", err)
	}
	log.Println("Connected to PostgreSQL successfully!")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        role VARCHAR(50) NOT NULL CHECK (role IN ('manager', 'worker', 'inventory_manager', 'admin')),
        email VARCHAR(150) UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,

		`CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        category VARCHAR(50),
        unit VARCHAR(20) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,

		`CREATE TABLE IF NOT EXISTS work_centers (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        type VARCHAR(50) NOT NULL CHECK (type IN ('machine','team','location')),
        capacity INT,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,

		`CREATE TABLE IF NOT EXISTS manufacturing_orders (
        id SERIAL PRIMARY KEY,
        product_id INT NOT NULL,
        quantity INT NOT NULL,
        status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'in_progress', 'done')),
        start_date DATE NOT NULL,
        due_date DATE,
        assigned_manager_id INT,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,

		`CREATE TABLE IF NOT EXISTS bom (
        id SERIAL PRIMARY KEY,
        product_id INT NOT NULL,
        component_id INT NOT NULL,
        quantity DECIMAL(10,2) NOT NULL,
        operation_name VARCHAR(100),
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,

		`CREATE TABLE IF NOT EXISTS work_orders (
        id SERIAL PRIMARY KEY,
        mo_id INT NOT NULL,
        step_name VARCHAR(50) NOT NULL,
        status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'in_progress', 'completed')),
        start_time TIMESTAMP,
        end_time TIMESTAMP,
        assigned_worker_id INT,
        work_center_id INT,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,

		`CREATE TABLE IF NOT EXISTS inventory (
        id SERIAL PRIMARY KEY,
        product_id INT NOT NULL,
        movement_type VARCHAR(10) NOT NULL CHECK (movement_type IN ('IN','OUT')),
        quantity DECIMAL(10,2) NOT NULL,
        date TIMESTAMP DEFAULT NOW(),
        reference_type VARCHAR(20),
        reference_id INT,
        current_balance DECIMAL(10,2),
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("could not execute query: %w", err)
		}
	}

	log.Println("All tables are ready!")
	return &Postgres{db: db}, nil

}

func (p *Postgres) CreateUser(name, role, email, password string) (*types.User, error) {
	query := `
		INSERT INTO users (name, role, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, name, role, email, password_hash, created_at, updated_at
	`

	var newUser types.User
	err := p.db.QueryRow(query, name, role, email, password).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.Role,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}
