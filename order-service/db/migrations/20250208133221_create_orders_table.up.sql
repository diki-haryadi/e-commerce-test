CREATE TABLE orders (
    "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status ENUM('pending', 'paid', 'cancelled', 'completed') DEFAULT 'pending',
    total_amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    payment_deadline TIMESTAMP
);

CREATE INDEX idx_order_status ON orders(status);