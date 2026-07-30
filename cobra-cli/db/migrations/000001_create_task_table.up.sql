-- 000001_create_tasks_table.up.sql
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ← SQLite
    name TEXT NOT NULL,                    -- ← SQLite
    status TEXT DEFAULT 'pending',         -- ← SQLite
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- ← SQLite
);
