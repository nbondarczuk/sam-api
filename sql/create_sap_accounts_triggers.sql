--------------------------------------------------------
--  DDL for Triggers
--------------------------------------------------------
CREATE OR REPLACE TRIGGER sap_accounts_bef_ins
BEFORE INSERT
   ON sap_accounts
   FOR EACH ROW
BEGIN
   IF :new.STATUS = 'P' AND :new.VALID_FROM_DATE < SYSDATE THEN
	  RAISE_APPLICATION_ERROR(-20000, 'Valid Date in the past, validation  error');
   END IF;
   :new.VALID_FROM_DATE := ROUND(TO_DATE(:new.VALID_FROM_DATE),'DAY');
END;
/   

SHOW ERROR

CREATE OR REPLACE TRIGGER sap_accounts_aft_ins
AFTER INSERT
   ON sap_accounts
   FOR EACH ROW
BEGIN
   INSERT INTO sap_accounts_log
   (
       OPCODE,
	   OPDATE,
 	   STATUS,
	   RELEASE_ID,
	   BSCS_ACCOUNT,
	   OFI_SAP_ACCOUNT,
	   VALID_FROM_DATE,
	   VAT_CODE_IND,
	   OFI_SAP_WBS_CODE,
	   CIT_MARKER_VAT_FLAG,
	   ENTRY_DATE,
	   ENTRY_OWNER,
	   UPDATE_DATE,
	   UPDATE_OWNER,
	   RELEASE_DATE,
	   RELEASE_OWNER,
	   REC_VERSION
   )
   VALUES
   (
       'I',
	   SYSDATE,
 	   :new.STATUS,
	   :new.RELEASE_ID,
	   :new.BSCS_ACCOUNT,
	   :new.OFI_SAP_ACCOUNT,
	   :new.VALID_FROM_DATE,
	   :new.VAT_CODE_IND,
	   :new.OFI_SAP_WBS_CODE,
	   :new.CIT_MARKER_VAT_FLAG,
	   :new.ENTRY_DATE,
	   :new.ENTRY_OWNER,
	   :new.UPDATE_DATE,
	   :new.UPDATE_OWNER,
	   :new.RELEASE_DATE,
	   :new.RELEASE_OWNER,
	   :new.REC_VERSION
   );
END;
/

SHOW ERROR

CREATE OR REPLACE TRIGGER sap_accounts_bef_upd
BEFORE UPDATE
   ON sap_accounts
   FOR EACH ROW
BEGIN
   IF :new.STATUS = 'P' AND :new.VALID_FROM_DATE < SYSDATE THEN
	  RAISE_APPLICATION_ERROR(-20000, 'Valid Date in the past, validation  error');
   END IF;
   :new.REC_VERSION := :new.REC_VERSION + 1;
   :new.VALID_FROM_DATE := ROUND(TO_DATE(:new.VALID_FROM_DATE),'DAY');
END;
/

SHOW ERROR

CREATE OR REPLACE TRIGGER sap_accounts_aft_upd
AFTER UPDATE
   ON sap_accounts
   FOR EACH ROW
BEGIN
   INSERT INTO sap_accounts_log
   (
       OPCODE,
	   OPDATE,
 	   STATUS,
	   RELEASE_ID,
	   BSCS_ACCOUNT,
	   OFI_SAP_ACCOUNT,
	   VALID_FROM_DATE,
	   VAT_CODE_IND,
	   OFI_SAP_WBS_CODE,
	   CIT_MARKER_VAT_FLAG,
	   ENTRY_DATE,
	   ENTRY_OWNER,
	   UPDATE_DATE,
	   UPDATE_OWNER,
	   RELEASE_DATE,
	   RELEASE_OWNER,
	   REC_VERSION
   )
   VALUES
   (
       'U',
	   SYSDATE,
 	   :new.STATUS,
	   :new.RELEASE_ID,
	   :new.BSCS_ACCOUNT,
	   :new.OFI_SAP_ACCOUNT,
	   :new.VALID_FROM_DATE,
	   :new.VAT_CODE_IND,
	   :new.OFI_SAP_WBS_CODE,
	   :new.CIT_MARKER_VAT_FLAG,
	   :new.ENTRY_DATE,
	   :new.ENTRY_OWNER,
	   :new.UPDATE_DATE,
	   :new.UPDATE_OWNER,
	   :new.RELEASE_DATE,
	   :new.RELEASE_OWNER,
	   :new.REC_VERSION
   );
END;
/

SHOW ERROR

CREATE OR REPLACE TRIGGER sap_accounts_bef_del
BEFORE DELETE
   ON sap_accounts
   FOR EACH ROW
BEGIN
   IF :old.STATUS = 'P' AND :old.VALID_FROM_DATE < SYSDATE THEN
	  RAISE_APPLICATION_ERROR(-20000, 'Valid Date in the past, validation  error');
   END IF;	        
END;
/

SHOW ERROR

CREATE OR REPLACE TRIGGER sap_accounts_aft_del
AFTER DELETE
   ON sap_accounts
   FOR EACH ROW
BEGIN
   INSERT INTO sap_accounts_log
   (
       OPCODE,
	   OPDATE,
 	   STATUS,
	   RELEASE_ID,
	   BSCS_ACCOUNT,
	   OFI_SAP_ACCOUNT,
	   VALID_FROM_DATE,
	   VAT_CODE_IND,
	   OFI_SAP_WBS_CODE,
	   CIT_MARKER_VAT_FLAG,
	   ENTRY_DATE,
	   ENTRY_OWNER,
	   UPDATE_DATE,
	   UPDATE_OWNER,
	   RELEASE_DATE,
	   RELEASE_OWNER,
	   REC_VERSION
   )
   VALUES
   (
       'D',
	   SYSDATE,
 	   :old.STATUS,
	   :old.RELEASE_ID,
	   :old.BSCS_ACCOUNT,
	   :old.OFI_SAP_ACCOUNT,
	   :old.VALID_FROM_DATE,
	   :old.VAT_CODE_IND,
	   :old.OFI_SAP_WBS_CODE,
	   :old.CIT_MARKER_VAT_FLAG,
	   :old.ENTRY_DATE,
	   :old.ENTRY_OWNER,
	   :old.UPDATE_DATE,
	   :old.UPDATE_OWNER,
	   :old.RELEASE_DATE,
	   :old.RELEASE_OWNER,
	   :old.REC_VERSION
   );
END;
/

SHOW ERROR

QUIT
/

