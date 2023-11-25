-- Создание схемы "auth"
CREATE SCHEMA auth;

-- Создание таблицы "auth.users"
CREATE TABLE auth.users (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    hash_password VARCHAR(255) NOT NULL,
    fullname VARCHAR(255) NOT NULL
);

-- Создание таблицы "auth.emails" 
CREATE TABLE auth.emails (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- Создание таблицы "auth.phone_numbers" 
CREATE TABLE auth.phone_numbers (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    phone_number VARCHAR(11) UNIQUE NOT NULL
);

-- Создание таблицы "auth.verification_codes"
CREATE TABLE auth.verification_codes (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL
);

-- Создание таблицы "auth.tokens"
CREATE TABLE auth.tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expiration_time TIMESTAMP NOT NULL
);

-- Создание таблицы "auth.reset_pwd_tokens"
CREATE TABLE auth.reset_pwd_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users(user_id) ON DELETE CASCADE,
    reset_token VARCHAR(255) NOT NULL
);

-- Создание триггера для заполнения таблицы "auth.emails" при внесении записи в "auth.users"
-- CREATE OR REPLACE FUNCTION insert_into_emails() RETURNS TRIGGER AS $$
-- BEGIN
--     IF NEW.email <> '' THEN
--         INSERT INTO auth.emails (user_id, email) VALUES (NEW.user_id, NEW.email);
--     END IF;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- Создание триггера, связывающего "auth.users" с "auth.emails"
-- CREATE TRIGGER users_insert_email
-- AFTER INSERT
-- ON auth.users
-- FOR EACH ROW
-- EXECUTE FUNCTION insert_into_emails();

-- -- Создание триггера для заполнения таблицы "auth.phone_numbers" при внесении записи в "auth.users"
-- CREATE OR REPLACE FUNCTION insert_into_phones() RETURNS TRIGGER AS $$
-- BEGIN
--     IF NEW.phone_number <> '' THEN
--         INSERT INTO auth.phone_numbers (user_id, phone_number) VALUES (NEW.user_id, NEW.phone_number);
--     END IF;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- Создание триггера, связывающего "auth.users" с "auth.phone_numbers"
-- CREATE TRIGGER users_insert_phone
-- AFTER INSERT
-- ON auth.users
-- FOR EACH ROW
-- EXECUTE FUNCTION insert_into_phones();



