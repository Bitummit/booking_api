-- +goose Up
-- +goose StatementBegin
CREATE TYPE  status_enum AS ENUM ('created', 'submitted', 'closed');

CREATE TABLE IF NOT EXISTS my_user(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(1024) NOT NULL,
    birthday DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS tag(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS city(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS hotel(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    city_id INT REFERENCES city (id)
);

CREATE TABLE IF NOT EXISTS tag_hotel(
    id SERIAL PRIMARY KEY,
    hotel_id INT REFERENCES hotel (id),
    tag_id INT REFERENCES tag (id)
);

CREATE TABLE IF NOT EXISTS room_category(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL NOT NULL,
    сapacity INT NOT NULL,
    description TEXT,
    size INT NOT NULL,
    hotel_id INT REFERENCES hotel (id)
);

CREATE TABLE IF NOT EXISTS room(
    id SERIAL PRIMARY KEY,
    number INT NOT NULL,
    category_id INT REFERENCES room_category (id)
);

CREATE TABLE IF NOT EXISTS booking(
    id SERIAL PRIMARY KEY,
    entry_date DATE NOT NULL,
    leave_date DATE NOT NULL,
    price DECIMAL NOT NULL,
    current_status status_enum,
    guests_count INT NOT NULL,
    user_id INT REFERENCES my_user (id),
    room_id INT REFERENCES room (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE my_user;
DROP TABLE tag;
DROP TABLE city;
DROP TABLE hotel;
DROP TABLE tag_hotel;
DROP TABLE room_category;
DROP TABLE room;
DROP TABLE booking;
DROP TYPE status_enum;
-- +goose StatementEnd
