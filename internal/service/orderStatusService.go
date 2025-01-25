package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/internal/utils"
	"net/http"
)

type OrderStatusService interface {
	GetAllOrderStatus(w http.ResponseWriter) (int, error)
	GetOrderStatus(w http.ResponseWriter, id int) (int, error)
}

type DefaultOrderStatService struct {
	repo dal.NewOrderStatusRepo
}

func DefaultOrderStatusService(repo dal.NewOrderStatusRepo) *DefaultOrderStatService {
	return &DefaultOrderStatService{repo: repo}
}

func (serv *DefaultOrderStatService) GetAllOrderStatus(w http.ResponseWriter) (int, error) {
	data, err := serv.repo.GetAllOrderStatus()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(data, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultOrderStatService) GetOrderStatus(w http.ResponseWriter, id int) (int, error) {
	exist, err := serv.repo.IsOrderStatusExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("order status not found")
	}
	data, err := serv.repo.GetOrderStatus(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(data, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
