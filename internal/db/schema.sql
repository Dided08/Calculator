-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

-- Таблица выражений (запросов на вычисление)
CREATE TABLE IF NOT EXISTS expressions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    expression TEXT NOT NULL,
    result TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'pending', -- pending, in_progress, done, error
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Таблица задач для агентов
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expression_id INTEGER NOT NULL,
    operation TEXT NOT NULL,
    arg1 TEXT NOT NULL,
    arg2 TEXT NOT NULL,
    result TEXT,
    status TEXT DEFAULT 'pending', -- pending, in_progress, done
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expression_id) REFERENCES expressions(id)
);