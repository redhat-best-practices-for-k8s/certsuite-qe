package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang/glog"
	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const (
	// FileConfigPath path to config file.
	FileConfigPath = "config/config.yaml"
)

// Config type keeps general configuration.
type Config struct {
	General struct {
		ReportDirAbsPath     string `yaml:"report" envconfig:"REPORT_DIR_NAME"`
		CnfNodeLabel         string `yaml:"cnf_worker_label" envconfig:"ROLE_WORKER_CNF"`
		WorkerNodeLabel      string `yaml:"worker_label" envconfig:"ROLE_WORKER"`
		TestImage            string `yaml:"test_image" envconfig:"TEST_IMAGE"`
		VerificationLogLevel string `yaml:"verification_log_level" envconfig:"VERIFICATION_LOG_LEVEL"`
		DebugTnf             string `envconfig:"DEBUG_TNF"`
		TnfConfigDir         string `yaml:"tnf_config_dir" envconfig:"TNF_CONFIG_DIR"`
		TnfRepoPath          string `envconfig:"TNF_REPO_PATH"`
		TnfEntryPointScript  string `yaml:"tnf_entry_point_script" envconfig:"TNF_ENTRY_POINT_SCRIPT"`
		TnfReportDir         string `yaml:"tnf_report_dir" envconfig:"TNF_REPORT_DIR"`
		TnfImage             string `yaml:"tnf_image" envconfig:"TNF_IMAGE"`
		TnfImageTag          string `yaml:"tnf_image_tag" envconfig:"TNF_IMAGE_TAG"`
	} `yaml:"general"`
}

// NewConfig returns instance Config type.
func NewConfig() (*Config, error) {
	var conf Config

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filepath.Dir(filepath.Join(filepath.Dir(filename), "..")))
	confFile, err := checkFileExists(baseDir, FileConfigPath)

	if err != nil {
		glog.Fatal(err)
	}

	err = readFile(&conf, confFile)

	if err != nil {
		return nil, err
	}

	conf.General.ReportDirAbsPath = filepath.Join(baseDir, conf.General.ReportDirAbsPath)

	err = readEnv(&conf)
	if err != nil {
		return nil, err
	}

	err = conf.deployTnfConfigDir(confFile)

	if err != nil {
		return nil, err
	}

	err = conf.deployTnfReportDir(confFile)
	if err != nil {
		return nil, err
	}

	conf.General.TnfRepoPath, err = conf.defineTnfRepoPath()

	if err != nil {
		glog.Fatal(err)
	}

	return &conf, nil
}

func (c *Config) DebugTnf() (bool, error) {
	if c.General.DebugTnf == "true" {
		err := os.Setenv("TNF_LOG_LEVEL", "trace")

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (c *Config) CreateLogFile(testSuite string, tcName string) *os.File {
	folderPath := filepath.Join(c.General.ReportDirAbsPath, "Debug", testSuite)

	err := os.MkdirAll(folderPath, 0755)
	if err != nil && !os.IsExist(err) {
		// we only panic in case the error is different than "folder already exists".
		panic(err)
	}

	tcFile := filepath.Join(folderPath, tcName+".log")

	// if the log file already exists, remove it and create a new one.
	if _, err := os.Stat(tcFile); err == nil {
		err = os.Remove(tcFile)
		if err != nil {
			panic(err)
		}
	}

	outfile, err := os.OpenFile(tcFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	return outfile
}

func readFile(cfg *Config, cfgFile string) error {
	openedCnfFile, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer openedCnfFile.Close()

	decoder := yaml.NewDecoder(openedCnfFile)
	err = decoder.Decode(&cfg)

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

func (c *Config) deployTnfConfigDir(configFileName string) error {
	return deployTnfDir(configFileName, c.General.TnfConfigDir, "tnf_config_dir", "TNF_CONFIG_DIR")
}

func (c *Config) deployTnfReportDir(configFileName string) error {
	return deployTnfDir(configFileName, c.General.TnfReportDir, "tnf_report_dir", "TNF_REPORT_DIR")
}

// GetReportPath returns full path to the report file.
func (c *Config) GetReportPath(file string) string {
	reportFileName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))

	return fmt.Sprintf("%s.xml", filepath.Join(c.General.ReportDirAbsPath, reportFileName))
}

func (c *Config) defineTnfRepoPath() (string, error) {
	if c.General.TnfRepoPath == "" {
		return "", fmt.Errorf("TNF_REPO_PATH env variable is not set. Please export TNF_REPO_PATH")
	}

	_, err := checkFileExists(c.General.TnfRepoPath, c.General.TnfEntryPointScript)

	if err != nil {
		return "", err
	}

	return c.General.TnfRepoPath, nil
}

// DefineClients sets client and return it's instance.
func DefineClients() (*testclient.ClientSet, error) {
	clients := testclient.New("")
	if clients == nil {
		return nil, fmt.Errorf("client is not set please check KUBECONFIG env variable")
	}

	return clients, nil
}

func checkFileExists(filePath, name string) (string, error) {
	if !filepath.IsAbs(filePath) {
		return "", fmt.Errorf(
			"make sure env var TNF_REPO_PATH is configured with absolute path instead of relative",
		)
	}

	fullPath, _ := filepath.Abs(filepath.Join(filePath, name))
	_, err := os.Stat(fullPath)

	if err == nil {
		return fullPath, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("path to %s file not found: %s , Exiting", name, fullPath)
	}

	return "", fmt.Errorf("path to %s file not valid: %s , err=%w, exiting", name, fullPath, err)
}

func deployTnfDir(confFileName string, dirName string, yamlTag string, envVar string) error {
	_, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		glog.V(4).Info(fmt.Sprintf("%s directory is not present. Creating directory", dirName))

		return os.MkdirAll(dirName, 0777)
	}

	if err != nil {
		return fmt.Errorf(
			"error in verifying the %s directory. Check if either %s is present in %s or "+
				"%s env var is set", dirName, yamlTag, envVar, confFileName)
	}

	return err
}
