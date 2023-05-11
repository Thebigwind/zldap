package client

import (
	"github.com/go-ldap/ldap/v3"
	. "zldap/common"
	. "zldap/manager"
)

type userManager interface {
	GetAllUsers() (error, map[string]UserEntry)
	GetUser(name string) (error, *ldap.SearchResult)
	AddUser(name, uid, gid, passwd, shell, home, shadowMax, shadowWarn string) (error, string)
	DeleteUser(name string) error
	ModifyUser(name, uid, gid, home, shell string) error
	Auth(name string, passwd string) error
	ChangePasswd(name string, old string, new string, force bool) error
}

type groupManager interface {
	GetAllGroups() (error, map[string]GroupEntry)
	GetGroup(name string) (error, *ldap.SearchResult)
	AddGroup(name string, gid string) (error, string)
	DeleteGroup(name string) error
	ModifyGroup(name string, newName string, gid string) error
	AddMember(name, add string) error
	DeleteMember(name, delete string) error
}

type client interface {
	userManager
	groupManager
}

type Client struct {
	UserManager
	GroupManager
}

func NewClient(servers []string) *Client {
	ldapdb := &LdapDB{Servers: servers, Conn: nil}
	return &Client{
		UserManager{LdapDB: ldapdb},
		GroupManager{LdapDB: ldapdb},
	}
}

func (c *Client) Close() {
	c.UserManager.Close()
	c.GroupManager.Close()
}
