package conf

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"log"
)

func MustLoad(path string, v any) {
	f := file.Provider(path)
	k := koanf.New(".")
	if err := k.Load(f, yaml.Parser()); err != nil {
		log.Fatalf("loading config err: %v", err.Error())
	}
}
