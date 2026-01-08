package config

import (
	"os"
	"path"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	Configuration struct {
		App      App      `yaml:"app"`
		Worker   Worker   `yaml:"worker"`
		Consumer Consumer `yaml:"consumer"`
	}

	App struct {
		Name            string        `yaml:"name"`
		GracefulTimeout time.Duration `yaml:"graceful_timeout"`
		TimeOutAPI      time.Duration `yaml:"time_out_api"`
		Port            string        `yaml:"port"`
	}

	Worker struct {
		UploadWorker struct {
			Size int `yaml:"size"`
		} `yaml:"upload_worker"`
	}

	Consumer struct {
		ReconciliationConsumer struct {
			Size int `yaml:"size"`
		} `yaml:"reconciliation_consumer"`
	}
)

func New() (*Configuration, error) {
	_, filename, _, _ := runtime.Caller(1)

	envPath := path.Join(path.Dir(filename), "../../config.yaml")

	_, err := os.Stat(envPath)
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(envPath)
	if err != nil {
		return nil, err
	}

	var c *Configuration
	err = yaml.Unmarshal([]byte(os.ExpandEnv(string(yamlFile))), &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
