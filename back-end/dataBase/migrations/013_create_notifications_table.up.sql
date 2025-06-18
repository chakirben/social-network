CREATE TABLE IF NOT EXISTS Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    senderId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiverId INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('follow_request', 'group_invite', 'group_join_request', 'new_event')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'refused')),
    notificationDate DATETIME DEFAULT CURRENT_TIMESTAMP,
    groupId INTEGER REFERENCES Groups(id) ON DELETE CASCADE,
    eventId INTEGER REFERENCES Events(id)  ON DELETE CASCADE
);
