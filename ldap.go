package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func RegisterUser(username string, password string, email string) (string, error) {

	l, err := ConnectToLdap()
	if err != nil {
		message := "Failed to bind with LDAP."
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

func DeleteLdapUser(username string) (string, error) {
	l, err := ConnectToLdap()
	if err != nil {
		message := "Failed to delete your account."
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

func GetLdapUsers() ([]string, error) {
	l, err := ConnectToLdap()
	if err != nil {
		return []string{}, err
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

func IsMemberOf(username string, group string) (bool, error) {
	groups, err := GetGroupMembership(username)
	if err != nil {
		return false, err
	}

	for _, g := range groups {
		if g == group {
			return true, nil
		}
	}

	return false, nil
}

func GetGroupMembership(username string) ([]string, error) {
	l, err := ConnectToLdap()
	if err != nil {
		return []string{}, err
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
		return []string{}, err
	}

	var groups []string
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}

	return groups, nil
}

func GetGroupMembers(group string) ([]string, error) {

	l, err := ConnectToLdap()
	if err != nil {
		return []string{}, err
	}

	searchFilter := fmt.Sprintf("(&(objectClass=groupOfUniqueNames)(cn=%s))", group)
	searchRequest := ldap.NewSearchRequest(
		"ou=groups,dc=skinny,dc=wsso",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"uniqueMember"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return []string{}, err
	}

	var members []string
	for _, entry := range sr.Entries {
		for _, member := range entry.GetAttributeValues("uniqueMember") {
			member = strings.Split(member, ",")[0]
			member = strings.Split(member, "=")[1]
			members = append(members, member)
		}
	}

	return members, nil

}

func AddUserToGroup(username string, group string) (string, error) {

	l, err := ConnectToLdap()
	if err != nil {
		message := "Failed to add user to group."
		return message, err
	}

	modifyRequest := ldap.NewModifyRequest("cn="+group+",ou=groups,dc=skinny,dc=wsso", nil)
	modifyRequest.Add("uniqueMember", []string{"uid=" + username + ",ou=users,dc=skinny,dc=wsso"})
	err = l.Modify(modifyRequest)

	if err != nil {
		message := fmt.Sprintf("Failed to add %s to %s!", username, group)
		return message, err
	}

	message := fmt.Sprintf("Successfully added %s to %s!", username, group)
	return message, nil

}

func RemoveUserFromGroup(username string, group string) (string, error) {

	l, err := ConnectToLdap()
	if err != nil {
		message := "Failed to remove user from group."
		return message, err
	}

	modifyRequest := ldap.NewModifyRequest("cn="+group+",ou=groups,dc=skinny,dc=wsso", nil)
	modifyRequest.Delete("uniqueMember", []string{"uid=" + username + ",ou=users,dc=skinny,dc=wsso"})
	err = l.Modify(modifyRequest)

	if err != nil {
		message := fmt.Sprintf("Failed to remove %s from %s!", username, group)
		return message, err
	}

	message := fmt.Sprintf("Successfully removed %s from %s!", username, group)
	return message, nil

}

func ConnectToLdap() (*ldap.Conn, error) {
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		return nil, err
	}

	err = l.Bind("cn=admin,dc=skinny,dc=wsso", os.Getenv("LDAP_ADMIN_PASSWORD"))

	return l, nil
}
