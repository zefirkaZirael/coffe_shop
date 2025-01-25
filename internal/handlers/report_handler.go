package handlers

import (
	"errors"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) Report_handler(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path[1:], "/")
	if len(splitted) < 2 {
		utils.Log_Err_Handler(errors.New("error URL adress"), http.StatusBadRequest, w)
		return
	}
	switch {
	case splitted[1] == "total-sales":
		code, err := h.service.Total_Sales(w)
		if err != nil {
			slog.Error("Failed to Handle Total sales Report", "Total Sales function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Total sales received succesfully")
		return
	case splitted[1] == "popular-items":
		code, err := h.service.Popular_Menu_Items(w)
		if err != nil {
			slog.Error("Failed to Handle Popular Items Report", "Popular Menu Items function: ", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Popular items received succesfully")
		return
	case splitted[1] == "numberOfOrderedItems":
		startDate := r.URL.Query().Get("startDate")
		endDate := r.URL.Query().Get("endDate")
		if startDate == "" {
			startDate = "1991-12-16"
		}
		if endDate == "" {
			currentTime := time.Now()
			endDate = currentTime.Format("2006-01-02")
		}
		code, err := h.service.GetOrderedItems(w, startDate, endDate)
		if err != nil {
			slog.Error("Failed to Handle Number of Ordered Items Report", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Number of Ordered Items Retrieved succesfully")
		return
	case splitted[1] == "search":
		var q string
		q = r.URL.Query().Get("q")
		if q == "" {
			slog.Error("Failed to Handle Full Text Search Report", errors.New("the required q key is missing "))
			utils.Log_Err_Handler(errors.New("the required q key is missing "), http.StatusBadRequest, w)
			return
		}
		filter := r.URL.Query().Get("filter")
		minPrice := r.URL.Query().Get("minPrice")
		maxPrice := r.URL.Query().Get("maxPrice")
		code, err := h.service.FullSearchReport(w, q, filter, minPrice, maxPrice)
		if err != nil {
			slog.Error("Failed to Handle Full Text Search Report", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Full Text Search Report succesfully completed")
	case splitted[1] == "orderedItemsByPeriod":
		period := r.URL.Query().Get("period")
		var month, year string
		switch period {
		case "day":
			month = r.URL.Query().Get("month")
			if month == "" {
				month = "january"
			}
		case "month":
			year = r.URL.Query().Get("year")
			if year == "" {
				year = "2024"
			}
		default:
			slog.Error("Failed to Handle Ordered Items Period", errors.New("period parameter is empty or invalid"))
			utils.Log_Err_Handler(errors.New("period parameter is empty or invalid"), http.StatusBadRequest, w)
			return
		}
		code, err := h.service.OrderedItemsPeriod(w, period, month, year)
		if err != nil {
			slog.Error("Failed to Handle Ordered Items Period", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Ordered Items By Period retrieved succesfully")
	case splitted[1] == "getLeftOvers":
		sortBy := r.URL.Query().Get("sortBy")
		pageStr := r.URL.Query().Get("page")
		var page int
		if pageStr == "" {
			page = 1
		} else {
			var err error
			page, err = strconv.Atoi(pageStr)
			if err != nil || page <= 0 {
				slog.Error("Failed to Handle LeftOvers report", errors.New("invalid page number"))
				utils.Log_Err_Handler(errors.New("invalid page number"), http.StatusBadRequest, w)
				return
			}
		}
		pageSizeStr := r.URL.Query().Get("pageSize")
		var pageSize int
		if pageSizeStr == "" {
			pageSize = 10
		} else {
			var err error
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil || pageSize <= 0 {
				slog.Error("Failed to Handle LeftOvers report", errors.New("invalid page size"))
				utils.Log_Err_Handler(errors.New("invalid page size"), http.StatusBadRequest, w)
				return
			}
		}
		if sortBy == "" {
			sortBy = "price"
		} else if sortBy != "price" && sortBy != "quantity" {
			slog.Error("Sorting Can be either: price (Sort by item price) or quantity (Sort by item quantity).")
			utils.Log_Err_Handler(errors.New("invalid sortBy value. sorting can only be 'price' or 'quantity'"), http.StatusBadRequest, w)
			return
		}
		code, err := h.service.GetLeftOvers(w, sortBy, page, pageSize)
		if err != nil {
			slog.Error("Failed to Handle Report", err)
			utils.Log_Err_Handler(err, code, w)
			return
		}
		slog.Info("Leftovers retrieved successfully")
		return
	}
}
