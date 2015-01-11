package config

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"ktkr.us/pkg/gas/auth"
)

var OnSave func(*Config)

var (
	configChan   = make(chan *Config)
	reloadChan   = make(chan struct{})
	errChan      = make(chan error)
	sharedConfig *Config
	confPath     string
	Default      Config
)

func Init(filePath string) error {
	confPath = filePath
	var err error
	sharedConfig, err = Load()
	return err
}

type Config struct {
	Host      string
	Port      int
	Password  []byte
	Salt      []byte
	Directory string
	MaxAge    int   // max age of uploads in days
	MaxSize   int64 // max total size of uploads in MB
}

// satisfies gas.User interface
func (c Config) Secrets() (pass, salt []byte, err error) {
	return c.Password, c.Salt, nil
}

func (c Config) Username() string {
	return ""
}

// Update the config with the new password hash, generating a new random salt
func (c *Config) SetPass(pass string) {
	c.Salt = make([]byte, 32)
	rand.Read(c.Salt)
	c.Password = auth.Hash([]byte(pass), c.Salt)
}

func Serve() {
	for {
		if sharedConfig == nil {
			sharedConfig = &Default
		}
		select {
		case conf := <-configChan:
			err := Save(conf)
			errChan <- err
			if err != nil {
				log.Printf("Failed to write config: %v", err)
			} else {
				log.Printf("Config updated on disk.")
				sharedConfig = conf
				if OnSave != nil {
					OnSave(sharedConfig)
				}
			}
		case configChan <- sharedConfig:
		case <-reloadChan:
			conf, err := Load()
			errChan <- err
			if err != nil {
				log.Printf("Failed to reload config: %v", err)
			} else {
				log.Print("Reloaded config.")
				sharedConfig = conf
			}
		}
	}
}

func Get() *Config {
	return <-configChan
}

func Set(c *Config) error {
	configChan <- c
	return <-errChan
}

func Reload() error {
	reloadChan <- struct{}{}
	return <-errChan
}

func Load() (*Config, error) {
	var conf Config

	confFile, err := os.Open(confPath)
	if err != nil {
		if os.IsNotExist(err) {
			conf = Default
			err = Save(&conf)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("reading config: %v", err)
		}
	} else {
		b, err := ioutil.ReadAll(confFile)
		if err != nil {
			return nil, fmt.Errorf("reading config: %v", err)
		}
		err = json.Unmarshal(b, &conf)
		if err != nil {
			return nil, fmt.Errorf("decoding config: %v", err)
		}
	}

	return &conf, nil
}

func Save(conf *Config) error {
	b, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return fmt.Errorf("encoding config: %v", err)
	}
	file, err := os.OpenFile(confPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return fmt.Errorf("writing config: %v", err)
	}
	defer file.Close()
	file.Write(b)
	return nil
}