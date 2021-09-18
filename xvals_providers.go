package xvals

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// XvalProvider is the interface all providers of values must implement
type XvalProvider interface {

	// Get Value or error
	Value(key string) (val string, err error)

	// Dump all values from the source
	Dump() map[string]string

	// Reload the values from the source
	Reload() error
}

// WithEnvironment adds environmental variables to the xval context.
// Will panic in case of errors
func WithEnvironment() (XvalProvider, error) {
	p := &envValProvider{}
	if e := p.Reload(); e != nil {
		return nil, e
	}
	ctxt = append(ctxt, p)
	return p, nil
}

// EnvValProvider provides values from the environment variables
type envValProvider struct {
	mapProvider
}

func (c *envValProvider) Reload() error {
	res := make(map[string]string)
	for _, v := range os.Environ() {
		s := strings.Split(v, "=")
		if len(s) == 1 {
			res[s[0]] = ""
		} else if len(s) > 1 {
			res[s[0]] = strings.Join(s[1:], "")
		}
	}
	c.vals = res
	return nil
}

// WithConfigFile adds a config file to the xval context. More than one file can be added.
// First added file has highest priority. Last added least priority.
func WithConfigFile(filename string) (XvalProvider, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	c := &configFileProvider{filename: absPath}
	e := c.readFile()
	if e != nil {
		return nil, e
	}
	ctxt = append(ctxt, c)
	return c, nil
}

// CfgFile is the structure of files the ConfigFileProvider uses.
type CfgFile struct {
	Values map[string]string `yaml:"values,inline"`
}

// A configFileProvider provides values from a configuration file.
type configFileProvider struct {
	filename string
	ctx      *CfgFile
}

func (c *configFileProvider) readFile() error {
	d, e := os.ReadFile(c.filename)
	if e != nil {
		return e
	}
	c.ctx = &CfgFile{Values: make(map[string]string)}
	e = yaml.Unmarshal(d, c.ctx)
	return e
}

func (c *configFileProvider) Value(key string) (val string, err error) {
	if v, ok := c.ctx.Values[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("failed to retrieve key %s from config file %s", key, c.filename)
}

func (c *configFileProvider) Dump() map[string]string {
	return c.ctx.Values
}

func (c *configFileProvider) Reload() error {
	return c.readFile()
}

// WithEnvironment adds environmental variables to the xval context.
func WithMap(src map[string]string) (XvalProvider, error) {
	p := &mapProvider{vals: src}
	ctxt = append(ctxt, p)
	return p, nil
}

// mapProvider provides values from a map
type mapProvider struct {
	vals map[string]string
}

func (c *mapProvider) Value(key string) (value string, err error) {
	if v, ok := c.vals[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("failed to retrieve key %s from src map", key)
}

func (c *mapProvider) Dump() map[string]string {
	return c.vals
}

func (c *mapProvider) Reload() error {
	return nil
}
