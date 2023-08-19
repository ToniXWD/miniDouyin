\c douyin;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    nickname VARCHAR(255),
    passwd VARCHAR(255) NOT NULL,
    follow_count INT DEFAULT 0,
    follower_count INT DEFAULT 0,
    work_count INT DEFAULT 0,
    favorite_count INT DEFAULT 0,
    token VARCHAR(255),
    avatar VARCHAR(255),
    background_image VARCHAR(255),
    signature VARCHAR(255),
    total_favorited INT DEFAULT 0,
    deleted date DEFAULT null,
);
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author INT,
    play_url VARCHAR(255) NOT NULL UNIQUE,
    cover_url VARCHAR(255) NOT NULL,
    favorite_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    deleted date DEFAULT null,
    FOREIGN KEY (author) REFERENCES users(id)
);
CREATE TABLE favorited_videos (
    id SERIAL PRIMARY KEY,
    user_id INT,
    video_id INT,
    deleted date DEFAULT null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
);
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id INT,
    video_id INT,
    content TEXT NOT NULL,
    date DATE NOT NULL,
    deleted date DEFAULT null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
);
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    from_id INT,
    to_id INT,
    content VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    deleted date DEFAULT null,
    FOREIGN KEY (from_id) REFERENCES users(id),
    FOREIGN KEY (to_id) REFERENCES users(id)
);
CREATE TABLE followers (
    id SERIAL PRIMARY KEY,
    user_id INT,
    follower_id INT,
    deleted date DEFAULT null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (follower_id) REFERENCES users(id)
);
CREATE TABLE follows (
    id SERIAL PRIMARY KEY,
    user_id INT,
    follow_id INT,
    deleted date DEFAULT null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (follow_id) REFERENCES users(id)
);
CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INT,
    friend_id INT,
    deleted date DEFAULT null,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (friend_id) REFERENCES users(id)
);