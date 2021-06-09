--------------------------------------------------------
--  DDL for Table
--------------------------------------------------------

CREATE OR REPLACE FORCE VIEW CGSYSADM.SAP_ACC_SEGM_ORDER_NUMBERS_P
(
	STATUS,
	RELEASE_ID,
	BSCS_ACCOUNT,
	SEGMENT_CODE,
	VALID_FROM_DATE,
	ORDER_NUMBER,
	ENTRY_DATE,
	ENTRY_OWNER,
	UPDATE_DATE,
	UPDATE_OWNER,
	RELEASE_DATE,
	RELEASE_OWNER
)
AS
SELECT
	STATUS,
	RELEASE_ID,
	BSCS_ACCOUNT,
	SEGMENT_CODE,
	VALID_FROM_DATE,
	ORDER_NUMBER,
	ENTRY_DATE,
	ENTRY_OWNER,
	UPDATE_DATE,
	UPDATE_OWNER,
	RELEASE_DATE,
	RELEASE_OWNER
FROM CGSYSADM.SAP_ACC_SEGM_ORDER_NUMBERS Q
WHERE STATUS = 'P'
AND RELEASE_ID = (SELECT MAX(RELEASE_ID) FROM CGSYSADM.SAP_ACC_SEGM_ORDER_NUMBERS SQ WHERE SQ.BSCS_ACCOUNT = Q.BSCS_ACCOUNT AND SQ.SEGMENT_CODE = Q.SEGMENT_CODE);

SHOW ERROR

GRANT SELECT ON CGSYSADM.SAP_ACC_SEGM_ORDER_NUMBERS_P TO SAMAPI;

--------------------------------------------------------
--  DDL for Synoyms
--------------------------------------------------------

CREATE OR REPLACE PUBLIC SYNONYM SAP_ACC_SEGM_ORDER_NUMBERS_P FOR CGSYSADM.SAP_ACC_SEGM_ORDER_NUMBERS_P;

QUIT
/

