package conf

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

func MustLoad(path string, v any) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, v)
	if err != nil {
		log.Fatal(err)
	}
}
