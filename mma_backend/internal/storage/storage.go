package storage

import "mma_api/internal/types"

type Storage interface {
	CreateUser(name, role, email, password string) (*types.User, error)
	UpdateUser(id int, name, role, email, password string) (*types.User, error)
	GetUsers() ([]types.User, error)
	GetUserByID(id int) (*types.User, error)
	DeleteUser(id int) error
	GetUserByEmail(email string) (*types.User, error)
	CreateProduct(name, description, category, unit string) (*types.Product, error)
	GetProductById(id int) (*types.Product, error)
	CreateBoM(productID, componentID int, quantity float64, operationName string) (*types.BoM, error)
}
