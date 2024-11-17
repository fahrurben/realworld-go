CREATE TABLE `users` (
    id BIGINT NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(100) NOT NULL,
    bio VARCHAR(1000),
    image VARCHAR(100),
    PRIMARY KEY (id)
)

CREATE TABLE `following`
(
    id             BIGINT NOT NULL AUTO_INCREMENT,
    user_id        BIGINT NOT NULL,
    follow_user_id BIGINT NOT NULL,
    PRIMARY KEY (id)
)

CREATE INDEX idx_following ON following (user_id, follow_user_id)