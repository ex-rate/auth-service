-- Создание схемы "auth"
CREATE SCHEMA auth;

-- Создание таблицы "auth.users"
CREATE TABLE auth.users (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR(255),
    hash_password VARCHAR(255),
    email VARCHAR(255),
    phone_number VARCHAR(11),
    fullname VARCHAR(255)
);

-- Создание таблицы "auth.emails" (добавлена отсутствовавшая таблица)
CREATE TABLE auth.emails (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id),
    email VARCHAR(255) NOT NULL
);

-- Создание таблицы "auth.verification_codes"
CREATE TABLE auth.verification_codes (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id),
    code VARCHAR(255) NOT NULL
);

-- Создание таблицы "auth.tokens"
CREATE TABLE auth.tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id),
    token VARCHAR(255) NOT NULL,
    expiration_time TIMESTAMP NOT NULL
);

-- Создание таблицы "auth.reset_pwd_tokens"
CREATE TABLE auth.reset_pwd_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id),
    reset_token VARCHAR(255) NOT NULL
);

-- Создание триггера для заполнения таблицы "auth.emails" при внесении записи в "auth.users"
CREATE OR REPLACE FUNCTION insert_into_emails() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO auth.emails (user_id, email) VALUES (NEW.user_id, NEW.email);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггера, связывающего "auth.users" с "auth.emails"
CREATE TRIGGER users_insert
AFTER INSERT
ON auth.users
FOR EACH ROW
EXECUTE FUNCTION insert_into_emails();
