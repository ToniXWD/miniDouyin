USE douyin;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    nickname VARCHAR(255),
    passwd VARCHAR(255) NOT NULL,
    follow_count INT,
    follower_count INT,
    work_count INT,
    favorite_count INT,
    token INT,
    avatar VARCHAR(255),
    background_image VARCHAR(255),
    signature VARCHAR(255),
    total_favorited INT
);

CREATE TABLE IF NOT EXISTS videos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author INT,
    play_url VARCHAR(255) NOT NULL UNIQUE,
    cover_url VARCHAR(255) NOT NULL,
    favorite_count INT,
    comment_count INT,
    FOREIGN KEY (author) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS favorited_videos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    video_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
);

CREATE TABLE IF NOT EXISTS comment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    video_id INT,
    content TEXT NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (video_id) REFERENCES videos(id)
);

CREATE TABLE IF NOT EXISTS message (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_id INT,
    to_id INT,
    content VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (from_id) REFERENCES users(id),
    FOREIGN KEY (to_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS followers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    follower_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (follower_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS follows (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    follow_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (follow_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS friends (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    friend_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (friend_id) REFERENCES users(id)
);