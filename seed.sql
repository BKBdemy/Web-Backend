DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS user_purchases CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS video CASCADE;
DROP TABLE IF EXISTS product_comments CASCADE;
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

CREATE TABLE product_comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
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
ALTER TABLE product_comments
ADD CONSTRAINT fk_product_comments
FOREIGN KEY (course_id)
REFERENCES video (id)
ON DELETE CASCADE;

/* User deleted -> delete comments */
ALTER TABLE product_comments
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

/* PHP */

INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('PHP Fundament', 'Dieser Kurs führt Sie in die Grundlagen von PHP ein und zeigt Ihnen, wie Sie eine Website strukturieren und gestalten können. Sie lernen auch, wie man Daten verarbeitet, Arrays verwendet und Probleme selbstständig löst. Der Kurs endet mit einem Projekt zur CAESAR-Verschlüsselung und einem Ausblick auf die nächsten Schritte.', 1000, '/static/php.jpeg', '/static/previews/1.mp4', 2);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Einführung', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video1.jpg', 'video1.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Variablen', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video2.jpg', 'video2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('PHP Fundament - Arrays', 'In diesem Video lernen Sie die Grundlagen von PHP kennen.', 100, 1, 'video3.jpg', 'video3.mp4');

/* Python */

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

/* HTML & CSS */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('Erstelle deine eigene Webseite', 'Tauche ein in die faszinierende Welt der Webentwicklung und erschaffe deine eigene atemberaubende Webseite. Lerne die Grundlagen von HTML und CSS, um deine kreativen Ideen zum Leben zu erwecken. Egal, ob du Anfänger bist oder bereits erste Erfahrungen hast, dieser Kurs bietet dir das nötige Wissen, um eine beeindruckende Homepage zu erstellen.', 500, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 1.mp4', 1);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Einführung in die aufregende Welt von HTML', 'Tauchen Sie ein in die spannende Welt von HTML und entdecken Sie die Grundlagen, um eine beeindruckende Webseite zu erstellen. Lernen Sie, wie Sie eine HTML-Datei von Grund auf erstellen und die mächtigen HTML-Tags nutzen, um Ihre Vision zum Leben zu erwecken.', 200, 3, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Der ultimative HTML Editor', 'Lassen Sie sich von uns in die Geheimnisse des Visual Studio Code einweihen und machen Sie sich bereit, Ihre Entwicklungsreise auf die nächste Stufe zu heben. Entdecken Sie die zahlreichen Features und Tools, die Ihnen zur Verfügung stehen, um Ihre HTML-Kreationen zu perfektionieren.', 200, 3, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 3.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Meistern Sie die Kunst der HTML Architektur', 'Tauchen Sie ein in die faszinierende Welt der HTML-Architektur und lernen Sie, wie Sie die Grundstruktur eines HTML-Dokuments gestalten. Entdecken Sie die Bausteine, die Ihre Webseite zusammenhalten, und lernen Sie, wie Sie sie effektiv einsetzen, um Ihre Inhalte optimal zur Geltung zu bringen.', 200, 3, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 4.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('CSS Basics: Meistern Sie die Kunst der Gestaltung', 'Erweitern Sie Ihre Webentwicklungsfähigkeiten und beherrschen Sie die Grundlagen von CSS. Lernen Sie, wie Sie CSS in Ihre HTML-Dateien einbinden und verwenden Sie es, um das Aussehen Ihrer Website nach Ihren Vorstellungen anzupassen. Tauchen Sie ein in eine Welt voller Stil und Kreativität.', 200, 3, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 5.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('HTML & CSS Layouts: Kreativität ohne Grenzen', 'Heben Sie Ihr Webdesign auf ein neues Niveau und lernen Sie, wie Sie mit HTML und CSS faszinierende Layouts erstellen. Entdecken Sie Techniken, um Ihre Inhalte zu organisieren und eine optimale Benutzererfahrung zu bieten. Lassen Sie Ihrer Kreativität freien Lauf und gestalten Sie beeindruckende Webseiten.', 200, 3, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 6.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('CSS Styling & HTML Embeds: Machen Sie Ihre Website zum Blickfang', 'Entdecken Sie die faszinierende Welt des CSS-Stylings und erfahren Sie, wie Sie das Aussehen Ihrer Website aufregend und ansprechend gestalten können. Lernen Sie, wie Sie Texte, Bilder und Links formatieren und beeindruckende Videos und Karten in Ihre Webseite einbetten. Bringen Sie Ihre Website zum Strahlen!', 200, 3, '/static/htmlcss/thumbnail.jpg', '/static/htmlcss/Abschnitt 7.mp4');


/* Java */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('Java für Anfänger', 'Java ist eine der beliebtesten Programmiersprachen der Welt. In diesem Kurs lernen Sie die Grundlagen von Java und werden in der Lage sein, Ihre eigenen Programme zu schreiben. Egal, ob Sie ein Anfänger sind oder bereits erste Erfahrungen haben, dieser Kurs bietet Ihnen das nötige Wissen, um Ihre eigenen Java-Programme zu erstellen.', 500, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 1.mp4', 2);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Einführung in Java-Konzepte', 'Starten Sie Ihre Reise in die aufregende Welt der Java-Programmierung. Erlernen Sie die grundlegenden Konzepte und Methoden, um Ihre eigenen Java-Anwendungen zu erstellen und Ihre Ideen zum Leben zu erwecken.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Einrichten von Java & IntelliJ für optimale Entwicklung', 'Lernen Sie die Vorteile und Funktionen von IntelliJ kennen und nehmen Sie Ihre Java-Entwicklung auf eine neue Ebene. Entdecken Sie die Vielfalt an Werkzeugen, die Ihnen zur Verfügung stehen, um Ihre Java-Programme zu optimieren und zu perfektionieren.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 3.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Erstes Java-Programm: Hallo Welt', 'Erstellen Sie Ihr erstes Java-Programm! Beginnen Sie Ihre Coding-Reise mit dem klassischen "Hallo Welt"-Programm und erleben Sie den Stolz der Code-Erstellung.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 4.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Macht der Variablen in Java', 'Erfahren Sie, wie Sie Variablen in Java effizient einsetzen, um Ihre Programme vielseitiger und dynamischer zu gestalten.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 5.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Java-Verzweigungen: Der Weg zur Entscheidungsfindung', 'Tauchen Sie ein in die Logik und die Kraft der Verzweigungen in Java. Lernen Sie, wie Sie diese verwenden, um die Entscheidungsfindung in Ihren Programmen zu steuern.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 6.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Funktionen in Java: Werkzeuge zur Code-Wiederverwendung', 'Erfahren Sie, wie Sie Funktionen in Java erstellen und nutzen, um Ihre Code-Effizienz zu steigern und die Wiederverwendbarkeit zu verbessern.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 7.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Java in Aktion: Währungsrechner Projekt', 'Wenden Sie Ihre bisherigen Kenntnisse an und lernen Sie, wie Sie ein funktionaler Währungsrechner mit Java erstellen können. Dieses Projekt wird Ihnen dabei helfen, die erlernten Konzepte zu festigen und praktische Erfahrungen zu sammeln.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 8.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Spielerisches Lernen: Zahlenraten-Spiel in Java', 'Lassen Sie uns spielerisch lernen! Erstellen Sie ein unterhaltsames Zahlenraten-Spiel in Java und üben Sie dabei Ihre Programmierkenntnisse.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 9.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Erstellen Sie ansprechende Benutzeroberflächen: GUIs in Java', 'Lernen Sie, wie Sie mit Java ansprechende und intuitive Benutzeroberflächen erstellen können. Erfahren Sie, wie GUIs die Benutzererfahrung Ihrer Programme verbessern können.', 200, 4, '/static/java/thumbnail.jpg', '/static/java/Abschnitt 10.mp4');


/* C++ */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('Objekt-Orientiertes C++', 'C++ ist eine der beliebtesten Programmiersprachen der Welt. In diesem Kurs lernen Sie die Grundlagen von C++ und werden in der Lage sein, Ihre eigenen Programme zu schreiben. Egal, ob Sie ein Anfänger sind oder bereits erste Erfahrungen haben, dieser Kurs bietet Ihnen das nötige Wissen, um Ihre eigenen C++-Programme zu erstellen.', 500, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 1.mp4', 1);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Klassen und Objekte in C++', 'Lernen Sie, wie Sie Klassen und Objekte in C++ erstellen und verwenden. Entdecken Sie die Vorteile der Objektorientierung und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 1.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('UML-Diagramme in C++', 'Erfahren Sie, wie Sie UML-Diagramme verwenden, um Ihre C++-Programme zu planen und zu entwerfen. Lernen Sie, wie Sie Klassen, Attribute und Methoden in UML-Diagrammen darstellen.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('UML-Diagramme in C++ #2', 'Erfahren Sie, wie Sie UML-Diagramme verwenden, um Ihre C++-Programme zu planen und zu entwerfen. Lernen Sie, wie Sie Klassen, Attribute und Methoden in UML-Diagrammen darstellen.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 3.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Erste Schritte in C++', 'Lernen Sie die Grundlagen von C++ kennen und erstellen Sie Ihr erstes Programm. Entdecken Sie die Syntax von C++ und lernen Sie, wie Sie Variablen, Datentypen und Operatoren verwenden.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 4.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('for-Schleifen in C++', 'Erfahren Sie, wie Sie for-Schleifen in C++ verwenden, um Ihre Programme effizienter zu gestalten. Lernen Sie, wie Sie for-Schleifen verwenden, um Code zu wiederholen und Ihre Programme zu optimieren.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 5.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Container-Klasse std::array in C++', 'Lernen Sie, wie Sie Container-Klassen in C++ verwenden, um Ihre Programme effizienter zu gestalten. Entdecken Sie die Vorteile von Container-Klassen und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 6.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Container-Klasse std::vector in C++', 'Lernen Sie, wie Sie Container-Klassen in C++ verwenden, um Ihre Programme effizienter zu gestalten. Entdecken Sie die Vorteile von Container-Klassen und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 7.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Iteratoren und deren Verwendung in <algorithm> in C++', 'Lernen Sie, wie Sie Iteratoren in C++ verwenden, um Ihre Programme effizienter zu gestalten. Entdecken Sie die Vorteile von Iteratoren und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 8.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Beispielprojekt Tic-Tac-Toe: Klassendiagramm', 'Lernen Sie, wie Sie UML-Diagramme verwenden, um Ihre C++-Programme zu planen und zu entwerfen. Lernen Sie, wie Sie Klassen, Attribute und Methoden in UML-Diagrammen darstellen.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 9.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Beispielprojekt Tic-Tac-Toe: Enumerationen', 'Lernen Sie, wie Sie Enumerationen in C++ verwenden, um Ihre Programme effizienter zu gestalten. Entdecken Sie die Vorteile von Enumerationen und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 10.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Beispielprojekt Tic-Tac-Toe: Klassenimplementierung', 'Implementieren Sie die Klassen für das Tic-Tac-Toe-Beispielprojekt. Lernen Sie, wie Sie Klassen in C++ implementieren.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 11.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Beispielprojekt Tic-Tac-Toe: Klasseninstanziierung', 'Instanziieren Sie die Klassen für das Tic-Tac-Toe-Beispielprojekt. Lernen Sie, wie Sie Klassen in C++ instanziieren.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 12.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Beispielprojekt Tic-Tac-Toe: Aggregation & Komposition', 'Lernen Sie, wie Sie Aggregation und Komposition in C++ verwenden, um Ihre Programme effizienter zu gestalten. Entdecken Sie die Vorteile von Aggregation und Komposition und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 13.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Beispielprojekt Tic-Tac-Toe: Konstruktor & Destruktor', 'Lernen Sie, wie Sie Konstruktoren und Destruktoren in C++ verwenden, um Ihre Programme effizienter zu gestalten. Entdecken Sie die Vorteile von Konstruktoren und Destruktoren und wie Sie diese nutzen können, um Ihre Programme zu verbessern.', 200, 5, '/static/cpp/thumbnail.jpg', '/static/cpp/Abschnitt 14.mp4');

/* One-Page Website */
/* SQL */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('One-Page Website mit HTML & CSS', 'Lernen Sie, wie Sie eine One-Page Website mit HTML & CSS erstellen. In diesem Kurs werden Sie die Grundlagen von HTML & CSS kennen lernen und Ihre eigene Website erstellen. Egal, ob Sie ein Anfänger sind oder bereits erste Erfahrungen haben, dieser Kurs bietet Ihnen das nötige Wissen, um Ihre eigene Website zu erstellen.', 500, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Einführung & Struktur.mp4', 3);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Einführung & Struktur', 'In dieser Lektion lernen Sie die Struktur einer One-Page Website kennen und welche Schritte Sie für deren Erstellung benötigen.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Einführung & Struktur.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Grundgerüst in HTML', 'In dieser Lektion lernen Sie, wie Sie das Grundgerüst für Ihre Website in HTML erstellen.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Grundgerüst in HTML.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Vorbereitungen treffen', 'In dieser Lektion bereiten wir alles vor, um mit dem Styling der Website beginnen zu können.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Vorbereitungen treffen.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Das Stylen beginnt', 'In dieser Lektion beginnen wir mit dem Styling der Website. Wir lernen, wie man CSS effektiv einsetzt, um die Website ansprechend zu gestalten.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Das Stylen beginnt.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Header & Navigation', 'In dieser Lektion gestalten wir den Header und die Navigation unserer Website. Diese Elemente sind wichtig für die Benutzerführung und das Gesamtbild der Website.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Header & Navigation.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Home & About', 'In dieser Lektion gestalten wir die Home- und About-Bereiche unserer Website. Wir lernen, wie man Inhalte ansprechend darstellt und den Benutzer auf der Website hält.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Home & About.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Work & Contact', 'In dieser Lektion gestalten wir den Work- und Contact-Bereich unserer Website. Wir lernen, wie man sein Portfolio präsentiert und den Benutzer dazu ermutigt, Kontakt aufzunehmen.', 200, 6, '/static/advancedhtml/thumbnail.jpg', '/static/advancedhtml/Work & Contact.mp4');


/* API Security */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('API Sicherheit', 'Lernen Sie, wie Sie APIs sicher gestalten und mögliche Sicherheitslücken vermeiden. Dieser Kurs behandelt verschiedene Aspekte der API-Sicherheit, einschließlich Bearer Tokens, OAuth 2.0, XSS-Injection, SQL-Injection und mehr.', 500, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 1.mp4', 3);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Was ist ein Bearer Token?', 'In dieser Lektion lernen Sie, was ein Bearer Token ist und wie es in der API-Sicherheit verwendet wird.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 1.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('OAuth 2.0 im Detail', 'Diese Lektion bietet eine detaillierte Erläuterung von OAuth 2.0 anhand eines Beispiels.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 2.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('APIs hacken, wie geht das?', 'In dieser Lektion lernen Sie, wie man APIs hackt und wie man die Sicherheit von APIs beurteilt.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 3.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('API Sicherheit - Excessive Data Exposure verhindern', 'Lernen Sie, wie Sie übermäßige Datenexposition in APIs verhindern.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 4.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Cross Site Scripting (XSS) Injection', 'Diese Lektion behandelt XSS-Injections und wie man sie in APIs verhindert.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 5.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Mass Assignment Autobinding Vulnerability', 'Erfahren Sie mehr über Mass Assignment und Autobinding-Schwachstellen in APIs und wie man sie vermeidet.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 6.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Broken Function Level Authorization', 'In dieser Lektion lernen Sie, was Broken Function Level Authorization ist und wie man es in APIs vermeidet.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 7.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('SQL Injection', 'Lernen Sie, was eine SQL-Injection ist und wie man sie in APIs vermeidet.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 8.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Rate Limiting Brute Force Angriffe', 'In dieser Lektion lernen Sie, was Rate Limiting und Brute-Force-Angriffe sind und wie man sie verhindert.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 9.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Insufficient Monitoring Logging', 'In dieser letzten Lektion lernen Sie, wie man unzureichendes Monitoring und Logging in APIs vermeidet.', 200, 6, '/static/apisecurity/thumbnail.jpg', '/static/apisecurity/Abschnitt 10.mp4');


/* JS */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('Javascript Tutorial für Anfänger', 'Ein Anfängerfreundliches Javascript Tutorial, das von den Grundlagen bis zu fortgeschrittenen Konzepten reicht.', 300, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Abschnitt_1_Einfuehrung_und_erstes_Programm.mp4', 1);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Einführung und erstes Programm', 'Der erste Teil des Javascript-Tutorials führt in die Grundlagen von Javascript ein und hilft Ihnen, Ihr erstes Programm zu schreiben.', 100, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Abschnitt_1_Einfuehrung_und_erstes_Programm.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Variablen', 'Der zweite Teil des Javascript-Tutorials führt in das Konzept der Variablen ein.', 100, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Abschnitt_2_Variablen.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Operatoren', 'Der dritte Teil des Javascript-Tutorials führt in das Konzept der Operatoren ein.', 100, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Abschnitt_3_Operatoren.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Bedingte Anweisungen und das DOM', 'Der vierte Teil des Javascript-Tutorials führt in das Konzept der bedingten Anweisungen und des Document Object Model (DOM) ein.', 100, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Abschnitt_4_Bedingte_Aweisungen_und_das_DOM.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Arrays und Schleifen', 'Der fünfte Teil des Javascript-Tutorials führt in das Konzept von Arrays und Schleifen ein.', 100, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Abschnitt_5_Arrays_und_Schleifen.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Übungsvideo zu Operatoren', 'Dieses Übungsvideo bietet zusätzliche Übungen zu den in Abschnitt 3 behandelten Konzepten.', 50, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Uebungsvideo_zu_Abschnitt_3_Operatoren.mp4');

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Übungsvideo zu Bedingten Anweisungen und Arrays', 'Dieses Übungsvideo bietet zusätzliche Übungen zu den in Abschnitt 4 und 5 behandelten Konzepten.', 50, 7, '/static/js_tutorial/thumbnail.jpg', '/static/js_tutorial/Uebungsvideo_zu_Abschnitten_4_und_5.mp4');

/* C# */
INSERT INTO products (name, description, price, image, preview_url, difficulty)
VALUES ('C# im Detail', 'Ein detailliertes C#-Tutorial, das die Grundlagen bis hin zu fortgeschrittenen Konzepten abdeckt.', 400, '/static/csharp_tutorial/thumbnail.jpg', '/static/csharp_tutorial/Abschnitt_1_Installation_und_erstes_Programm.mp4', 2);

INSERT INTO video (name, description, points, parent_product_id, thumbnail, filename)
VALUES ('Installation und erstes Programm', 'Der erste Abschnitt des C#-Tutorials, der in die Installation und das Schreiben des ersten Programms einführt.', 100, 8, '/static/csharp_tutorial/thumbnail.jpg', '/static/csharp_tutorial/Abschnitt_1_Installation_und_erstes_Programm.mp4'),
       ('Variablen und Datentypen', 'Der zweite Abschnitt des C#-Tutorials, der das Konzept von Variablen und Datentypen einführt.', 100, 8, '/static/csharp_tutorial/thumbnail.jpg', '/static/csharp_tutorial/Abschnitt_2_Variablen_und_Datentypen.mp4'),
       ('Mathematische Operatoren', 'Der dritte Abschnitt des C#-Tutorials, der das Konzept von mathematischen Operatoren einführt.', 100, 8, '/static/csharp_tutorial/thumbnail.jpg', '/static/csharp_tutorial/Abschnitt_3_Mathematische_Operatoren.mp4'),
       ('If Abfragen', 'Der vierte Abschnitt des C#-Tutorials, der das Konzept von if-Abfragen einführt.', 100, 8, '/static/csharp_tutorial/thumbnail.jpg', '/static/csharp_tutorial/Abschnitt_4_If_Abfragen.mp4'),
       ('Switch Blöcke', 'Der fünfte Abschnitt des C#-Tutorials, der das Konzept von switch Blöcken einführt.', 100, 8, '/static/csharp_tutorial/thumbnail.jpg', '/static/csharp_tutorial/Abschnitt_5_Switch_Bloecke.mp4');

-- Add similar entries for the rest of the videos in the series...


INSERT INTO user_purchases (user_id, product_id)
VALUES (1, 1);
INSERT INTO user_purchases (user_id, product_id)
VALUES (1, 2);

/* Sample comments */
INSERT INTO product_comments (user_id, course_id, comment)
VALUES (1, 1, 'Das ist ein Kommentar');

INSERT INTO product_comments (user_id, course_id, comment)
VALUES (1, 1, 'Das ist ein Kommentar');

INSERT INTO product_comments (user_id, course_id, comment)
VALUES (2, 1, 'Das ist auch ein Kommentar');

INSERT INTO product_comments (user_id, course_id, comment)
VALUES (1, 1, 'das ist ja krass bro');

/* --Indexes-- */
CREATE INDEX idx_user_purchases_user_id ON user_purchases (user_id);
CREATE INDEX idx_user_purchases_product_id ON user_purchases (product_id);

CREATE INDEX idx_video_comments_user_id ON product_comments (user_id);
CREATE INDEX idx_video_comments_video_id ON product_comments (course_id);

CREATE INDEX idx_video_parent_product_id ON video (parent_product_id);

