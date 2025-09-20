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

// ------------------users--------Radiator-------------------------//

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

func (p *Postgres) UpdateUser(id int, name, role, email, password string) (*types.User, error) {
	query := `
		UPDATE users
		SET name = $1,
		    role = $2,
		    email = $3,
		    password_hash = $4,
		    updated_at = NOW()
		WHERE id = $5
		RETURNING id, name, role, email, password_hash, created_at, updated_at
	`

	var updatedUser types.User
	err := p.db.QueryRow(query, name, role, email, password, id).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Role,
		&updatedUser.Email,
		&updatedUser.PasswordHash,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

func (p *Postgres) GetUsers() ([]types.User, error) {
	query := `SELECT id, name, role, email, password_hash, created_at, updated_at FROM users`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []types.User

	for rows.Next() {
		var u types.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Role,
			&u.Email,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}
func (p *Postgres) GetUserByID(id int) (*types.User, error) {
	query := `SELECT id, name, role, email, password_hash, created_at, updated_at
	          FROM users WHERE id = $1`

	var u types.User
	err := p.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Role,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &u, nil
}

func (p *Postgres) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	return nil
}

func (p *Postgres) GetUserByEmail(email string) (*types.User, error) {
	query := `SELECT id, name, role, email, password_hash, created_at, updated_at
	          FROM users
	          WHERE email = $1`

	var u types.User
	err := p.db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Role,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &u, nil
}

//------------------users--------Radiator-------------------------//

// -----------------products-------Radiator------------------------//
func (p *Postgres) CreateProduct(name, description, category, unit string) (*types.Product, error) {
	query := `
        INSERT INTO products (name, description, category, unit)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, description, category, unit, created_at, updated_at
    `

	var product types.Product
	err := p.db.QueryRow(query, name, description, category, unit).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Unit,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("could not create product: %w", err)
	}

	return &product, nil
}
func (p *Postgres) GetProducts() ([]types.Product, error) {
	query := `
        SELECT id, name, description, category, unit, created_at, updated_at
        FROM products
        ORDER BY id ASC
    `

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not fetch products: %w", err)
	}
	defer rows.Close()

	var products []types.Product

	for rows.Next() {
		var prod types.Product
		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Description,
			&prod.Category,
			&prod.Unit,
			&prod.CreatedAt,
			&prod.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan product: %w", err)
		}
		products = append(products, prod)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return products, nil
}
func (p *Postgres) GetProductById(id int) (*types.Product, error) {
	query := `
		SELECT id, name, description, category, unit, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product types.Product
	err := p.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Unit,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, fmt.Errorf("could not fetch product: %w", err)
	}

	return &product, nil
}
func (p *Postgres) CreateBoM(productID, componentID int, quantity float64, operationName string) (*types.BoM, error) {

	var exists bool
	err := p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking product existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("product with id %d does not exist", productID)
	}

	err = p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", componentID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking component existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("component product with id %d does not exist", componentID)
	}

	query := `
		INSERT INTO bom (product_id, component_id, quantity, operation_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, component_id, quantity, operation_name, created_at, updated_at
	`

	var bom types.BoM
	err = p.db.QueryRow(query, productID, componentID, quantity, operationName).Scan(
		&bom.ID,
		&bom.ProductID,
		&bom.ComponentID,
		&bom.Quantity,
		&bom.OperationName,
		&bom.CreatedAt,
		&bom.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create bom: %w", err)
	}

	return &bom, nil
}
func (p *Postgres) GetBoM(productID int) ([]types.BoM, error) {
	query := `
		SELECT id, product_id, component_id, quantity, operation_name, created_at, updated_at
		FROM bom
		WHERE product_id = $1
	`

	rows, err := p.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch BoM: %w", err)
	}
	defer rows.Close()

	var boms []types.BoM
	for rows.Next() {
		var bom types.BoM
		if err := rows.Scan(
			&bom.ID,
			&bom.ProductID,
			&bom.ComponentID,
			&bom.Quantity,
			&bom.OperationName,
			&bom.CreatedAt,
			&bom.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan bom row: %w", err)
		}
		boms = append(boms, bom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return boms, nil
}

//-----------------products-------Radiator------------------------//

//-----------------MO------------Radiator-------------------------//

//-----------------MO------------Radiator-------------------------//
