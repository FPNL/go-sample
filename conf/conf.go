package conf

import (
	"errors"
	"fmt"

	"github.com/fpnl/go-sample/pkg/tools"
)

const debug = "debug"

// debugMode -ldflags="-X 'full_package_path.variable=value'"
var debugMode = ""

func InitAPI(confPath string) (*Bootstrap, error) {
	var c = &Bootstrap{}

	err := loadBootstrapConfig(confPath, c)
	if err != nil {
		return nil, fmt.Errorf("load Bootstrap Config fail: %w", err)
	}

	if debugMode == debug {
		c.Project.IsDebug = true
	}

	if err = validate(c); err != nil {
		return nil, fmt.Errorf("validate config fail: %w", err)
	}

	return c, nil
}

// LoadBootstrapConfig 由於 devops 給設定檔的方式為，把設定內容 encode to json 然後設定到環境變數 envConfig
// 但是在本地開發還是以 env.json 為主，所以這邊要做一個判斷，如果有 envConfig 就用 envConfig，沒有就用 env.json
func loadBootstrapConfig(confPath string, c any) (err error) {
	if err = tools.LoadFromEnv("envConfig", "json", c); err == nil {
		return nil
	}

	if err = tools.LoadFromFile(confPath, c); err != nil {
		return fmt.Errorf("load config from file fail: %w", err)
	}

	return nil
}

func validate(c *Bootstrap) error {
	if c.Server.HTTP.Addr == "" {
		return errors.New("config http's addr is empty")
	}

	return nil
}
