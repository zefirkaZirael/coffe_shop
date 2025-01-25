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

type InventHandler struct {
	service service.InventService
}

func NewInventHandle(service service.InventService) *InventHandler {
	return &InventHandler{service: service}
}

func (h *InventHandler) Inventory_Handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 {
		num, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Convertation error: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		id = num
	}
	switch {
	case r.Method == http.MethodPost && len(splitted) == 1:
		inventory, err := h.Get_Body_Inventory(r)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Get Body Inventory function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.Add_Inventory(inventory)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Add Inventory function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Inventory added succesfully")
		w.WriteHeader(code)
		return
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.Retrieve_All_Inventory(w)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Retrieve all inventory function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("All Inventory retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.Retrieve_Inventory(w, id)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Retrieve Inventory function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Inventory retrieved succesfully")
		return
	case r.Method == http.MethodPut && len(splitted) == 2:
		inventory, err := h.Get_Body_Inventory(r)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Get Body Inventory function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.Update_Inventory(inventory, id)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Update Inventory function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Inventory updated succesfully")
		return
	case r.Method == http.MethodDelete && len(splitted) == 2:
		code, err := h.service.Delete_Inventory(id)
		if err != nil {
			slog.Error("Failed to Handle Inventory", "Delete Inventory function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Inventory deleted succesfully")
		w.WriteHeader(code)
		return
	default:
		utils.Log_Err_Handler(errors.New("error method in inventory"), http.StatusMethodNotAllowed, w)
		return
	}
}

func (h *InventHandler) InventoryTransaction_Handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 {
		num, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("Failed to Handle Inventory Transaction", "convertation error: ", err)
			utils.Log_Err_Handler(err, http.StatusInternalServerError, w)
			return
		}
		id = num
	}
	switch {
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.GetAllTransactionData(w)
		if err != nil {
			slog.Error("Failed to Handle Inventory Transaction", "Get All Transaction function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("All Inventory Transaction retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.GetInventoryTransaction(w, id)
		if err != nil {
			slog.Error("Failed to Handle Inventory Transaction", "Get Inventory Transaction: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Inventory Transaction retrieved succesfully")
		return
	case r.Method == http.MethodPost && len(splitted) == 1:
		transaction, err := h.Get_Body_InventoryTransaction(r)
		if err != nil {
			slog.Error("Failed to Handle Inventory Transaction", "Get Body Inventory Transaction function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.FillInventory(transaction)
		if err != nil {
			slog.Error("Failed to Handle Inventory Transaction", "Fill Inventory function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Inventory filled succesfully")
		return
	default:
		utils.Log_Err_Handler(errors.New("error method in Inventory Transaction"), http.StatusMethodNotAllowed, w)
		return
	}
}

func (h *InventHandler) Get_Body_InventoryTransaction(r *http.Request) (models.InventoryTransaction, error) {
	var transaction models.InventoryTransaction
	if r.Body == nil {
		return transaction, errors.New("request body is empty")
	}

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return transaction, err
	}
	if transaction.ID != 0 {
		return transaction, errors.New("transaction id must be empty")
	}
	if transaction.Inventory_id <= 0 {
		return transaction, errors.New("transaction inventory ID field cannot be negative or empty")
	}
	if transaction.Price <= 0 {
		return transaction, errors.New("transaction price field cannot be negative or empty")
	}
	if transaction.Quantity <= 0 {
		return transaction, errors.New("transaction quantity field cannot be negative or empty")
	}
	if !transaction.Transaction_date.IsZero() {
		return transaction, errors.New("transaction date must be empty")
	}
	return transaction, nil
}

func (h *InventHandler) Get_Body_Inventory(r *http.Request) (models.InventoryItem, error) {
	var inventory models.InventoryItem

	if r.Body == nil {
		return inventory, errors.New("request body is empty")
	}

	if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
		return inventory, err
	}

	if inventory.ID != 0 {
		return inventory, errors.New("inventory id field must be empty")
	}
	if inventory.Name == "" {
		return inventory, errors.New("name field is missing")
	}
	if inventory.StockLevel <= 0 {
		return inventory, errors.New("stock_level field cannot be negative or empty")
	}
	if inventory.UnitType == "" {
		return inventory, errors.New("unit_type field is missing")
	}
	if inventory.ReorderLevel <= 0 {
		return inventory, errors.New("reorder_level field cannot be negative or empty")
	}
	if !inventory.LastUpdated.IsZero() {
		return inventory, errors.New("last_updated field must be empty")
	}
	return inventory, nil
}
