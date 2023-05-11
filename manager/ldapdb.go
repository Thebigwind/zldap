package manager

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strconv"
	. "zldap/common"
)

const (
	BASE_DN        string = "dc=zdlz,dc=com"
	ADM_DN         string = "cn=admin,dc=zdlz,dc=com"
	ADM_PASS       string = "123456"
	PEOPLEDN       string = "ou=People,dc=zdlz,dc=com"
	GROUPDN        string = "ou=Group,dc=zdlz,dc=com"
	lDAPMAILDOMAIN string = "zdlz.com"
)

var (
	userQueryString  = "(&(uid=%s)(objectClass=posixAccount))"
	groupQueryString = "(&(cn=%s)(objectClass=posixGroup))"
	SHADOWMAX        = "99999"
	SHADOWWARNING    = "14"
)

type LdapDB struct {
	Servers []string
	Conn    *ldap.Conn
}

func NewLdapDB(ldapserver []string) *LdapDB {
	db := &LdapDB{
		Servers: ldapserver,
		Conn:    nil,
	}

	return db
}

func (db *LdapDB) createConnection() (error, *ldap.Conn) {
	var err error
	for _, server := range db.Servers {
		conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", server, 389))
		if err == nil {
			return nil, conn
		}
	}
	return fmt.Errorf("Fail to dial ldap server, %s\n", err.Error()), nil
}

func (db *LdapDB) bindConnection(conn *ldap.Conn, userDn string, passwd string) (error, *ldap.Conn) {
	err := conn.Bind(userDn, passwd)
	if err != nil {
		return fmt.Errorf("Fail to bind to ldap server, %s\n", err.Error()), conn
	}
	return nil, conn
}

func (db *LdapDB) getConnection(userdn string, passwd string) error {
	if db.Conn == nil {
		err, conn := db.createConnection()
		if err != nil {
			return err
		}

		err, db.Conn = db.bindConnection(conn, userdn, passwd)
		if err != nil {
			db.Close()
			return err
		}
	}

	return nil
}

func (db *LdapDB) search(dn string, fliter string, attr []string) (error, *ldap.SearchResult) {
	err := db.getConnection(ADM_DN, ADM_PASS)
	if err != nil {
		return err, nil
	}

	searchRequest := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fliter,
		attr,
		nil,
	)

	sr, err := db.Conn.Search(searchRequest)
	if err != nil {
		db.Close()
		return fmt.Errorf(err.Error()), nil
	}

	return nil, sr
}

func (db *LdapDB) add(addRequest *ldap.AddRequest) error {
	err := db.getConnection(ADM_DN, ADM_PASS)
	if err != nil {
		return err
	}

	return db.Conn.Add(addRequest)
}

func (db *LdapDB) userAdd(attr *UserAttr) error {
	name := attr.Name[0]
	a := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", name, PEOPLEDN), nil)
	a.Attribute("cn", attr.Name)
	a.Attribute("sn", attr.Name)
	a.Attribute("objectClass", attr.ObjectClass)
	a.Attribute("shadowMax", attr.ShadowMax)
	a.Attribute("shadowWarning", attr.ShadowWarning)
	a.Attribute("loginShell", attr.LoginShell)
	a.Attribute("uidNumber", attr.UidNumber)
	a.Attribute("gidNumber", attr.GidNumber)
	a.Attribute("userPassword", attr.UserPassword)
	a.Attribute("homeDirectory", attr.HomeDirectory)
	a.Attribute("mail", attr.Mail)

	return db.add(a)
}

func (db *LdapDB) groupAdd(attr *GroupAttr) error {
	name := attr.Name[0]
	a := ldap.NewAddRequest(fmt.Sprintf("cn=%s,%s", name, GROUPDN), nil)
	a.Attribute("cn", attr.Name)
	a.Attribute("objectClass", attr.ObjectClass)
	a.Attribute("gidNumber", attr.GidNumber)

	return db.add(a)
}

func (db *LdapDB) delete(delRequest *ldap.DelRequest) error {
	err := db.getConnection(ADM_DN, ADM_PASS)
	if err != nil {
		return err
	}

	return db.Conn.Del(delRequest)
}

func (db *LdapDB) modify(modifyRequest *ldap.ModifyRequest) error {
	err := db.getConnection(ADM_DN, ADM_PASS)
	if err != nil {
		return err
	}

	return db.Conn.Modify(modifyRequest)
}

func (db *LdapDB) changePasswd(username, old, new string) error {
	err := db.getConnection(ADM_DN, ADM_PASS)
	if err != nil {
		return err
	}

	userdn := fmt.Sprintf("uid=%s,%s", username, PEOPLEDN)
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userdn, old, new)
	_, err = db.Conn.PasswordModify(passwordModifyRequest)

	return err
}

func (db *LdapDB) getNextID(subtree string) (error, string) {
	big := 10000

	switch subtree {
	case "user":
		sfilter := "(&(uidNumber=*)(objectClass=posixAccount))"
		err, res := db.search(BASE_DN, sfilter, []string{})
		if err != nil {
			return err, ""
		}

		for _, entry := range res.Entries {
			suid := entry.GetAttributeValue("uidNumber")
			uid, err := strconv.Atoi(suid)
			if err != nil {
				return err, ""
			}

			if uid > big && uid < 60000 {
				big = uid
			}
		}
		nextId := big + 1
		return nil, strconv.Itoa(nextId)

	case "group":
		sfilter := "(&(gidNumber=*)(objectClass=posixGroup))"
		err, res := db.search(BASE_DN, sfilter, []string{})
		if err != nil {
			return err, ""
		}

		for _, entry := range res.Entries {
			sgid := entry.GetAttributeValue("gidNumber")
			gid, err := strconv.Atoi(sgid)
			if err != nil {
				return err, ""
			}

			if gid > big && gid < 60000 {
				big = gid
			}
		}
		nextId := big + 1
		return nil, strconv.Itoa(nextId)

	default:
		return fmt.Errorf("unknow subtree type, must be 'user' or 'group'"), ""
	}
}

func (db *LdapDB) Close() {
	if db.Conn != nil {
		db.Conn.Close()
		db.Conn = nil
	}
}
