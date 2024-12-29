CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    text_id TEXT NOT NULL, -- unique string used for lookups after oauth
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    picture_url TEXT NOT NULL,
    deleted INTEGER DEFAULT FALSE,

    UNIQUE(text_id)
);

CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    deleted INTEGER DEFAULT FALSE
);

CREATE TABLE event_organizers (
    id INTEGER PRIMARY KEY,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,

    FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE participants (
    id INTEGER PRIMARY KEY,
    event_id INTEGER NOT NULL,
    email TEXT NOT NULL,

    FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE TABLE email_templates (
    id INTEGER PRIMARY KEY,
    event_id INTEGER NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,

    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    expires_at TEXT NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);