package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"city2city/api/models"
	"city2city/storage"
	"github.com/google/uuid"
)

type tripCustomerRepo struct {
	db *sql.DB
}

func NewTripCustomerRepo(db *sql.DB) storage.ITripCustomerRepo {
	return &tripCustomerRepo{
		db: db,
	}
}

func (c *tripCustomerRepo) Create(req models.CreateTripCustomer) (string, error) {
	// Generate a new UUID
	uid := uuid.New().String()

	// Prepare the SQL query with a placeholder for the UUID
	query := `INSERT INTO trip_customers (id,trip_id, customer_id) VALUES ($1, $2, $3) RETURNING id`

	// Execute the query, passing the generated UUID as a parameter
	rows, err := c.db.Query(
		query,
		uid,
		req.TripID,
		req.CustomerID,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// No need to scan for the ID as it's already generated and returned
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", err
		}
		return "", errors.New("failed to create city")
	}

	return uid, nil
}

func (c *tripCustomerRepo) Get(id string) (models.TripCustomer, error) {
	query := `
        SELECT id, trip_id, customer_id, created_at
        FROM trip_customers
        WHERE id = $1
    `
	row := c.db.QueryRow(query, id)
	var tripCustomer models.TripCustomer
	if err := row.Scan(&tripCustomer.ID, &tripCustomer.TripID, &tripCustomer.CustomerID, &tripCustomer.CustomerData, &tripCustomer.CreatedAt); err != nil {
		return models.TripCustomer{}, fmt.Errorf("failed to get trip customer: %w", err)
	}
	return tripCustomer, nil
}

func (c *tripCustomerRepo) GetList(req models.GetListRequest) (models.TripCustomersResponse, error) {
	query := `
        SELECT id, trip_id, customer_id, created_at
        FROM trip_customers
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
	rows, err := c.db.Query(query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return models.TripCustomersResponse{}, fmt.Errorf("failed to get trip customers list: %w", err)
	}
	defer rows.Close()

	var tripCustomers []models.TripCustomer
	for rows.Next() {
		var tripCustomer models.TripCustomer
		if err := rows.Scan(&tripCustomer.ID, &tripCustomer.TripID, &tripCustomer.CustomerID, &tripCustomer.CustomerData, &tripCustomer.CreatedAt); err != nil {
			return models.TripCustomersResponse{}, fmt.Errorf("failed to scan trip customer: %w", err)
		}
		tripCustomers = append(tripCustomers, tripCustomer)
	}

	if err := rows.Err(); err != nil {
		return models.TripCustomersResponse{}, fmt.Errorf("failed to iterate trip customers: %w", err)
	}

	countQuery := `
        SELECT COUNT(*) FROM trip_customers
    `
	row := c.db.QueryRow(countQuery)
	var count int
	if err := row.Scan(&count); err != nil {
		return models.TripCustomersResponse{}, fmt.Errorf("failed to count trip customers: %w", err)
	}

	return models.TripCustomersResponse{
		TripCustomers: tripCustomers,
		Count:         count,
	}, nil
}

func (c *tripCustomerRepo) Update(req models.TripCustomer) (string, error) {
	query := `
        UPDATE trip_customers
        SET trip_id = $1, customer_id = $2
        WHERE id = $3
        RETURNING id
    `
	row := c.db.QueryRow(query, req.TripID, req.CustomerID, req.ID)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("failed to update trip customer: %w", err)
	}
	return id, nil
}

func (c *tripCustomerRepo) Delete(id string) error {
	query := `
        DELETE FROM trip_customers
        WHERE id = $1
    `
	_, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trip customer: %w", err)
	}
	return nil
}
