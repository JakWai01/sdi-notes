package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

func main() {
	uid := flag.String("uid", "", "User-Id e.g. jw163")
	password := flag.String("pw", "", "Password e.g. password")
	flag.Parse()

	ldapURL := "ldaps://ldap1.mi.hdm-stuttgart.de"
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	err = l.Bind(fmt.Sprintf("uid=%s, ou=userlist,dc=hdm-stuttgart,dc=de", *uid), *password)
	if err != nil {
		log.Fatal(err)
	}

	// Until here, everything should work
	BaseDN := "dc=hdm-stuttgart,dc=de"
	Filter := fmt.Sprintf("(uid=%s)", *uid)
	searchReq := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		Filter,
		[]string{},
		nil,
	)
	result, err := l.Search(searchReq)
	if err != nil {
		log.Fatal(err)
	}

	if len(result.Entries) > 0 {
		result.PrettyPrint(4)
	} else {
		fmt.Println("Couldn't fetch search entries")
	}
}
