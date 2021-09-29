package xvals

import (
	"fmt"
	"log"
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
	// In case of an error, the variables will simply not be available
	// An error message should be logged in that case.
	Reload()
}

// WithEnvironment adds environmental variables to the xval context.
// Will panic in case of errors
func WithEnvironment() XvalProvider {
	p := &envValProvider{}
	p.Reload()
	ctxt = append(ctxt, p)
	return p
}

// EnvValProvider provides values from the environment variables
type envValProvider struct {
	mapProvider
}

func (c *envValProvider) Reload() {
	res := make(map[string]string)
	for _, v := range os.Environ() {
		s := strings.Split(v, "=")
		if len(s) == 1 {
			res[strings.ToLower(s[0])] = ""
		} else if len(s) > 1 {
			res[strings.ToLower(s[0])] = strings.Join(s[1:], "")
		}
	}
	c.vals = res
}

// WithConfigFile adds a config file to the xval context. More than one file can be added.
// First added file has highest priority. Last added least priority.
func WithConfigFile(filename string) XvalProvider {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil
	}
	c := &configFileProvider{filename: absPath}
	c.Reload()
	ctxt = append(ctxt, c)
	return c
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
	vals := make(map[string]string)
	c.ctx = &CfgFile{Values: make(map[string]string)}
	e = yaml.Unmarshal(d, &vals)
	for k, v := range vals {
		c.ctx.Values[strings.ToLower(k)] = v
	}
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

func (c *configFileProvider) Reload() {
	e := c.readFile()
	if e != nil {
		log.Printf("failed to reload configFileProvider %v", e)
	}
}

// WithMap adds a map to the xval context.
func WithMap(src map[string]string) XvalProvider {
	p := &mapProvider{vals: src}
	ctxt = append(ctxt, p)
	return p
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

func (c *mapProvider) Reload() {
}
