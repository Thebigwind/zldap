package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"os"
	"zldap/client"
	"zldap/common"
)

var (
	servers = kingpin.Flag("server", "client server address").
		Default("10.10.10.125:389").Strings()
	//ldapaddr         = kingpin.Flag("addr", "ldap addr").Default("10.10.10.125").String()
	//ldapport         = kingpin.Flag("port", "ldap connect port").Default("389").Int()
	_ = kingpin.Command("userls", "list all users from ldap server.")
	_ = kingpin.Command("groupls", "list all groups from ldap server.")
)

func main() {
	subcmd := kingpin.Parse()
	//ldap := client.NewLdapDB(*servers)
	ldap := client.NewClient(*servers)

	switch subcmd {
	case "userls":
		err, userMap := ldap.GetAllUsers()
		if err != nil {
			fmt.Println("Get all users from ldap server fail.")
			fmt.Printf("  Reason: %s \n", err.Error())
			os.Exit(1)
		}
		common.ShowUserList(userMap)

	case "groupls":
		err, groupMap := ldap.GetAllGroups()
		if err != nil {
			fmt.Println("Get all groups from ldap server fail.")
			fmt.Printf("  Reason: %s \n", err.Error())
			os.Exit(1)
		}
		common.ShowGroupList(groupMap)
	}
}
