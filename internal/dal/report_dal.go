package dal

import (
	"database/sql"
	"frappuccino/models"
	"strings"
)

type ReportRepo interface {
	Get_Popular_List() ([]models.PopularItems, error)
	GetTotalSales() (models.TotalSale, error)
	GetOrderedItems(startDate, endDate string) ([]models.OrderedItemsNum, error)
	GetLeftOvers(sortBy string, page, pageSize int) ([]models.LeftoverItem, int, error)
	FullSearchMenu(q string, minPrice, maxPrice float64) ([]models.Menu, error)
	FullSearchOrder(q string, minPrice, maxPrice float64) ([]models.OrderSearchResult, error)
	GetDayPeriod(month int, orderRequest *models.OrderByDayRequest) error
	GetMonthPeriod(year int, orderRequest *models.OrderByMonthRequest) error
}

type DefReportRepo struct {
	DB *sql.DB
}

func DefaultReportRepo(DB *sql.DB) *DefReportRepo {
	return &DefReportRepo{DB: DB}
}

// Gets Popular Ordered items list by ID
func (repo *DefReportRepo) Get_Popular_List() ([]models.PopularItems, error) {
	var Popular_Items []models.PopularItems
	rows, err := repo.DB.Query(`SELECT menu_item_id,name ,SUM(quantity) as sale_count 
	FROM order_items INNER JOIN menu_items using (menu_item_id)
    GROUP BY menu_item_id,name
	ORDER BY sale_count DESC
    LIMIT 10;`)
	if err != nil {
		return Popular_Items, err
	}
	defer rows.Close()
	for rows.Next() {
		var popular_item models.PopularItems
		err := rows.Scan(&popular_item.Menu_item_id, &popular_item.Name, &popular_item.Sale_count)
		if err != nil {
			return Popular_Items, err
		}
		Popular_Items = append(Popular_Items, popular_item)
	}
	return Popular_Items, nil
}

// Gets summary of total amount From Database
func (repo *DefReportRepo) GetTotalSales() (models.TotalSale, error) {
	var totalsale models.TotalSale
	err := repo.DB.QueryRow(`SELECT SUM(total_amount) FROM orders WHERE status='closed'`).Scan(&totalsale.TotalSaleAmount)
	return totalsale, err
}

// Retrieves Infromation about ordered items by period
func (repo *DefReportRepo) GetOrderedItems(startDate, endDate string) ([]models.OrderedItemsNum, error) {
	var orderedItems []models.OrderedItemsNum
	rows, err := repo.DB.Query(`SELECT 
    mi.name, 
    COALESCE(SUM(oi.quantity), 0) AS total_quantity
FROM 
    menu_items mi
INNER JOIN order_items oi 
    ON oi.menu_item_id = mi.menu_item_id
INNER JOIN order_status_history osh 
    ON oi.order_id = osh.order_id
    AND osh.status = 'closed'
    AND osh.changed_at BETWEEN $1 AND $2
GROUP BY 
    mi.menu_item_id, mi.name
ORDER BY 
    total_quantity DESC;`, startDate, endDate)
	if err != nil {
		return orderedItems, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderedItem models.OrderedItemsNum
		err := rows.Scan(&orderedItem.Name, &orderedItem.Quantity)
		if err != nil {
			return orderedItems, err
		}
		orderedItems = append(orderedItems, orderedItem)
	}
	return orderedItems, nil
}

// Search menu item(s) by filter
func (repo *DefReportRepo) FullSearchMenu(q string, minPrice, maxPrice float64) ([]models.Menu, error) {
	var menus []models.Menu
	if maxPrice == 0 {
		err := repo.DB.QueryRow(`SELECT MAX(price) FROM menu_items`).Scan(&maxPrice)
		if err != nil {
			return menus, err
		}
	}
	q = strings.ToLower(q)
	rows, err := repo.DB.Query(`
SELECT 
    menu_item_id, 
    description, 
    name, 
    price,  
    ROUND(
        CASE 
            WHEN LENGTH(REPLACE(LOWER(CONCAT(name, description)), $1, '')) = 0
            THEN 1 
            WHEN LENGTH(REPLACE(CONCAT(name, description),$1,''))*1.0<LENGTH($1)
            THEN LENGTH(REPLACE(LOWER(CONCAT(name, description)),$1,''))*1.0/LENGTH($1)
            ELSE LENGTH($1)*1.0/LENGTH(REPLACE(LOWER(CONCAT(name, description)),$1,''))
        END,
        3
    ) AS relevance
FROM 
    menu_items
WHERE 
    (LOWER(name) LIKE '%'||$1||'%' OR LOWER(description) LIKE '%'||$1||'%') 
    AND price BETWEEN $2 AND $3
ORDER BY 
	relevance DESC
	`, q, minPrice, maxPrice)
	if err != nil {
		return menus, err
	}
	defer rows.Close()
	for rows.Next() {
		var menu models.Menu
		err := rows.Scan(&menu.ID, &menu.Description, &menu.Name, &menu.Price, &menu.Relevance)
		if err != nil {
			return menus, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

// Search Order information by Filter
func (repo *DefReportRepo) FullSearchOrder(q string, minPrice, maxPrice float64) ([]models.OrderSearchResult, error) {
	var orders []models.OrderSearchResult
	if maxPrice == 0 {
		err := repo.DB.QueryRow(`SELECT MAX(total_amount) FROM orders`).Scan(&maxPrice)
		if err != nil {
			return orders, err
		}
	}
	q = strings.ToLower(q)
	rows, err := repo.DB.Query(`WITH relevance as (SELECT ROUND(
        CASE 
            WHEN LENGTH(REPLACE(LOWER(string_agg(mi.name, '')||string_agg(description,'')||string_agg(c.name,'')), $1, '')) = 0
            THEN 1 
			WHEN LENGTH($1)>LENGTH(REPLACE(LOWER(string_agg(mi.name,'')||string_agg(description,'')||string_agg(c.name,'')),$1,''))
			THEN LENGTH(REPLACE(LOWER(string_agg(mi.name, '')||string_agg(description,'')||string_agg(c.name,'')),$1,''))*1.0/LENGTH($1)
            ELSE
            LENGTH($1)*1.0/LENGTH(REPLACE(LOWER(string_agg(mi.name, '')||string_agg(description,'')||string_agg(c.name,'')),$1,''))
        END,
        3
    ) 
AS relevance,
order_id
	FROM order_items oi
		INNER JOIN menu_items mi USING(menu_item_id)
		INNER JOIN orders o USING(order_id)
		INNER JOIN customers c USING(customer_id)
					WHERE (LOWER(mi.name) LIKE '%'||$1||'%' OR LOWER(mi.description) LIKE '%'||$1||'%' OR LOWER(c.name) LIKE '%'||$1||'%') 
					GROUP BY order_id)

  SELECT 
		o.order_id,c.name,o.total_amount ,r.relevance
	FROM 
		orders o INNER JOIN customers c 
	ON 
		o.customer_id=c.customer_id
  INNER JOIN relevance r ON o.order_id = r.order_id
   WHERE total_amount BETWEEN $2  and $3
   ORDER BY relevance DESC;
	`, q, minPrice, maxPrice)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	orderRepo := DefaultOrderRepo(repo.DB)
	for rows.Next() {
		var order models.OrderSearchResult
		err := rows.Scan(&order.ID, &order.CustomerName, &order.Total, &order.Relevance)
		if err != nil {
			return orders, err
		}
		order.Items, err = orderRepo.GetOrderItemsArray(order.ID)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Gets paginated inventory leftovers
func (repo *DefReportRepo) GetLeftOvers(sortBy string, page, pageSize int) ([]models.LeftoverItem, int, error) {
	var totalItems int
	var leftovers []models.LeftoverItem

	// Query to get paginated inventory leftovers
	/*
			SELECT
		    menu_item_id,
		    description,
		    name,
		    price,
		    ts_rank(
		        to_tsvector('english', name || ' ' || description),
		        plainto_tsquery('english', $1)
		    ) AS relevance
		FROM
		    menu_items
		WHERE
		    to_tsvector('english', name || ' ' || description) @@ plainto_tsquery('english', $1)
		    AND price BETWEEN $2 AND $3
		ORDER BY
		    relevance DESC;

	*/
	rows, err := repo.DB.Query(`WITH ingredient_prices AS (
    SELECT 
        it.inventory_id,
        COALESCE(SUM(it.price) / NULLIF(SUM(it.quantity), 0), 0) AS avg_price
    FROM 
        inventory_transactions it
    GROUP BY 
        it.inventory_id
),
sorted_inventory AS (
    SELECT 
        i.name, 
        i.stock_level AS quantity,
        COALESCE(ip.avg_price, 0) AS price,
        ROW_NUMBER() OVER (
            ORDER BY 
                CASE 
                    WHEN $1 = 'price' THEN COALESCE(ip.avg_price, 0)
                    WHEN $1 = 'quantity' THEN i.stock_level
                END ASC
        ) AS row_num
    FROM 
        inventory i
    LEFT JOIN 
        ingredient_prices ip 
        ON i.inventory_id = ip.inventory_id
)
SELECT 
    name, 
    quantity, 
    price
FROM 
    sorted_inventory
WHERE 
    row_num > (($2 - 1) * $3) 
    AND row_num <= ($2 * $3)
ORDER BY row_num;
	`, sortBy, page, pageSize)
	if err != nil {
		return leftovers, 0, err
	}
	defer rows.Close()

	// Parse query results
	for rows.Next() {
		var leftover models.LeftoverItem
		err := rows.Scan(&leftover.Name, &leftover.Quantity, &leftover.Price)
		if err != nil {
			return leftovers, 0, err
		}
		leftovers = append(leftovers, leftover)
	}

	err = repo.DB.QueryRow(`SELECT COUNT(*) FROM inventory`).Scan(&totalItems)
	if err != nil {
		return leftovers, 0, err
	}

	return leftovers, totalItems, nil
}

// Retrieve of sum of totally ordered items by day period
func (repo *DefReportRepo) GetDayPeriod(month int, orderRequest *models.OrderByDayRequest) error {
	rows, err := repo.DB.Query(`SELECT EXTRACT(day FROM order_date) ,SUM(quantity)
    FROM orders 
    INNER JOIN order_items USING(order_id)
    WHERE EXTRACT(month FROM order_date)=$1 AND EXTRACT (YEAR FROM order_date)=2024 AND status='closed'
    GROUP BY EXTRACT(day FROM order_date)
    ORDER BY EXTRACT(day FROM order_date);`, month)
	if err != nil {
		return err
	}
	defer rows.Close()
	var orderedItems []models.OrderedItemDay
	for rows.Next() {
		var orderedItem models.OrderedItemDay
		err := rows.Scan(&orderedItem.Day, &orderedItem.Quantity)
		if err != nil {
			return err
		}
		orderedItems = append(orderedItems, orderedItem)
	}
	orderRequest.OrderedItems = orderedItems
	return nil
}

// Retrieve of sum of totally ordered items by Month period
func (repo *DefReportRepo) GetMonthPeriod(year int, orderRequest *models.OrderByMonthRequest) error {
	rows, err := repo.DB.Query(`SELECT TO_CHAR(order_date, 'FMMonth') AS month_name, SUM(quantity)
	FROM orders
	INNER JOIN order_items USING(order_id)
	WHERE EXTRACT(year FROM order_date) = $1 AND status='closed'
	GROUP BY EXTRACT(month FROM order_date), TO_CHAR(order_date, 'FMMonth')
	ORDER BY EXTRACT(month FROM order_date);`, year)
	if err != nil {
		return err
	}
	defer rows.Close()
	var orderedItems []models.OrderedItemMonth
	for rows.Next() {
		var orderedItem models.OrderedItemMonth
		err := rows.Scan(&orderedItem.Month, &orderedItem.Quantity)
		if err != nil {
			return err
		}
		orderedItems = append(orderedItems, orderedItem)
	}
	orderRequest.OrderedItems = orderedItems
	return nil
}
