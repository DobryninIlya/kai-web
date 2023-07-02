package main

//import "github.com/BurntSushi/toml"
import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	app "main/internal/app/vk_app"
	"path"
)

var (
	configPath     string
	configPathTest string
)

func init() {
	flag.StringVar(&configPath, "config-path", path.Join("configs", "vk_app.toml"), "path to config file")
	flag.StringVar(&configPathTest, "config-path-test", path.Join("..", "..", "configs", "vk_app.toml"), "path to config file")
}

func main() {
	flag.Parse()
	config := app.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Start(config); err != nil {
		log.Fatal(err)
	}
}
