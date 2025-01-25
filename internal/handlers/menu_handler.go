package handlers

import (
	"encoding/json"
	"errors"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type HandlerMenu struct {
	service service.DefaultMenuService
}

func NewMenuHandle(service service.DefaultMenuService) *HandlerMenu {
	return &HandlerMenu{service: service}
}

func (h *HandlerMenu) Menu_Handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 {
		num, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("Failed to Handle Menu", "convertation error ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		id = num
	}
	switch {
	case r.Method == http.MethodPost && len(splitted) == 1:
		menu, err := h.Get_Body_Menu(r)
		if err != nil {
			slog.Error("Failed to Handle Menu", ": ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.Add_Menu(menu, menu.ItemIngredient)
		if err != nil {
			slog.Error("Failed to Handle Menu", "Add Menu function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Menu added succesfully")
		w.WriteHeader(code)
		return
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.Retrieve_All_Menu(w)
		if err != nil {
			slog.Error("Failed to Handle Menu", "Retrieve Menu function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("All menu retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.Retrieve_Menu(w, id)
		if err != nil {
			slog.Error("Failed to Handle Menu", "Retrieve Menu function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Menu retrieved succesfully")
		return
	case r.Method == http.MethodPut && len(splitted) == 2:
		menu, err := h.Get_Body_Menu(r)
		if err != nil {
			slog.Error("Failed to Handle Menu", "Get Body Menu function: ", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		code, err := h.service.Update_Menu(menu, id)
		if err != nil {
			slog.Error("Failed to Handle Menu", "Update Menu function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Menu updated succesfully")
		return
	case r.Method == http.MethodDelete && len(splitted) == 2:
		code, err := h.service.Delete_Menu(id)
		if err != nil {
			slog.Error("Failed to Handle Menu", "Delete Menu function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Menu deleted succesfully")
		w.WriteHeader(code)
		return
	default:
		utils.Log_Err_Handler(errors.New("error method in Menu"), http.StatusMethodNotAllowed, w)
		return
	}
}

func (h *HandlerMenu) Menu_Price_History_Handle(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 1 {
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	var id int
	if len(splitted) == 2 {
		num, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			slog.Error("Failed to Handle Menu Price History", "convertation error", err)
			utils.Log_Err_Handler(err, http.StatusBadRequest, w)
			return
		}
		id = num
	}
	switch {
	case r.Method == http.MethodGet && len(splitted) == 1:
		code, err := h.service.GetAllMenuPriceHistory(w)
		if err != nil {
			slog.Error("Failed to Handle Menu Price History", "Get All menu Price History fucntion: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("All Menu Price History retrieved succesfully")
		return
	case r.Method == http.MethodGet && len(splitted) == 2:
		code, err := h.service.GetMenuPriceHistory(w, id)
		if err != nil {
			slog.Error("Failed to Handle Menu Price History", "Get Menu Price History function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Menu Price History retrieved succesfully")
		return
	default:
		utils.Log_Err_Handler(errors.New("error method in Menu Price History"), http.StatusMethodNotAllowed, w)
		return
	}
}

func (h *HandlerMenu) Get_Body_Menu(r *http.Request) (models.Menu, error) {
	var menu models.Menu
	if r.Body == nil {
		return menu, errors.New("request body is empty")
	}

	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		return menu, err
	}
	if menu.ID != 0 {
		return menu, errors.New("menu id must be empty")
	}
	if menu.Name == "" {
		return menu, errors.New("menu name is missing")
	}
	if menu.Description == "" {
		return menu, errors.New("menu description is missing")
	}
	if menu.Price <= 0 {
		return menu, errors.New("menu price must be greater than 0")
	}
	if len(menu.ItemIngredient) == 0 {
		return menu, errors.New("menu ingredients field is missing")
	}
	for _, menuItem := range menu.ItemIngredient {
		if menuItem.ID != 0 {
			return menu, errors.New("menu item id field must be empty")
		}
		if menuItem.InventoryID <= 0 {
			return menu, errors.New("menu item inventory id field is missing or invalid")
		}
		if menuItem.MenuItemID != 0 {
			return menu, errors.New("menu id must be empty")
		}
		if menuItem.Quantity <= 0 {
			return menu, errors.New("menu item quantity field is missing or invalid")
		}
	}

	return menu, nil
}
