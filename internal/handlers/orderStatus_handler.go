package handlers

import (
	"errors"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type OrderStatusHandler struct {
	service service.OrderStatusService
}

func NewOrderStatusHandle(service service.OrderStatusService) *OrderStatusHandler {
	return &OrderStatusHandler{service: service}
}

func (h *OrderStatusHandler) OrderStatus_handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 {
		num, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("Failed to Handle Order Status", "convertation error", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		id = num
	}
	switch {
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.GetAllOrderStatus(w)
		if err != nil {
			slog.Error("Failed to Handle Order Status", "Get All Order Status function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("All Order Status retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.GetOrderStatus(w, id)
		if err != nil {
			slog.Error("Failed to Handle Order Status", "Get Order Status function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Order Status retrieved succesfully")
		return
	default:
		utils.Log_Err_Handler(errors.New("error method in order status"), http.StatusMethodNotAllowed, w)
		return
	}
}
