package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"

	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

const (
	// FileConfigPath path to config file.
	FileConfigPath = "config/config.yaml"
)

// Config type keeps general GetConfiguration().
type Config struct {
	General struct {
		ReportDirAbsPath          string `yaml:"report" envconfig:"REPORT_DIR_NAME"`
		CnfNodeLabel              string `yaml:"cnf_worker_label" envconfig:"ROLE_WORKER_CNF"`
		WorkerNodeLabel           string `yaml:"worker_label" envconfig:"ROLE_WORKER"`
		TestImage                 string `yaml:"test_image" envconfig:"TEST_IMAGE"`
		VerificationLogLevel      string `yaml:"verification_log_level" envconfig:"VERIFICATION_LOG_LEVEL"`
		DebugCertsuite            string `envconfig:"DEBUG_CERTSUITE"`
		CertsuiteConfigDir        string `yaml:"certsuite_config_dir" envconfig:"CERTSUITE_CONFIG_DIR"`
		CertsuiteRepoPath         string `envconfig:"CERTSUITE_REPO_PATH"`
		CertsuiteEntryPointBinary string `yaml:"certsuite_entry_point_binary" envconfig:"CERTSUITE_ENTRY_POINT_BINARY"`
		CertsuiteReportDir        string `yaml:"certsuite_report_dir" envconfig:"CERTSUITE_REPORT_DIR"`
		DockerConfigDir           string `yaml:"docker_config_dir" envconfig:"DOCKER_CONFIG_DIR"`
		CertsuiteImage            string `yaml:"certsuite_image" envconfig:"CERTSUITE_IMAGE"`
		CertsuiteImageTag         string `yaml:"certsuite_image_tag" envconfig:"CERTSUITE_IMAGE_TAG"`
		DisableIntrusiveTests     string `yaml:"disable_intrusive_tests" envconfig:"DISABLE_INTRUSIVE_TESTS"`
		ContainerEngine           string `default:"docker" yaml:"container_engine" envconfig:"CONTAINER_ENGINE"`
		UseBinary                 string `default:"false" yaml:"use_binary" envconfig:"USE_BINARY"`
	} `yaml:"general"`
}

// DefineClients sets client and return it's instance.
func DefineClients() (*testclient.ClientSet, error) {
	clients := testclient.New("")
	if clients == nil {
		return nil, fmt.Errorf("client is not set please check KUBECONFIG env variable")
	}

	return clients, nil
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

	err = conf.deployCertsuiteConfigDir(confFile)

	if err != nil {
		return nil, err
	}

	err = conf.deployCertsuiteReportDir(confFile)
	if err != nil {
		return nil, err
	}

	err = conf.makeDockerConfig()
	if err != nil {
		return nil, err
	}

	conf.General.CertsuiteRepoPath, err = conf.defineCertsuiteRepoPath()

	if err != nil {
		glog.Fatal(err)
	}

	return &conf, nil
}

// DebugCertsuite activates debug mode.
func (c *Config) DebugCertsuite() (bool, error) {
	if c.General.DebugCertsuite == "true" {
		err := os.Setenv("CERTSUITE_LOG_LEVEL", "debug")
		if err != nil {
			return false, fmt.Errorf("failed to set env var CERTSUITE_LOG_LEVEL: %w", err)
		}

		return true, nil
	}

	return false, nil
}

// CreateLogFile creates log file for testSuite.
func (c *Config) CreateLogFile(testSuite string, tcName string) *os.File {
	folderPath := filepath.Join(c.General.ReportDirAbsPath, "Debug", testSuite, tcName)

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

// GetReportPath returns full path to the report file.
func (c *Config) GetReportPath(file string) string {
	reportFileName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))

	return fmt.Sprintf("%s.xml", filepath.Join(c.General.ReportDirAbsPath, reportFileName))
}

func (c *Config) defineCertsuiteRepoPath() (string, error) {
	if c.General.CertsuiteRepoPath == "" {
		return "", fmt.Errorf("CERTSUITE_REPO_PATH env variable is not set. Please export CERTSUITE_REPO_PATH")
	}

	return c.General.CertsuiteRepoPath, nil
}

func readFile(cfg *Config, cfgFile string) error {
	openedCnfFile, err := os.Open(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to open cfg file: %w", err)
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

func (c *Config) deployCertsuiteConfigDir(configFileName string) error {
	return deployCertsuiteDir(configFileName, c.General.CertsuiteConfigDir, "certsuite_config_dir", "CERTSUITE_CONFIG_DIR")
}

func (c *Config) deployCertsuiteReportDir(configFileName string) error {
	return deployCertsuiteDir(configFileName, c.General.CertsuiteReportDir, "certsuite_report_dir", "CERTSUITE_REPORT_DIR")
}

func checkFileExists(filePath, name string) (string, error) {
	if !filepath.IsAbs(filePath) {
		return "", fmt.Errorf(
			"make sure env var CERTSUITE_REPO_PATH is configured with absolute path instead of relative",
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

func deployCertsuiteDir(confFileName string, dirName string, yamlTag string, envVar string) error {
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

func (c *Config) makeDockerConfig() error {
	var configFile *os.File

	err := os.MkdirAll(c.General.DockerConfigDir, 0777)

	if err != nil {
		return err
	}

	err = os.Chdir(c.General.DockerConfigDir)

	if err != nil {
		return err
	}

	configFile, err = os.Create("config")

	if err != nil {
		return err
	}

	_, err = configFile.Write([]byte("{ \"auths\": {} }"))

	return err
}
