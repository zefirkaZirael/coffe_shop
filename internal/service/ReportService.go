package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"net/http"
	"strconv"
	"strings"
)

type ReportService interface {
	Total_Sales(w http.ResponseWriter) (int, error)
	Popular_Menu_Items(w http.ResponseWriter) (int, error)
	GetOrderedItems(w http.ResponseWriter, startDate, endDate string) (int, error)
	FullSearchReport(w http.ResponseWriter, q, filter, minPricestr, maxPricestr string) (int, error)
	GetLeftOvers(w http.ResponseWriter, sortBy string, page, pageSize int) (int, error)
	OrderedItemsPeriod(w http.ResponseWriter, period, monthstr, yearstr string) (int, error)
}

type DefaultReportService struct {
	repo dal.DefReportRepo
}

func NewReportService(repo dal.DefReportRepo) *DefaultReportService {
	return &DefaultReportService{repo: repo}
}

func (serv *DefaultReportService) Total_Sales(w http.ResponseWriter) (int, error) {
	total_sales, err := serv.repo.GetTotalSales()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = utils.Send_Request(total_sales, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultReportService) Popular_Menu_Items(w http.ResponseWriter) (int, error) {
	popularitems, err := serv.repo.Get_Popular_List()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(popularitems) == 0 {
		return http.StatusNotFound, errors.New("popular items list is empty")
	}
	err = utils.Send_Request(popularitems, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultReportService) GetOrderedItems(w http.ResponseWriter, startDate, endDate string) (int, error) {
	orderedItems, err := serv.repo.GetOrderedItems(startDate, endDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(orderedItems) == 0 {
		return http.StatusNotFound, errors.New("ordered items list is empty")
	}
	err = utils.Send_Request(orderedItems, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultReportService) FullSearchReport(w http.ResponseWriter, q, filter, minPricestr, maxPricestr string) (int, error) {
	var menuExist, orderExist bool
	var minPrice, maxPrice float64
	var result models.SearchResponse
	if filter != "" {
		filters := strings.Split(filter, ",")
		for i := 0; i < len(filters); i++ {
			switch filters[i] {
			case "menu":
				menuExist = true
			case "orders":
				orderExist = true
			default:
				return http.StatusBadRequest, errors.New("filter value is incorrect")
			}
		}
	}
	if !menuExist && !orderExist {
		menuExist, orderExist = true, true
	}
	var err error
	if minPricestr != "" {
		minPrice, err = strconv.ParseFloat(minPricestr, 256)
		if err != nil {
			return http.StatusBadRequest, err
		}
	} else {
		minPrice = 0
	}
	if maxPricestr != "" {
		maxPrice, err = strconv.ParseFloat(maxPricestr, 256)
		if err != nil {
			return http.StatusBadRequest, err
		}
	} else {
		maxPrice = 0
	}
	if minPrice < 0 || maxPrice < 0 {
		return http.StatusBadRequest, errors.New("min price or max price cannot be negative")
	}
	if menuExist {
		menu, err := serv.repo.FullSearchMenu(q, minPrice, maxPrice)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		result.MenuItems = menu
	}
	if orderExist {
		orders, err := serv.repo.FullSearchOrder(q, minPrice, maxPrice)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		result.Orders = orders
	}
	result.TotalMatches = len(result.MenuItems) + len(result.Orders)
	if result.TotalMatches == 0 {
		return http.StatusNotFound, errors.New("total matches is 0")
	}
	err = utils.Send_Request(result, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultReportService) GetLeftOvers(w http.ResponseWriter, sortBy string, page, pageSize int) (int, error) {
	leftovers, totalItems, err := serv.repo.GetLeftOvers(sortBy, page, pageSize)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	if page > totalPages {
		return http.StatusBadRequest, errors.New("page does not exist")
	}

	if len(leftovers) == 0 {
		return http.StatusNotFound, errors.New("no inventory leftovers found")
	}

	hasNextPage := page < totalPages

	// Prepare response
	response := models.LeftoversResponse{
		CurrentPage: page,
		HasNextPage: hasNextPage,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		Data:        leftovers,
	}

	// Send the response
	err = utils.Send_Request(response, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (serv *DefaultReportService) OrderedItemsPeriod(w http.ResponseWriter, period, monthstr, yearstr string) (int, error) {
	var year, month int
	var err error
	if period == "day" {
		monthMap := map[string]int{
			"january":   1,
			"february":  2,
			"march":     3,
			"april":     4,
			"may":       5,
			"june":      6,
			"july":      7,
			"august":    8,
			"september": 9,
			"october":   10,
			"november":  11,
			"december":  12,
		}
		var exists bool
		month, exists = monthMap[monthstr]
		if !exists {
			return http.StatusBadRequest, errors.New("month parameter is incorrect")
		}
		var orderByDay models.OrderByDayRequest
		orderByDay.Period = "day"
		orderByDay.Month = monthstr
		orderByDay.OrderedItems = []models.OrderedItemDay{}
		err = serv.repo.GetDayPeriod(month, &orderByDay)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = utils.Send_Request(orderByDay, w)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else if period == "month" {
		year, err = strconv.Atoi(yearstr)
		if err != nil {
			return http.StatusBadRequest, err
		}
		if year <= 0 {
			return http.StatusBadRequest, errors.New("year parameter must be greater than zero")
		}
		var orderByMonth models.OrderByMonthRequest
		orderByMonth.Period = "month"
		orderByMonth.Year = yearstr
		orderByMonth.OrderedItems = []models.OrderedItemMonth{}
		err = serv.repo.GetMonthPeriod(year, &orderByMonth)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = utils.Send_Request(orderByMonth, w)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}
