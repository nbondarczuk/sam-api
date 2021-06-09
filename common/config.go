package common

import (
	"encoding/json"
	"log"
	"os"
)

type (
	configuration struct {
		ServerIPAddress,
		ServerPort,
		RunPath,
		KeyPath,
		Debug,
		OracleDBUser,
		OracleDBPassword,
		OracleServiceName,
		AlertMailAddress,
		AlertMailServerAddress,
		AlertMailSenderAddress,
		JWTTokenValidHours,
		LdapBase,
		LdapHost,
		LdapPort,
		LdapBindDN,
		Testing string
	}
)

// AppConfig holds the configuration values from config.json file
var AppConfig configuration

// Initialize AppConfig
func InitConfig(config string) {
	loadAppConfig(config)
}

//
// Reads config.json and decode into AppConfig
//
func loadAppConfig(config string) {
	// config file name may be overriden by explicit invocation parameter
	log.Printf("Using config file: %s", config)

	// file must exist
	file, err := os.Open(config)
	defer file.Close()
	if err != nil {
		log.Fatalf("Can not open file: %s\n", err)
	}

	// decode json contents
	decoder := json.NewDecoder(file)
	AppConfig = configuration{}
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("Can not decode config file: %s\n", err)
	}
}
