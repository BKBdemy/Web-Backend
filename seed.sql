TRUNCATE licenses RESTART IDENTITY CASCADE;
TRUNCATE product_licenses RESTART IDENTITY CASCADE;
TRUNCATE products RESTART IDENTITY CASCADE;
TRUNCATE user_purchases RESTART IDENTITY CASCADE;
TRUNCATE users RESTART IDENTITY CASCADE;

INSERT INTO users (username, password, balance)
VALUES ('admin', '$argon2id$v=19$m=256000,t=6,p=1$dGVzdHRlc3Q$MMMzLViNOBi+zmhnFWj4y1y6TqYfRvmUAI6BiH30mIk', 1000); /* password is admin */

INSERT INTO users (username, password)
VALUES ('user', '$argon2id$v=19$m=256000,t=6,p=1$dGVzdHRlc3Q$MMMzLViNOBi+zmhnFWj4y1y6TqYfRvmUAI6BiH30mIk'); /* password is admin */

INSERT INTO products (name, description, price, image, mpd_url)
VALUES ('Product 1', 'Product 1 description', 1000, 'product1.jpg', 'product1.mpd');

INSERT INTO products (name, description, price, image, mpd_url)
VALUES ('Product 2', 'Product 2 description', 2000, 'product2.jpg', 'product2.mpd');

INSERT INTO licenses (key_id, encryption_key, product_id) VALUES ('f3c5e0361e6654b28f8049c778b23946', 'a4631a153a443df9eed0593043db7519', 1);
INSERT INTO licenses (key_id, encryption_key, product_id) VALUES ('abba271e8bcf552bbd2e86a434a9a5d9', '69eaa802a6763af979e8d1940fb88392', 1);

INSERT INTO licenses (key_id, encryption_key, product_id) VALUES ('f3c5e0361e6654b28f8049c778b23946', 'a4631a153a443df9eed0593043db7519', 2);
INSERT INTO licenses (key_id, encryption_key, product_id) VALUES ('abba271e8bcf552bbd2e86a434a9a5d9', '69eaa802a6763af979e8d1940fb88392', 2);

INSERT INTO product_licenses (product_id, license_id)
VALUES (1, 1);
INSERT INTO product_licenses (product_id, license_id)
VALUES (1, 2);
INSERT INTO product_licenses (product_id, license_id)
VALUES (2, 3);
INSERT INTO product_licenses (product_id, license_id)
VALUES (2, 4);

INSERT INTO user_purchases (user_id, product_id)
VALUES (1, 1);
INSERT INTO user_purchases (user_id, product_id)
VALUES (1, 2);
