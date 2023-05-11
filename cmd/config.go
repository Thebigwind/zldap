package main

import (
	"os/user"
	"path"
)

func ClientConfFile() string {
	dir := ClientConfFile()
	if dir == "" {
		return ""
	}
	return path.Join(dir, "zldap.conf")
}

func ClientConfDir() string {
	user, err := user.Current()
	if err != nil {
		//Logger.Errorf("Fail to get current user: %s\n", err.Error())
		return ""
	}
	return user.HomeDir
}
