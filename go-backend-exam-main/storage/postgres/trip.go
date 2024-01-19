package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"city2city/api/models"
	"city2city/storage"
	"github.com/google/uuid"
)

type tripRepo struct {
	db *sql.DB
}

func NewTripRepo(db *sql.DB) storage.ITripRepo {
	return &tripRepo{
		db: db,
	}
}

func (c *tripRepo) Create(trip models.CreateTrip) (string, error) {
	// Generate a new UUID
	tripID := uuid.New().String()

	// Prepare the SQL query with a placeholder for the UUID
	query := `INSERT INTO trips (id, trip_number_id, from_city_id, to_city_id, driver_id, price) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	// Execute the query, passing the generated UUID as a parameter
	rows, err := c.db.Query(
		query,
		tripID,
		trip.TripNumberID,
		trip.FromCityID,
		trip.ToCityID,
		trip.DriverID,
		trip.Price,
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

	return tripID, nil
}

func (c *tripRepo) Get(id string) (models.Trip, error) {
	stmt, err := c.db.Prepare("SELECT id, trip_number_id, from_city_id, to_city_id, driver_id, price, created_at FROM trips WHERE id = $1")
	if err != nil {
		return models.Trip{}, fmt.Errorf("failed to get Trip: %w", err)
	}
	defer stmt.Close()

	var trip models.Trip
	err = stmt.QueryRow(id).Scan(&trip.ID, &trip.TripNumberID, &trip.FromCityID, &trip.ToCityID, &trip.DriverID, &trip.Price, &trip.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Trip{}, fmt.Errorf("failed to get Trip: %w", err)
		} else {
			return models.Trip{}, err
		}
	}

	return trip, nil
}

func (c *tripRepo) GetList(req models.GetListRequest) (models.TripsResponse, error) {
	query := "SELECT id, trip_number_id, from_city_id, to_city_id, driver_id, price, created_at FROM trips"

	if req.Page > 0 && req.Limit > 0 {
		offset := (req.Page - 1) * req.Limit
		query += fmt.Sprintf(" ORDER BY created_at DESC OFFSET %d LIMIT %d", offset, req.Limit)
	}

	rows, err := c.db.Query(query)
	if err != nil {
		return models.TripsResponse{}, fmt.Errorf("failed to getList Trips: %w", err)
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var trip models.Trip
		err = rows.Scan(&trip.ID, &trip.TripNumberID, &trip.FromCityID, &trip.ToCityID, &trip.DriverID, &trip.Price, &trip.CreatedAt)
		if err != nil {
			return models.TripsResponse{}, err
		}
		trips = append(trips, trip)
	}

	countQuery := "SELECT COUNT(*) FROM trips"
	row := c.db.QueryRow(countQuery)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return models.TripsResponse{}, err
	}

	return models.TripsResponse{Trips: trips, Count: count}, nil
}

func (c *tripRepo) Update(trip models.Trip) (string, error) {
	stmt, err := c.db.Prepare("UPDATE trips SET trip_number_id = $1, from_city_id = $2, to_city_id = $3, driver_id = $4, price = $5 WHERE id = $6")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	result, err := stmt.Exec(trip.TripNumberID, trip.FromCityID, trip.ToCityID, trip.DriverID, trip.Price, trip.ID)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", errors.New("failed to update Trip")
	}

	return trip.ID, nil
}

func (c *tripRepo) Delete(id string) error {
	stmt, err := c.db.Prepare("DELETE FROM trips WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete Trip")
	}

	return nil
}
