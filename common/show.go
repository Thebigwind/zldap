package common

import (
	"fmt"
	"github.com/crackcell/gotabulate"
	"github.com/xlab/treeprint"
)

type showEntry struct {
}

func ShowUserList(users map[string]UserEntry) {
	tabulator := gotabulate.NewTabulator()
	tabulator.SetFirstRowHeader(true)
	tabulator.SetFormat("grid")

	var table [][]string
	table = append(table, []string{"User", "UID", "GID", "Home", "Shell"})
	for u, e := range users {
		table = append(table, []string{u, e.Uid, e.Gid, e.Home, e.Shell})
	}

	fmt.Print(tabulator.Tabulate(table))
}

func ShowGroupList(groups map[string]GroupEntry) {
	tree := treeprint.New()
	for g, e := range groups {
		if len(e.Users) > 0 {
			group := tree.AddBranch(fmt.Sprintf("%s[%s]", g, e.Gid))
			for _, u := range e.Users {
				group.AddNode(u)
			}
		} else {
			tree.AddNode(fmt.Sprintf("%s[%s]", g, e.Gid))
		}
	}
	fmt.Print(tree.String())
}
