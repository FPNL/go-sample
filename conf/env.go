package conf

import (
	"oltp/pkg/tools"
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
		Timeout tools.Duration
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
		Pwd    string
		DbName string
		Conn   int
	}
}
