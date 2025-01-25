package dal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"frappuccino/models"
	"net/http"
	"strconv"
	"time"
)

type OrderRepo interface {
	Get_Need_Inventory(order models.Order) (map[string]float64, error)
	Is_Enough(need_invent map[string]float64) (bool, error)
	ProductID_Exists(product_id string) (bool, error)
	IsOrderExist(order_id int) (bool, error)
	SaveOrder(order models.Order) error
	Get_Orders() ([]models.Order, error)
	GetOrder(order_id int) (models.Order, error)
	Get_Order_Items(order_id int) ([]models.OrderItem, error)
	GetOrderItemsArray(order_id int) ([]string, error)
	CheckOrder(newOrder models.Order) (int, error)
	GetPriceAtOrderItems(order *models.Order) error
	GetTotalAmount(order *models.Order) error
	UpdateOrder(order models.Order, id int) (int, error)
	DeleteOrder(order_id int) error
	CloseOrder(order_id int) error
	UseInventory(needInventory map[string]float64) error
	GetLastOrderId() (int, error)
}

type NewOrderRepo struct {
	DB *sql.DB
}

func DefaultOrderRepo(db *sql.DB) *NewOrderRepo {
	return &NewOrderRepo{DB: db}
}

// Get_Need_Inventory calculates the inventory required for an order
func (repo *NewOrderRepo) Get_Need_Inventory(order models.Order) (map[string]float64, error) {
	needInventory := make(map[string]float64)

	for _, item := range order.Items {
		rows, err := repo.DB.Query(`
			SELECT inventory_id,
			quantity * $1
			FROM menu_item_ingredients 
			WHERE menu_item_id = $2
		`, item.Quantity, item.MenuItemID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var ingredientID string
			var quantity float64
			if err := rows.Scan(&ingredientID, &quantity); err != nil {
				return nil, err
			}
			needInventory[ingredientID] += quantity
		}
	}
	return needInventory, nil
}

// Is_Enough checks if the inventory is sufficient for an order
func (repo *NewOrderRepo) Is_Enough(needInventory map[string]float64) (bool, error) {
	for ingredientID, requiredQuantity := range needInventory {
		var availableQuantity float64
		err := repo.DB.QueryRow(`
			SELECT stock_level
			FROM inventory
			WHERE inventory_id = $1
		`, ingredientID).Scan(&availableQuantity)
		if err != nil {
			return false, err
		}
		if availableQuantity <= requiredQuantity {
			return false, nil
		}
	}
	return true, nil
}

// ProductID_Exists checks if a product exists in the menu
func (repo *NewOrderRepo) ProductID_Exists(productID string) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`
		SELECT COUNT(*)
		FROM menu_items
		WHERE menu_item_id = $1
	`, productID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SaveOrder inserts a new order into the database
func (repo *NewOrderRepo) SaveOrder(order models.Order) error {
	time := time.Now()
	items := order.Items
	tx, err := repo.DB.Begin() // Start a transaction
	if err != nil {
		return err
	}

	// Marshal `special_instructions` to JSON
	specialInstructionsJSON, err := json.Marshal(order.SpecialInstructions)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert order
	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (customer_id, total_amount, status, special_Instructions)
		VALUES ($1, $2, $3, $4) RETURNING order_id
	`, order.CustomerID, order.TotalAmount, order.Status, specialInstructionsJSON).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert order items
	for _, item := range items {
		customizationsJSON, err := json.Marshal(item.Customizations)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`
			INSERT INTO order_items (menu_item_id, order_id, customizations, price_at_order_time, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`, item.MenuItemID, orderID, customizationsJSON, item.PriceAtOrderTime, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	order.ID = orderID
	_, err = tx.Exec(`INSERT INTO order_status_history(order_id,status,changed_at)
	VALUES($1,'active', $2)
	`, order.ID, time)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit() // Commit the transaction
}

// Get_Orders retrieves all orders from the database
func (repo *NewOrderRepo) Get_Orders() ([]models.Order, error) {
	rows, err := repo.DB.Query(`
		SELECT order_id, customer_id,order_date, status, total_amount, special_instructions
		FROM orders
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var specialInstructions []byte // To store the raw JSONB data
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.CreatedAt, &order.Status, &order.TotalAmount, &specialInstructions); err != nil {
			return nil, err
		}

		// Unmarshal the `special_instructions` JSONB data
		if len(specialInstructions) > 0 {
			if err := json.Unmarshal(specialInstructions, &order.SpecialInstructions); err != nil {
				return nil, err
			}
		}

		orderItems, err := repo.Get_Order_Items(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = orderItems
		orders = append(orders, order)
	}

	return orders, nil
}

// Get_Order_Items finds order items that belong to a specific order.
func (repo *NewOrderRepo) Get_Order_Items(order_id int) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	var customizations []byte
	rows, err := repo.DB.Query(`SELECT order_item_id, menu_item_id, order_id, customizations, price_at_order_time, quantity 
	FROM order_items
	WHERE order_id=$1`, order_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderItem models.OrderItem
		err := rows.Scan(&orderItem.ID, &orderItem.MenuItemID, &orderItem.OrderID, &customizations, &orderItem.PriceAtOrderTime, &orderItem.Quantity)
		if err != nil {
			return nil, err
		}
		if len(customizations) > 0 {
			err = json.Unmarshal(customizations, &orderItem.Customizations)
			if err != nil {
				return nil, err
			}
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

// Gets array of ordered items by order_id from database
func (repo *NewOrderRepo) GetOrderItemsArray(order_id int) ([]string, error) {
	var orderItems []string
	rows, err := repo.DB.Query(`SELECT name
	FROM order_items
	INNER JOIN menu_items USING(menu_item_id)
	WHERE order_id=$1`, order_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderItem string
		err := rows.Scan(&orderItem)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

// CheckOrder validates an order before saving
func (repo *NewOrderRepo) CheckOrder(newOrder models.Order) (int, error) {
	// Check if products exist
	for _, item := range newOrder.Items {
		exists, err := repo.ProductID_Exists(strconv.Itoa(item.MenuItemID))
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if !exists {
			return http.StatusBadRequest, errors.New("ordered item does not exist: " + strconv.Itoa(item.MenuItemID))
		}
	}

	// Check inventory requirements
	needInventory, err := repo.Get_Need_Inventory(newOrder)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	enough, err := repo.Is_Enough(needInventory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !enough {
		return http.StatusBadRequest, errors.New("not enough inventory")
	}
	return http.StatusOK, nil
}

// GetPriceAtOrderItems sets the value of price_at_order_items for the order
func (repo *NewOrderRepo) GetPriceAtOrderItems(order *models.Order) error {
	for index, orderItem := range order.Items {
		var PriceAtOrderTime float64
		err := repo.DB.QueryRow("SELECT price*$1 FROM menu_items WHERE menu_item_id=$2", orderItem.Quantity, orderItem.MenuItemID).
			Scan(&PriceAtOrderTime)
		if err != nil {
			return err
		}
		order.Items[index].PriceAtOrderTime = PriceAtOrderTime
	}
	return nil
}

// GetTotalAmount Sets the value of TotalAmount for the order
func (repo *NewOrderRepo) GetTotalAmount(order *models.Order) error {
	var total float64
	for _, orderItem := range order.Items {
		total += orderItem.PriceAtOrderTime
	}
	order.TotalAmount = total
	return nil
}

// Checks if the order exists in the database
func (repo *NewOrderRepo) IsOrderExist(order_id int) (bool, error) {
	var count int
	repo.DB.QueryRow(`SELECT COUNT(*) FROM orders
	WHERE order_id=$1
	`, order_id).Scan(&count)
	return count > 0, nil
}

// Finds the order by order_id in the database
func (repo *NewOrderRepo) GetOrder(order_id int) (models.Order, error) {
	var order models.Order
	var specialInstructions []byte // To handle JSONB data

	err := repo.DB.QueryRow(`SELECT order_id,customer_id,order_date,status,total_amount,special_instructions
	FROM orders
	WHERE order_id=$1
	`, order_id).Scan(&order.ID, &order.CustomerID, &order.CreatedAt, &order.Status, &order.TotalAmount, &specialInstructions)
	if err != nil {
		return order, err
	}

	// Unmarshal the JSONB data into the SpecialInstructions field
	if len(specialInstructions) > 0 {
		err = json.Unmarshal(specialInstructions, &order.SpecialInstructions)
		if err != nil {
			return order, err
		}
	}

	order.Items, err = repo.Get_Order_Items(order.ID)
	if err != nil {
		return order, err
	}
	return order, nil
}

// Deletes information about order from database
func (repo *NewOrderRepo) DeleteOrder(order_id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM orders
		WHERE order_id=$1
	`, order_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Closes the order and deducts required inventory from database
func (repo *NewOrderRepo) CloseOrder(order_id int) error {
	order, err := repo.GetOrder(order_id)
	if err != nil {
		return err
	}
	need_invent, err := repo.Get_Need_Inventory(order)
	if err != nil {
		return err
	}
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE orders
	SET status='closed'
	WHERE order_id=$1
	`, order_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`INSERT INTO order_status_history (order_id, status)
	VALUES($1,'closed')
	`, order_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = DefaultInventRepo(repo.DB).Use_Inventory(need_invent)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Update order information from database
func (repo *NewOrderRepo) UpdateOrder(order models.Order, id int) (int, error) {
	date := time.Now()
	code, err := repo.CheckOrder(order)
	if err != nil {
		return code, err
	}
	order.CreatedAt = date
	order.Status = "active"

	// Marshal `special_instructions` to JSON
	specialInstructionsJSON, err := json.Marshal(order.SpecialInstructions)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = repo.GetPriceAtOrderItems(&order)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = repo.GetTotalAmount(&order)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	exist, err := repo.IsOrderExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("order not found")
	}
	if exist := DefaultCustomerRepo(repo.DB).IsCustomerExist(order.CustomerID); !exist {
		return http.StatusBadRequest, errors.New("customer id is not exist")
	}
	tx, err := repo.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = tx.Exec(`DELETE FROM orders
		WHERE order_id=$1
	`, id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	order.ID = id
	items := order.Items

	// Insert order
	_, err = tx.Exec(`
		INSERT INTO orders (order_id, customer_id, total_amount, status, special_instructions)
		VALUES ($1, $2, $3, $4, $5)
	`, order.ID, order.CustomerID, order.TotalAmount, order.Status, specialInstructionsJSON)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// Insert order items
	for _, item := range items {
		customizationsJSON, err := json.Marshal(item.Customizations)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
		_, err = tx.Exec(`
			INSERT INTO order_items (menu_item_id, order_id, customizations, price_at_order_time, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`, item.MenuItemID, order.ID, customizationsJSON, item.PriceAtOrderTime, item.Quantity)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}
	_, err = tx.Exec(`INSERT INTO order_status_history(order_id,status,changed_at)
	VALUES($1,'active', $2)
	`, order.ID, date)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if tx.Commit() != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil // Commit the transaction
}

// Deducts are required inventory items from database
func (repo *NewOrderRepo) UseInventory(needInventory map[string]float64) error {
	for ingredientID, quantityUsed := range needInventory {
		_, err := repo.DB.Exec(`
			UPDATE inventory
			SET stock_level = stock_level - $1
			WHERE inventory_id = $2 AND stock_level >= $1
		`, quantityUsed, ingredientID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *NewOrderRepo) GetLastOrderId() (int, error) {
	var lastID int
	err := repo.DB.QueryRow(`SELECT order_id FROM orders
	ORDER BY order_date DESC
	LIMIT 1;
	`).Scan(&lastID)
	return lastID, err
}
