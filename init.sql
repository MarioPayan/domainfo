DROP DATABASE IF EXISTS domainfo CASCADE;
CREATE USER IF NOT EXISTS domainfouser;
CREATE DATABASE domainfo;
GRANT ALL ON DATABASE domainfo TO domainfouser;