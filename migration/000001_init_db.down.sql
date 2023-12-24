DROP TABLE IF EXISTS auth.users CASCADE;
DROP TABLE IF EXISTS auth.emails CASCADE;
DROP TABLE IF EXISTS auth.verification_codes CASCADE;
DROP TABLE IF EXISTS auth.tokens CASCADE;
DROP TABLE IF EXISTS auth.reset_pwd_tokens CASCADE;

DROP TRIGGER IF EXISTS users_insert_email ON auth.users CASCADE;
DROP TRIGGER IF EXISTS users_insert_phone ON auth.users CASCADE;

DROP TABLE IF EXISTS auth.phone_numbers CASCADE;
DROP TABLE IF EXISTS auth.refresh_tokens;

DROP SCHEMA IF EXISTS auth;