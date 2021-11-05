package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	v.OnConfigChange(func(in fsnotify.Event) {
		// OUTPUT:
		// "~learn-by-test/viper/config.toml": WRITE
		fmt.Println(in.String())
	})

	v.WatchConfig()

	name := v.GetString("owner.name")

	fmt.Println("origin:", name)

	// modify config.toml by hand
	for name == "Tom Preston-Werner" {
		name = v.GetString("owner.name")
	}

	fmt.Println("modified:", name)
}
