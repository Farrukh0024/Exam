package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"city2city/api/models"
)

func (h Handler) City(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCity(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCityList(w)
		} else {
			h.GetCityByID(w, r)
		}
	case http.MethodPut:
		h.UpdateCity(w, r)
	case http.MethodDelete:
		h.DeleteCity(w, r)
	}
}

func (h Handler) CreateCity(w http.ResponseWriter, r *http.Request) {
	createCity := models.CreateCity{}

	if err := json.NewDecoder(r.Body).Decode(&createCity); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.City().Create(createCity)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.storage.City().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, user)
}

func (h Handler) GetCityByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	var err error

	city, err := h.storage.City().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, city)
}

func (h Handler) GetCityList(w http.ResponseWriter) {
	var (
		page, limit = 1, 10
		err         error
	)

	resp, err := h.storage.City().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (h Handler) UpdateCity(w http.ResponseWriter, r *http.Request) {
	city := models.City{}

	if err := json.NewDecoder(r.Body).Decode(&city); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.City().Update(city)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.storage.City().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, user)
}

func (h Handler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.City().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data successfully deleted")
}
