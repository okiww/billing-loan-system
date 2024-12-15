-- +goose Up
CREATE TABLE users
(
    id         INTEGER PRIMARY KEY,
    name       VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE loans
(
    id                  INTEGER PRIMARY KEY,
    user_id             INTEGER,
    name                VARCHAR(255),
    loan_amount         INTEGER,
    interest_percentage DECIMAL,
    status              ENUM('ACTIVE', 'DELINQUENT', 'CLOSED'),
    start_date          DATE,
    due_date            DATE,
    loan_terms_per_week INTEGER,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE loan_bills
(
    id             INTEGER PRIMARY KEY,
    loan_id        INTEGER,
    billing_date   DATE,
    billing_amount INTEGER,
    billing_number INTEGER,
    status         ENUM('PENDING', 'PAID', 'OVERDUE'),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE billing_configs
(
    id    INTEGER PRIMARY KEY,
    name  VARCHAR(255),
    value TEXT
);

ALTER TABLE loans
    ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE loan_bills
    ADD FOREIGN KEY (loan_id) REFERENCES loans (id);

-- +goose Down
ALTER TABLE loan_bills DROP FOREIGN KEY loan_id;
ALTER TABLE loans DROP FOREIGN KEY user_id;

DROP TABLE billing_configs;
DROP TABLE loan_bills;
DROP TABLE loans;
DROP TABLE users;

