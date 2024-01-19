package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"city2city/api/models"
)

func (h Handler) Car(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCar(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCarList(w, r)
		} else {
			h.GetCarByID(w, r)
		}
	case http.MethodPut:
		values := r.URL.Query()
		if _, ok := values["route"]; ok {
			h.UpdateCarRoute(w, r)
		} else if _, ok := values["status"]; ok {
			h.UpdateCarStatus(w, r)
		} else {
			h.UpdateCar(w, r)
		}
	case http.MethodDelete:
		h.DeleteCar(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h Handler) CreateCar(w http.ResponseWriter, r *http.Request) {
	createCar := models.CreateCar{}

	if err := json.NewDecoder(r.Body).Decode(&createCar); err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	id, err := h.storage.Car().Create(createCar)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, car)
}

func (h Handler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	id := values["id"][0]

	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, car)
}

func (h Handler) GetCarList(w http.ResponseWriter, r *http.Request) {
	var (
		cars  = []models.Car{}
		count int
	)

	// Data query
	req := models.GetListRequest{
		Page:  1,
		Limit: 10,
	}

	resp, err := h.storage.Car().GetList(req)
	if err != nil {
		fmt.Println("Error in GetList:", err) // Hatanın ayrıntılarını yazdır
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	cars = resp.Cars
	count = resp.Count

	handleResponse(w, http.StatusOK, models.CarsResponse{
		Cars:  cars,
		Count: count,
	})
}

func (h Handler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	updateCar := models.Car{}

	if err := json.NewDecoder(r.Body).Decode(&updateCar); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Car().Update(updateCar)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	car, err := h.storage.Car().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, car)
}

func (h Handler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	id := values["id"][0]

	if err := h.storage.Car().Delete(id); err != nil {
		handleResponse(w, http.StatusNotFound, err)
		return
	} else {
		handleResponse(w, http.StatusInternalServerError, err)
	}

	handleResponse(w, http.StatusOK, "data succesfully deleted")
}

func (h Handler) UpdateCarRoute(w http.ResponseWriter, r *http.Request) {
	updateCarRoute := models.UpdateCarRoute{}

	if err := json.NewDecoder(r.Body).Decode(&updateCarRoute); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(updateCarRoute.CarID) == 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("CarID is required"))
		return
	}

	if err := h.storage.Car().UpdateCarRoute(updateCarRoute); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "Car route successfully updated")
}

func (h Handler) UpdateCarStatus(w http.ResponseWriter, r *http.Request) {
	updateCarStatus := models.UpdateCarStatus{}

	if err := json.NewDecoder(r.Body).Decode(&updateCarStatus); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(updateCarStatus.ID) == 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("ID is required"))
		return
	}

	if err := h.storage.Car().UpdateCarStatus(updateCarStatus); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "Car status updated successfully")
}
