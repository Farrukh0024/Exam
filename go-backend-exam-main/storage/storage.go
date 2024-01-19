package storage

import (
	"city2city/api/models"
)

type IStorage interface {
	CloseDB()
	City() ICityRepo
	Customer() ICustomerRepo
	Driver() IDriverRepo
	Car() ICarRepo
	Trip() ITripRepo
	TripCustomer() ITripCustomerRepo
}

type ICityRepo interface {
	Create(models.CreateCity) (string, error)
	Get(id string) (models.City, error)
	GetList(models.GetListRequest) (models.CitiesResponse, error)
	Update(models.City) (string, error)
	Delete(id string) error
}

type ICustomerRepo interface {
	Create(customer models.CreateCustomer) (string, error)
	Get(id string) (models.Customer, error)
	GetList(models.GetListRequest) (models.CustomersResponse, error)
	Update(models.Customer) (string, error)
	Delete(id string) error
}

type IDriverRepo interface {
	Create(driver models.CreateDriver) (string, error)
	Get(id string) (models.Driver, error)
	GetList(models.GetListRequest) (models.DriversResponse, error)
	Update(models.Driver) (string, error)
	Delete(id string) error
}

type ICarRepo interface {
	Create(models.CreateCar) (string, error)
	Get(id string) (models.Car, error)
	GetList(models.GetListRequest) (models.CarsResponse, error)
	Update(models.Car) (string, error)
	Delete(id string) error
	UpdateCarStatus(updateCarStatus models.UpdateCarStatus) error
	UpdateCarRoute(updateCarRoute models.UpdateCarRoute) error
}

type ITripRepo interface {
	Create(models.CreateTrip) (string, error)
	Get(id string) (models.Trip, error)
	GetList(models.GetListRequest) (models.TripsResponse, error)
	Update(models.Trip) (string, error)
	Delete(id string) error
}

type ITripCustomerRepo interface {
	Create(models.CreateTripCustomer) (string, error)
	Get(id string) (models.TripCustomer, error)
	GetList(models.GetListRequest) (models.TripCustomersResponse, error)
	Update(models.TripCustomer) (string, error)
	Delete(id string) error
}
