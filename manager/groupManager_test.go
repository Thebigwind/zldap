package manager

import (
	"fmt"
	"os"
	"testing"
	. "zldap/common"
)

var gm *GroupManager
var um *UserManager
var server = []string{"10.46.221.141"}
var testGroup = "unitestGroup"
var groupadd = "groupadd"
var groupadd1 = "groupadd1"
var groupadd2 = "groupadd2"

func TestMain(m *testing.M) {
	db := NewLdapDB(server)
	gm = NewGroupManager(db)
	um = NewUserManager(db)
	exitCode := m.Run()
	gm.Close()
	os.Exit(exitCode)
}

func TestGroupADD(t *testing.T) {
	err, _ := gm.AddGroup(testGroup, "")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGroupModify(t *testing.T) {
	err := gm.ModifyGroup(testGroup, "unitestGroup", "13140")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGroupAddUser(t *testing.T) {
	t.Run("add testUser", testGroupAddUserFunc(testGroup, groupadd))
	t.Run("add testUser1", testGroupAddUserFunc(testGroup, groupadd1))
	t.Run("add testUser2", testGroupAddUserFunc(testGroup, groupadd2))
}

func testGroupAddUserFunc(name, add string) func(t *testing.T) {
	return func(t *testing.T) {
		err := gm.AddMember(name, add)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGetGroupMems(t *testing.T) {
	err, mems := gm.getGroupMems(testGroup)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Print("          ")
	fmt.Println(mems)
}

func TestGetGroup(t *testing.T) {
	err, group := gm.GetGroup(testGroup)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(group.Entries) != 1 {
		t.Errorf("get group %s failure", testGroup)
	}
	group.PrettyPrint(10)
}

func TestGroupDelUser(t *testing.T) {
	t.Run("delete groupadd", testGroupDelUserFunc(testGroup, groupadd))
	t.Run("delete groupadd1", testGroupDelUserFunc(testGroup, groupadd1))
	t.Run("delete groupadd2", testGroupDelUserFunc(testGroup, groupadd2))
}

func testGroupDelUserFunc(groupname, username string) func(t *testing.T) {
	return func(t *testing.T) {
		err := gm.DeleteMember(groupname, username)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGetAllGroups(t *testing.T) {
	err, groupMap := gm.GetAllGroups()
	if err != nil {
		t.Errorf(err.Error())
	}
	ShowGroupList(groupMap)
}

func TestGroupDel(t *testing.T) {
	err := gm.DeleteGroup(testGroup)
	if err != nil {
		t.Error(err.Error())
	}
}
