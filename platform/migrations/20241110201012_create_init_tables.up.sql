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

CREATE TABLE article (
    id BIGINT NOT NULL AUTO_INCREMENT,
    author_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(1000) NOT NULL,
    body TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT now(),
    updated_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (author_id) REFERENCES users(id)
)

CREATE TABLE tag (
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (name)
)

CREATE TABLE article_tags
(
    id BIGINT NOT NULL AUTO_INCREMENT,
    article_id     BIGINT NOT NULL,
    tag_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (article_id) REFERENCES article(id)
)


CREATE TABLE `favorites_article`
(
    id         BIGINT NOT NULL AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    user_id    BIGINT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (article_id) REFERENCES article (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
)