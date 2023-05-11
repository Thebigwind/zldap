package manager

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strconv"
	"strings"
	. "zldap/common"
)

type UserManager struct {
	*LdapDB
}

func NewUserManager(db *LdapDB) *UserManager {
	return &UserManager{
		LdapDB: db,
	}
}

func (mgr *UserManager) GetAllUsers() (error, map[string]UserEntry) {
	err, sr := mgr.search(BASE_DN, "(&(objectClass=posixAccount))", []string{})
	if err != nil {
		return err, nil
	}

	userMap := make(map[string]UserEntry)
	for _, entry := range sr.Entries {
		userEntry := UserEntry{
			Pass:  entry.GetAttributeValue("userPassword"),
			Uid:   entry.GetAttributeValue("uidNumber"),
			Gid:   entry.GetAttributeValue("gidNumber"),
			Gecos: "",
			Home:  entry.GetAttributeValue("homeDirectory"),
			Shell: entry.GetAttributeValue("loginShell"),
		}
		username := entry.GetAttributeValue("uid")
		userMap[username] = userEntry
	}

	return nil, userMap
}

func (mgr *UserManager) GetUser(username string) (error, *ldap.SearchResult) {
	if len(strings.TrimSpace(username)) == 0 {
		return fmt.Errorf("user name can not be empty when get user"), nil
	}

	fliter := fmt.Sprintf(userQueryString, username)
	err, sr := mgr.search(BASE_DN, fliter, []string{})
	if err != nil {
		return err, nil
	}

	return nil, sr
}

func (mgr *UserManager) AddUser(username, uid, gid, passwd, shell, home, shadowMax, shadowWarn string) (error, string) {
	var err error

	if len(strings.TrimSpace(username)) == 0 {
		return fmt.Errorf("user name can not be empty when add user"), ""
	}

	if len(strings.TrimSpace(uid)) == 0 {
		err, uid = mgr.getNextID("user")
		if err != nil {
			return err, ""
		}
	} else {
		err = mgr.verifyId(uid)
		if err != nil {
			return err, ""
		}
	}

	if len(strings.TrimSpace(gid)) == 0 {
		err, gid = mgr.getNextID("group")
		if err != nil {
			return err, ""
		}
	}

	if len(strings.TrimSpace(passwd)) == 0 {
		passwd = "123456"
	}

	if len(strings.TrimSpace(shell)) == 0 {
		shell = "/bin/bash"
	}

	if len(strings.TrimSpace(home)) == 0 {
		home = fmt.Sprintf("/home/%s", username)
	}

	if len(strings.TrimSpace(shadowMax)) == 0 {
		shadowMax = SHADOWMAX
	} else {
		_, err := strconv.Atoi(shadowMax)
		if err != nil {
			return fmt.Errorf("shadowMax can not convent to int"), ""
		}
	}

	if len(strings.TrimSpace(shadowWarn)) == 0 {
		shadowWarn = SHADOWWARNING
	} else {
		_, err := strconv.Atoi(shadowWarn)
		if err != nil {
			return fmt.Errorf("shadowMax can not convent to int"), ""
		}
	}

	attr := &UserAttr{
		Name:          []string{username},
		ObjectClass:   []string{"inetOrgPerson", "posixAccount", "top", "shadowAccount"},
		UidNumber:     []string{uid},
		GidNumber:     []string{gid},
		UserPassword:  []string{passwd},
		ShadowMax:     []string{shadowMax},
		ShadowWarning: []string{shadowWarn},
		LoginShell:    []string{shell},
		HomeDirectory: []string{home},
		Mail:          []string{fmt.Sprintf("%s@%s", username, lDAPMAILDOMAIN)},
	}

	// TODO: need optimize
	groupManager := NewGroupManager(mgr.LdapDB)
	groupManager.AddGroup(username, gid)
	return mgr.userAdd(attr), uid
}

func (mgr *UserManager) DeleteUser(username string) error {
	if len(strings.TrimSpace(username)) == 0 {
		return fmt.Errorf("user name can not be empty when delete user")
	}

	userdn := fmt.Sprintf("uid=%s,%s", username, PEOPLEDN)
	d := ldap.NewDelRequest(userdn, nil)

	// TODO: need optimize
	mgr.groupDelete(username)
	return mgr.delete(d)
}

func (mgr *UserManager) ModifyUser(username, uid, gid, home, shell string) error {
	if len(strings.TrimSpace(username)) == 0 {
		return fmt.Errorf("user name can not be empty when modify user")
	}

	if len(strings.TrimSpace(uid)) == 0 && len(strings.TrimSpace(gid)) == 0 &&
		len(strings.TrimSpace(home)) == 0 && len(strings.TrimSpace(shell)) == 0 {
		return fmt.Errorf("Parameters can not both be empty")
	}

	userdn := fmt.Sprintf("uid=%s,%s", username, PEOPLEDN)
	modify := ldap.NewModifyRequest(userdn, nil)

	if len(strings.TrimSpace(uid)) != 0 {
		if err := mgr.verifyId(uid); err != nil {
			return err
		}
		modify.Replace("uidNumber", []string{uid})
	}

	if len(strings.TrimSpace(gid)) != 0 {
		groupManager := NewGroupManager(mgr.LdapDB)
		if err := groupManager.verifyId(gid); err != nil {
			return err
		}
		modify.Replace("gidNumber", []string{gid})
	}

	if len(strings.TrimSpace(home)) != 0 {
		modify.Replace("homeDirectory", []string{home})
	}

	if len(strings.TrimSpace(shell)) != 0 {
		modify.Replace("loginShell", []string{shell})
	}

	return mgr.modify(modify)
}

func (mgr *UserManager) Auth(username string, passwd string) error {
	userdn := fmt.Sprintf("uid=%s,%s", username, PEOPLEDN)

	err, conn := mgr.createConnection()
	defer conn.Close()
	if err != nil {
		return err
	}

	err, _ = mgr.bindConnection(conn, userdn, passwd)

	return err
}

func (mgr *UserManager) ChangePasswd(username string, old string, new string, force bool) error {
	if !force {
		err := mgr.Auth(username, old)
		if err != nil {
			return fmt.Errorf("old password error")
		}
	}

	return mgr.changePasswd(username, old, new)
}

func (mgr *UserManager) groupDelete(groupname string) error {
	if len(strings.TrimSpace(groupname)) == 0 {
		return fmt.Errorf("group name can not be empty when delete group")
	}

	groupdn := fmt.Sprintf("cn=%s,%s", groupname, GROUPDN)
	d := ldap.NewDelRequest(groupdn, nil)

	return mgr.delete(d)
}

func (mgr *UserManager) isAssigned(id string) (error, bool) {
	err, allUser := mgr.GetAllUsers()
	if err != nil {
		return err, true
	}
	for _, entry := range allUser {
		if entry.Uid == id {
			return nil, true
		}
	}

	return nil, false
}

func (mgr *UserManager) verifyId(id string) error {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("uid can not convent to int")
	}

	if intId < 10000 || intId > 60000 {
		return fmt.Errorf("uid must between 10000 and 60000")
	}

	err, assigned := mgr.isAssigned(id)
	if err != nil {
		return err
	}
	if assigned {
		return fmt.Errorf("uid %s already be assigned, can not assign again", id)
	}

	return nil
}
