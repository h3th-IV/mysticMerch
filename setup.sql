    CREATE TABLE users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id VARCHAR(255) NOT NULL,
        first_name VARCHAR(50) NOT NULL,
        last_name VARCHAR(50) NOT NULL,
        email VARCHAR(255) NOT NULL,
        phone_number VARCHAR(20) NOT NULL,
        password_hash VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

    --add unique email constraints
    ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

    CREATE TABLE products (
        id INT AUTO_INCREMENT PRIMARY KEY,
        product_id VARCHAR(255) NOT NULL,
        product_name VARCHAR(255),
        description LONGTEXT,  
        image VARCHAR(255),
        price INT,
        rating INT
    );

    CREATE TABLE carts (
        cart_id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        product_id VARCHAR(255),
        product_name VARCHAR(255),
        description LONGTEXT,
        price INT,
        rating INT,
        image VARCHAR(255),
        quantity INT,
        color VARCHAR(50),
        size VARCHAR(50),
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (product_id) REFERENCES products(id)
    );


    CREATE TABLE orders (
        order_id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        ordered_at TIMESTAMP,
        price INT,
        discount INT,
        payment_type ENUM('Electronic', 'Cash'),
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE TABLE address (
        address_id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        house_no VARCHAR(50),
        street VARCHAR(255),
        city VARCHAR(100),
        postal_code VARCHAR(20),
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    --added delete cascade so when clumns can be removed along side userfir