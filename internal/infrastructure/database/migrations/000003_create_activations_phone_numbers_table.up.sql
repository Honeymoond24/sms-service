CREATE TABLE IF NOT EXISTS activations
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    sum_price  REAL    NOT NULL DEFAULT 0,
    status     INTEGER NOT NULL REFERENCES activations_statuses (id),
    phone_id   INTEGER NOT NULL REFERENCES phone_numbers (id),
    service_id INTEGER NOT NULL REFERENCES services (id)
);

CREATE TABLE IF NOT EXISTS activations_statuses
(
    id   INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO activations_statuses (id, name)
VALUES (-1, 'ACCESS_CANCEL'),
       (1, 'ACCESS_READY'),
       (3, 'ACCESS_RETRY_GET'),
       (6, 'ACCESS_ACTIVATION'),
       (8, 'ACCESS_CANCEL');

CREATE TABLE IF NOT EXISTS phone_numbers
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    country_id INTEGER     NOT NULL,
    operator   VARCHAR(255) DEFAULT NULL,
    number     VARCHAR(31) NOT NULL,
    FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE CASCADE
);

INSERT INTO phone_numbers (country_id, operator, number)
VALUES (1, NULL, '79185556601'),
       (1, NULL, '79185556602'),
       (2, NULL, '79185556603'),
       (2, NULL, '79185556604'),
       (3, NULL, '79185556605'),
       (3, NULL, '79185556606'),
       (4, NULL, '79185556607'),
       (5, NULL, '79185556608'),
       (5, NULL, '79185556609'),
       (5, NULL, '79185556610');
