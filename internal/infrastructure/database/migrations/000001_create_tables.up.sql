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
    name         VARCHAR(255) NOT NULL
);

INSERT INTO services (service_code, name)
VALUES ('vk', 'vk.com'),
       ('tg', 'telegram'),
       ('fb', 'facebook'),
       ('wa', 'whatsapp'),
       ('ok', 'ok.ru'),
       ('tw', 'twitter'),
       ('in', 'instagram'),
       ('yt', 'youtube'),
       ('gd', 'google'),
       ('ms', 'microsoft')
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
    country_id INTEGER            NOT NULL,
    number     INTEGER(31) UNIQUE NOT NULL,
    FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE CASCADE
);

INSERT INTO phone_numbers (country_id, number)
VALUES
(1, 71000001),
(1, 71000002),
(1, 71000003),
(1, 71000004),
(1, 71000005),
(1, 71000006),
(1, 71000007),
(1, 71000008),
(1, 71000009),
(2, 72000001),
(2, 72000002),
(2, 72000003),
(2, 72000004),
(2, 72000005),
(2, 72000006),
(2, 72000007),
(2, 72000008),
(3, 73000001),
(3, 73000002),
(3, 73000003),
(3, 73000004),
(3, 73000005),
(3, 73000006),
(3, 73000007),
(3, 73000008),
(3, 73000009),
(3, 73000010),
(3, 73000011),
(3, 73000012),
(4, 74000001),
(4, 74000002),
(4, 74000003),
(4, 74000004),
(4, 74000005),
(4, 74000006),
(4, 74000007),
(5, 75000001),
(5, 75000002),
(5, 75000003),
(5, 75000004),
(5, 75000005);



------------------------------------------------------------

CREATE TABLE IF NOT EXISTS sms
(
    id         INTEGER PRIMARY KEY,
    sms_id     INTEGER      NOT NULL,
    phone_id   INTEGER      NOT NULL REFERENCES phone_numbers (id),
    phone_from VARCHAR(255) NOT NULL,
    text       TEXT         NOT NULL
);