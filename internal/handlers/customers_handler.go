package handlers

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandle(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Customers_handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		slog.Error("Failed to Handle Customer", "error", "invalid URL adress")
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 {
		num, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("Failed to Handle Customer", "Convertation error: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		id = num
	}
	switch {
	case r.Method == http.MethodPost && len(splitted) == 1:
		customer, err := GetCustomersBody(r)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Get Customers Body function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.CreateCustomer(customer)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Create Customer function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Customer created succesfully")
		w.WriteHeader(code)
		return
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.GetAllCustomers(w)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Get All Customers Function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Customers retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.GetCustomer(w, id)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Get Customer function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Customer retrieved succesfully")
		return
	case r.Method == http.MethodPut && len(splitted) == 2:
		customer, err := GetCustomersBody(r)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Get Customers Body function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.UpdateCustomer(customer, id)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Update Customer function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Customer updated succesfully")
		return
	case r.Method == http.MethodDelete && len(splitted) == 2:
		code, err := h.service.DeleteCustomer(id)
		if err != nil {
			slog.Error("Failed to Handle Customer", "Delete Customer function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		if code == http.StatusNotFound {
			w.WriteHeader(code)
			slog.Info("Customer not found")
			return
		}
		slog.Info("Customer deleted succesfully")
		w.WriteHeader(code)
		return
	default:
		slog.Info("error method in customers")
		utils.Log_Err_Handler(errors.New("error method in customers"), http.StatusMethodNotAllowed, w)
		return
	}
}

func GetCustomersBody(r *http.Request) (models.Customer, error) {
	var customer models.Customer

	if r.Body == nil {
		return customer, errors.New("request Body is empty")
	}
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		return customer, err
	}
	if customer.Customer_id != 0 {
		return customer, errors.New("customer id must be empty")
	}
	if customer.Name == "" {
		return customer, errors.New("customer name is missing")
	}
	if customer.Email == "" {
		return customer, errors.New("customer email is missing")
	}
	return customer, nil
}
