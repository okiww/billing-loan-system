-- +goose Up
CREATE TABLE users
(
    id         INTEGER PRIMARY KEY AUTO_INCREMENT,
    name       VARCHAR(255),
    is_delinquent BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE loans
(
    id                  INTEGER PRIMARY KEY AUTO_INCREMENT,
    user_id             INTEGER,
    name                VARCHAR(255),
    loan_amount         INT,
    loan_total_amount   INT,
    outstanding_amount  INT,
    interest_percentage DECIMAL,
    status              ENUM('ACTIVE', 'DELINQUENT', 'CLOSED'),
    start_date          DATE,
    due_date            DATE,
    loan_terms_per_week INTEGER,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_loans_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE loan_bills
(
    id                   INTEGER PRIMARY KEY AUTO_INCREMENT,
    loan_id              INTEGER,
    billing_date         DATE,
    billing_amount       INT,
    billing_total_amount INT,
    billing_number       INT,
    status               ENUM('PENDING', 'PAID', 'BILLED', 'OVERDUE'),
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_loans_loan_id FOREIGN KEY (loan_id) REFERENCES loans (id)
);

CREATE TABLE billing_configs
(
    id    INTEGER PRIMARY KEY AUTO_INCREMENT,
    name  VARCHAR(255),
    value TEXT
);

-- +goose Down
ALTER TABLE loan_bills DROP FOREIGN KEY fk_loans_loan_id;
ALTER TABLE loans DROP FOREIGN KEY fk_loans_user_id;

DROP TABLE billing_configs;
DROP TABLE loan_bills;
DROP TABLE loans;
DROP TABLE users;

