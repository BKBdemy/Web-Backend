DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS user_purchases CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS video CASCADE;
DROP TABLE IF EXISTS video_comments CASCADE;
DROP TABLE IF EXISTS user_watched_videos CASCADE;
DROP TABLE IF EXISTS user_tokens CASCADE;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    points INTEGER DEFAULT 0
);

CREATE TABLE user_tokens (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    token      VARCHAR NOT NULL,
    expiry    TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    price INTEGER NOT NULL,
    image VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE video (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL,
    parent_product_id INTEGER NOT NULL,
    thumbnail VARCHAR NOT NULL,
    filename VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_purchases (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE video_comments (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    video_id INTEGER NOT NULL,
    comment VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_watched_videos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    video_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

/* --Constraints-- */

/* Video deleted -> delete watched videos */
ALTER TABLE user_watched_videos
ADD CONSTRAINT fk_video_watched
FOREIGN KEY (video_id)
REFERENCES video (id)
ON DELETE CASCADE;

/* User deleted -> delete watched videos */
ALTER TABLE user_watched_videos
ADD CONSTRAINT fk_user_watched
FOREIGN KEY (user_id)
REFERENCES users (id)
ON DELETE CASCADE;

/* Product deleted -> delete videos */
ALTER TABLE video
ADD CONSTRAINT fk_product_videos
FOREIGN KEY (parent_product_id)
REFERENCES products (id)
ON DELETE CASCADE;

/* User deleted -> delete purchases */
ALTER TABLE user_purchases
ADD CONSTRAINT fk_user_purchases
FOREIGN KEY (product_id)
REFERENCES products (id)
ON DELETE CASCADE;

/* Product deleted -> delete comments */
ALTER TABLE video_comments
ADD CONSTRAINT fk_video_comments
FOREIGN KEY (video_id)
REFERENCES video (id)
ON DELETE CASCADE;

/* User deleted -> delete comments */
ALTER TABLE video_comments
ADD CONSTRAINT fk_comment_author
FOREIGN KEY (user_id)
REFERENCES users (id)
ON DELETE CASCADE;

/* Prevent negative balance */
ALTER TABLE users
ADD CONSTRAINT check_balance
CHECK (balance >= 0);

/* Sample data */

INSERT INTO users (username, password, balance)
VALUES ('admin', '$argon2id$v=19$m=256000,t=6,p=1$dGVzdHRlc3Q$MMMzLViNOBi+zmhnFWj4y1y6TqYfRvmUAI6BiH30mIk', 1000);
/* password is admin */

INSERT INTO users (username, password)
VALUES ('user', '$argon2id$v=19$m=256000,t=6,p=1$dGVzdHRlc3Q$MMMzLViNOBi+zmhnFWj4y1y6TqYfRvmUAI6BiH30mIk');
/* password is admin */

INSERT INTO products (name, description, price, image)
VALUES ('PHP Fundament', 'Dieser Kurs führt Sie in die Grundlagen von PHP ein und zeigt Ihnen, wie Sie eine Website strukturieren und gestalten können. Sie lernen auch, wie man Daten verarbeitet, Arrays verwendet und Probleme selbstständig löst. Der Kurs endet mit einem Projekt zur CAESAR-Verschlüsselung und einem Ausblick auf die nächsten Schritte.', 1000, '/static/php.jpeg');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Einführung', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video1.jpg', 'video1.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Variablen', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video2.jpg', 'video2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Arrays', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video3.jpg', 'video3.mp4');

INSERT INTO products (name, description, price, image)
VALUES ('Python Grundlagen', 'In unserem Kurs “Python lernen in unter 4 Stunden” führen wir dich in die Grundlagen von Python ein. Du lernst, wie man Python installiert und einrichtet, die Grundlagen von Python, die Verwendung von Schleifen, Listen und Funktionen und vieles mehr. Während des Kurses baust du auch zwei Projekte - einen Geburtstagskarten-Generator und ein Number Guessing Spiel - die dir helfen werden, das Gelernte anzuwenden und zu vertiefen.', 2000, '/static/python.jpeg');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Product 2 - Video 1', 'Product 2 - Video 1 description', 200, 2, 'video4.jpg', 'video4.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Product 2 - Video 2', 'Product 2 - Video 2 description', 200, 2, 'video5.jpg', 'video5.mp4');

INSERT INTO user_purchases (user_id, product_id)
VALUES (1, 1);
INSERT INTO user_purchases (user_id, product_id)
VALUES (1, 2);

INSERT INTO video_comments (user_id, video_id, comment)
VALUES (1, 1, 'This is a comment');

INSERT INTO video_comments (user_id, video_id, comment)
VALUES (1, 2, 'This is a comment');

INSERT INTO video_comments (user_id, video_id, comment)
VALUES (1, 3, 'This is a comment');

/* --Indexes-- */
CREATE INDEX idx_user_purchases_user_id ON user_purchases (user_id);
CREATE INDEX idx_user_purchases_product_id ON user_purchases (product_id);

CREATE INDEX idx_video_comments_user_id ON video_comments (user_id);
CREATE INDEX idx_video_comments_video_id ON video_comments (video_id);

CREATE INDEX idx_video_parent_product_id ON video (parent_product_id);

