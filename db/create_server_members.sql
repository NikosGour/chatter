CREATE TABLE IF NOT EXISTS server_members
(
    server_id TEXT NOT NULL REFERENCES servers (id) ON DELETE CASCADE,
    user_id   TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (server_id, user_id)
);