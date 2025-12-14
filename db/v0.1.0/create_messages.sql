CREATE TABLE IF NOT EXISTS messages
(
    id           BIGSERIAL PRIMARY KEY,
    "text"       TEXT      NOT NULL,
    sender_id    TEXT      NOT NULL REFERENCES channels (id),
    recipient_id TEXT      NOT NULL REFERENCES channels (id),
    date_sent    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
