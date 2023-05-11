package common

/*An Entry contains all the fields for a specific user*/
type UserEntry struct {
	Pass  string `json:"Pass"`
	Uid   string `json:"Uid"`
	Gid   string `json:"Gid"`
	Gecos string `json:"Gecos"`
	Home  string `json:"Home"`
	Shell string `json:"Shell"`
}

/*An Entry contains all the fields for a specific group*/
type GroupEntry struct {
	Pass  string   `json:"Pass"`
	Gid   string   `json:"Gid"`
	Users []string `json:"Users"`
}

type UserAttr struct {
	Name          []string
	ObjectClass   []string
	UidNumber     []string
	GidNumber     []string
	UserPassword  []string
	ShadowMax     []string
	ShadowWarning []string
	LoginShell    []string
	HomeDirectory []string
	Mail          []string
}

type GroupAttr struct {
	Name        []string
	ObjectClass []string
	GidNumber   []string
}
