package block

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Configuration struct {
	Server        ServerConfig             `json:"server" toml:"server"`
	Authenticator *AuthenticatorConfig     `json:"authenticator,omitempty" toml:"authenticator"`
	Handlers      map[string]HandlerConfig `json:"handlers,omitempty" toml:"handler"`
	ConfigPath    string                   `json:"-" toml:"-"`
}

type AuthenticatorConfig struct {
	Protocol   string `json:"protocol" toml:"protocol" validate:"required"`
	Entrypoint string `json:"entrypoint" toml:"entrypoint" validate:"required"`
}

type HandlerConfig struct {
	Addr       string `json:"addr" toml:"addr"`
	Port       uint16 `json:"port" toml:"port" validate:"required"`
	Path       string `json:"path" toml:"path" validate:"required"`
	Entrypoint string `json:"entrypoint" toml:"entrypoint" validate:"required"`
}

type ServerConfig struct {
	Addr string `json:"addr" toml:"addr"`
	Port uint16 `json:"port" toml:"port" validate:"required"`
}

func NewConfigurationFromFile(path string) *Configuration {
	p, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := &Configuration{
		ConfigPath: p,
	}
	_, err = toml.DecodeFile(p, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return c
}

func (c *Configuration) HandlerConfig(name string) HandlerConfig {
	return c.Handlers[name]
}

func (c *Configuration) String() string {
	bs, err := json.Marshal(c)
	if err != nil {
		log.WithError(err).Fatalln("json.Marshal")
	}
	return string(bs)
}

func (c *Configuration) Valid() error {
	validate := validator.New()
	for key := range c.Handlers {
		if err := validate.Struct(c.Handlers[key]); err != nil {
			return err
		}
	}
	return validate.Struct(c)
}
