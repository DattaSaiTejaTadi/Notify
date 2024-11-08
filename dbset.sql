CREATE DATABASE user_management;

USE user_management;

CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,       
    name VARCHAR(255) NOT NULL,
    domainid VARCHAR(255) NOT NULL UNIQUE, 
    password VARCHAR(255) NOT NULL, 
    role varchar(50)      NOT NULL,   
    active_state BOOLEAN DEFAULT TRUE, 
    INDEX (domainid)                   
    );