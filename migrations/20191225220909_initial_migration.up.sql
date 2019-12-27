CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    role VARCHAR NOT NULL DEFAULT "user",
    archived BOOLEAN NOT NULL DEFAULT 0,
    archived_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE integrations (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    username VARCHAR NOT NULL UNIQUE,
    api_key VARCHAR NOT NULL,
    auth_token VARCHAR NOT NULL,
    
    archived BOOLEAN NOT NULL DEFAULT 0,
    archived_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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
    avatar_blob_id VARCHAR,

    archived BOOLEAN NOT NULL DEFAULT 0,
    archived_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (integration_id, vrchat_id)
);

CREATE TABLE attendance (
    timestamp INT NOT NULL,
    integration_id INT NULL NULL REFERENCES integrations(id),
    friend_id INT REFERENCES friends(id),
    teacher_id INT REFERENCES friends(id),
    location VARCHAR NOT NULL,
    
    archived BOOLEAN NOT NULL DEFAULT 0,
    archived_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (timestamp, friend_id, integration_id)
);

CREATE TABLE blobs (
    id INTEGER PRIMARY KEY,
    file_name VARCHAR NOT NULL,
    mime_type VARCHAR NOT NULL,
    file_size_bytes INT NOT NULL,
    EXTENSION VARCHAR NOT NULL,
    file BLOB NOT NULL,
    views INTEGER DEFAULT 0,

    archived BOOLEAN NOT NULL DEFAULT 0,
    archived_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
