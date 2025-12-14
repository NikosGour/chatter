CREATE TABLE IF NOT EXISTS group_members
(
    group_id TEXT NOT NULL REFERENCES groups (id),
    user_id  TEXT NOT NULL REFERENCES users (id),
    PRIMARY KEY (group_id, user_id)
);
