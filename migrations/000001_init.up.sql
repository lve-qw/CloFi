-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL,
    password VARCHAR(255) NOT NULL, -- хеш bcrypt
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица лайков (составной первичный ключ)
CREATE TABLE IF NOT EXISTS likes (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, product_id)
);

-- Индекс для быстрого поиска лайков пользователя
CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);

