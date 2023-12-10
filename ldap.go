package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func registerUser(username string, password string, email string) (string, error) {
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		message := "Failed to connect to LDAP server."
		return message, err
	}
	defer l.Close()

	// Bind with Admin
	err = l.Bind("cn=admin,dc=skinny,dc=wsso", os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != nil {
		message := "Failed to bind with LDAP server."
		return message, err
	}

	addRequest := ldap.NewAddRequest("uid="+username+",ou=users,dc=skinny,dc=wsso", nil)
	addRequest.Attribute("objectClass", []string{"top", "posixAccount", "shadowAccount", "inetOrgPerson"})
	addRequest.Attribute("uid", []string{username})
	addRequest.Attribute("cn", []string{username})
	addRequest.Attribute("sn", []string{username})
	addRequest.Attribute("userPassword", []string{password})
	addRequest.Attribute("loginShell", []string{"/bin/bash"})
	addRequest.Attribute("uidNumber", []string{"10000"})
	addRequest.Attribute("gidNumber", []string{"10000"})
	addRequest.Attribute("homeDirectory", []string{"/home/" + username})
	addRequest.Attribute("shadowLastChange", []string{"0"})
	addRequest.Attribute("shadowMax", []string{"99999"})
	addRequest.Attribute("shadowWarning", []string{"7"})
	addRequest.Attribute("mail", []string{email})
	err = l.Add(addRequest)

	if err != nil {
		if strings.Contains(err.Error(), "68") {
			message := fmt.Sprintf("Username %s is not available!", username)
			return message, err
		}
		message := "Failed to register your account. Please contact an administrator."
		return message, err
	}

	message := "Account created successfully!"
	return message, nil
}

func deleteLdapUser(username string) (string, error) {
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		message := "Failed to connect to LDAP server."
		return message, err
	}
	defer l.Close()

	// Bind with Admin
	err = l.Bind("cn=admin,dc=skinny,dc=wsso", os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != nil {
		message := "Failed to bind with LDAP server."
		return message, err
	}

	delRequest := ldap.NewDelRequest("uid="+username+",ou=users,dc=skinny,dc=wsso", nil)
	err = l.Del(delRequest)

	if err != nil {
		message := "Failed to delete your account."
		return message, err
	}

	message := "Account deleted successfully!"
	return message, nil
}

func getLdapUsers() ([]string, error) {
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		message := "Failed to connect to LDAP server."
		return []string{message}, err
	}
	defer l.Close()

	// Bind with Admin
	err = l.Bind("cn=admin,dc=skinny,dc=wsso", os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != nil {
		message := "Failed to bind with LDAP server."
		return []string{message}, err
	}

	searchRequest := ldap.NewSearchRequest(
		"ou=users,dc=skinny,dc=wsso",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=posixAccount)",
		[]string{"uid"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		message := "Failed to list users."
		return []string{message}, err
	}

	var users []string
	for _, entry := range sr.Entries {
		users = append(users, entry.GetAttributeValue("uid"))
	}

	return users, nil
}

func isMemberOf(username string, group string) (bool, error) {
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return false, err
	}
	defer l.Close()

	// Bind with Admin
	err = l.Bind("cn=admin,dc=skinny,dc=wsso", os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != nil {
		return false, err
	}

	uniqueMember := fmt.Sprintf("uid=%s,ou=users,dc=skinny,dc=wsso", username)
	searchFilter := fmt.Sprintf("(&(objectClass=groupOfUniqueNames)(uniqueMember=%s))", uniqueMember)
	searchRequest := ldap.NewSearchRequest(
		"ou=groups,dc=skinny,dc=wsso",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, err
	}

	for _, entry := range sr.Entries {
		if entry.GetAttributeValue("cn") == group {
			return true, nil
		}
	}

	return false, nil
}
