package main

import (
	"flag"
	"frappuccino/internal/dal"
	"frappuccino/internal/handlers"
	"frappuccino/internal/service"
	"frappuccino/internal/utils"
	"frappuccino/models"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db, err := dal.CheckDB()
	if err != nil {
		slog.Error("Failed to start program", "CheckDB err:", err)
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	if *models.Prt_help {
		utils.Help()
		return
	}
	err = utils.CheckPort()
	if err != nil {
		slog.Error("Failed to start program", "CheckPort err:", err)
		return
	}

	customerRepo := dal.DefaultCustomerRepo(db)
	customerServ := service.NewDefaultServiceCustomer(*customerRepo)
	customerHandler := handlers.NewCustomerHandle(customerServ)

	orderRepo := dal.DefaultOrderRepo(db)
	orderService := service.NewDefaultOrderService(*orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	orderStatusRepo := dal.DefaultOrderStatusRepo(db)
	orderStatusService := service.DefaultOrderStatusService(*orderStatusRepo)
	orderStatusHandler := handlers.NewOrderStatusHandle(orderStatusService)

	menuRepo := dal.DefaultMenuRepo(db)
	menuService := service.NewDefaultMenuService(menuRepo)
	menuHandler := handlers.NewMenuHandle(*menuService)

	inventRepo := dal.DefaultInventRepo(db)
	inventService := service.NewDefaultInventService(inventRepo)
	inventHandler := handlers.NewInventHandle(inventService)

	reportRepo := dal.DefaultReportRepo(db)
	reportService := service.NewReportService(*reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	mux := http.NewServeMux()
	// Customers mux
	mux.HandleFunc("/customers", customerHandler.Customers_handle)
	mux.HandleFunc("/customers/{id}", customerHandler.Customers_handle)

	// Orders mux
	mux.HandleFunc("/orders", orderHandler.Order_Handle)
	mux.HandleFunc("/orders/{id}", orderHandler.Order_Handle)
	mux.HandleFunc("/orders/{id}/close", orderHandler.Order_Handle)
	mux.HandleFunc("/orders/batch-process", orderHandler.Order_Handle)

	// Orders Status mux
	mux.HandleFunc("/order-status", orderStatusHandler.OrderStatus_handle)
	mux.HandleFunc("/order-status/{id}", orderStatusHandler.OrderStatus_handle)

	// Menu mux
	mux.HandleFunc("/menu", menuHandler.Menu_Handle)
	mux.HandleFunc("/menu/{id}", menuHandler.Menu_Handle)
	mux.HandleFunc("/menu-price", menuHandler.Menu_Price_History_Handle)
	mux.HandleFunc("/menu-price/{id}", menuHandler.Menu_Price_History_Handle)

	// Inventory mux
	mux.HandleFunc("/inventory", inventHandler.Inventory_Handle)
	mux.HandleFunc("/inventory/{id}", inventHandler.Inventory_Handle)
	mux.HandleFunc("/inventory-transaction", inventHandler.InventoryTransaction_Handle)
	mux.HandleFunc("/inventory-transaction/{id}", inventHandler.InventoryTransaction_Handle)

	// Report mux
	mux.HandleFunc("/reports/total-sales", reportHandler.Report_handler)
	mux.HandleFunc("/reports/popular-items", reportHandler.Report_handler)
	mux.HandleFunc("/reports/numberOfOrderedItems", reportHandler.Report_handler)
	mux.HandleFunc("/reports/search", reportHandler.Report_handler)
	mux.HandleFunc("/reports/getLeftOvers", reportHandler.Report_handler)
	mux.HandleFunc("/reports/orderedItemsByPeriod", reportHandler.Report_handler)
	// Error mux
	mux.HandleFunc("/", utils.Err_Handler)

	slog.Info("Server is running on port " + *models.Port)
	if err := http.ListenAndServe(":"+*models.Port, mux); err != nil {
		slog.Error("Error starting server:", err)
		return
	}
}
