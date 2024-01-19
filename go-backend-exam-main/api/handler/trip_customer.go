package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"city2city/api/models"
)

func (h Handler) TripCustomer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTripCustomer(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetTripCustomerList(w)
		} else {
			h.GetTripCustomerByID(w, r)
		}
	case http.MethodPut:
		h.UpdateTripCustomer(w, r)
	case http.MethodDelete:
		h.DeleteTripCustomer(w, r)
	}
}

func (h Handler) CreateTripCustomer(w http.ResponseWriter, r *http.Request) {
	createTrip := models.CreateTripCustomer{}

	if err := json.NewDecoder(r.Body).Decode(&createTrip); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.TripCustomer().Create(createTrip)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	trip, err := h.storage.TripCustomer().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, trip)
}

func (h Handler) GetTripCustomerByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	var err error

	tripCustumer, err := h.storage.TripCustomer().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, tripCustumer)
}

func (h Handler) GetTripCustomerList(w http.ResponseWriter) {
	var (
		page, limit = 1, 10
		err         error
	)

	resp, err := h.storage.TripCustomer().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (h Handler) UpdateTripCustomer(w http.ResponseWriter, r *http.Request) {
	tripCustomer := models.TripCustomer{}

	if err := json.NewDecoder(r.Body).Decode(&tripCustomer); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.TripCustomer().Update(tripCustomer)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	t, err := h.storage.TripCustomer().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, t)
}

func (h Handler) DeleteTripCustomer(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.TripCustomer().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data successfully deleted")
}
