CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL
);
CREATE TABLE integrations (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    username VARCHAR NOT NULL UNIQUE,
    api_key VARCHAR NOT NULL,
    auth_token VARCHAR NOT NULL 
);
CREATE TABLE friends (
    id INTEGER PRIMARY KEY,
    integration_id INT NOT NULL REFERENCES integrations(id),
    is_teacher BOOLEAN NOT NULL DEFAULT 0,
    vrchat_id VARCHAR UNIQUE NOT NULL,
    vrchat_username VARCHAR NOT NULL,
    vrchat_display_name VARCHAR NOT NULL,
    vrchat_avatar_image_url VARCHAR NOT NULL,
    vrchat_avatar_thumbnail_image_url VARCHAR NOT NULL,
    vrchat_location VARCHAR NOT NULL,
    UNIQUE (integration_id, vrchat_id)
);

CREATE TABLE attendance (
    timestamp INT NOT NULL,
    integration_id INT NULL NULL REFERENCES integrations(id),
    friend_id INT REFERENCES friends(id),
    teacher_id INT REFERENCES friends(id),
    location VARCHAR NOT NULL,
    PRIMARY KEY (timestamp, friend_id, integration_id)
);