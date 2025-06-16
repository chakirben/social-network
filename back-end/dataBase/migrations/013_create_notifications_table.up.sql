CREATE TABLE IF NOT EXISTS Notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    senderId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiverId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('follow_request', 'group_invite', 'group_join_request', 'new_event')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'refused')),
    notificationDate DATETIME NOT NULL,
    userTargetId INTEGER REFERENCES users(id),
    groupTargetId INTEGER REFERENCES Groups(id),
    eventTargetId INTEGER REFERENCES Events(id)
);
