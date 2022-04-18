package main

type Config struct {
	AppId string
	Comet string
	Logic string
	Users []User
}

type User struct {
	Token  string
	Uid    string
	Groups []string
	Text   string
}
