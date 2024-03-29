--------------------------------------------------------
--  DDL for Table
--------------------------------------------------------

CREATE TABLE "CGSYSADM"."SAP_ACCOUNTS" (
	   STATUS CHAR(1),	
	   RELEASE_ID INTEGER,
	   BSCS_ACCOUNT VARCHAR2(32),
	   OFI_SAP_ACCOUNT VARCHAR2(32),	
	   VALID_FROM_DATE DATE,
	   VAT_CODE_IND VARCHAR(32),
	   OFI_SAP_WBS_CODE VARCHAR2(32),
	   CIT_MARKER_VAT_FLAG INTEGER,
	   ENTRY_DATE DATE,
	   ENTRY_OWNER VARCHAR2(16),
	   UPDATE_DATE DATE,
	   UPDATE_OWNER VARCHAR2(16),	   
	   RELEASE_DATE DATE,
	   RELEASE_OWNER VARCHAR(16),
       REC_VERSION INTEGER  DEFAULT 0	   
) SEGMENT CREATION IMMEDIATE 
PCTFREE 10 PCTUSED 40 INITRANS 1 MAXTRANS 255 
NOCOMPRESS NOLOGGING
STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
TABLESPACE "DATA_ROT" ;

COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."STATUS" IS 'Status of the ntry, W like Working, C like Controller, P like Production';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."RELEASE_ID" IS 'Running identifier of the release, 0 for Work, else for Production release versions';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."BSCS_ACCOUNT" IS 'BSCS GL account code used for booking, source of the mapping';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."OFI_SAP_ACCOUNT" IS 'OFI SAP account number, destination of the mapping';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."VALID_FROM_DATE" IS 'Validity of the mapping, must be runded down to the 1st day of the month';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."VAT_CODE_IND" IS 'Business property: VAR code indicator';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."OFI_SAP_WBS_CODE" IS 'Business property: SAP WBS code';
COMMENT ON COLUMN "CGSYSADM"."SAP_ACCOUNTS"."CIT_MARKER_VAT_FLAG" IS 'Business property:  SAP marker vAt flag';
COMMENT ON TABLE "CGSYSADM"."SAP_ACCOUNTS"  IS 'BSCS GL account to SAP IFI account mapping table';

--------------------------------------------------------
--  DDL for Index
--------------------------------------------------------

CREATE UNIQUE INDEX "CGSYSADM"."PK_SAP_ACCT_IDX" ON "CGSYSADM"."SAP_ACCOUNTS" ("STATUS", "RELEASE_ID", "BSCS_ACCOUNT") 
PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS NOLOGGING 
STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
TABLESPACE "DATA_ROT" ;

--------------------------------------------------------
--  DDL for Constraints
--------------------------------------------------------

ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS"
ADD CONSTRAINT "PK_SAP_ACCT_IDX" PRIMARY KEY ("STATUS", "RELEASE_ID", "BSCS_ACCOUNT")
USING INDEX PCTFREE 10 INITRANS 2 MAXTRANS 255 COMPUTE STATISTICS NOLOGGING 
STORAGE(INITIAL 65536 NEXT 1048576 MINEXTENTS 1 MAXEXTENTS 2147483645
PCTINCREASE 0 FREELISTS 1 FREELIST GROUPS 1
BUFFER_POOL DEFAULT FLASH_CACHE DEFAULT CELL_FLASH_CACHE DEFAULT)
TABLESPACE "DATA_ROT" ENABLE;

ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS" MODIFY ("STATUS" NOT NULL ENABLE);
ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS" MODIFY ("RELEASE_ID" NOT NULL ENABLE);
ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS" MODIFY ("BSCS_ACCOUNT" NOT NULL ENABLE);
ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS" MODIFY ("ENTRY_DATE" NOT NULL ENABLE);
ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS" MODIFY ("ENTRY_OWNER" NOT NULL ENABLE);

ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS"
ADD CONSTRAINT FK_SAP_ACC_OFI
FOREIGN KEY (OFI_SAP_ACCOUNT)
REFERENCES SAP_OFI_ACCOUNTS (SAP_OFI_ACCOUNT);

ALTER TABLE "CGSYSADM"."SAP_ACCOUNTS"
ADD CONSTRAINT FK_SAP_ACC_BSCS
FOREIGN KEY (BSCS_ACCOUNT)
REFERENCES GLACCOUNTS (GLACODE);

--------------------------------------------------------
--  DDL for Grants
--------------------------------------------------------

GRANT SELECT, INSERT, UPDATE, DELETE ON "CGSYSADM"."SAP_ACCOUNTS" TO SAMAPI;

--------------------------------------------------------
--  DDL for Synoyms
--------------------------------------------------------

CREATE OR REPLACE PUBLIC SYNONYM SAP_ACCOUNTS FOR "CGSYSADM"."SAP_ACCOUNTS";

QUIT
/
