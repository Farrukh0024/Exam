package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"city2city/api/models"
)

func (h Handler) Customer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCustomer(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			h.GetCustomerList(w, r)
		} else {
			h.GetCustomerByID(w, r)
		}
	case http.MethodPut:
		h.UpdateCustomer(w, r)
	case http.MethodDelete:
		h.DeleteCustomer(w, r)
	}
}

func (h Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	createCustomer := models.CreateCustomer{}

	if err := json.NewDecoder(r.Body).Decode(&createCustomer); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.Customer().Create(createCustomer)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	customer, err := h.storage.Customer().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, customer)
}

func (h Handler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]
	var err error

	customer, err := h.storage.Customer().Get(id)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, customer)
}

func (h Handler) GetCustomerList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10
		err         error
	)
	values := r.URL.Query()

	if len(values["page"]) > 0 {
		page, err = strconv.Atoi(values["page"][0])
		if err != nil {
			page = 1
		}
	}

	if len(values["limit"]) > 0 {
		limit, err = strconv.Atoi(values["limit"][0])
		if err != nil {
			fmt.Println("limit", values["limit"])
			limit = 10
		}
	}

	resp, err := h.storage.Customer().GetList(models.GetListRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (h Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customer := models.Customer{}

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.Customer().Update(customer)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	c, err := h.storage.Customer().Get(pKey)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, c)
}

func (h Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if len(values["id"]) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	id := values["id"][0]

	if err := h.storage.Customer().Delete(id); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, "data successfully deleted")
}
