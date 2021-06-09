--------------------------------------------------------
--  DDL for Table
--------------------------------------------------------

CREATE OR REPLACE FORCE VIEW "CGSYSADM"."GLACCOUNTS" ("GLACODE", "GLADESC", "GLATYPE", "GLACTIVE", "ENTRY_DATE", "ENTRY_OWNER", "UPDATE_DATE", "UPDATE_OWNER") AS
SELECT GLACODE, GLADESC, GLATYPE, GLAACTIVE, GLAENTDATE, 'SYSADM', GLAMODDATE, 'SYSADM'
FROM glaccount_all@BSCSDB.WORLD;

SHOW ERROR

GRANT SELECT ON "CGSYSADM"."GLACCOUNTS" TO SAMAPI;

--------------------------------------------------------
--  DDL for Synoyms
--------------------------------------------------------

CREATE OR REPLACE PUBLIC SYNONYM GLACCOUNTS FOR "CGSYSADM"."GLACCOUNTS";

alter user SAMAPI grant connect through CGSYSADM;

QUIT
/

