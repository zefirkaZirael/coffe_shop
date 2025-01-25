CREATE TYPE order_status_enum as ENUM('active','closed');


CREATE TABLE customers(
    customer_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    number VARCHAR(20)
);

CREATE TABLE orders(
    order_id SERIAL PRIMARY KEY, 
    customer_id INT REFERENCES customers(customer_id) ON DELETE CASCADE,
    order_date TIMESTAMPTZ DEFAULT NOW(),
    status order_status_enum NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL CHECK(total_amount>0),
    special_instructions JSONB
);

CREATE TABLE menu_items(
    menu_item_id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    name VARCHAR(100) NOT NULL ,
    price DECIMAL(10,2) NOT NULL CHECK(price>0),
    tags TEXT[]
);

CREATE TABLE order_items(
    order_item_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    customizations JSONB,
    price_at_order_time DECIMAL(10,2) NOT NULL CHECK(price_at_order_time>0),
    quantity INT NOT NULL CHECK (quantity >0)
);

CREATE TABLE inventory(
    inventory_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    stock_level DECIMAL(10,2) NOT NULL CHECK(stock_level>0),
    unit_type VARCHAR(50) NOT NULL,
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    reorder_level DECIMAL(10,2) NOT NULL CHECK(reorder_level>0)
);

CREATE TABLE menu_item_ingredients(
    id SERIAL PRIMARY KEY,
    inventory_id INT REFERENCES inventory(inventory_id) ON DELETE CASCADE,
    menu_item_id INT REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL CHECK(quantity>0)
);

CREATE TABLE inventory_transactions(
    transaction_id SERIAL PRIMARY KEY,
    inventory_id INT REFERENCES inventory(inventory_id) ON DELETE CASCADE,
    price DECIMAL(10,2) NOT NULL CHECK(price>0),
    quantity DECIMAL(10,2) NOT NULL CHECK(quantity>0),
    transaction_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE order_status_history(
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    status order_status_enum NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
    old_price DECIMAL(10,2) CHECK(old_price>0),
    new_price DECIMAL(10,2) CHECK(new_price>0),
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- orders
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_order_date ON orders(order_date);

-- customers
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_name_ft ON customers USING GIN (to_tsvector('english', name));

-- order_items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item_id ON order_items(menu_item_id);

-- menu_items
CREATE INDEX idx_menu_items_name_ft ON menu_items USING GIN (to_tsvector('english', name));
CREATE INDEX idx_menu_items_description_ft ON menu_items USING GIN (to_tsvector('english', description));
CREATE INDEX idx_menu_items_tags ON menu_items USING GIN (tags);

-- inventory
CREATE INDEX idx_inventory_name ON inventory(name);
CREATE INDEX idx_inventory_stock_level ON inventory(stock_level);

-- inventory_transactions
CREATE INDEX idx_inventory_transactions_inventory_id ON inventory_transactions(inventory_id);
CREATE INDEX idx_inventory_transactions_date ON inventory_transactions(transaction_date);

-- order_status_history
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_order_status_history_composite ON order_status_history(order_id, changed_at);

-- price_history
CREATE INDEX idx_price_history_menu_item_id ON price_history(menu_item_id);
CREATE INDEX idx_price_history_changed_at ON price_history(changed_at);

INSERT INTO customers(name,email,number)
VALUES
    ('John Doe', 'john.doe@example.com', '1234567890'),
    ('Jane Smith', 'jane.smith@example.com', '0987654321'),
    ('Alice Johnson', 'alice.j@example.com', '5556667777'),
    ('Bob Brown', 'bob.brown@example.com', '4445556666'),
    ('Charlie Black', 'charlie.black@example.com', '1112223333'),
    ('Emily White', 'emily.white@example.com', '2223334444'),
    ('Daniel Green', 'daniel.green@example.com', '3334445555'),
    ('Sophia Blue', 'sophia.blue@example.com', '4445556666'),
    ('Michael Grey', 'michael.grey@example.com', '5556667777'),
    ('Olivia Purple', 'olivia.purple@example.com', '6667778888'),
    ('Ethan Orange', 'ethan.orange@example.com', '7778889999'),
    ('Isabella Pink', 'isabella.pink@example.com', '8889990000'),
    ('Liam Yellow', 'liam.yellow@example.com', '9990001111'),
    ('Mia Cyan', 'mia.cyan@example.com', '0001112222'),
    ('Noah Indigo', 'noah.indigo@example.com', '1112223333'),
    ('Ava Violet', 'ava.violet@example.com', '2223334444'),
    ('Lucas Silver', 'lucas.silver@example.com', '3334445555'),
    ('Sophia Gold', 'sophia.gold@example.com', '4445556666'),
    ('James Brown', 'james.brown@example.com', '5556667777'),
    ('Charlotte Pink', 'charlotte.pink@example.com', '6667778888'),
    ('Benjamin Red', 'benjamin.red@example.com', '7778889999'),
    ('Amelia White', 'amelia.white@example.com', '8889990000'),
    ('Elijah Green', 'elijah.green@example.com', '9990001111'),
    ('Harper Blue', 'harper.blue@example.com', '0001112222'),
    ('William Grey', 'william.grey@example.com', '1112223333'),
    ('Evelyn Purple', 'evelyn.purple@example.com', '2223334444'),
    ('Henry Orange', 'henry.orange@example.com', '3334445555'),
    ('Abigail Pink', 'abigail.pink@example.com', '4445556666'),
    ('Alexander Yellow', 'alexander.yellow@example.com', '5556667777'),
    ('Sofia Cyan', 'sofia.cyan@example.com', '6667778888');

    INSERT INTO orders (customer_id, total_amount, status, special_instructions,order_date)
VALUES
    (1, 10.50, 'active', '{"notes": "extra sugar"}','2020-08-16'),
    (2, 20.75, 'closed', '{"notes": "no milk"}','2022-01-17'),
    (3, 15.25, 'active', '{"notes": "decaf"}','2024-02-12'),
    (4, 30.00, 'closed', '{"notes": "extra hot"}','2024-01-16'),
    (5, 12.50, 'active', '{"notes": "half foam"}','2024-02-19'),
    (6, 25.00, 'closed', '{"notes": "syrup on the side"}','2024-07-10'),
    (7, 17.50, 'active', '{"notes": "with ice"}','2024-03-01'),
    (8, 19.00, 'closed', '{"notes": "no cream"}','2020-01-16'),
    (9, 22.50, 'active', '{"notes": "small size"}','2024-04-28'),
    (10, 18.00, 'closed', '{"notes": "extra cinnamon"}','2024-05-18'),
    (11, 14.75, 'active', '{"notes": "with oat milk"}','2024-06-27'),
    (12, 20.00, 'closed', '{"notes": "light sugar"}','2024-08-22'),
    (13, 16.00, 'active', '{"notes": "double shot"}','2024-07-16'),
    (14, 27.50, 'closed', '{"notes": "no foam"}','2024-08-24'),
    (15, 21.00, 'active', '{"notes": "extra whipped cream"}','2024-09-19'),
    (16, 11.50, 'closed', '{"notes": "hot, no ice"}','2024-10-12'),
    (17, 28.00, 'active', '{"notes": "no whipped cream"}','2024-11-15'),
    (18, 23.50, 'closed', '{"notes": "extra shot"}','2024-01-01'),
    (19, 13.75, 'active', '{"notes": "no syrup"}','2021-01-02'),
    (20, 29.50, 'closed', '{"notes": "extra caramel"}','2024-01-03'),    
    (21, 24.00, 'active', '{"notes": "extra hot water"}','2024-01-05'),
    (22, 12.00, 'closed', '{"notes": "no sugar"}','2024-01-06'),
    (23, 16.25, 'active', '{"notes": "small cup"}','2024-01-07'),
    (24, 17.00, 'closed', '{"notes": "larger size"}','2024-01-08'),
    (25, 14.25, 'active', '{"notes": "light on milk"}','2024-01-09'),
    (26, 30.00, 'closed', '{"notes": "decaf, no sugar"}','2024-01-10'),
    (27, 15.50, 'active', '{"notes": "extra foam"}','2024-01-11'),
    (28, 22.00, 'closed', '{"notes": "extra caramel"}','2024-01-12'),
    (29, 18.75, 'active', '{"notes": "no ice"}','2024-01-13'),
    (30, 20.50, 'closed', '{"notes": "decaf, no milk"}','2024-01-14');


INSERT INTO menu_items (description, name, price, tags) 
VALUES
    ('A strong, black coffee made by forcing steam through ground coffee beans', 'Espresso', 2.50, ARRAY['coffee', 'hot']),
    ('Espresso topped with steamed milk foam', 'Cappuccino', 3.00, ARRAY['coffee', 'hot']),
    ('Espresso with steamed milk', 'Latte', 3.50, ARRAY['coffee', 'hot']),
    ('Espresso with hot water', 'Americano', 2.75, ARRAY['coffee', 'hot']),
    ('Espresso with a small amount of steamed milk', 'Macchiato', 2.85, ARRAY['coffee', 'hot']),
    ('Cold coffee served over ice', 'Iced Coffee', 3.00, ARRAY['coffee', 'cold']),
    ('Espresso with caramel and steamed milk', 'Caramel Macchiato', 4.00, ARRAY['coffee', 'hot']),
    ('Slow-steeped cold coffee served chilled', 'Cold Brew', 3.50, ARRAY['coffee', 'cold']),
    ('Freshly baked bagel with cream cheese', 'Bagel', 2.00, ARRAY['food', 'breakfast']),
    ('Blueberry muffin with a crunchy topping', 'Muffin', 2.50, ARRAY['food', 'dessert']);

INSERT INTO order_items (menu_item_id, order_id, customizations, price_at_order_time, quantity)
VALUES
    (1, 1, '{"size": "large1"}', 2.50, 3), 
    (4, 4,        '{}',         2.75, 5),
    (1, 1, '{"size": "large1"}', 2.50, 3), 
    (1, 1, '{"size": "large1"}', 2.50, 3), 
    (2, 2, '{"foam": "none"}', 3.00, 4),
    (3, 3, '{"milk": "soy"}', 3.50, 2), 
    (4, 4, '{"water": "extra"}', 2.75, 5), 
    (5, 5, '{"milk": "half"}', 2.85, 3), 
    (6, 6, '{"ice": "light"}', 3.00, 2), 
    (7, 7, '{"caramel": "extra"}', 4.00, 4), 
    (8, 8, '{"brew": "strong"}', 3.50, 3), 
    (9, 9, '{"spread": "butter"}', 2.00, 1),
    (10, 10, '{"topping": "extra crunch"}', 2.50, 2),
    (10, 10, '{"topping": "extra crunch"}', 2.50, 2),
    (1, 11, '{"shot": "double"}', 2.50, 2),
    (2, 12, '{"milk": "none"}', 3.00, 3),
    (3, 13, '{"decaf": true}', 3.50, 2), 
    (4, 14, '{"water": "less"}', 2.75, 4),
    (5, 15, '{"foam": "extra"}', 2.85, 5),
    (6, 16, '{"size": "medium"}', 3.00, 2),
    (7, 17, '{"caramel": "none"}', 4.00, 3), 
    (8, 18, '{"ice": "no"}', 3.50, 4),
    (9, 19, '{"spread": "none"}', 2.00, 1),
    (10, 20, '{"topping": "extra sugar"}', 2.50, 2),
    (1, 21, '{"brew": "bold"}', 2.50, 3), 
    (2, 22, '{"milk": "skim"}', 3.00, 3),
    (3, 23, '{"shot": "single"}', 3.50, 2),
    (4, 24, '{"water": "extra hot"}', 2.75, 5), 
    (5, 25, '{"milk": "light"}', 2.85, 3),
    (6, 26, '{"size": "small"}', 3.00, 1),
    (7, 27, '{"caramel": "light"}', 4.00, 4), 
    (8, 28, '{"brew": "mild"}', 3.50, 3),
    (9, 29, '{"spread": "cream cheese"}', 2.00, 2),
    (10, 30, '{"topping": "extra cinnamon"}', 2.50, 1);

INSERT INTO inventory (name, stock_level, unit_type, reorder_level) 
VALUES
    ('Espresso Beans', 50000.00, 'kg', 10.00),
    ('Milk', 30000.00, 'liters', 5.00),
    ('Sugar', 100000.00, 'kg', 20.00),
    ('Flavored Syrup', 25000.00, 'liters', 5.00),
    ('Cinnamon', 5000.00, 'kg', 1.00),
    ('Caramel Syrup', 15000.00, 'liters', 3.00),
    ('Butter', 10000.00, 'kg', 2.00),
    ('Cream Cheese', 20000.00, 'kg', 4.00),
    ('Bagels', 6000.00, 'pieces', 10.00),
    ('Muffins', 4000.00, 'pieces', 8.00),
    ('Cold Brew Coffee', 35000.00, 'liters', 7.00),
    ('Iced Coffee', 4000.00, 'liters', 6.00),
    ('Macchiato Syrup', 1000.00, 'liters', 2.00),
    ('Milk Foam', 2500.00, 'liters', 4.00),
    ('Vanilla Syrup', 1200.00, 'liters', 2.00),
    ('Coffee Cups', 5000.00, 'pieces', 50.00),
    ('Straws', 10000.00, 'pieces', 100.00),
    ('Coffee Filters', 2000.00, 'pieces', 20.00),
    ('Napkins', 100000.00, 'pieces', 100.00);

INSERT INTO menu_item_ingredients (inventory_id, menu_item_id, quantity)
VALUES
    (1, 1, 10), 
    (1, 2, 8), 
    (14, 2, 10), 
    (1, 3, 8),  
    (2, 3, 15), 
    (1, 4, 10), 
    (1, 5, 8), 
    (14, 5, 5),  
    (1, 6, 8), 
    (2, 6, 10), 
    (1, 7, 8), 
    (6, 7, 5), 
    (1, 8, 8),
    (9, 9, 1),  
    (8, 9, 2), 
    (10, 10, 1), 
    (3, 10, 1); 


INSERT INTO inventory_transactions (inventory_id, price, quantity, transaction_date) 
VALUES
    (1, 2.50, 5.00, '2024-12-01'),
    (2, 3.00, 10.00, '2024-12-01'),
    (3, 1.00, 50.00, '2024-12-02'),
    (4, 2.00, 15.00, '2024-12-03'),
    (5, 3.00, 20.00, '2024-12-04'),
    (6, 4.00, 30.00, '2024-12-05');

INSERT INTO order_status_history (order_id, status, changed_at) 
VALUES
    (1, 'active', '2020-08-16'),
    (2, 'active', '2022-01-17'),
    (3, 'active', '2024-02-12'),
    (4, 'active', '2024-01-16'),
    (5, 'active', '2024-02-19'),
    (6, 'active', '2024-07-10'),
    (7, 'active', '2024-03-01'),
    (8, 'active', '2020-01-16'),
    (9, 'active', '2024-04-28'),
    (10, 'active', '2024-05-18'),
    (11, 'active', '2024-06-27'),
    (12, 'active', '2024-08-22'),
    (13, 'active', '2024-07-16'),
    (14, 'active', '2024-08-24'),
    (15, 'active', '2024-09-19'),
    (16, 'active', '2024-10-12'),
    (17, 'active', '2024-11-15'),
    (18, 'active', '2024-01-01'),
    (19, 'active', '2021-01-02'),
    (20, 'active', '2024-01-03'),
    (21, 'active', '2024-01-05'),
    (22, 'active', '2024-01-06'),
    (23, 'active', '2024-01-07'),
    (24, 'active', '2024-01-08'),
    (25, 'active', '2024-01-09'),
    (26, 'active', '2024-01-10'),
    (27, 'active', '2024-01-11'),
    (28, 'active', '2024-01-12'),
    (29, 'active', '2024-01-13'),
    (30, 'active', '2024-01-14');
  
  INSERT INTO order_status_history (order_id, status, changed_at) 
VALUES
    (2, 'closed', '2022-01-18'),
    (4, 'closed', '2024-01-17'),
    (6, 'closed', '2024-07-11'),
    (8, 'closed', '2020-01-17'),
    (10, 'closed', '2024-05-19'),
    (12, 'closed', '2024-08-23'),
    (14, 'closed', '2024-08-25'),
    (16, 'closed', '2024-10-13'),
    (18, 'closed', '2024-01-02'),
    (20, 'closed', '2024-01-04'),
    (22, 'closed', '2024-01-07'),
    (24, 'closed', '2024-01-09'),
    (26, 'closed', '2024-01-11'),
    (28, 'closed', '2024-01-13'),
    (30, 'closed', '2024-01-15');

INSERT INTO price_history (menu_item_id, old_price, new_price, changed_at) 
VALUES
    (1, 2.50, 2.75, '2024-06-01'),
    (2, 3.00, 3.20, '2024-07-01'),
    (3, 3.50, 3.60, '2024-08-01'),
    (4, 2.75, 3.00, '2024-09-01'),
    (5, 2.85, 3.00, '2024-10-01'),
    (6, 3.00, 3.20, '2024-11-01');