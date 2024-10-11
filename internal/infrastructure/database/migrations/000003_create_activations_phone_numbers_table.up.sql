CREATE TABLE IF NOT EXISTS activations
(
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    status INTEGER NOT NULL REFERENCES activations_statuses (id)
);

CREATE TABLE IF NOT EXISTS activations_statuses
(
    id   INTEGER PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

INSERT INTO activations_statuses (id, name)
VALUES (-1, 'ACCESS_CANCEL '),
       (1, 'ACCESS_READY'),
       (3, 'ACCESS_RETRY_GET'),
       (6, 'ACCESS_ACTIVATION'),
       (8, 'ACCESS_CANCEL');

CREATE TABLE IF NOT EXISTS phone_numbers
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    country_id    INTEGER     NOT NULL,
    operator      VARCHAR(255) DEFAULT NULL,
    number        VARCHAR(31) NOT NULL,
    activation_id INTEGER      DEFAULT NULL REFERENCES activations (id),
    FOREIGN KEY (country_id) REFERENCES countries (id) ON DELETE CASCADE
);

INSERT INTO activations (id, status)
VALUES (1, 1),
       (2, 1),
       (3, 1),
       (4, 1),
       (5, 1),
       (6, 1),
       (7, 1),
       (8, 1),
       (9, 1),
       (10, 1);

INSERT INTO phone_numbers (country_id, operator, number, activation_id)
VALUES (1, NULL, '79185556601', 1),
       (1, NULL, '79185556602', 2),
       (2, NULL, '79185556603', 3),
       (2, NULL, '79185556604', 4),
       (3, NULL, '79185556605', 5),
       (3, NULL, '79185556606', 6),
       (4, NULL, '79185556607', 7),
       (5, NULL, '79185556608', 8),
       (5, NULL, '79185556609', 9),
       (5, NULL, '79185556610', 10);
