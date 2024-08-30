package config

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const envFilePath = "env.test.json"

type Bootstrap struct {
	Project *struct {
		Name string `json:"name" yaml:"name"`
	} `json:"project" yaml:"project"`
	Mysql *struct {
		Host        string   `json:"host" yaml:"host"`
		Port        int      `json:"port" yaml:"port"`
		Timeout     Duration `json:"timeout" yaml:"timeout"`
		IpWhitelist []string `json:"ip_whitelist" yaml:"ip_whitelist"`
	} `json:"mysql" yaml:"mysql"`
}

func TestLoadFromEnvironment(t *testing.T) {
	err := os.Setenv("envConfig", `---
mysql:
  host: localhost
  port: 3306
  timeout: 5s
  ip_whitelist:
  - '1'
  - '2'
  - '3'
  - '4'
`)
	require.NoError(t, err)

	var c = &Bootstrap{}

	err = LoadFromEnv("envConfig", "yaml", c)
	require.NoError(t, err)

	require.NotNil(t, c.Mysql)

	assert.Equal(t, c.Mysql.Host, "localhost")
	assert.Equal(t, c.Mysql.Port, 3306)
	assert.Equal(t, c.Mysql.Timeout.Duration, time.Second*5)
	assert.Len(t, c.Mysql.IpWhitelist, 4)
	assert.ElementsMatch(t, c.Mysql.IpWhitelist, []string{"1", "2", "3", "4"})
}

func TestLoadFromFile(t *testing.T) {
	var c = &Bootstrap{}

	err := LoadFromFile(envFilePath, c)
	require.NoError(t, err)

	require.NotNil(t, c.Mysql)

	assert.Equal(t, c.Mysql.Host, "localhost")
	assert.Equal(t, c.Mysql.Port, 3306)
	assert.Equal(t, c.Mysql.Timeout.Duration, time.Second*5)
	assert.Len(t, c.Mysql.IpWhitelist, 4)
	assert.ElementsMatch(t, c.Mysql.IpWhitelist, []string{"1", "2", "3", "goconfig"})
}

func TestWatchFile(t *testing.T) {
	var c = &Bootstrap{}
	err := LoadFromFile(envFilePath, c, WithWatcher(
		"mysql.host",
		func(key string, value config.Value) {
			s, err := value.String()
			require.NoError(t, err)

			c.Mysql.Host = s
		}, func(err error) {
			require.NoError(t, err)
		}))
	require.NoError(t, err)

	file, err := os.OpenFile(envFilePath, os.O_RDWR, 0666)
	require.NoError(t, err)
	defer file.Close()

	var offset int64 = 0
	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096) // Adjust buffer size as needed

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			require.NoError(t, err)
		}
		if n == 0 {
			break
		}

		index := bytes.Index(buffer[:n], []byte("localhost"))
		if index != -1 {
			// Calculate the position for WriteAt
			offset = offset + int64(index)
			_, err = file.WriteAt([]byte("hostlocal"), offset)
			require.NoError(t, err)
			break
		}

		offset += int64(n)
	}
	time.Sleep(time.Second * 2)
	require.NotNil(t, c.Mysql)

	assert.Equal(t, c.Mysql.Host, "hostlocal")
	assert.Equal(t, c.Mysql.Port, 3306)
	assert.Equal(t, c.Mysql.Timeout.Duration, time.Second*5)
	assert.Len(t, c.Mysql.IpWhitelist, 4)
	assert.ElementsMatch(t, c.Mysql.IpWhitelist, []string{"1", "2", "3", "goconfig"})

	_, err = file.WriteAt([]byte("localhost"), offset)
	require.NoError(t, err)
}
