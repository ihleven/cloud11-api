package auth

type Account struct {
	ID       uint32 `json:"-"`
	Username string `json:"username"`
	Password string `json:"-"`
	// Uid, Gid      uint32
	// Username      string
	// Name          string
	HomeDir string
	// Authenticated bool
	Groups []string `json:"groups"`
}

var Matthias Account = Account{
	1,
	"matt",
	"pwd",
	"/Users/mi/tmp",
	[]string{},
}
var Vicky Account = Account{2, "vicky", "pwd2", "/home/vicky", []string{}}
var Anonymous Account = Account{0, "anonymous", "none", "", []string{}}

var CurrentUser *Account = &Matthias
