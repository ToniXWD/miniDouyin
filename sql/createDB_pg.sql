\c
douyin;
CREATE TABLE users
(
    id               SERIAL PRIMARY KEY,           -- 自增主键
    username         VARCHAR(255) NOT NULL UNIQUE, -- 用户名
    nickname         VARCHAR(255),                 -- 昵称
    passwd           VARCHAR(255) NOT NULL,        -- 密码
    follow_count     bigint  DEFAULT 0,               --关注数
    follower_count   bigint  DEFAULT 0,               --粉丝数
    work_count       bigint  DEFAULT 0,               --作品数
    favorite_count   bigint  DEFAULT 0,               --点赞数
    token            VARCHAR(255),                 --鉴权
    avatar           VARCHAR(255),                 --头像url
    background_image VARCHAR(255),                 --背景图url
    signature        VARCHAR(255),                 --个人介绍
    total_favorited  bigint  DEFAULT 0,               --获赞数
    deleted          date DEFAULT null             --软删除
);
CREATE TABLE videos
(
    id             serial PRIMARY KEY,                                    -- 自增主键
    title          character varying(255) NOT NULL,                       -- 标题
    author         biginteger,                                               -- 作者id，外键
    play_url       character varying(255) NOT NULL UNIQUE,                --视频url
    cover_url      character varying(255) NOT NULL,                       --封面url
    favorite_count biginteger DEFAULT 0,                                     --获赞数
    comment_count  biginteger DEFAULT 0,                                     --评论数
    created_at     timestamp without time zone DEFAULT CURRENT_TIMESTAMP, --投稿时间
    deleted        date,                                                  --软删除
    FOREIGN KEY (author) REFERENCES users (id)
);
CREATE TABLE favorited_videos
(
    id       SERIAL PRIMARY KEY,   -- 自增主键
    user_id  bigint,
    video_id bigint,               -- 视频id，外键
    deleted  date DEFAULT null, -- 软删除
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (video_id) REFERENCES videos (id)
);
CREATE TABLE comments
(
    id       SERIAL PRIMARY KEY, -- 自增主键
    user_id  bigint,                --
    video_id bigint,                --
    content  TEXT NOT NULL,      --
    date     DATE NOT NULL,      --
    deleted  date DEFAULT null,  --
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (video_id) REFERENCES videos (id)
);
CREATE TABLE messages
(
    id      SERIAL PRIMARY KEY,    -- 自增主键
    from_id bigint,                   -- 发送者id，外键
    to_id   bigint,                   -- 接受者id，外键
    content VARCHAR(255) NOT NULL, -- 内容
    date    DATE         NOT NULL, -- 日期
    deleted date DEFAULT null,     -- 软删除
    FOREIGN KEY (from_id) REFERENCES users (id),
    FOREIGN KEY (to_id) REFERENCES users (id)
);
CREATE TABLE followers
(
    id          SERIAL PRIMARY KEY, -- 自增主键
    user_id     bigint,                -- 用户id，外键
    follower_id bigint,                -- 粉丝id，外键
    deleted     date DEFAULT null,  -- 软删除
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (follower_id) REFERENCES users (id)
);
CREATE TABLE follows
(
    id        SERIAL PRIMARY KEY,                 -- 自增主键
    user_id   bigint,                                -- 用户id，外键
    follow_id bigint,                                -- 被关注者id，外键
    deleted   date DEFAULT null,                  -- 软删除
    FOREIGN KEY (user_id) REFERENCES users (id),  --
    FOREIGN KEY (follow_id) REFERENCES users (id) --
);
CREATE TABLE friends
(
    id        SERIAL PRIMARY KEY, -- 自增主键
    user_id   bigint,                -- 用户id，外键
    friend_id bigint,                -- 好友id，外键
    deleted   date DEFAULT null,  -- 软删除
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (friend_id) REFERENCES users (id)
);