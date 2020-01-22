CREATE USER docker WITH
    LOGIN
    SUPERUSER
    INHERIT
    PASSWORD 'docker';

CREATE DATABASE docker
    WITH
    OWNER = docker
    ENCODING = 'UTF8';
