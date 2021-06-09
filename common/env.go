package common

import (
	"flag"
	"fmt"
	"log"
	"os"

	. "sam-api/utl/str"
)

// process environment
var (
	fversion                bool = false
	fconfig                 string
	fserveripaddress        string
	fserverport             string
	frunpath                string
	fkeypath                string
	fdebug                  string
	foracledbuser           string
	foracledbpassword       string
	foracleservicename      string
	falertmailaddress       string
	falertmailserveraddress string
	falertmailsenderaddress string
	fjwttokenvalidhours     string
	fldapbase               string
	fldaphost               string
	fldapport               string
	fldapbinddn             string
	TestRun                 bool = false
)

// get invocation flags
func FlagsInit() {
	flag.BoolVar(&fversion, "v", false, "Version check")
	flag.StringVar(&fconfig, "config", "config.json", "Config json file if not in $RUNPATH/config.json")
	flag.StringVar(&fserveripaddress, "serveripaddress", "0.0.0.0", "Server address")
	flag.StringVar(&fserverport, "serverport", "8000", "Server port")
	flag.StringVar(&frunpath, "runpath", ".", "Run path")
	flag.StringVar(&fkeypath, "keypath", "keys", "Key path")
	flag.StringVar(&fdebug, "debug", "0", "Debug level")
	flag.StringVar(&foracledbuser, "oracledbuser", "", "Oracle DB user")
	flag.StringVar(&foracledbpassword, "oracledbpassword", "", "Oracle DB password")
	flag.StringVar(&foracleservicename, "oracleservicename", "", "Oracle service name")
	flag.StringVar(&falertmailaddress, "alertmailaddress", "root@localhost", "Alert mail address")
	flag.StringVar(&falertmailserveraddress, "alertmailserveraddress", "", "Alert SNMP server address")
	flag.StringVar(&falertmailsenderaddress, "alertmailsenderaddress", "samapi@localhost", "Alert mail sender address")
	flag.StringVar(&fjwttokenvalidhours, "jwttokenvalidhours", "1", "JWT token validity period in hours")
	flag.StringVar(&fldapbase, "ldapbase", "", "LDAP base")
	flag.StringVar(&fldaphost, "ldaphost", "", "LDAP host")
	flag.StringVar(&fldapport, "ldapport", "", "LDAP port")
	flag.StringVar(&fldapbinddn, "ldapbinddn", "", "LDAP bind DN")
}

// load env variables if they are set otherwise use default values or config file
func EnvInit(version, build, level string) {
	flag.Parse()

	// shortcut, anyway we would print it out
	setVersion(version, build, level)
	if fversion {
		fmt.Printf("The version is: %s\n", GetVersion())
		os.Exit(0)
	}

	if version == "test" {
		TestRun = true
	}

	config := os.Getenv("CONFIG")
	if Empty(config) {
		config = fconfig
	}
	InitConfig(config)

	//override loaded config values with env variables or command line switches
	AppConfig.ServerIPAddress = Nvl(Nvl(os.Getenv("SERVERIPADDRESS"), fserveripaddress), AppConfig.ServerIPAddress)
	AppConfig.ServerPort = Nvl(Nvl(os.Getenv("SERVERPORT"), fserverport), AppConfig.ServerPort)
	AppConfig.RunPath = Nvl(Nvl(os.Getenv("RUNPATH"), frunpath), AppConfig.RunPath)
	AppConfig.KeyPath = Nvl(Nvl(os.Getenv("KEYPATH"), fkeypath), AppConfig.KeyPath)
	AppConfig.Debug = Nvl(Nvl(os.Getenv("DEBUG"), fdebug), AppConfig.Debug)
	AppConfig.OracleDBUser = Nvl(Nvl(os.Getenv("ORACLEDBUSER"), foracledbuser), AppConfig.OracleDBUser)
	AppConfig.OracleDBPassword = Nvl(Nvl(os.Getenv("ORACLEDBPASSWORD"), foracledbpassword), AppConfig.OracleDBPassword)
	AppConfig.OracleServiceName = Nvl(Nvl(os.Getenv("ORACLESERVICENAME"), foracleservicename), AppConfig.OracleServiceName)
	AppConfig.AlertMailAddress = Nvl(Nvl(os.Getenv("ALERTMAILADDRESS"), falertmailaddress), AppConfig.AlertMailAddress)
	AppConfig.AlertMailServerAddress = Nvl(Nvl(os.Getenv("ALERTMAILSERVERADDRESS"), falertmailserveraddress), AppConfig.AlertMailServerAddress)
	AppConfig.AlertMailSenderAddress = Nvl(Nvl(os.Getenv("ALERTMAILSENDERADDRESS"), falertmailsenderaddress), AppConfig.AlertMailSenderAddress)
	AppConfig.JWTTokenValidHours = Nvl(Nvl(os.Getenv("JWTTOKENVALIDHOURS"), fjwttokenvalidhours), AppConfig.JWTTokenValidHours)
	AppConfig.LdapBase = Nvl(Nvl(os.Getenv("LDAPBASE"), fldapbase), AppConfig.LdapBase)
	AppConfig.LdapHost = Nvl(Nvl(os.Getenv("LDAPHOST"), fldaphost), AppConfig.LdapHost)
	AppConfig.LdapPort = Nvl(Nvl(os.Getenv("LDAPPORT"), fldapport), AppConfig.LdapPort)
	AppConfig.LdapBindDN = Nvl(Nvl(os.Getenv("LDAPBINDDN"), fldapbinddn), AppConfig.LdapBindDN)

	EnvLog()
}

// show in log working environment
func EnvLog() {
	log.Printf("API Server execution environment")
	log.Printf("%s: %s", "ServerIPAddress       ", AppConfig.ServerIPAddress)
	log.Printf("%s: %s", "ServerPort            ", AppConfig.ServerPort)
	log.Printf("%s: %s", "RunPath               ", AppConfig.RunPath)
	log.Printf("%s: %s", "KeyPath               ", AppConfig.KeyPath)
	log.Printf("%s: %s", "Debug                 ", AppConfig.Debug)
	log.Printf("%s: %s", "OracleDBUser          ", AppConfig.OracleDBUser)
	log.Printf("%s: %s", "OracleDBPassword      ", AppConfig.OracleDBPassword)
	log.Printf("%s: %s", "OracleServiceName     ", AppConfig.OracleServiceName)
	log.Printf("%s: %s", "AlertMailAddress      ", AppConfig.AlertMailAddress)
	log.Printf("%s: %s", "AlertMailServerAddress", AppConfig.AlertMailServerAddress)
	log.Printf("%s: %s", "AlertMailSenderAddress", AppConfig.AlertMailSenderAddress)
	log.Printf("%s: %s", "JWTTokenValidHours    ", AppConfig.JWTTokenValidHours)
	log.Printf("%s: %s", "LdapBase              ", AppConfig.LdapBase)
	log.Printf("%s: %s", "LdapHost              ", AppConfig.LdapHost)
	log.Printf("%s: %s", "LdapPort              ", AppConfig.LdapPort)
	log.Printf("%s: %s", "LdapBindDN            ", AppConfig.LdapBindDN)
}
