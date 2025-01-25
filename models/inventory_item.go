package models

import "time"

type InventoryItem struct {
	ID           int       `json:"id"`            // Matches inventory_id
	Name         string    `json:"name"`          // Matches name
	StockLevel   float64   `json:"stock_level"`   // Matches stock_level
	UnitType     string    `json:"unit_type"`     // Matches unit_type
	LastUpdated  time.Time `json:"last_updated"`  // Matches last_updated
	ReorderLevel float64   `json:"reorder_level"` // Matches reorder_level
}

type InventoryTransaction struct {
	ID               int       `json:"id"`               // Matches transaction_id
	Inventory_id     int       `json:"inventory_id"`     // Matches inventory_id
	Price            float64   `json:"price"`            // Matches price
	Quantity         float64   `json:"Quantity"`         // Matches quantity
	Transaction_date time.Time `json:"transaction_date"` // Matches transaction_date
}
