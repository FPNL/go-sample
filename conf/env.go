package conf

import (
	"time"
)

type Bootstrap struct {
	Project *Project
	Server  *Server
	Data    *Data
}

type Project struct {
	Name          string
	Env           string
	IsDebug       bool
	DecryptorPath string
}

type Server struct {
	HTTP struct {
		Network string
		Addr    string
		Timeout time.Duration
	}
	IPWhitelist struct {
		Internal []string
		Uu       []string
	}
}

type Data struct {
	Mysql struct {
		URL    string
		Port   string
		User   string
		DbName string
		Pwd    string
		Conn   int
	}
}
