CREATE TABLE IF NOT EXISTS Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    content TEXT NOT NULL,
    image TEXT,
    privacy TEXT CHECK (privacy IN ('public', 'almostPrivate', 'private' , 'inGroup') OR privacy IS NULL),
    groupId INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    creatorId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);