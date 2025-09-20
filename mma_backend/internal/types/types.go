package types

import "time"

type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Description string    `json:"description,omitempty" db:"description"`
	Category    string    `json:"category,omitempty" db:"category"`
	Unit        string    `json:"unit" db:"unit" validate:"required,min=1,max=20"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Role         string    `json:"role" db:"role"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type WorkCenter struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	Capacity  int       `json:"capacity,omitempty" db:"capacity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ManufacturingOrder struct {
	ID                int        `json:"id" db:"id"`
	ProductID         int        `json:"product_id" db:"product_id"`
	Quantity          int        `json:"quantity" db:"quantity"`
	Status            string     `json:"status" db:"status"`
	StartDate         time.Time  `json:"start_date" db:"start_date"`
	DueDate           *time.Time `json:"due_date,omitempty" db:"due_date"`
	AssignedManagerID *int       `json:"assigned_manager_id,omitempty" db:"assigned_manager_id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

type BoM struct {
	ID            int       `json:"id" db:"id"`
	ProductID     int       `json:"product_id" db:"product_id"`
	ComponentID   int       `json:"component_id" db:"component_id"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	OperationName string    `json:"operation_name,omitempty" db:"operation_name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type WorkOrder struct {
	ID               int        `json:"id" db:"id"`
	MOID             int        `json:"mo_id" db:"mo_id"`
	StepName         string     `json:"step_name" db:"step_name"`
	Status           string     `json:"status" db:"status"`
	StartTime        *time.Time `json:"start_time,omitempty" db:"start_time"`
	EndTime          *time.Time `json:"end_time,omitempty" db:"end_time"`
	AssignedWorkerID *int       `json:"assigned_worker_id,omitempty" db:"assigned_worker_id"`
	WorkCenterID     *int       `json:"work_center_id,omitempty" db:"work_center_id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}
type Inventory struct {
	ID             int       `json:"id" db:"id"`
	ProductID      int       `json:"product_id" db:"product_id"`
	MovementType   string    `json:"movement_type" db:"movement_type"`
	Quantity       float64   `json:"quantity" db:"quantity"`
	Date           time.Time `json:"date" db:"date"`
	ReferenceType  *string   `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID    *int      `json:"reference_id,omitempty" db:"reference_id"`
	CurrentBalance float64   `json:"current_balance" db:"current_balance"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
