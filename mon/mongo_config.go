package mon

type ConfigMongo struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Uri    string `yaml:"uri" json:"uri"`
}
