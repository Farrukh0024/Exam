package postgres

import (
	"database/sql"
	"fmt"

	"city2city/api/models"
	"city2city/storage"
	"github.com/google/uuid"
)

type carRepo struct {
	db *sql.DB
}

func NewCarRepo(db *sql.DB) storage.ICarRepo {
	return carRepo{db: db}
}

func (c carRepo) Create(car models.CreateCar) (string, error) {
	uid := uuid.New()

	result, err := c.db.Exec("INSERT INTO cars (id, model, brand, number, driver_id) VALUES ($1, $2, $3, $4, $5)",
		uid,
		car.Model,
		car.Brand,
		car.Number,
		car.DriverID,
	)
	if err != nil {
		return "", fmt.Errorf("error while inserting data: %w", err)
	}

	// LastInsertId ve RowsAffected değerlerini kontrol et
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting LastInsertId: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("error getting RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("no rows affected")
	}

	return fmt.Sprintf("Last Insert ID: %d", lastInsertID), nil
}

func (c carRepo) Get(id string) (models.Car, error) {
	query := `SELECT * FROM cars WHERE id = $1`
	row := c.db.QueryRow(query, id)
	var car models.Car
	if err := row.Scan(&car.ID, &car.Model, &car.Brand, &car.Number, &car.DriverID, &car.DriverData, &car.CreatedAt); err != nil {
		return models.Car{}, fmt.Errorf("error getting car: %w", err)
	}
	return car, nil
}

func (c carRepo) GetList(req models.GetListRequest) (models.CarsResponse, error) {
	var (
		cars  = []models.Car{}
		count int
	)

	// Count query
	countQuery := `SELECT COUNT(*) FROM cars`
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		return models.CarsResponse{}, fmt.Errorf("error getting car count: %w", err)
	}

	// Data query
	query := `SELECT c.id, c.model, c.brand, c.number, c.driver_id, c.created_at,
	                 d.id as driver_id, d.full_name, d.phone,
	                 d.from_city_id as driver_from_city_id,
	                 d.to_city_id as driver_to_city_id,
	                 d.created_at as driver_created_at
              FROM cars c
              LEFT JOIN drivers d ON c.driver_id = d.id
              ORDER BY c.created_at DESC
              LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return models.CarsResponse{}, fmt.Errorf("error getting car list: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		var driver models.Driver
		if err := rows.Scan(&car.ID, &car.Model, &car.Brand, &car.Number, &car.DriverID, &car.CreatedAt,
			&driver.ID, &driver.FullName, &driver.Phone, &driver.FromCityID, &driver.ToCityID, &driver.CreatedAt); err != nil {
			return models.CarsResponse{}, fmt.Errorf("error scanning car row: %w", err)
		}
		car.DriverData = models.Driver{
			ID:           driver.ID,
			FullName:     driver.FullName,
			Phone:        driver.Phone,
			FromCityID:   driver.FromCityID,
			FromCityData: models.City{ID: driver.FromCityID}, // Bu kısmı veritabanınıza uygun şekilde doldurun
			ToCityID:     driver.ToCityID,
			ToCityData:   models.City{ID: driver.ToCityID}, // Bu kısmı veritabanınıza uygun şekilde doldurun
			CreatedAt:    driver.CreatedAt,
		}
		cars = append(cars, car)
	}

	return models.CarsResponse{
		Cars:  cars,
		Count: count,
	}, nil
}

func (c carRepo) Update(car models.Car) (string, error) {
	query := `UPDATE cars
		              SET model = $1,
		                  brand = $2,
		                  number = $3
		                  WHERE id = $4 `
	result, err := c.db.Exec(query, car.Model, car.Brand, car.Number, car.ID)
	if err != nil {
		return "", fmt.Errorf("error updating car: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("error getting affected rows: %w", err)
	}

	if affectedRows == 0 {
		return "", fmt.Errorf("no rows affected")
	}

	return car.ID, nil
}

func (c carRepo) Delete(id string) error {
	query := `DELETE FROM cars WHERE id = $1`
	result, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting car: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if affectedRows == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (c carRepo) UpdateCarRoute(models.UpdateCarRoute) error {

	return nil
}
func (c carRepo) UpdateCarStatus(updateCarStatus models.UpdateCarStatus) error {

	return nil
}
