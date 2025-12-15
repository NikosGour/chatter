CREATE TABLE IF NOT EXISTS messages
(
    id        bigserial UNIQUE PRIMARY KEY,
    "text"    TEXT      NOT NULL,
    sender_id    TEXT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tab_id    TEXT      NOT NULL REFERENCES tabs (id) ON DELETE CASCADE,
    date_sent TIMESTAMP NOT NULL
);