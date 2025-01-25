package dal

import (
	"database/sql"
	"frappuccino/models"
)

type OrderStatusRepo interface {
	GetAllOrderStatus() ([]models.OrderStatus, error)
	GetOrderStatus(id int) (models.OrderStatus, error)
	IsOrderStatusExist(id int) (bool, error)
}

type NewOrderStatusRepo struct {
	DB *sql.DB
}

func DefaultOrderStatusRepo(db *sql.DB) *NewOrderStatusRepo {
	return &NewOrderStatusRepo{DB: db}
}

// Retrieves the information about all order status histories from database
func (repo *NewOrderStatusRepo) GetAllOrderStatus() ([]models.OrderStatus, error) {
	var data []models.OrderStatus
	rows, err := repo.DB.Query(`SELECT id, order_id, status, changed_at 
		FROM order_status_history
	`)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderStatus models.OrderStatus
		err := rows.Scan(&orderStatus.Id, &orderStatus.Order_id, &orderStatus.Status, &orderStatus.CreatedAt)
		if err != nil {
			return data, err
		}
		data = append(data, orderStatus)
	}
	return data, nil
}

// Retrieves the order status history by id from database
func (repo *NewOrderStatusRepo) GetOrderStatus(id int) ([]models.OrderStatus, error) {
	var data []models.OrderStatus
	rows, err := repo.DB.Query(`SELECT id, order_id, status, changed_at 
	FROM order_status_history
	WHERE order_id=$1
	`, id)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderStatus models.OrderStatus
		err := rows.Scan(&orderStatus.Id, &orderStatus.Order_id, &orderStatus.Status, &orderStatus.CreatedAt)
		if err != nil {
			return data, err
		}
		data = append(data, orderStatus)
	}
	return data, nil
}

// Checks is order status exist by ID from database
func (repo *NewOrderStatusRepo) IsOrderStatusExist(id int) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`SELECT COUNT(*) FROM order_status_history
	WHERE order_id=$1
	`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
