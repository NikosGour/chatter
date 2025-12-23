CREATE TABLE IF NOT EXISTS users
(
    id           TEXT UNIQUE PRIMARY KEY,
    username     TEXT      NOT NULL,
    password     TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    is_test      boolean DEFAULT FALSE,
   
    CONSTRAINT username_it_test_unique UNIQUE (username, is_test)
);
