DROP USER SAMAPI CASCADE;

CREATE USER SAMAPI IDENTIFIED BY "SAMAPI";

GRANT CONNECT TO SAMAPI;

ALTER USER "SAMAPI"
DEFAULT TABLESPACE "DATA_ROT"
TEMPORARY TABLESPACE "TEMP"
ACCOUNT UNLOCK ;

ALTER USER "SAMAPI" QUOTA UNLIMITED ON "DATA_ROT";

QUIT
/
