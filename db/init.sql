-- for auth service
CREATE USER auth WITH ENCRYPTED PASSWORD 'authpwd';
CREATE DATABASE auth
    WITH
    OWNER = auth
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

GRANT TEMPORARY, CONNECT ON DATABASE auth TO PUBLIC;
GRANT ALL ON DATABASE auth TO auth;
GRANT ALL ON DATABASE auth TO auth;


-- for auth accural
CREATE USER accural WITH ENCRYPTED PASSWORD 'accuralpwd';
CREATE DATABASE accural
    WITH
    OWNER = accural
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

GRANT TEMPORARY, CONNECT ON DATABASE accural TO PUBLIC;
GRANT ALL ON DATABASE accural TO accural;
GRANT ALL ON DATABASE accural TO accural;


-- for auth gophermart
CREATE USER gophermart WITH ENCRYPTED PASSWORD 'gophermartpwd';
CREATE DATABASE gophermart
    WITH
    OWNER = gophermart
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

GRANT TEMPORARY, CONNECT ON DATABASE gophermart TO PUBLIC;
GRANT ALL ON DATABASE gophermart TO gophermart;
GRANT ALL ON DATABASE gophermart TO gophermart;