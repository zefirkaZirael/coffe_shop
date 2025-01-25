package models

type Customer struct {
	Customer_id int    `json:"id"`     // Matches customer_id
	Name        string `json:"name"`   // Matches name
	Email       string `json:"email"`  // Matches email
	Number      string `json:"number"` // Matches number
}
