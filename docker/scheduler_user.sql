-- User: scheduler_user
-- DROP USER scheduler_user;

CREATE USER scheduler_user WITH
    LOGIN
    NOSUPERUSER
    INHERIT
    CREATEDB
    CREATEROLE
    NOREPLICATION;