CREATE TABLE IF NOT EXISTS messages
(
    id        bigserial UNIQUE PRIMARY KEY,
    "text"    TEXT      NOT NULL,
    tab_id    TEXT      NOT NULL REFERENCES tabs (id) ON DELETE CASCADE,
    date_sent TIMESTAMP NOT NULL
);