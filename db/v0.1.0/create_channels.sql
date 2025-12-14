CREATE TABLE IF NOT EXISTS channels
(
    id           TEXT PRIMARY KEY,
    channel_type text NOT NULL CHECK ( channel_type IN ('user', 'group'))
);
