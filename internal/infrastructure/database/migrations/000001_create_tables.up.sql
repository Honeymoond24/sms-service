BEGIN;

CREATE TABLE IF NOT EXISTS countries
(
    id           SERIAL PRIMARY KEY,
    country_code INT          NOT NULL,
    name         VARCHAR(255) NOT NULL
);
INSERT INTO countries (country_code, name)
VALUES (1, 'russia'),
       (2, 'ukraine'),
       (3, 'kazakhstan'),
       (4, 'usa'),
       (5, 'germany')
ON CONFLICT DO NOTHING;



CREATE TABLE IF NOT EXISTS services
(
    id           SERIAL PRIMARY KEY,
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



CREATE TABLE IF NOT EXISTS activations_statuses
(
    id   INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO activations_statuses (id, name)
VALUES (-1, 'ACCESS_CANCEL'),
       (1, 'ACCESS_READY'),
       (3, 'ACCESS_RETRY_GET'),
       (6, 'ACCESS_ACTIVATION'),
       (8, 'ACCESS_CANCEL')
ON CONFLICT DO NOTHING;



CREATE TABLE IF NOT EXISTS phone_numbers
(
    id         SERIAL PRIMARY KEY,
    country_id INT        NOT NULL,
    number     INT UNIQUE NOT NULL,
    CONSTRAINT phone_numbers_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS activations
(
    id         SERIAL PRIMARY KEY,
    sum_price  REAL NOT NULL DEFAULT 0,
    status     INT  NOT NULL REFERENCES activations_statuses (id),
    phone_id   INT  NOT NULL REFERENCES phone_numbers (id),
    service_id INT  NOT NULL REFERENCES services (id),
    CONSTRAINT activations_phone_id_service_id_unique UNIQUE (phone_id, service_id)
);



CREATE TABLE IF NOT EXISTS sms
(
    id         SERIAL PRIMARY KEY,
    sms_id     INT          NOT NULL,
    phone_id   INT          NOT NULL REFERENCES phone_numbers (id),
    phone_from VARCHAR(255) NOT NULL,
    text       TEXT         NOT NULL
);

COMMIT;