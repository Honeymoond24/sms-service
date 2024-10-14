CREATE TABLE IF NOT EXISTS countries
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    country_code INTEGER      NOT NULL,
    name         VARCHAR(255) NOT NULL
);
INSERT INTO countries (country_code, name)
VALUES (1, 'russia'),
       (2, 'ukraine'),
       (3, 'kazakhstan'),
       (4, 'usa'),
       (5, 'germany')
ON CONFLICT DO NOTHING;

------------------------------------------------------------

CREATE TABLE IF NOT EXISTS services
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    service_code VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    operator_id  INTEGER DEFAULT 0 NOT NULL,
    amount       INTEGER DEFAULT 0 NOT NULL
);

INSERT INTO services (service_code, name, operator_id, amount)
VALUES ('vk', 'vk.com', 0, 1688),
       ('tg', 'telegram', 0, 100),
       ('ms', 'microsoft', 0, 32)
ON CONFLICT DO NOTHING;

------------------------------------------------------------

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

------------------------------------------------------------

CREATE TABLE IF NOT EXISTS phone_numbers
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    country_id INTEGER     NOT NULL,
    operator   VARCHAR(255) DEFAULT NULL,
    number     VARCHAR(31) UNIQUE NOT NULL,
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

------------------------------------------------------------

CREATE TABLE IF NOT EXISTS sms (
    id INTEGER PRIMARY KEY,
    sms_id INTEGER NOT NULL,
    phone VARCHAR(31) NOT NULL,
    phone_from VARCHAR(255) NOT NULL,
    text TEXT NOT NULL
);