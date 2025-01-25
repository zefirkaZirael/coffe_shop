package service

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type OrderService interface {
	Create_Order(newOrder models.Order) (int, error)
	Retrieve_All_Orders(w http.ResponseWriter) (int, error)
	Retrieve_Order(w http.ResponseWriter, id int) (int, error)
	Update_Order(order models.Order, id int) (int, error)
	Delete_Order(id int) (int, error)
	Close_Order(id int, w http.ResponseWriter) (int, error)
	BatchProcessOrders(w http.ResponseWriter, r *http.Request)
}

type DefaultOrderService struct {
	repo dal.NewOrderRepo
}

func NewDefaultOrderService(repo dal.NewOrderRepo) *DefaultOrderService {
	return &DefaultOrderService{repo: repo}
}

func (s *DefaultOrderService) Create_Order(newOrder models.Order) (int, error) {
	date := time.Now()
	code, err := s.repo.CheckOrder(newOrder)
	if err != nil {
		return code, err
	}
	newOrder.CreatedAt = date
	newOrder.Status = "active"
	err = s.repo.GetPriceAtOrderItems(&newOrder)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = s.repo.GetTotalAmount(&newOrder)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if exist := dal.DefaultCustomerRepo(s.repo.DB).IsCustomerExist(newOrder.CustomerID); !exist {
		return http.StatusBadRequest, errors.New("customer id is not exist")
	}

	err = s.repo.SaveOrder(newOrder)
	if err != nil {
		slog.Error("Failed to create order", "Save Order function", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *DefaultOrderService) Retrieve_All_Orders(w http.ResponseWriter) (int, error) {
	orders, err := s.repo.Get_Orders()
	if err != nil {
		slog.Error("Failed to Retrieve orders", "Get order function", err)
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(orders, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *DefaultOrderService) Retrieve_Order(w http.ResponseWriter, id int) (int, error) {
	exist, err := s.repo.IsOrderExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("order not found")
	}
	order, err := s.repo.GetOrder(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(order, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *DefaultOrderService) Update_Order(order models.Order, id int) (int, error) {
	return s.repo.UpdateOrder(order, id)
}

func (s *DefaultOrderService) Delete_Order(id int) (int, error) {
	exist, err := s.repo.IsOrderExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("order not found")
	}
	err = s.repo.DeleteOrder(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (s *DefaultOrderService) Close_Order(id int, w http.ResponseWriter) (int, error) {
	exist, err := s.repo.IsOrderExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("order not found")
	}
	order, err := s.repo.GetOrder(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if order.Status == "closed" {
		return http.StatusBadRequest, errors.New("order is already closed")
	}
	need_invent, err := s.repo.Get_Need_Inventory(order)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	enough, err := s.repo.Is_Enough(need_invent)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !enough {
		return http.StatusConflict, errors.New("not enough inventory")
	}
	err = s.repo.CloseOrder(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	reorderItems, err := dal.DefaultInventRepo(s.repo.DB).CheckReordering()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if reorderItems == nil {
		return http.StatusOK, nil
	}
	response := struct {
		Items   []models.InventoryItem `json:"items"`
		Comment string                 `json:"comment"`
	}{
		Items:   reorderItems,
		Comment: "these items need to be replenished",
	}
	err = utils.Send_Request(response, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *DefaultOrderService) BatchProcessOrders(w http.ResponseWriter, r *http.Request) {
	var request models.BatchProcessRequest
	var response models.BatchProcessResponse
	var totalRevenue float64
	var accepted, rejected int
	var inventoryUpdates []models.InventoryUpdate

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.Log_Err_Handler(errors.New("failed to parse request body "+err.Error()), http.StatusBadRequest, w)
		return
	}

	// Start a transaction
	tx, err := s.repo.DB.Begin()
	if err != nil {
		utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, order := range request.Orders {
		status := "accepted"
		var reason string
		if len(order.Items) == 0 {
			status = "rejected"
			reason = "invalid order:order items is empty"
		}
		for _, item := range order.Items {
			exists, err := s.repo.ProductID_Exists(strconv.Itoa(item.MenuItemID))
			if err != nil {
				status = "rejected"
				reason = "error occured: " + err.Error()
			}
			if !exists {
				status = "rejected"
				reason = "ordered item does not exist: " + strconv.Itoa(item.MenuItemID)
			}
		}

		// Validate and check inventory
		needInventory, err := s.repo.Get_Need_Inventory(order)
		if err != nil {
			status = "rejected"
			reason = "invalid_order"
		} else {
			enough, err := s.repo.Is_Enough(needInventory)
			if err != nil {
				status = "rejected"
				reason = "invalid order"
			}
			if !enough {
				status = "rejected"
				reason = "insufficient_inventory"
			}
		}
		if exist := dal.DefaultCustomerRepo(s.repo.DB).IsCustomerExist(order.CustomerID); !exist {
			status = "rejected"
			reason = "customer id is not exist"
		}
		// Process the order if valid
		if status == "accepted" {
			date := time.Now()
			code, err := s.repo.CheckOrder(order)
			if err != nil {
				tx.Rollback()
				utils.Log_Err_Handler(err, code, w)
				return
			}
			order.CreatedAt = date
			order.Status = "active"
			err = s.repo.GetPriceAtOrderItems(&order)
			if err != nil {
				tx.Rollback()
				utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
				return
			}
			err = s.repo.GetTotalAmount(&order)
			if err != nil {
				tx.Rollback()
				utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
				return
			}
			err = s.repo.SaveOrder(order)
			if err != nil {
				tx.Rollback()
				utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
				return
			}
			order.ID, err = s.repo.GetLastOrderId()
			if err != nil {
				tx.Rollback()
				utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
				return
			}
			// Update inventory
			err = s.repo.CloseOrder(order.ID)
			if err != nil {
				tx.Rollback()
				utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
				return
			}

			// Update revenue and inventory updates
			totalRevenue += order.TotalAmount
			for ingredientID, quantityUsed := range needInventory {
				num, err := strconv.Atoi(ingredientID)
				if err != nil {
					tx.Rollback()
					utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
					return
				}
				inventory, err := dal.DefaultInventRepo(s.repo.DB).GetInventory(num)
				if err != nil {
					tx.Rollback()
					utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
					return
				}
				inventoryUpdates = append(inventoryUpdates, models.InventoryUpdate{
					IngredientID: ingredientID,
					QuantityUsed: quantityUsed,
					Remaining:    inventory.StockLevel, // Retrieve the updated stock from DB if needed
				})
			}
			accepted++
		} else {
			rejected++
		}

		// Add to processed orders
		response.ProcessedOrders = append(response.ProcessedOrders, models.ProcessedOrder{
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Status:     status,
			Total:      order.TotalAmount,
			Reason:     reason,
		})
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
		return
	}

	// Prepare the response summary
	response.Summary = models.BatchProcessSummary{
		TotalOrders:      len(request.Orders),
		Accepted:         accepted,
		Rejected:         rejected,
		TotalRevenue:     totalRevenue,
		InventoryUpdates: inventoryUpdates,
	}

	// Send response
	err = utils.Send_Request(response, w)
	if err != nil {
		utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
		return
	}
	slog.Info("Order batch process has been succesfully completed")
}
