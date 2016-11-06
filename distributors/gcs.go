package distributors

type GCS struct {
	AuthKey      string              `yaml:"auth_key"`
	Source       map[string][]string `yaml:"source"`
	Target       string              `yaml:"target"`
	Ignore       string              `yaml:"ignore"`
	ACL          []string            `yaml:"acl"`
	CacheControl string              `yaml:"cache_control"`
	Metadata     map[string]string   `yaml:"metadata"`
}

func (g *GCS) Valid() error {
	return nil
}
