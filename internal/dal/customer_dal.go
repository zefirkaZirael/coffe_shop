package dal

import (
	"database/sql"
	"frappuccino/models"
)

type CustomerRepo interface {
	SaveCustomer(customer models.Customer) error
	IsEmailUnique(email string) (bool, error)
	GetAllCustomers() ([]models.Customer, error)
	GetCustomerByID(id int) (models.Customer, error)
	UpdateCustomer(customer models.Customer, id int) error
	DeleteCustomer(id int) error
	IsCustomerExist(id int) bool
}

type NewCustomerRepo struct {
	DB *sql.DB
}

func DefaultCustomerRepo(db *sql.DB) *NewCustomerRepo {
	return &NewCustomerRepo{DB: db}
}

// Saves Customer information to the Database
func (repo *NewCustomerRepo) SaveCustomer(customer models.Customer) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO customers(name, email, number)
	VALUES($1,$2,$3)`, customer.Name, customer.Email, customer.Number)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Checks is customer's email unique in the database
func (repo *NewCustomerRepo) IsEmailUnique(email string) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`SELECT COUNT(*)
	FROM customers
	WHERE email=$1
	`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Retrieves all Customers information from database
func (repo *NewCustomerRepo) GetAllCustomers() ([]models.Customer, error) {
	rows, err := repo.DB.Query("SELECT customer_id, name, email, number FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Number); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

// Retrieve information about customer by ID from database
func (repo *NewCustomerRepo) GetCustomerByID(id int) (models.Customer, error) {
	var customer models.Customer
	err := repo.DB.QueryRow("SELECT customer_id, name, email, number FROM customers WHERE customer_id=$1", id).
		Scan(&customer.Customer_id, &customer.Name, &customer.Email, &customer.Number)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

// Updates Customer information
func (repo *NewCustomerRepo) UpdateCustomer(customer models.Customer, id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE customers
		SET name=$1,email=$2,number=$3
		WHERE customer_id=$4
	`, customer.Name, customer.Email, customer.Number, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Delete information about Customer
func (repo *NewCustomerRepo) DeleteCustomer(id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM customers
	WHERE customer_id=$1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Check is Customer exist by ID
func (repo *NewCustomerRepo) IsCustomerExist(id int) bool {
	var count int
	repo.DB.QueryRow("SELECT COUNT(*) FROM customers WHERE customer_id=$1", id).Scan(&count)
	return count != 0
}
