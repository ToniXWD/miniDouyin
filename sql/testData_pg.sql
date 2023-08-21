\c
tmp;
INSERT INTO users (username, passwd, token, work_count)
VALUES ('toni111', '123456', 'toni123456', 2);

INSERT INTO videos (title, author, play_url, cover_url)
VALUES ('bear', 1, 'videos/bear.mp4', 'bgs/pandas.png'),
       ('city', 1, 'videos/city.mp4', 'bgs/pandas.png');