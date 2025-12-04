CREATE TABLE IF NOT EXISTS messages
(
    id           TEXT PRIMARY KEY,
    sender_id    TEXT      NOT NULL REFERENCES channels (id),
    recipient_id TEXT      NOT NULL REFERENCES channels (id),
    date_sent    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
