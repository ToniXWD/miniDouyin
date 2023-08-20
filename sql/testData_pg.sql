\c
tmp;
INSERT INTO users (username, passwd, token)
VALUES ('toni111', '123456', 'toni123456');

INSERT INTO videos (title, author, play_url, cover_url)
VALUES ('bear', 1, 'videos/bear.mp4', 'bgs/pandas.png'),
       ('city', 1, 'videos/city.mp4', 'bgs/pandas.png');