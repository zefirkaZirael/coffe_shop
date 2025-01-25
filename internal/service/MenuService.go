package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"net/http"
)

type MenuService interface {
	Add_Menu(Menu models.Menu) (int, error)
	Retrieve_All_Menu(w http.ResponseWriter) (int, error)
	Retrieve_Menu(w http.ResponseWriter, id int) (int, error)
	Update_Menu(Menu models.Menu, id int) (int, error)
	Delete_Menu(id int) (int, error)
	GetAllMenuPriceHistory(w http.ResponseWriter) (int, error)
	GetMenuPriceHistory(w http.ResponseWriter, id int) (int, error)
}

type DefaultMenuService struct {
	repo dal.MenuRepo
}

func NewDefaultMenuService(repo dal.MenuRepo) *DefaultMenuService {
	return &DefaultMenuService{repo: repo}
}

// Add_Menu adds a new menu item
func (serv *DefaultMenuService) Add_Menu(menu models.Menu, ingredients []models.MenuItemIngredient) (int, error) {
	// Check if menu is unique
	isUnique, err := serv.repo.Check_UniqueMenu(menu.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !isUnique {
		return http.StatusBadRequest, errors.New("menu name is not unique")
	}
	exist, err := serv.repo.Check_Menu_Inventory(menu)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusBadRequest, errors.New("menu item ingredient is not exist in inventory")
	}
	// Save the menu and its ingredients
	err = serv.repo.Save_Menu(menu, ingredients)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

// Retrieve_All_Menu retrieves all menu items
func (serv *DefaultMenuService) Retrieve_All_Menu(w http.ResponseWriter) (int, error) {
	menu, err := serv.repo.Get_AllMenu()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(menu, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Retrieve_Menu retrieves a specific menu item by ID
func (serv *DefaultMenuService) Retrieve_Menu(w http.ResponseWriter, id int) (int, error) {
	exist, err := serv.repo.IsMenuExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("menu not found")
	}
	menu, err := serv.repo.GetMenu(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(menu, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// Update_Menu updates a menu item
func (serv *DefaultMenuService) Update_Menu(menu models.Menu, id int) (int, error) {
	return serv.repo.Update_Menu(menu, id)
}

// Delete_Menu deletes a menu item by ID
func (serv *DefaultMenuService) Delete_Menu(id int) (int, error) {
	exist, err := serv.repo.IsMenuExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("menu not found")
	}
	err = serv.repo.Delete_Menu(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (serv *DefaultMenuService) GetAllMenuPriceHistory(w http.ResponseWriter) (int, error) {
	data, err := serv.repo.GetAllMenuPriceHistory()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(data, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultMenuService) GetMenuPriceHistory(w http.ResponseWriter, id int) (int, error) {
	exist, err := serv.repo.IsPriceHistoryExist(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exist {
		return http.StatusNotFound, errors.New("menu price history is not exist")
	}
	data, err := serv.repo.GetMenuPriceHistory(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(data, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
