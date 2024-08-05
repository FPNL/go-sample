package conf

import (
	"github.com/fpnl/go-sample/pkg/tools"
)

type Bootstrap struct {
	Project *Project
	Server  *Server
	Data    *Data
	Log     *Log
}

type Log struct {
	OutPath    string
	AccessPath string
	PanicPath  string
	Stdout     bool
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
		Outsider []string
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
