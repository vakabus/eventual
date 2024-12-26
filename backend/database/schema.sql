CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    deleted INTEGER DEFAULT FALSE
);

CREATE TABLE events (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    readme TEXT NOT NULL DEFAULT "",
    deleted INTEGER DEFAULT FALSE
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