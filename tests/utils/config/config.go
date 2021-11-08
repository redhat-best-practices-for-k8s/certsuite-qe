package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const (
	// FileConfigPath path to config file
	FileConfigPath = "config/config.yaml"
)

// Config type keeps general configuration
type Config struct {
	General struct {
		ReportDirAbsPath string `yaml:"report" envconfig:"REPORT_DIR_NAME"`
		CnfNodeLabel     string `yaml:"cnf_worker_label" envconfig:"ROLE_WORKER_CNF"`
	} `yaml:"general"`
}

// NewConfig returns instance Config type
func NewConfig() (*Config, error) {
	var c Config
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filepath.Dir(filepath.Join(filepath.Dir(filename), "..")))
	confFile := filepath.Join(baseDir, FileConfigPath)
	err := readFile(&c, confFile)
	if err != nil {
		return nil, err
	}
	c.General.ReportDirAbsPath = filepath.Join(baseDir, c.General.ReportDirAbsPath)

	err = readEnv(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func readFile(c *Config, cfgFile string) error {
	f, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		return err
	}
	return nil
}

func readEnv(c *Config) error {
	err := envconfig.Process("", c)
	if err != nil {
		return err
	}
	return nil
}

// GetReportPath returns full path to the report file
func (c *Config) GetReportPath(file string) string {
	reportFileName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))
	return fmt.Sprintf("%s.xml", filepath.Join(c.General.ReportDirAbsPath, reportFileName))
}
