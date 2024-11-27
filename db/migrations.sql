-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY,          -- 用户ID
                                     username VARCHAR(50) NOT NULL UNIQUE,       -- 登录名
    password VARCHAR(100) NOT NULL,             -- 密码
    nickname VARCHAR(50) DEFAULT 'User',        -- 昵称
    gender ENUM('Male', 'Female') DEFAULT 'Male' -- 性别
    );

-- 创建书籍表
CREATE TABLE IF NOT EXISTS books (
                                     id INT AUTO_INCREMENT PRIMARY KEY,          -- 书籍ID
                                     user_id INT NOT NULL,                       -- 用户ID（外键）
                                     isbn VARCHAR(50) NOT NULL UNIQUE,           -- 书号（ISBN）
    title VARCHAR(100) NOT NULL,                -- 书名
    author VARCHAR(100) NOT NULL,               -- 作者
    publisher VARCHAR(100) NOT NULL,            -- 出版社
    description TEXT,                           -- 简介
    cover_image VARCHAR(255),                   -- 封面图片URL
    FOREIGN KEY (user_id) REFERENCES users(id)  -- 外键关联用户表
    );
