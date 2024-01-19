package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"city2city/api/models"
	"city2city/storage"
	"github.com/google/uuid"
)

type cityRepo struct {
	db *sql.DB
}

func NewCityRepo(db *sql.DB) storage.ICityRepo {
	return cityRepo{db: db}
}

func (c cityRepo) Create(city models.CreateCity) (string, error) {
	// Generate a new UUID
	cityID := uuid.New().String()

	// Prepare the SQL query with a placeholder for the UUID
	query := `INSERT INTO cities (id, name) VALUES ($1, $2) RETURNING id`

	// Execute the query, passing the generated UUID as a parameter
	rows, err := c.db.Query(query, cityID, city.Name)
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

	return cityID, nil
}

func (c cityRepo) Get(id string) (models.City, error) {
	var city models.City
	err := c.db.QueryRow("SELECT id, name, created_at FROM cities WHERE id = $1", id).Scan(&city.ID, &city.Name, &city.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("error while getting city", err.Error())
			return models.City{}, err
		}
		return models.City{}, err
	}

	return city, nil
}

func (c cityRepo) GetList(req models.GetListRequest) (models.CitiesResponse, error) {
	limit := req.Limit
	offset := (req.Page - 1) * limit

	rows, err := c.db.Query(
		`SELECT id, name, created_at FROM cities ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return models.CitiesResponse{}, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.ID, &city.Name, &city.CreatedAt)
		if err != nil {
			return models.CitiesResponse{}, err
		}
		cities = append(cities, city)
	}

	count, err := c.countCities()
	if err != nil {
		return models.CitiesResponse{}, err
	}

	return models.CitiesResponse{Cities: cities, Count: count}, nil
}

func (c cityRepo) Update(city models.City) (string, error) {
	result, err := c.db.Exec("UPDATE cities SET name = $1 WHERE id = $2", city.Name, city.ID)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		fmt.Println("error while updating city", err.Error())
		return "", err
	}

	return city.ID, nil
}

func (c cityRepo) Delete(id string) error {
	result, err := c.db.Exec("DELETE FROM cities WHERE id = $1", id)
	if err != nil {
		fmt.Println("error while deleting city", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (c cityRepo) countCities() (int, error) {
	var count int
	err := c.db.QueryRow("SELECT COUNT(*) FROM cities").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
