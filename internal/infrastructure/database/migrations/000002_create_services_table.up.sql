CREATE TABLE IF NOT EXISTS services
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    service_code VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    operator_id  INTEGER DEFAULT 0 NOT NULL,
    amount       INTEGER DEFAULT 0 NOT NULL
);

INSERT INTO services (service_code, name, operator_id, amount)
VALUES ('hn', '1688', 0, 1688),
       ('wj', '1хbet', 0, 100),
       ('qi', '32red', 0, 32)
ON CONFLICT DO NOTHING;