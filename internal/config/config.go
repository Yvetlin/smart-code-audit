package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Target string `yaml:"target"`
}

type Config struct {
	Project ProjectConfig `yaml:"project"`
	LLM     LLMConfig     `yaml:"llm"`
	Policy  PolicyConfig  `yaml:"policy"`
}

type LLMConfig struct {
	Provider   string `yaml:"provider"`
	APIKeyEnv  string `yaml:"api_key_env"`
	BaseURLEnv string `yaml:"base_url_env"`
	Model      string `yaml:"model"`
}

type PolicyConfig struct {
	MinSeverity    string  `yaml:"min_severity"`
	MinProbability float64 `yaml:"min_probability"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
