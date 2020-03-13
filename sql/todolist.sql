CREATE TABLE IF NOT EXISTS todo_list
(
    id      serial PRIMARY KEY,
    user_id INTEGER     NOT NULL,
    todo    VARCHAR(50) NOT NULL,
    status  boolean
);