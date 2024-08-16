DROP DATABASE IF EXISTS donation_rest;
CREATE DATABASE donation_rest;

\c donation_rest;

CREATE TABLE user_registed (
    id SERIAL PRIMARY KEY,
    login VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    token VARCHAR NOT NULL
);

CREATE TABLE donation_list (
    id SERIAL PRIMARY KEY,
    loginstrimer VARCHAR NOT NULL,
    namesub VARCHAR NOT NULL,
    val INTEGER,
    text VARCHAR NOT NULL
);

