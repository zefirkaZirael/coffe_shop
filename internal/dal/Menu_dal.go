package dal

import (
	"database/sql"
	"errors"
	"frappuccino/models"
	"net/http"

	"github.com/lib/pq"
)

type MenuRepo interface {
	Get_AllMenu() ([]models.Menu, error)
	GetMenu(id int) (models.Menu, error)
	Check_UniqueMenu(name string) (bool, error)
	IsMenuExist(id int) (bool, error)
	Save_Menu(menu models.Menu, ingredients []models.MenuItemIngredient) error
	Get_Ingredients(menuItemID int) ([]models.MenuItemIngredient, error)
	Check_Menu_Inventory(Menu models.Menu) (bool, error)
	GetOldPrice(menu_id int) (float64, error)
	Update_Menu(menu models.Menu, id int) (int, error)
	Delete_Menu(id int) error
	GetAllMenuPriceHistory() ([]models.MenuPriceHistory, error)
	GetMenuPriceHistory(id int) (models.MenuPriceHistory, error)
	IsPriceHistoryExist(id int) (bool, error)
}

type NewMenuRepo struct {
	DB *sql.DB
}

func DefaultMenuRepo(db *sql.DB) *NewMenuRepo {
	return &NewMenuRepo{DB: db}
}

// Get_Menu retrieves all menu items from the database
func (repo *NewMenuRepo) Get_AllMenu() ([]models.Menu, error) {
	rows, err := repo.DB.Query("SELECT menu_item_id, name, description, price, tags FROM menu_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var menu models.Menu
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Price, pq.Array(&menu.Tags)); err != nil {
			return nil, err
		}
		menu.ItemIngredient, err = repo.Get_Ingredients(menu.ID)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menus, nil
}

// Retrieves information about menu by ID from database
func (repo *NewMenuRepo) GetMenu(id int) (models.Menu, error) {
	var menu models.Menu
	err := repo.DB.QueryRow(`SELECT menu_item_id, name, description, price, tags 
	FROM menu_items
	WHERE menu_item_id=$1`, id).Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Price, pq.Array(&menu.Tags))
	if err != nil {
		return menu, err
	}
	menu.ItemIngredient, err = repo.Get_Ingredients(id)
	if err != nil {
		return menu, err
	}
	return menu, nil
}

// Retrieves the previous price of the menu item
func (repo *NewMenuRepo) GetOldPrice(menu_id int) (float64, error) {
	var oldprice float64
	err := repo.DB.QueryRow(`SELECT price FROM menu_items
	WHERE menu_item_id=$1
	`, menu_id).Scan(&oldprice)
	if err != nil {
		return oldprice, err
	}
	return oldprice, nil
}

// Check_UniqueMenu checks if a menu item name is unique
func (repo *NewMenuRepo) Check_UniqueMenu(name string) (bool, error) {
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM menu_items WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// IsMenuExist checks is menu exist in database
func (repo *NewMenuRepo) IsMenuExist(id int) (bool, error) {
	var count int
	err := repo.DB.QueryRow("SELECT COUNT(*) FROM menu_items WHERE menu_item_id = $1", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Save_Menu saves a menu item and its ingredients
func (repo *NewMenuRepo) Save_Menu(menu models.Menu, ingredients []models.MenuItemIngredient) error {
	tx, err := repo.DB.Begin() // Start a transaction
	if err != nil {
		return err
	}

	// Insert menu item
	var menuItemID int
	err = tx.QueryRow(`
		INSERT INTO menu_items (name, description, price, tags)
		VALUES ($1, $2, $3, $4) RETURNING menu_item_id
	`, menu.Name, menu.Description, menu.Price, pq.Array(menu.Tags)).Scan(&menuItemID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert ingredients
	for _, ingredient := range ingredients {
		_, err = tx.Exec(`
			INSERT INTO menu_item_ingredients (menu_item_id, inventory_id, quantity)
			VALUES ($1, $2, $3)
		`, menuItemID, ingredient.InventoryID, ingredient.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit() // Commit the transaction
}

// Checks if the ingredients of menu item are available in inventory
func (repo *NewMenuRepo) Check_Menu_Inventory(Menu models.Menu) (bool, error) {
	inventrepo := DefaultInventRepo(repo.DB)
	inventory, err := inventrepo.Get_AllInventory()
	if err != nil {
		return false, err
	}
	menuInventory := Menu.ItemIngredient
	for i := 0; i < len(menuInventory); i++ {
		var contain bool
		for j := 0; j < len(inventory); j++ {
			if inventory[j].ID == menuInventory[i].InventoryID {
				contain = true
				break
			}
		}
		if !contain {
			return false, nil
		}
	}
	return true, nil
}

// Get_Ingredients retrieves ingredients for a menu item
func (repo *NewMenuRepo) Get_Ingredients(menuItemID int) ([]models.MenuItemIngredient, error) {
	rows, err := repo.DB.Query(`
		SELECT id, inventory_id, menu_item_id, quantity
		FROM menu_item_ingredients
		WHERE menu_item_id = $1
	`, menuItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.MenuItemIngredient
	for rows.Next() {
		var ingredient models.MenuItemIngredient
		if err := rows.Scan(&ingredient.ID, &ingredient.InventoryID, &ingredient.MenuItemID, &ingredient.Quantity); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

// Updates infromation about menu item from database
func (repo *NewMenuRepo) Update_Menu(menu models.Menu, id int) (int, error) {
	menu.ID = id
	exist, err := repo.IsMenuExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusBadRequest, errors.New("menu id is not exist")
	}
	oldprice, err := repo.GetOldPrice(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	tx, err := repo.DB.Begin()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = tx.Exec(`DELETE FROM menu_items
		WHERE menu_item_id=$1
	`, id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	unique, err := repo.Check_UniqueMenu(menu.Name)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if !unique {
		tx.Rollback()
		return http.StatusBadRequest, errors.New("menu name must be unique")
	}
	exist, err = repo.Check_Menu_Inventory(menu)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if !exist {
		tx.Rollback()
		return http.StatusBadRequest, errors.New("menu item ingredient is not exist in inventory")
	}

	// Insert menu item
	_, err = tx.Exec(`
		INSERT INTO menu_items (menu_item_id, name, description, price, tags)
		VALUES ($1, $2, $3, $4, $5)
	`, menu.ID, menu.Name, menu.Description, menu.Price, pq.Array(menu.Tags))
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// Insert ingredients
	for _, ingredient := range menu.ItemIngredient {
		_, err = tx.Exec(`
			INSERT INTO menu_item_ingredients (menu_item_id, inventory_id, quantity)
			VALUES ($1, $2, $3)
		`, menu.ID, ingredient.InventoryID, ingredient.Quantity)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}
	if oldprice != menu.Price {
		_, err := tx.Exec(`INSERT INTO price_history (menu_item_id,old_price,new_price)
		VALUES ($1, $2, $3)
		`, menu.ID, oldprice, menu.Price)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}
	if tx.Commit() != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil // Commit the transaction
}

// Deletes information about menu item from database
func (repo *NewMenuRepo) Delete_Menu(id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM menu_items
		WHERE menu_item_id=$1
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Retrieves price historyof all menu items from database
func (repo *NewMenuRepo) GetAllMenuPriceHistory() ([]models.MenuPriceHistory, error) {
	var data []models.MenuPriceHistory
	rows, err := repo.DB.Query(`SELECT id, menu_item_id, old_price, new_price, changed_at
	FROM price_history
`)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var price_history models.MenuPriceHistory
		err := rows.Scan(&price_history.ID, &price_history.MenuItemID, &price_history.OldPrice, &price_history.NewPrice, &price_history.ChangedAt)
		if err != nil {
			return data, err
		}
		data = append(data, price_history)
	}
	return data, nil
}

// Retrieves price histore of menu item from database by ID
func (repo *NewMenuRepo) GetMenuPriceHistory(id int) (models.MenuPriceHistory, error) {
	var price_history models.MenuPriceHistory
	err := repo.DB.QueryRow(`SELECT id, menu_item_id, old_price, new_price, changed_at
	FROM price_history
	WHERE id=$1
	`, id).Scan(&price_history.ID, &price_history.MenuItemID, &price_history.OldPrice, &price_history.NewPrice, &price_history.ChangedAt)
	if err != nil {
		return price_history, err
	}
	return price_history, nil
}

// Checks is price history information exists by ID in database
func (repo *NewMenuRepo) IsPriceHistoryExist(id int) (bool, error) {
	var count int
	err := repo.DB.QueryRow(`SELECT COUNT(*)
	FROM price_history
	WHERE menu_item_id=$1
	`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
