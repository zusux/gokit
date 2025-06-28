package mongo

type ConfigMongo struct {
	Enable  bool              `yaml:"enable" json:"enable"`
	Mapping map[string]string `yaml:"mapping" json:"mapping"`
}
