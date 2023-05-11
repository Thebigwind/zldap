package manager

import (
	"testing"
	. "zldap/common"
)

var testUser = "unitestUser"
var testUser1 = "unitestUser1"
var testUser2 = "unitestUser2"

func TestUserAdd(t *testing.T) {
	t.Run(testUser, testUserAddFunc(testUser, "", "", "", "", "", "", ""))
	t.Run(testUser1, testUserAddFunc(testUser1, "", "", "", "", "", "", ""))
	t.Run(testUser2, testUserAddFunc(testUser2, "", "", "", "", "", "", ""))
}

func testUserAddFunc(name, uid, gid, passwd, shell, home, shadowMax, shadowWarn string) func(t *testing.T) {
	return func(t *testing.T) {
		err, _ := um.AddUser(name, uid, gid, passwd, shell, home, shadowMax, shadowWarn)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestAuth(t *testing.T) {
	um.getConnection(ADM_DN, ADM_PASS)

	err := um.Auth(testUser, "123456")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestChangePasswd(t *testing.T) {
	t.Run("123456", testChangePasswdFunc(testUser, "123456", "111111", false))
	t.Run("111111", testChangePasswdFunc(testUser, "111111", "123456", true))
}

func testChangePasswdFunc(name string, old string, new string, force bool) func(t *testing.T) {
	return func(t *testing.T) {
		err := um.ChangePasswd(name, old, new, force)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestGetUser(t *testing.T) {
	err, user := um.GetUser(testUser)
	if err != nil {
		t.Errorf(err.Error())
	}
	user.PrettyPrint(10)
}

func TestIsAssgined(t *testing.T) {
	t.Run("uid10001", testIsAssigned("uid", "10001", true))
	t.Run("gid10001", testIsAssigned("gid", "10001", true))
	t.Run("uid40000", testIsAssigned("uid", "40000", false))
	t.Run("gid50000", testIsAssigned("gid", "50000", false))
}

func testIsAssigned(subtype string, id string, expected bool) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		var actual bool
		if subtype == "uid" {
			err, actual = um.isAssigned(id)
		} else {
			err, actual = um.isAssigned(id)
		}

		if err != nil {
			t.Error(err.Error())
		}
		if actual != expected {
			t.Errorf("Expected the testIsAssigned %s of %s to be %t but instead got %t!", subtype, id, expected, actual)
		}
	}

}

func TestUserModify(t *testing.T) {
	t.Run("modify uid", testUserModFunc(testUser, "11111", "11111", "", "/bin/sh"))
}

func testUserModFunc(name, uid, gid, home, shell string) func(t *testing.T) {
	return func(t *testing.T) {
		err := um.ModifyUser(name, uid, gid, home, shell)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGetNextID(t *testing.T) {
	t.Run("user", testGetNextIDFunc("user"))
	t.Run("group", testGetNextIDFunc("group"))
	//t.Run("other", testGetNextIDFunc("other"))
}

func testGetNextIDFunc(subtree string) func(*testing.T) {
	return func(t *testing.T) {
		err, id := um.getNextID(subtree)

		if err != nil {
			t.Errorf(err.Error())
		} else {
			t.Logf("the resutl is ok, next %s id is %s ", subtree, id)
		}
	}
}

func TestGetAllUsers(t *testing.T) {
	err, userMap := um.GetAllUsers()
	if err != nil {
		t.Errorf(err.Error())
	}
	ShowUserList(userMap)
}

func TestUserDel(t *testing.T) {
	t.Run(testUser, testUserDelFunc(testUser))
	t.Run(testUser1, testUserDelFunc(testUser1))
	t.Run(testUser2, testUserDelFunc(testUser2))
}

func testUserDelFunc(name string) func(t *testing.T) {
	return func(t *testing.T) {
		err := um.DeleteUser(name)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
