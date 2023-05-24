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
    difficulty INTEGER NOT NULL DEFAULT 1,
    preview_url VARCHAR NOT NULL DEFAULT '',
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

INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('PHP Fundament', 'Dieser Kurs führt Sie in die Grundlagen von PHP ein und zeigt Ihnen, wie Sie eine Website strukturieren und gestalten können. Sie lernen auch, wie man Daten verarbeitet, Arrays verwendet und Probleme selbstständig löst. Der Kurs endet mit einem Projekt zur CAESAR-Verschlüsselung und einem Ausblick auf die nächsten Schritte.', 1000, '/static/php.jpeg', '/static/previews/1.mp4', 2);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Einführung', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video1.jpg', 'video1.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Variablen', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video2.jpg', 'video2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Arrays', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video3.jpg', 'video3.mp4');



INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('Python Grundlagen', 'In unserem Kurs “Python lernen in unter 4 Stunden” führen wir dich in die Grundlagen von Python ein. Du lernst, wie man Python installiert und einrichtet, die Grundlagen von Python, die Verwendung von Schleifen, Listen und Funktionen und vieles mehr. Während des Kurses baust du auch zwei Projekte - einen Geburtstagskarten-Generator und ein Number Guessing Spiel - die dir helfen werden, das Gelernte anzuwenden und zu vertiefen.', 2000, '/static/python.jpeg', '/static/previews/2.mp4', 1);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Installation & Setup', 'Dieses Kapitel behandelt die Grundlagen der Python-Installation auf Ihrem Computer. Es führt Sie durch den Prozess des Herunterladens und Installierens der neuesten Python-Version, das Einrichten Ihrer Programmierumgebung und das Testen, um sicherzustellen, dass alles korrekt eingerichtet ist.', 200, 2, 'video4.jpg', '/static/python/Abschnitt 2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Das erste Programm', 'Dieses Kapitel führt Sie durch das Schreiben, Ausführen und Verstehen Ihres allerersten Python-Skripts. Es stellt grundlegende Konzepte wie Skriptstruktur und -ausführung vor.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 3.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Zahlen & Operatoren', 'In diesem Kapitel lernen Sie, wie Sie Zahlen in Python verwenden und manipulieren. Es werden sowohl Ganzzahlen als auch Gleitkommazahlen behandelt, und es werden die verschiedenen mathematischen Operatoren vorgestellt, die in Python verfügbar sind.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 4.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Zeichenketten', 'Dieser Abschnitt behandelt Zeichenketten (Strings) in Python. Es werden Themen wie das Erstellen, Manipulieren und Kombinieren von Zeichenketten behandelt.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 5.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Variablen', 'Hier lernen Sie, wie Sie Variablen in Python definieren und verwenden. Das Kapitel behandelt die Regeln zur Benennung von Variablen, den Prozess der Variablendeklaration und -zuweisung und die Verwendung von Variablen in Ihrem Code.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 6.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Datentypen', 'Dieses Kapitel behandelt die verschiedenen Datentypen, die in Python verfügbar sind, einschließlich Zahlen, Zeichenketten, Listen, Tupeln, Sets und Wörterbüchern.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 7.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Type-Casting', 'In diesem Kapitel lernen Sie, wie Sie Datentypen in Python umwandeln können. Es behandelt die Konzepte der impliziten und expliziten Typumwandlung.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 8.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Erstes Beispielprojekt', 'Hier setzen Sie das bisher Gelernte in die Praxis um, indem Sie ein Beispielprojekt erstellen. Dieses Projekt ermöglicht es Ihnen, die Konzepte, die Sie bisher gelernt haben, zu verstehen und anzuwenden.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 9.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Vergleiche und Booleans', 'Dieses Kapitel behandelt Vergleichsoperatoren und boolesche Werte in Python. Sie lernen, wie Sie Bedingungen mit Vergleichsoperatoren erstellen und wie Sie boolesche Werte verwenden können, um den Fluss Ihres Programms zu steuern.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 10.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Die for-Schleife', 'Das letzte Kapitel führt Sie in die Verwendung von ''for''-Schleifen in Python ein. Sie lernen, wie Sie eine Schleife erstellen, durchlaufen und steuern und wie Sie mit Schleifen verschiedene Arten von Problemen lösen können.', 200, 2, 'video5.jpg', '/static/python/Abschnitt 11.mp4');

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

