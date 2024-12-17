-- +goose Up
INSERT INTO `users` (`id`, `is_delinquent`, `name`, `created_at`, `updated_at`)
VALUES
    (1, 1, 'John Doe', '2024-12-16 19:18:13', '2024-12-17 00:12:00');

INSERT INTO `billing_configs` (`id`, `name`, `value`)
VALUES
    (1, 'loan_interest_percentage', '{"is_active":true,"value":10}'),
    (2, 'loan_term_per_week', '{"is_active":true,"value":50}');



-- +goose Down
DELETE FROM `users` WHERE `id` = 123;
DELETE FROM `billing_configs` WHERE `id` IN (1, 2);
