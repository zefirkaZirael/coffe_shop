package models

type PopularItems struct {
	Menu_item_id int    `json:"menu_item_id"`
	Name         string `json:"name"`
	Sale_count   int    `json:"sale_count"`
}

type TotalSale struct {
	TotalSaleAmount float64 `json:"Total_Sale_Amount"`
}

type OrderedItemsNum struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type SearchResponse struct {
	MenuItems    []Menu              `json:"menu_items"`
	Orders       []OrderSearchResult `json:"orders"`
	TotalMatches int                 `json:"total_matches"`
}

type LeftoverItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type LeftoversResponse struct {
	CurrentPage int            `json:"currentPage"`
	HasNextPage bool           `json:"hasNextPage"`
	PageSize    int            `json:"pageSize"`
	TotalPages  int            `json:"totalPages"`
	Data        []LeftoverItem `json:"data"`
}

// Структура для запроса по дням
type OrderByDayRequest struct {
	Period       string           `json:"period"`
	Month        string           `json:"month"`
	OrderedItems []OrderedItemDay `json:"orderedItems"`
}

type OrderedItemDay struct {
	Day      string `json:"day"`
	Quantity int    `json:"quantity"`
}

// Структура для запроса по месяцам
type OrderByMonthRequest struct {
	Period       string             `json:"period"`
	Year         string             `json:"year"`
	OrderedItems []OrderedItemMonth `json:"orderedItems"`
}

type OrderedItemMonth struct {
	Month    string `json:"month"`
	Quantity int    `json:"quantity"`
}
