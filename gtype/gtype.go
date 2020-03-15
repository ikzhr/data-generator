package gtype

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Method string

const (
	RangeMethod Method = "range"
	EnumMethod  Method = "enum"
)

type Gtype struct {
	Name    string    `yaml:"name"`
	DBtype  string    `yaml:"dbtype"`
	Num     int       `yaml:"num"`
	Method  Method    `yaml:"method"`
	Range   Range     `yaml:"range"`
	TimeFmt string    `yaml:"timefmt"`
	Enum    []string  `yaml:"enum"`
	Weights []float32 `yaml:"weights"`
}

type Range struct {
	Min interface{} `yaml:"min"`
	Max interface{} `yaml:"max"`
}

type GtypeConfig struct {
	Gtypes []Gtype `yaml:"gtypes"`
}

func LoadConfig(path string) *GtypeConfig {
	var parsed GtypeConfig
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(buf, &parsed); err != nil {
		panic(err)
	}

	return &parsed
}

func (gconf *GtypeConfig) GenerateGtypesMap() map[string]Gtype {
	gtype_map := map[string]Gtype{}
	for _, gtype := range gconf.Gtypes {
		gtype_map[gtype.Name] = gtype
	}

	return gtype_map
}
