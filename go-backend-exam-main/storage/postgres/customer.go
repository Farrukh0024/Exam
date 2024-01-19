package postgres

import (
	"database/sql"
	"fmt"

	"city2city/api/models"
	"github.com/google/uuid"
)

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) customerRepo {
	return customerRepo{
		db: db,
	}
}

func (c customerRepo) Create(customer models.CreateCustomer) (string, error) {
	uid := uuid.New().String()

	if _, err := c.db.Exec("INSERT INTO customers (id, full_name, phone, email) VALUES ($1, $2, $3, $4)",
		uid,
		customer.FullName,
		customer.Phone,
		customer.Email,
	); err != nil {
		fmt.Println("error while inserting data ", err.Error())
		return "", err
	}
	return "", nil
}

func (c customerRepo) Get(id string) (models.Customer, error) {
	query := `
        SELECT id, full_name, phone, email, created_at
        FROM customers
        WHERE id = $1
    `

	row := c.db.QueryRow(query, id)
	var customer models.Customer
	if err := row.Scan(&customer.ID, &customer.FullName, &customer.Phone, &customer.Email, &customer.CreatedAt); err != nil {
		return models.Customer{}, fmt.Errorf("error getting customer: %w", err)
	}

	return customer, nil
}

func (c customerRepo) GetList(req models.GetListRequest) (models.CustomersResponse, error) {
	query := `
        SELECT id, full_name, phone, email, created_at
        FROM customers
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := c.db.Query(query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return models.CustomersResponse{}, fmt.Errorf("error getting customer list: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.FullName, &customer.Phone, &customer.Email, &customer.CreatedAt); err != nil {
			return models.CustomersResponse{}, fmt.Errorf("error scanning customer: %w", err)
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return models.CustomersResponse{}, fmt.Errorf("error iterating customers: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM customers`
	var count int
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		return models.CustomersResponse{}, fmt.Errorf("error getting customer count: %w", err)
	}

	return models.CustomersResponse{
		Customers: customers,
		Count:     count,
	}, nil
}

func (c customerRepo) Update(customer models.Customer) (string, error) {
	query := `
        UPDATE customers
        SET full_name = $2, phone = $3, email = $4
        WHERE id = $1
        RETURNING id
    `

	row := c.db.QueryRow(query, customer.ID, customer.FullName, customer.Phone, customer.Email)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("error updating customer: %w", err)
	}

	return id, nil
}

func (c customerRepo) Delete(id string) error {
	query := `
        DELETE FROM customers
        WHERE id = $1
    `

	if _, err := c.db.Exec(query, id); err != nil {
		return fmt.Errorf("error deleting customer: %w", err)
	}

	return nil
}
