# frappuccino

## Abstract

Hello everyone. We present our project - Frappuccino. In short, it is a [hot-coffee](https://github.com/alem-platform/foundation/tree/main/hot-coffee) project, but instead of a json database, postgresql is used. Systems become more complex and modernized over time, and the hot-coffee project is no exception. While using a JSON-based database is convenient, it is not scalable and makes maintenance difficult for other developers. Moving to a PostgreSQL database is a better solution.

### Containerization Guide

Since adding a database dependency to your project makes it more challenging for auditor to run and test it, you need to containerize your service and database into separate containers.
Don't worry about it, everything made up for you. Just use command:
```bash
docker-compose up --build
```

### New Endpoints
- **Orders:**

    - `POST /orders`: Create a new order.
    - `GET /orders`: Retrieve all orders.
    - `GET /orders/{id}`: Retrieve a specific order by ID.
    - `PUT /orders/{id}`: Update an existing order.
    - `DELETE /orders/{id}`: Delete an order.
    - `POST /orders/{id}/close`: Close an order.
- **Menu Items:**

    - `POST /menu`: Add a new menu item.
    - `GET /menu`: Retrieve all menu items.
    - `GET /menu/{id}`: Retrieve a specific menu item.
    - `PUT /menu/{id}`: Update a menu item.
    - `DELETE /menu/{id}`: Delete a menu item.
- **Inventory:**

    - `POST /inventory`: Add a new inventory item.
    - `GET /inventory`: Retrieve all inventory items.
    - `GET /inventory/{id}`: Retrieve a specific inventory item.
    - `PUT /inventory/{id}`: Update an inventory item.
    - `DELETE /inventory/{id}`: Delete an inventory item.
- **Customers:**

    - `POST /customers`: Add a new customer
    - `GET /customers`: Retrieve all customers
    - `GET /customers/{id}`: Retrieve a specific customer by id
    - `PUT /customers/{id}`: Update a customer information
    - `DELETE /customers/{id}`: Delete a specific customer by id
- **Order-status:**

    - `GET /order-status` :  Retrieve all orders status history
    - `GET /order-status/{id}`: Retrieve specific order status history by ID
- **Menu-price:**

    - `GET /menu-price`: Retrieve all menu price history
    - `GET /menu-price/{id}`: Retrieve specific menu price history
- **Inventory-transaction:**

    - `POST /inventory-transaction`: Fill inventory item
    - `GET /inventory-transaction`: Retrieve all transactions
    - `GET /inventory-transaction/{id}`: Retrieve specific transaction by id

- **Aggregations:**

    - `GET /reports/total-sales`: Get the total sales amount.
    - `GET /reports/popular-items`: Get a list of popular menu items.

### We also have the following new endpoints:

#### 1. Number of ordered items
`GET /reports/numberOfOrderedItems?startDate={startDate}&endDate={endDate}`: Returns a list of ordered items and their quantities for a specified time period. If the `startDate` and `endDate` parameters are not provided, the endpoint should return data for the entire time span.
##### **Parameters**:
- `startDate` _(optional)_: The start date of the period in `YYYY-MM-DD` format.
- `endDate` _(optional)_: The end date of the period in `YYYY-MM-DD` format.

Response example:
```json
GET /reports/numberOfOrderedItems?startDate=10.11.2024&endDate=11.11.2024
HTTP/1.1 200 OK
Content-Type: application/json

{
  "latte": 109,
  "muffin": 56,
  "espresso": 120,
  "raff": 0,
  ...
}
```

#### 2. Full Text Search Report
`GET /reports/search`: Search through orders, menu items, and customers with partial matching and ranking.

##### Parameters:
- `q` _(required)_: Search query string
- `filter` _(optional)_: What to search through, can be multiple values comma-separated:
   * orders (search in customer names and order details)
   * menu (search in item names and descriptions)
   * all (default, search everywhere)
- `minPrice` _(optional)_: Minimum order/item price to include
- `maxPrice` _(optional)_: Maximum order/item price to include

Response example:
```json
GET /reports/search?q=chocolate cake&filter=menu,orders&minPrice=10
HTTP/1.1 200 OK
Content-Type: application/json

{
   "menu_items": [
       {
           "id": "12",
           "name": "Double Chocolate Cake",
           "description": "Rich chocolate layer cake",
           "price": 15.99,
           "relevance": 0.89
       },
       {
           "id": "15",
           "name": "Chocolate Cheesecake",
           "description": "Creamy cheesecake with chocolate",
           "price": 12.99,
           "relevance": 0.75
       }
   ],
   "orders": [
       {
           "id": "1234",
           "customer_name": "Alice Brown",
           "items": ["Chocolate Cake", "Coffee"],
           "total": 18.99,
           "relevance": 0.68
       }
   ],
   "total_matches": 3
}
```

#### 3. Ordered items by period
`GET /reports/orderedItemsByPeriod?period={day|month}&month={month}`: Returns the number of orders for the specified period, grouped by day within a month or by month within a year. The `period` parameter can take the value `day` or `month`. The `month` parameter is optional and used only when `period=day`.
##### **Parameters**:
- `period` _(required)_:
    - `day`: Groups data by day within the specified month.
    - `month`: Groups data by month within the specified year.
- `month` _(optional)_: Specifies the month (e.g., `october`). Used only if `period=day`.
- `year` _(optional)_: Specifies the year. Used only if `period=month`.

Response example:
```json
GET /orderedItemsByPeriod?period=day&month=october
HTTP/1.1 200 OK
Content-Type: application/json

{
    "period": "day",
    "month": "october",
    "orderedItems": [
        { "1": 109 },
        { "2": 234 },
        { "3": 198 },
        { "4": 157 },
        { "5": 223 },
        { "6": 143 },
        { "7": 256 },
        { "8": 199 },
        { "9": 275 },
        { "10": 187 },
        { "11": 234 },
        { "12": 150 },
        { "13": 178 },
        { "14": 210 },
        { "15": 202 },
        { "16": 190 },
        { "17": 260 },
        { "18": 215 },
        { "19": 240 },
        { "20": 180 },
        { "21": 300 },
        { "22": 250 },
        { "23": 199 },
        { "24": 210 },
        { "25": 220 },
        { "26": 190 },
        { "27": 170 },
        { "28": 260 },
        { "29": 230 },
        { "30": 210 },
        { "31": 180 }
    ]
}

```

Response example:
```json
GET /orderedItemsByPeriod?period=month&year=2023
HTTP/1.1 200 OK
Content-Type: application/json

{
    "period": "month",
    "year": "2023",
    "orderedItems": [
        { "january": 6528 },
        { "february": 7324 },
        { "march": 8452 },
        { "april": 7890 },
        { "may": 9103 },
        { "june": 8675 },
        { "july": 9234 },
        { "august": 8820 },
        { "september": 9345 },
        { "october": 8901 },
        { "november": 8123 },
        { "december": 9576 }
    ]
}
```

#### 4. Get leftovers
`GET /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}`: Returns the inventory leftovers in the coffee shop, including sorting and pagination options.
##### **Parameters**:
- `sortBy` _(optional)_: Determines the sorting method. Can be either:
    - `price`: Sort by item price.
    - `quantity`: Sort by item quantity.
- `page` _(optional)_: Current page number, starting from 1.
- `pageSize` _(optional)_: Number of items per page. Default value: `10`.
##### **Response**:
- Includes:
    - A list of leftovers sorted and paginated.
    - `currentPage`: The current page number.
    - `hasNextPage`: Boolean indicating whether there is a next page.
    - `totalPages`: Total number of pages.

Response example:
```json
GET /getLeftOvers?sortBy=quantity?page=1&pageSize=4
HTTP/1.1 200 OK
Content-Type: application/json

{
    "currentPage": 1,
    "hasNextPage": true,
    "pageSize": 4,
    "totalPages": 10,
    "data": [
        {
            "name": "croissant",
            "quantity": 109,
            "price": 950
        },
        {
            "name": "sugar",
            "quantity": 93,
            "price": 50
        },
        {
            "name": "muffin",
            "quantity": 63,
            "price": 350
        },
        {
            "name": "milk",
            "quantity": 1,
            "price": 200
        }
    ]
}
```

#### 5. Bulk Order Processing
`POST /orders/batch-process`: Process multiple orders simultaneously while ensuring inventory consistency. This endpoint must handle concurrent orders and maintain data integrity using transactions.

Request Body:
```json
{
   "orders": [
       {
           "customer_name": "Alice",
           "items": [
               {
                   "menu_item_id": 1,
                   "quantity": 2
               },
               {
                   "menu_item_id": 3,
                   "quantity": 1
               }
           ]
       },
       {
           "customer_name": "Bob",
           "items": [
               {
                   "menu_item_id": 2,
                   "quantity": 1
               }
           ]
       }
   ]
}
```

Response example:

```json
{
    "processed_orders": [
        {
            "order_id": 123,
            "customer_name": "Alice",
            "status": "accepted",
            "total": 15.50
        },
        {
            "order_id": 124,
            "customer_name": "Bob",
            "status": "rejected",
            "reason": "insufficient_inventory"
        }
    ],
    "summary": {
        "total_orders": 2,
        "accepted": 1,
        "rejected": 1,
        "total_revenue": 15.50,
        "inventory_updates": [
            {
                "ingredient_id": 1,
                "name": "Coffee Beans",
                "quantity_used": 100,
                "remaining": 2400
            }
        ]
    }
}
```