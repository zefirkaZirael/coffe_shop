package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"net/http"
)

type InventService interface {
	Add_Inventory(inventory models.InventoryItem) (int, error)
	Retrieve_All_Inventory(w http.ResponseWriter) (int, error)
	Retrieve_Inventory(w http.ResponseWriter, id int) (int, error)
	Update_Inventory(inventory models.InventoryItem, id int) (int, error)
	Delete_Inventory(id int) (int, error)
	GetAllTransactionData(w http.ResponseWriter) (int, error)
	GetInventoryTransaction(w http.ResponseWriter, id int) (int, error)
	FillInventory(transaction models.InventoryTransaction) (int, error)
}

type DefaultInventService struct {
	repo dal.InventRepo
}

func NewDefaultInventService(repo dal.InventRepo) *DefaultInventService {
	return &DefaultInventService{repo: repo}
}

// Add_Inventory adds a new inventory item
func (serv *DefaultInventService) Add_Inventory(inventory models.InventoryItem) (int, error) {
	unique, err := serv.repo.IsInventUnique(inventory.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !unique {
		return http.StatusBadRequest, errors.New("inventory name must be unique")
	}
	err = serv.repo.Save_Inventory(inventory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

// Retrieve_All_Inventory retrieves all inventory items
func (serv *DefaultInventService) Retrieve_All_Inventory(w http.ResponseWriter) (int, error) {
	inventory, err := serv.repo.Get_AllInventory()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = utils.Send_Request(inventory, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Retrieve_Inventory retrieves a specific inventory item by ID
func (serv *DefaultInventService) Retrieve_Inventory(w http.ResponseWriter, id int) (int, error) {
	exist, err := serv.repo.IsInventExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("inventory item not found")
	}

	foundItem, err := serv.repo.GetInventory(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(foundItem, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Update_Inventory updates an existing inventory item by ID
func (serv *DefaultInventService) Update_Inventory(inventory models.InventoryItem, id int) (int, error) {
	unique, err := serv.repo.IsInventUnique(inventory.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !unique {
		return http.StatusBadRequest, errors.New("inventory name must be unique")
	}
	exist, err := serv.repo.IsInventExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("inventory item not found")
	}
	inventory.ID = id
	err = serv.repo.Update_Inventory(inventory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Delete_Inventory deletes an inventory item by ID
func (serv *DefaultInventService) Delete_Inventory(id int) (int, error) {
	exist, err := serv.repo.IsInventExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("inventory item not found")
	}
	err = serv.repo.Delete_Inventory(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (serv *DefaultInventService) GetAllTransactionData(w http.ResponseWriter) (int, error) {
	data, err := serv.repo.GetAllInventoryTransactions()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(data, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultInventService) GetInventoryTransaction(w http.ResponseWriter, id int) (int, error) {
	exist, err := serv.repo.IsInventTransactionExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("inventory transaction not found")
	}
	data, err := serv.repo.GetInventoryTransaction(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(data, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultInventService) FillInventory(transaction models.InventoryTransaction) (int, error) {
	exist, err := serv.repo.IsInventExist(transaction.Inventory_id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("inventory item is not found")
	}
	err = serv.repo.FillInventory(transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
