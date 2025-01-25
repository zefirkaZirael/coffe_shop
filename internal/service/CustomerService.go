package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"net/http"
)

type DefaultCustomerService struct {
	repo dal.NewCustomerRepo
}

type CustomerService interface {
	CreateCustomer(customer models.Customer) (int, error)
	GetAllCustomers(w http.ResponseWriter) (int, error)
	GetCustomer(w http.ResponseWriter, id int) (int, error)
	UpdateCustomer(customer models.Customer, id int) (int, error)
	DeleteCustomer(id int) (int, error)
}

func NewDefaultServiceCustomer(repo dal.NewCustomerRepo) *DefaultCustomerService {
	return &DefaultCustomerService{repo: repo}
}

func (serv *DefaultCustomerService) CreateCustomer(customer models.Customer) (int, error) {
	unique, err := serv.repo.IsEmailUnique(customer.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !unique {
		return http.StatusBadRequest, errors.New("email must be unique")
	}
	err = serv.repo.SaveCustomer(customer)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func (serv *DefaultCustomerService) GetAllCustomers(w http.ResponseWriter) (int, error) {
	customers, err := serv.repo.GetAllCustomers()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(customers, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultCustomerService) GetCustomer(w http.ResponseWriter, id int) (int, error) {
	if !serv.repo.IsCustomerExist(id) {
		return http.StatusNotFound, errors.New("customer is not exist")
	}
	customer, err := serv.repo.GetCustomerByID(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(customer, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultCustomerService) UpdateCustomer(customer models.Customer, id int) (int, error) {
	if !serv.repo.IsCustomerExist(id) {
		return http.StatusNotFound, errors.New("customer is not exist")
	}
	unique, err := serv.repo.IsEmailUnique(customer.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !unique {
		return http.StatusBadRequest, errors.New("email must be unique")
	}
	err = serv.repo.UpdateCustomer(customer, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultCustomerService) DeleteCustomer(id int) (int, error) {
	if !serv.repo.IsCustomerExist(id) {
		return http.StatusNotFound, errors.New("customer is not exist")
	}
	err := serv.repo.DeleteCustomer(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}
