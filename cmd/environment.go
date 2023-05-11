package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type envCommand struct {
}

func newEnvCommand() *envCommand {
	return &envCommand{}
}

func (cmd *envCommand) Load() (error, map[string]interface{}) {
	m := make(map[string]interface{})
	_, err := os.Stat(ClientConfFile())
	if os.IsNotExist(err) {
		content := []byte("{}")
		err = ioutil.WriteFile(ClientConfFile(), content, 644)
		return err, m
	}

	o, err := ioutil.ReadFile(ClientConfFile())
	if err != nil {
		return err, m
	}

	err = json.Unmarshal(o, &m)
	return err, m
}

func (cmd *envCommand) Set(key string, val string) {
	tmpfile, err := ioutil.TempFile(ClientConfDir(), "zldap_env")
	if err != nil {
		fmt.Println("Fail to create temp env file", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	err, m := cmd.Load()
	if err != nil {
		fmt.Println("Load config file error!", err)
		return
	}
	m[key] = val
	res, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Fail to encode user config!", err)
		return
	}

	if _, err := tmpfile.Write(res); err != nil {
		fmt.Println("Write temp config file error!", err)
		return
	}

	err = os.Rename(tmpfile.Name(), ClientConfFile())
	if err != nil {
		fmt.Println("Fail to rename file %s \n", err.Error())
	}

	tmpfile.Close()
}

func (cmd *envCommand) Get(key string) {
	err, m := cmd.Load()
	if err != nil {
		fmt.Println("Load config file error!", err)
		return
	}

	if key != "all" {
		if _, ok := m[key]; ok {
			fmt.Println(m[key].(string))
		} else {
			fmt.Printf("Unknown env key %s . \n", key)
		}
	} else {
		res, _ := json.Marshal(m)
		fmt.Println(string(res))
	}
}
