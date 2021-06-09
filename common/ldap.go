package common

import (
	"log"
	"os/exec"
	"syscall"
)

var bindPassword, groupFilter, serverName, userFilter string
var port int
var useSSL bool
var skipTLS bool

func LdapServerAuth(user, role, password string) (err error) {
	if AppConfig.LdapBase == "" || AppConfig.LdapHost == "" || AppConfig.LdapPort == "" || AppConfig.LdapBindDN == "" {
		log.Printf("Skip LDAP validation of the user: %s", user)
		return nil
	}

	log.Printf("Running LDAP check with: %s", "ldapcheck")
	cmd := exec.Command("ldapcheck", user, password, role)
	if err := cmd.Start(); err != nil {
		log.Printf("cmd.Start: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Printf("cmd.Wait: %v", err)
		}
	}

	log.Printf("LDAP check result: %s", err)

	return err
}
