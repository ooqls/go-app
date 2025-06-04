package app

import (
	"os"

	"gopkg.in/yaml.v2"
)

type DocsConfig struct {
	Enabled     bool   `yaml:"enabled"`
	DocsApiPath string `yaml:"docs_api_path"`
	DocsDir     string `yaml:"docs_dir"`
}

type ServerConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CaPath   string `yaml:"ca_path"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type JWTConfig struct {
	Enabled                bool   `yaml:"enabled"`
	RSAKeyPath             string `yaml:"rsa_key_path"`
	RSAPubKeyPath          string `yaml:"rsa_pub_key_path"`
	TokenConfigurationPaths []string `yaml:"token_configuration_paths"`
}

type SQLFilesConfig struct {
	Enabled          bool     `yaml:"enabled"`
	SQLFilesDir      string   `yaml:"sql_files_dir"`
	SQLFiles         []string `yaml:"sql_files"`
	CreateTableStmts []string `yaml:"create_table_stmts"`
	CreateIndexStmts []string `yaml:"create_index_stmts"`
}

type RegistryConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type LoggingAPIConfig struct {
	Enabled bool `yaml:"enabled"`
}

type AppConfig struct {
	LoggingAPI   LoggingAPIConfig `yaml:"logging_api"`
	DocsConfig   DocsConfig       `yaml:"docs"`
	ServerConfig ServerConfig     `yaml:"server"`
	TLS          TLSConfig        `yaml:"tls"`
	JWT          JWTConfig        `yaml:"jwt"`
	SQLFiles     SQLFilesConfig   `yaml:"sql"`
	Registry     RegistryConfig   `yaml:"registry"`
}

func LoadConfig(path string) (*AppConfig, error) {
	cfg := &AppConfig{}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
