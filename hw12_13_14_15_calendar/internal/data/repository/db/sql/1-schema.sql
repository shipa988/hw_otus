DROP DATABASE IF EXISTS calendar;
CREATE DATABASE calendar;
CREATE USER igor WITH encrypted password 'igor';
GRANT ALL PRIVILEGES ON DATABASE calendar to igor;