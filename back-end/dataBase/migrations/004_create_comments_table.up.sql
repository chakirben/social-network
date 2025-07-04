CREATE TABLE IF NOT EXISTS Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    postId INTEGER NOT NULL REFERENCES Posts(id) ON DELETE CASCADE,
    userId INTEGER NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    image TEXT,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);