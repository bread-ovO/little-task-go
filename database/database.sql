-- 创建数据库
CREATE DATABASE IF NOT EXISTS book_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE book_management;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    gender ENUM('Male', 'Female', 'Other') NOT NULL
);

-- 创建书籍表
CREATE TABLE IF NOT EXISTS books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    book_number VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(255) NOT NULL,
    description TEXT,
    cover_image VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 插入5个用户
INSERT INTO users (username, password, nickname, gender) VALUES
('user1', 'password1', '昵称1', 'Male'),
('user2', 'password2', '昵称2', 'Female'),
('user3', 'password3', '昵称3', 'Other'),
('user4', 'password4', '昵称4', 'Male'),
('user5', 'password5', '昵称5', 'Female');

-- 为每个用户插入100本书籍
DELIMITER //

CREATE PROCEDURE InsertBooks()
BEGIN
    DECLARE u INT DEFAULT 1;
    DECLARE i INT;
    WHILE u <= 5 DO
        SET i = 1;
        WHILE i <= 100 DO
            INSERT INTO books (user_id, book_number, title, author, publisher, description, cover_image)
            VALUES (
                u,
                CONCAT('BN', u, LPAD(i, 3, '0')),
                CONCAT('书名', u, '-', i),
                CONCAT('作者', u),
                CONCAT('出版社', u),
                CONCAT('这是书籍', u, '-', i, '的简介。'),
                'default_cover.jpg'
            );
            SET i = i + 1;
        END WHILE;
        SET u = u + 1;
    END WHILE;
END //

DELIMITER ;

CALL InsertBooks();

DROP PROCEDURE InsertBooks;

