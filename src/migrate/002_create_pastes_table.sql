CREATE TABLE pastes (
    id SERIAL PRI MARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    extension VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    slug VARCHAR(255) NOT NULL UNIQUE
);

