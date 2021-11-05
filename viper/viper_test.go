package main_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestViper(t *testing.T) {
	viper.SetConfigName("config")         // name of config file, Does not include extension
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	assert.NoError(t, err)

	assert.Equal(t, "TOML Example", viper.GetString("title"))
	assert.Equal(t, "192.168.1.1", viper.GetString("database.server"))
	assert.Equal(t, []int{8000, 8001, 8002}, viper.GetIntSlice("database.ports"))

	assert.Equal(t, "10.0.0.1", viper.GetString("servers.alpha.ip"))
	assert.Equal(t, []interface{}{[]interface{}{"gamma", "delta"}, []interface{}{int64(1), int64(2)}}, viper.Get("clients.data"))
}

func TestYAML(t *testing.T) {
	viper.SetConfigName("docker-compose")
	viper.AddConfigPath("../etcd/")

	err := viper.ReadInConfig()
	assert.NoError(t, err)

	assert.Equal(t, []string{"12379:2379", "12380:2380"}, viper.GetStringSlice("services.etcd1.ports"))
}
