-- +goose Up
-- Create the payments table
CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER,
    loan_id INTEGER,
    loan_bill_id INTEGER,
    amount INTEGER,
    status ENUM('PENDING', 'PROCESS', 'COMPLETED'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- Drop the payments table
DROP TABLE IF EXISTS payments;
