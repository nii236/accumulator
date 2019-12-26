CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL
);
CREATE TABLE integrations (
    id INTEGER PRIMARY KEY,
    api_key VARCHAR NOT NULL DEFAULT '',
    auth_token VARCHAR NOT NULL DEFAULT ''
);
CREATE TABLE friends (
    id VARCHAR PRIMARY KEY,
    integration_id INT NULL NULL REFERENCES integrations(id),
    is_teacher BOOLEAN NOT NULL DEFAULT 0,
    vrchat_username VARCHAR,
    vrchat_display_name VARCHAR,
    vrchat_avatar_image_url VARCHAR,
    vrchat_avatar_thumbnail_image_url VARCHAR
);

CREATE TABLE attendance (
    timestamp INT NOT NULL,
    integration_id INT NULL NULL REFERENCES integrations(id),
    friend_id INT REFERENCES friends(id),
    teacher_id INT REFERENCES friends(id),
    world_id VARCHAR NOT NULL,
    instance_id VARCHAR NOT NULL,
    PRIMARY KEY (timestamp, friend_id)
);