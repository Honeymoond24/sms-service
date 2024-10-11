CREATE TABLE IF NOT EXISTS services
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    service_code VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    country_id   INTEGER NOT NULL REFERENCES countries (id) ON DELETE CASCADE,
    operator_id  INTEGER DEFAULT 0 NOT NULL,
    amount       INTEGER DEFAULT 0 NOT NULL
);

INSERT INTO services (service_code, name, country_id, operator_id, amount)
VALUES ('hn', '1688', 1, 0, 1688),
       ('wj', '1хbet', 1, 0, 100),
       ('qi', '32red', 2, 0, 32)
ON CONFLICT DO NOTHING;