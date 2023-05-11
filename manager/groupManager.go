package manager

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strconv"
	"strings"
	. "zldap/common"
)

type GroupManager struct {
	*LdapDB
}

func NewGroupManager(db *LdapDB) *GroupManager {
	return &GroupManager{
		LdapDB: db,
	}
}

func (mgr *GroupManager) GetAllGroups() (error, map[string]GroupEntry) {
	err, sr := mgr.search(BASE_DN, "(&(objectClass=posixGroup))", []string{})
	if err != nil {
		return err, nil
	}

	groupMap := make(map[string]GroupEntry)
	for _, entry := range sr.Entries {
		greoupEntry := GroupEntry{
			Pass:  "",
			Gid:   entry.GetAttributeValue("gidNumber"),
			Users: entry.GetAttributeValues("memberUid"),
		}
		groupName := entry.GetAttributeValue("cn")
		groupMap[groupName] = greoupEntry
	}
	return nil, groupMap
}

func (mgr *GroupManager) GetGroup(groupname string) (error, *ldap.SearchResult) {
	if len(strings.TrimSpace(groupname)) == 0 {
		return fmt.Errorf("group name can not be empty when get group"), nil
	}

	fliter := fmt.Sprintf(groupQueryString, groupname)
	err, sr := mgr.search(BASE_DN, fliter, []string{})
	if err != nil {
		return err, nil
	}

	return nil, sr
}

func (mgr *GroupManager) AddGroup(groupname string, gid string) (error, string) {
	var err error
	if len(strings.TrimSpace(groupname)) == 0 {
		return fmt.Errorf("group name can not be empty when add group"), ""
	}

	if len(strings.TrimSpace(gid)) == 0 {
		err, gid = mgr.getNextID("group")
		if err != nil {
			return err, ""
		}
	} else {
		if err := mgr.verifyId(gid); err != nil {
			return err, ""
		}
	}

	attr := &GroupAttr{
		Name:        []string{groupname},
		ObjectClass: []string{"posixGroup", "top"},
		GidNumber:   []string{gid},
	}

	return mgr.groupAdd(attr), gid
}

func (mgr *GroupManager) DeleteGroup(groupname string) error {
	if len(strings.TrimSpace(groupname)) == 0 {
		return fmt.Errorf("group name can not be empty when delete group")
	}

	err, groupMems := mgr.getGroupMems(groupname)
	if err != nil {
		return err
	}

	if len(groupMems) > 0 {
		return fmt.Errorf("group %s not removed because it has other members.", groupname)
	}

	groupdn := fmt.Sprintf("cn=%s,%s", groupname, GROUPDN)
	d := ldap.NewDelRequest(groupdn, nil)

	return mgr.delete(d)
}

func (mgr *GroupManager) ModifyGroup(groupname string, newName string, gid string) error {
	if len(strings.TrimSpace(groupname)) == 0 {
		return fmt.Errorf("group name can not be empty when modify group")
	}

	if len(strings.TrimSpace(newName)) == 0 && len(strings.TrimSpace(gid)) == 0 {
		return fmt.Errorf("Parameters can not both be empty")
	}

	groupdn := fmt.Sprintf("cn=%s,%s", groupname, GROUPDN)
	modify := ldap.NewModifyRequest(groupdn, nil)
	// TODO: implement modify dn
	if len(strings.TrimSpace(newName)) != 0 {
	}

	if len(strings.TrimSpace(gid)) != 0 {
		if err := mgr.verifyId(gid); err != nil {
			return err
		}

		modify.Replace("gidNumber", []string{gid})
	}

	return mgr.modify(modify)
}

func (mgr *GroupManager) AddMember(groupname, username string) error {
	groupdn := fmt.Sprintf("cn=%s,%s", groupname, GROUPDN)
	modify := ldap.NewModifyRequest(groupdn, nil)

	if len(strings.TrimSpace(username)) != 0 {
		modify.Add("memberUid", []string{username})
	}

	return mgr.modify(modify)
}

func (mgr *GroupManager) DeleteMember(groupname, username string) error {
	groupdn := fmt.Sprintf("cn=%s,%s", groupname, GROUPDN)
	modify := ldap.NewModifyRequest(groupdn, nil)

	if len(strings.TrimSpace(username)) != 0 {
		modify.Delete("memberUid", []string{username})
	}

	return mgr.modify(modify)
}

func (mgr *GroupManager) getGroupMems(groupname string) (error, []string) {
	err, groupInfo := mgr.GetGroup(groupname)
	if err != nil {
		return err, nil
	}
	if len(groupInfo.Entries) == 0 {
		return nil, nil
	}

	return nil, groupInfo.Entries[0].GetAttributeValues("memberUid")
}

func (mgr *GroupManager) isAssigned(id string) (error, bool) {
	err, allGroups := mgr.GetAllGroups()
	if err != nil {
		return err, true
	}
	for _, entry := range allGroups {
		if entry.Gid == id {
			return nil, true
		}
	}

	return nil, false
}

func (mgr *GroupManager) verifyId(id string) error {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("gid can not convent to int")
	}

	if intId < 10000 || intId > 60000 {
		return fmt.Errorf("gid must between 10000 and 60000")
	}

	err, assigned := mgr.isAssigned(id)
	if err != nil {
		return err
	}
	if assigned {
		return fmt.Errorf("gid %s already be assigned, can not assign again", id)
	}

	return nil
}
