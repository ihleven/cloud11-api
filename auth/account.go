package auth

type Account struct {
	ID       uint32
	Username string
	Password string
	// Uid, Gid      uint32
	// Username      string
	// Name          string
	// HomeDir       string
	// Authenticated bool
}

var Matthias Account = Account{1, "matt", "pwd"}
var Vicky Account = Account{2, "vicky", "pwd2"}

var CurrentUser *Account = &Matthias
