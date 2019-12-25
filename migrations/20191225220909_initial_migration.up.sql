CREATE TABLE users (
    id INT PRIMARY KEY,
    email VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    api_key VARCHAR,
    auth_token VARCHAR
);

CREATE TABLE friends (
    id VARCHAR PRIMARY KEY,
    is_teacher BOOLEAN NOT NULL DEFAULT 0,
    vrchat_username VARCHAR,
    vrchat_display_name VARCHAR,
    vrchat_avatar_image_url VARCHAR,
    vrchat_avatar_thumbnail_image_url VARCHAR
);

CREATE TABLE attendance (
    timestamp INT NOT NULL,
    friend_id INT REFERENCES friends(id),
    teacher_id INT REFERENCES friends(id),
    world_id VARCHAR NOT NULL,
    instance_id VARCHAR NOT NULL,
    PRIMARY KEY (timestamp, friend_id)
);