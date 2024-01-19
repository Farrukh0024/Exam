package postgres

import (
	"database/sql"
	"fmt"

	"city2city/api/models"

	"github.com/google/uuid"
)

type driverRepo struct {
	db *sql.DB
}

func NewDriverRepo(db *sql.DB) driverRepo {
	return driverRepo{
		db: db,
	}
}

func (d driverRepo) Create(driver models.CreateDriver) (string, error) {
	uid := uuid.New().String()

	if _, err := d.db.Exec("INSERT INTO drivers (id, full_name, phone, from_city_id,to_city_id) VALUES ($1, $2, $3, $4, $5)",
		uid,
		driver.FullName,
		driver.Phone,
		driver.FromCityID,
		driver.ToCityID); err != nil {
		fmt.Println("error while inserting data drivers...", err.Error())
		return "", err
	}

	return "", nil
}

func (d driverRepo) Get(id string) (models.Driver, error) {
	stmt, err := d.db.Prepare("SELECT * FROM drivers WHERE id = $1")
	if err != nil {
		return models.Driver{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var driver models.Driver
	err = row.Scan(&driver.ID, &driver.FullName, &driver.Phone, &driver.FromCityID, &driver.ToCityID, &driver.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("error while getting data", err.Error())
			return models.Driver{}, err
		} else {
			return models.Driver{}, err
		}
	}

	return driver, nil
}

func (d driverRepo) GetList(req models.GetListRequest) (models.DriversResponse, error) {
	var (
		drivers           = []models.Driver{}
		count             = 0
		countQuery, query string
		page              = req.Page
		offset            = (page - 1) * req.Limit
	)

	countQuery = `
 SELECT count(1) from drivers `

	if err := d.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of drivers", err.Error())
		return models.DriversResponse{}, err
	}

	query = `
 SELECT id, full_name, phone, from_city_id, to_city_id, created_at
  FROM drivers
      `

	query += fmt.Sprintf("LIMIT %d OFFSET %d", req.Limit, offset)

	rows, err := d.db.Query(query)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.DriversResponse{}, err
	}

	for rows.Next() {
		driver := models.Driver{}

		if err = rows.Scan(&driver.ID, &driver.FullName, &driver.Phone, &driver.FromCityID, &driver.ToCityID, &driver.CreatedAt); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.DriversResponse{}, err
		}

		drivers = append(drivers, driver)
	}

	return models.DriversResponse{
		Drivers: drivers,
		Count:   count,
	}, nil
}

func (d driverRepo) Update(driver models.Driver) (string, error) {
	stmt, err := d.db.Prepare("UPDATE drivers SET full_name=$1, phone=$2, from_city_id=$3, to_city_id=$4 WHERE id=$5")
	if err != nil {
		fmt.Println("error while updating driver data", err.Error())
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(driver.FullName, driver.Phone, driver.FromCityID, driver.ToCityID, driver.ID)
	if err != nil {
		fmt.Println("error while updating driver data", err.Error())
		return "", err
	}

	return driver.ID, nil
}

func (d driverRepo) Delete(id string) error {
	stmt, err := d.db.Prepare("DELETE FROM drivers WHERE id=$1")
	if err != nil {
		fmt.Println("error while deleting data", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Println("error while deleting data", err.Error())
		return err
	}

	return nil
}

func (d driverRepo) countDrivers() (int, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM drivers").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
