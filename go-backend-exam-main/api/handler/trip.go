package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"city2city/api/models"
)

func (h Handler) Trip(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTrip(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetTripList(w)
		} else {
			h.GetTripByID(w, r)
		}
	case http.MethodPut:
		h.UpdateTrip(w, r)
	case http.MethodDelete:
		h.DeleteTrip(w, r)
	}
}

func (h Handler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	createTrip := models.CreateTrip{}

	if err := json.NewDecoder(r.Body).Decode(&createTrip); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.Trip().Create(createTrip)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	trip, err := h.storage.Trip().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, trip)
}

func (h Handler) GetTripByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	var err error

	trip, err := h.storage.Trip().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, trip)
}

func (h Handler) GetTripList(w http.ResponseWriter) {
	var (
		page, limit = 1, 10 // Adjust defaults as needed
		err         error
	)

	resp, err := h.storage.Trip().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (h Handler) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	trip := models.Trip{}

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.Trip().Update(trip)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	t, err := h.storage.Trip().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, t)
}

func (h Handler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Trip().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data successfully deleted")
}
