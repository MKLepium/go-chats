DROP DATABASE IF EXISTS chat_app;
CREATE DATABASE chat_app;
USE chat_app;

-- Drop tables if they exist to prevent errors during schema creation
DROP TABLE IF EXISTS room_participants;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chat_rooms;
DROP TABLE IF EXISTS users;


CREATE TABLE users (
    user_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL
);

-- Create the chat_rooms table with explicit type for room_id and created_by_user_id
CREATE TABLE chat_rooms (
    room_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    room_name VARCHAR(255) NOT NULL,
    created_by_user_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX (created_by_user_id)
);

-- Create the messages table with explicit types for IDs
CREATE TABLE messages (
    message_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    from_user_id BIGINT UNSIGNED NOT NULL,
    to_user_id BIGINT UNSIGNED NOT NULL,
    message_text TEXT,
    sent_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `read` BOOLEAN DEFAULT FALSE,
    INDEX (from_user_id),
    INDEX (to_user_id)
);

CREATE TABLE room_participants (
    room_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id),
    INDEX (room_id),
    INDEX (user_id)
);

-- Add foreign key constraints after all table creation statements
ALTER TABLE messages
ADD CONSTRAINT fk_from_user_id FOREIGN KEY (from_user_id) REFERENCES users(user_id),
ADD CONSTRAINT fk_to_user_id FOREIGN KEY (to_user_id) REFERENCES users(user_id);

ALTER TABLE room_participants
ADD CONSTRAINT fk_room_id FOREIGN KEY (room_id) REFERENCES chat_rooms(room_id),
ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id);

ALTER TABLE chat_rooms
ADD CONSTRAINT fk_created_by_user_id FOREIGN KEY (created_by_user_id) REFERENCES users(user_id);

INSERT INTO users (username, email, password_hash)
VALUES ('admin', 'admin@example.com', SHA2('password', 256));

INSERT INTO users (username, email, password_hash)
VALUES ('user1', 'user1@example.com', SHA2('password', 256));	