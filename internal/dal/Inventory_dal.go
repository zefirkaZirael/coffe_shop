package dal

import (
	//"encoding/json"
	"database/sql"

	//"errors"
	"frappuccino/models"
	//"io"
	//"os"
)

type InventRepo interface {
	Get_AllInventory() ([]models.InventoryItem, error)
	GetInventory(id int) (models.InventoryItem, error)
	Use_Inventory(need_inventory map[string]float64) error
	Save_Inventory(inventory models.InventoryItem) error
	IsInventExist(id int) (bool, error)
	IsInventUnique(name string) (bool, error)
	CheckReordering() ([]models.InventoryItem, error)
	Update_Inventory(inventory models.InventoryItem) error
	Delete_Inventory(id int) error
	IsInventTransactionExist(id int) (bool, error)
	GetAllInventoryTransactions() ([]models.InventoryTransaction, error)
	GetInventoryTransaction(inventory_id int) ([]models.InventoryTransaction, error)
	FillInventory(transaction models.InventoryTransaction) error
}

type NewInventRepo struct {
	DB *sql.DB
}

func DefaultInventRepo(db *sql.DB) *NewInventRepo {
	return &NewInventRepo{DB: db}
}

// Gets information about all inventory from Database
func (repo *NewInventRepo) Get_AllInventory() ([]models.InventoryItem, error) {
	rows, err := repo.DB.Query("SELECT inventory_id, name, stock_level,unit_type,last_updated,reorder_level FROM inventory")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ID, &item.Name, &item.StockLevel, &item.UnitType, &item.LastUpdated, &item.ReorderLevel); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// Get information about Inventory from database
func (repo *NewInventRepo) GetInventory(id int) (models.InventoryItem, error) {
	var item models.InventoryItem
	err := repo.DB.QueryRow(`SELECT inventory_id, name, stock_level,unit_type,last_updated,reorder_level 
	FROM inventory
	WHERE inventory_id=$1
	`, id).Scan(&item.ID, &item.Name, &item.StockLevel, &item.UnitType, &item.LastUpdated, &item.ReorderLevel)
	if err != nil {
		return item, err
	}
	return item, nil
}

// Changes the quantity of inventory based on the required ingredients
func (repo *NewInventRepo) Use_Inventory(need_inventory map[string]float64) error {
	tx, err := repo.DB.Begin() // Start a transaction
	if err != nil {
		return err
	}

	for ingredientID, quantity := range need_inventory {
		_, err := tx.Exec(
			"UPDATE inventory SET stock_level = stock_level - $1 WHERE inventory_id = $2 AND stock_level >= $3",
			quantity, ingredientID, quantity,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit() // Commit the transaction
}

// Inserts information about new inventory to the database
func (repo *NewInventRepo) Save_Inventory(inventory models.InventoryItem) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO inventory (name, stock_level, unit_type, reorder_level)
	VALUES ($1, $2, $3, $4)
	`, inventory.Name, inventory.StockLevel, inventory.UnitType, inventory.ReorderLevel)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Checks is inventory exist by ID
func (repo *NewInventRepo) IsInventExist(id int) (bool, error) {
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM inventory WHERE inventory_id = $1", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// Checks is inventory name unique from database
func (repo *NewInventRepo) IsInventUnique(name string) (bool, error) {
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM inventory WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Checks the reorder level of each inventory item in the database and returns the inventory items that need to be reordered.
func (repo *NewInventRepo) CheckReordering() ([]models.InventoryItem, error) {
	rows, err := repo.DB.Query("SELECT inventory_id,name,stock_level,reorder_level FROM inventory")
	if err != nil {
		return nil, err
	}
	var NeedToReorder []models.InventoryItem
	defer rows.Close()
	for rows.Next() {
		var inventoryItem models.InventoryItem
		if err := rows.Scan(&inventoryItem.ID, &inventoryItem.Name, &inventoryItem.StockLevel, &inventoryItem.ReorderLevel); err != nil {
			return nil, err
		}
		if inventoryItem.StockLevel <= inventoryItem.ReorderLevel {
			NeedToReorder = append(NeedToReorder, inventoryItem)
		}
	}
	if len(NeedToReorder) != 0 {
		return NeedToReorder, nil
	}
	return nil, nil
}

// Updates information about inventory from database
func (repo *NewInventRepo) Update_Inventory(inventory models.InventoryItem) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE inventory
	SET name=$1, stock_level=$2, unit_type=$3, last_updated=$4, reorder_level=$5
	WHERE inventory_id=$6
	`, inventory.Name, inventory.StockLevel, inventory.UnitType, inventory.LastUpdated, inventory.ReorderLevel, inventory.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Delete inventory information from database
func (repo *NewInventRepo) Delete_Inventory(id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM inventory
	WHERE inventory_id=$1
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Checks is inventory transaction exists by ID from database
func (repo *NewInventRepo) IsInventTransactionExist(id int) (bool, error) {
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM inventory_transactions WHERE inventory_id = $1", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// Gets information about inventory transaction by ID from database
func (repo *NewInventRepo) GetInventoryTransaction(inventory_id int) ([]models.InventoryTransaction, error) {
	var transactions []models.InventoryTransaction
	rows, err := repo.DB.Query(`SELECT transaction_id, inventory_id, price, quantity, transaction_date
	FROM inventory_transactions
	WHERE inventory_id=$1
	`, inventory_id)
	if err != nil {
		return transactions, err
	}
	for rows.Next() {
		var transaction models.InventoryTransaction
		err := rows.Scan(&transaction.ID, &transaction.Inventory_id, &transaction.Price, &transaction.Quantity, &transaction.Transaction_date)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// Gets information about all inventory transactions from database
func (repo *NewInventRepo) GetAllInventoryTransactions() ([]models.InventoryTransaction, error) {
	var transactions []models.InventoryTransaction
	rows, err := repo.DB.Query(`SELECT transaction_id, inventory_id, price, quantity, transaction_date
	FROM inventory_transactions
	`)
	if err != nil {
		return transactions, err
	}
	for rows.Next() {
		var transaction models.InventoryTransaction
		err := rows.Scan(&transaction.ID, &transaction.Inventory_id, &transaction.Price, &transaction.Quantity, &transaction.Transaction_date)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// Updates information about inventory
func (repo *NewInventRepo) FillInventory(transaction models.InventoryTransaction) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO inventory_transactions (inventory_id,price,quantity)
	VALUES($1, $2, $3)
	`, transaction.Inventory_id, transaction.Price, transaction.Quantity)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`UPDATE inventory
	SET stock_level=stock_level+$1
	WHERE inventory_id=$2
	`, transaction.Quantity, transaction.Inventory_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
