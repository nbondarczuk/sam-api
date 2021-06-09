--------------------------------------------------------
--  DDL for Table
--------------------------------------------------------

DROP TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS";

CREATE TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS" (
	   SAP_OFI_ACCOUNT VARCHAR2(32),
	   NAME VARCHAR2(255),
	   STATUS VARCHAR(8),
	   ENTRY_DATE DATE,
	   ENTRY_OWNER VARCHAR2(16),
	   UPDATE_DATE DATE,
	   UPDATE_OWNER VARCHAR2(16),
	   REC_VERSION INTEGER  DEFAULT 0
) SEGMENT CREATION IMMEDIATE 
PCTFREE 10 PCTUSED 40 INITRANS 1 MAXTRANS 255 
NOCOMPRESS NOLOGGING
STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
TABLESPACE "DATA_ROT" ;

COMMENT ON COLUMN "CGSYSADM"."SAP_OFI_ACCOUNTS"."SAP_OFI_ACCOUNT" IS 'SAP OFI account code';
COMMENT ON COLUMN "CGSYSADM"."SAP_OFI_ACCOUNTS"."NAME" IS 'SAP OFI account description';
COMMENT ON COLUMN "CGSYSADM"."SAP_OFI_ACCOUNTS"."STATUS" IS 'SAP OFI account status';
COMMENT ON TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS"  IS 'SAP OFI account dictionary';

--------------------------------------------------------
--  DDL for Index
--------------------------------------------------------

CREATE UNIQUE INDEX "CGSYSADM"."PK_SAP_OFI_ACCT_IDX" ON "CGSYSADM"."SAP_OFI_ACCOUNTS" ("SAP_OFI_ACCOUNT") 
PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS NOLOGGING 
STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
TABLESPACE "DATA_ROT" ;

--------------------------------------------------------
--  DDL for Constraints
--------------------------------------------------------

ALTER TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS"
ADD CONSTRAINT "PK_SAP_OFI_ACCT_IDX" PRIMARY KEY (SAP_OFI_ACCOUNT)
USING INDEX PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS NOLOGGING 
STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
TABLESPACE "DATA_ROT" ENABLE;

ALTER TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS" MODIFY ("SAP_OFI_ACCOUNT" NOT NULL ENABLE);
ALTER TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS" MODIFY ("ENTRY_DATE" NOT NULL ENABLE);
ALTER TABLE "CGSYSADM"."SAP_OFI_ACCOUNTS" MODIFY ("ENTRY_OWNER" NOT NULL ENABLE);

--------------------------------------------------------
--  DDL for Grants
--------------------------------------------------------

GRANT SELECT, INSERT, UPDATE, DELETE ON "CGSYSADM"."SAP_OFI_ACCOUNTS" TO SAMAPI;

--------------------------------------------------------
--  DDL for Synoyms
--------------------------------------------------------

CREATE OR REPLACE PUBLIC SYNONYM SAP_OFI_ACCOUNTS FOR "CGSYSADM"."SAP_OFI_ACCOUNTS";

QUIT
/

