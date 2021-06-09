#/bin/bash

# make user SAMAPI to handle the connection
ORA='SYSTEM/oracle@XE'
sqlplus ${ORA} @create_user_samapi.sql

# make tables in CGSYSADM schema to be accessed by SAMAPI
ORA="CGSYSADM/cgsysadm17@XE"
sqlplus ${ORA} @dropall.sql
sqlplus ${ORA} @create_sap_ofi_bscs_glaccounts_table.sql
sqlplus ${ORA} @create_sap_ofi_accounts.sql
sqlplus ${ORA} @create_customer_segment.sql
sqlplus ${ORA} @create_sap_accounts.sql
sqlplus ${ORA} @create_sap_accounts_log.sql
sqlplus ${ORA} @create_sap_accounts_triggers.sql
sqlplus ${ORA} @create_sap_accounts_p_view.sql
sqlplus ${ORA} @create_sap_acc_segm_order_numbers.sql
sqlplus ${ORA} @create_sap_acc_segm_order_numbers_log.sql
sqlplus ${ORA} @create_sap_acc_segm_order_numbers_triggers.sql
sqlplus ${ORA} @create_sap_acc_segm_order_numbers_p_view.sql






