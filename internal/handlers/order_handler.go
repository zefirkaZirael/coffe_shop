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

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Order_Handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		slog.Error("Failed to Handle Order: error URL adress")
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 || len(splitted) == 3 {
		if splitted[1] != "batch-process" {
			num, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				slog.Error("Failed to Handle Order", "convertation error: ", err)
				utils.Log_Err_Handler(err, http.StatusBadRequest, w)
				return
			}
			id = num
		} else {
			id = 1
		}
	}
	switch {
	case r.Method == http.MethodPost && len(splitted) == 1:
		order, err := h.Get_Body_Order(r)
		if err != nil {
			slog.Error("Failed to Handle Order", "Get Body Order function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.Create_Order(order)
		if err != nil {
			slog.Error("Failed to Handle Order", "Create Order function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Order created succesfullly")
		w.WriteHeader(http.StatusCreated)
		return
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.Retrieve_All_Orders(w)
		if err != nil {
			slog.Error("Failed to Handle Order", "Retrieve All Orders function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("All Orders retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.Retrieve_Order(w, id)
		if err != nil {
			slog.Error("Failed to Handle Order", "Retrieve Order function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		if code == http.StatusNotFound {
			w.WriteHeader(code)
			slog.Info("Order not found")
			return
		}
		slog.Info("Order retrieved succesfully")
		return
	case r.Method == http.MethodPut && len(splitted) == 2:
		order, err := h.Get_Body_Order(r)
		if err != nil {
			slog.Error("Failed to Handle Order", "Get Body Order function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.Update_Order(order, id)
		if err != nil {
			slog.Error("Failed to Handle Order", "Update Order function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Order updated succesfully")
		return
	case r.Method == http.MethodDelete && len(splitted) == 2:
		code, err := h.service.Delete_Order(id)
		if err != nil {
			slog.Error("Failed to Handle Order", "Delete Order function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Order deleted succesfully")
		w.WriteHeader(code)
		return
	case r.Method == http.MethodPost && len(splitted) == 3:
		code, err := h.service.Close_Order(id, w)
		if err != nil {
			slog.Error("Failed to Handle Order", "Close Order function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		if code == http.StatusNotFound {
			w.WriteHeader(code)
			slog.Info("Order not found")
			return
		}
		slog.Info("Order closed succesfully")
		return
	case r.Method == http.MethodPost && len(splitted) == 2 && splitted[1] == "batch-process":
		h.service.BatchProcessOrders(w, r)
		return
	default:
		slog.Error("Failed to Handle Order: error method in orders")
		utils.Log_Err_Handler(errors.New("error method in orders"), http.StatusMethodNotAllowed, w)
		return
	}
}

func (h *OrderHandler) Get_Body_Order(r *http.Request) (models.Order, error) {
	var order models.Order
	if r.Body == nil {
		return order, errors.New("request body is empty")
	}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		return order, err
	}
	if order.ID != 0 {
		return order, errors.New("order ID must be empty")
	}
	if order.CustomerID <= 0 {
		return order, errors.New("customer_id is missing or invalid")
	}

	if order.TotalAmount != 0 {
		return order, errors.New("total_amount must be empty")
	}

	if !order.CreatedAt.IsZero() {
		return order, errors.New("CreatedAt field must be empty")
	}

	if order.Status != "" {
		return order, errors.New("status must be empty")
	}

	if len(order.Items) == 0 {
		return order, errors.New("order must contain at least one item")
	}

	for _, item := range order.Items {
		if item.MenuItemID <= 0 {
			return order, errors.New("menu_item_id is missing or invalid in one of the items")
		}
		if item.OrderID != 0 {
			return order, errors.New("order_id must be empty")
		}
		if item.Quantity <= 0 {
			return order, errors.New("quantity must be greater than 0 in one of the items")
		}
		if item.PriceAtOrderTime != 0 {
			return order, errors.New("price_at_order_time must be empty")
		}
	}
	return order, nil
}
