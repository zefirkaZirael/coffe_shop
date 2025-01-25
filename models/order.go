package models

import "time"

type Order struct {
	ID                  int                    `json:"id"`                   // Matches order_id
	CustomerID          int                    `json:"customer_id"`          // Matches customer_id
	TotalAmount         float64                `json:"total_amount"`         // Matches total_amount
	Status              string                 `json:"status"`               // Matches status ENUM
	CreatedAt           time.Time              `json:"created_at"`           // Matches order_date
	SpecialInstructions map[string]interface{} `json:"special_instructions"` // Matches special_instructions (JSONB)
	Items               []OrderItem            `json:"items"`                // Linked items from order_items
}

type OrderItem struct {
	ID               int                    `json:"id"`                  // Matches order_item_id
	MenuItemID       int                    `json:"menu_item_id"`        // Matches menu_item_id
	OrderID          int                    `json:"order_id"`            // Matches order_id
	Customizations   map[string]interface{} `json:"customizations"`      // Matches customizations (JSONB)
	PriceAtOrderTime float64                `json:"price_at_order_time"` // Matches price_at_order_time
	Quantity         int                    `json:"quantity"`            // Matches quantity

}

type OrderStatus struct {
	Id        int       `json:"id"`
	Order_id  int       `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderSearchResult struct {
	ID           int      `json:"id"`
	CustomerName string   `json:"customer_name"`
	Total        float64  `json:"total_amount"`
	Items        []string `json:"items"`
	Relevance    float64  `json:"relevance"`
}

type ProcessedOrder struct {
	OrderID    int     `json:"order_id"`
	CustomerID int     `json:"customer_id"` // Matches customer_id
	Status     string  `json:"status"`
	Total      float64 `json:"total"`
	Reason     string  `json:"reason,omitempty"`
}

type InventoryUpdate struct {
	IngredientID string  `json:"ingredient_id"`
	Name         string  `json:"name,omitempty"`
	QuantityUsed float64 `json:"quantity_used"`
	Remaining    float64 `json:"remaining"`
}

type BatchProcessSummary struct {
	TotalOrders      int               `json:"total_orders"`
	Accepted         int               `json:"accepted"`
	Rejected         int               `json:"rejected"`
	TotalRevenue     float64           `json:"total_revenue"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}

type BatchProcessResponse struct {
	ProcessedOrders []ProcessedOrder    `json:"processed_orders"`
	Summary         BatchProcessSummary `json:"summary"`
}

type BatchProcessRequest struct {
	Orders []Order `json:"orders"` // Массив заказов
}
