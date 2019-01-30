CREATE DATABASE RedCoins;
USE RedCoins;
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    birthday VARCHAR(50) NOT NULL,
    balance FLOAT NOT NULL
);

CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR(50) NOT NULL,
    hour VARCHAR(50) NOT NULL,
    bitcoins FLOAT NOT NULL,
    convert_tx FLOAT NOT NULL,
    final_value FLOAT NOT NULL,
    user_id_1 INT NOT NULL,
    user_id_2 INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_id_1) REFERENCES users(id)
);