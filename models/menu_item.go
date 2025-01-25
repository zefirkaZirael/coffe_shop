package models

import "time"

type Menu struct {
	ID             int                  `json:"id"`          // Matches menu_item_id
	Name           string               `json:"name"`        // Matches name
	Description    string               `json:"description"` // Matches description
	Price          float64              `json:"price"`       // Matches price
	Tags           []string             `json:"tags"`        // Matches tags
	ItemIngredient []MenuItemIngredient `json:"menuitems"`   // Matches MenuItems
	Relevance      float64              `json:"relevance"`   // Matches Relevance
}

type MenuItemIngredient struct {
	ID          int     `json:"id"`           // Matches id in menu_item_ingredients
	MenuItemID  int     `json:"menu_item_id"` // Matches menu_item_id
	InventoryID int     `json:"inventory_id"` // Matches inventory_id
	Quantity    float64 `json:"quantity"`     // Matches quantity
}

type MenuPriceHistory struct {
	ID         int       `json:"id"`           // Matches id
	MenuItemID int       `json:"menu_item_id"` // Matches menu_item_id
	OldPrice   float64   `json:"old_price"`    // Matches old_price
	NewPrice   float64   `json:"new_price"`    // Matches new_price
	ChangedAt  time.Time `json:"changed_at"`   // Matches changed_at
}
