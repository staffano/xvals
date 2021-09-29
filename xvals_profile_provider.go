package xvals

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ProfileFile struct {
	CurrentProfile string                       `yaml:"current_profile"`
	Profiles       map[string]map[string]string `yaml:"profiles"`
}

// profileProvider is an xvals provider holding several profiles of values
type profileProvider struct {
	mapProvider
	filename string
}

// WithProfile adds a file that has a one or several profiles, which of one
// is the current profile. Each profile is a mapProvider
func WithProfile(profileFilePath string) XvalProvider {
	p := &profileProvider{filename: profileFilePath}
	p.Reload()
	ctxt = append(ctxt, p)
	return p
}

// Reload the profile file
func (c *profileProvider) Reload() {
	data, err := os.ReadFile(c.filename)
	if err != nil {
		log.Printf("failed to load file %s %v", c.filename, err)
	}
	content := &ProfileFile{}
	err = yaml.Unmarshal(data, content)
	if err != nil {
		log.Printf("failed to parse file %s %v", c.filename, err)
	}
	cp, ok := content.Profiles[content.CurrentProfile]
	if ok {
		c.mapProvider = mapProvider{vals: cp}
		return
	}
	log.Printf("current profile %s is not available in profile file %s", content.CurrentProfile, c.filename)
}
