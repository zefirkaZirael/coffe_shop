
##Project Context

The Frappuccino project is a coffee shop management system built with PostgreSQL for efficient data management and robust querying. The goal of the project is to handle a wide range of operations related to order management, menu items, inventory, and reporting. It utilizes PostgreSQL as the database system to take advantage of its advanced features like complex queries, aggregation capabilities, and relational structure.

The project is designed to provide a scalable and maintainable solution for managing orders, inventory, and menu items, while offering advanced reporting tools. Key features include dynamic order processing, customizable menu items, real-time inventory tracking, and detailed analytics to help make data-driven business decisions.
Key Features and Components
1. Database Design and Structure

The Frappuccino project uses a relational database schema in PostgreSQL, ensuring data integrity and efficient querying. The database schema is designed to handle various aspects of the coffee shop's operations:

    Orders: Stores detailed order information, including customer details, status updates, and timestamps.

    Menu Items: Contains information about the menu, such as item names, categories, prices, and ingredients.

    Inventory: Manages stock levels for ingredients, tracking quantities and prices.

    Order Items: Tracks individual items in each order, including quantity and customization details.

    Price History and Status History: Tracks changes to menu item prices and order status transitions over time.

2. PostgreSQL Data Types

To handle specific data needs, several PostgreSQL data types are employed:

    ENUM: Used for predefined values, such as order statuses (pending, accepted, rejected), item sizes (small, medium, large), and payment methods (credit_card, cash).

    JSONB: Utilized for storing flexible, structured data such as menu item customization options, order special instructions, and customer preferences.

    Arrays: Ideal for storing multiple values in a single field, such as tags, allergens, and ingredient substitutes.

    Timestamp with Time Zone: Used to track key events like order dates, inventory updates, and price changes, with full time zone support.

3. Endpoints and SQL Queries

The project features multiple endpoints that enable users to manage orders, menu items, inventory, and generate reports. These endpoints are backed by SQL queries to interact with the PostgreSQL database. The key endpoints include:

    Order Management:

        Create, retrieve, update, and delete orders.

        Close orders and track their statuses.

    Menu Management:

        Add, update, delete, and retrieve menu items, along with their customization options and ingredients.

    Inventory Management:

        Add, update, delete, and retrieve inventory items, including tracking stock levels and prices.

    Reporting and Analytics:

        Generate reports on total sales, popular items, order history, inventory leftovers, and more.

        Perform full-text search across orders, menu items, and customer data for advanced querying.

4. Aggregation and Reporting

Advanced reporting endpoints allow users to gather insights into their operations, such as:

    Total Sales: Retrieve the total sales amount over a specified period.

    Popular Items: Identify the most frequently ordered menu items.

    Ordered Items by Period: Group ordered items by day or month for analysis.

    Inventory Leftovers: Track inventory items with sorting and pagination options.

    Batch Order Processing: Process multiple orders at once while ensuring inventory consistency.

5. Containerization

The Frappuccino project is containerized using Docker to simplify deployment and testing. With Docker Compose, the project can be set up with a single command, which automatically initializes the PostgreSQL database and runs the application in separate containers.

The project includes:

    A docker-compose.yml file to manage both the application and database containers.

    An init.sql file to initialize the database schema, including table creation, relationships, and mock data insertion.

By using Docker, the project ensures consistency across different environments, making it easy to test and deploy.
6. Mock Data and Testing

The database is populated with mock data for realistic testing scenarios. This includes:

    Menu items with various categories and prices.

    Inventory items with different quantities and prices.

    Orders with various statuses and associated items.

    Price history and order status history to simulate real-time changes.

    Inventory transactions to track ingredient usage.

This mock data ensures comprehensive testing coverage, allowing for validation of:

    Full-text search functionality.

    Date range queries for sales and inventory.

    Order status transitions.

    Price history tracking.

    Inventory usage and calculations.
