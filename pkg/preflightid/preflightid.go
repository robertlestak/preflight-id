package preflightid

import (
	"encoding/json"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type IDProvider interface {
	Run() error
}

type Provider string

const (
	ProviderGCP  Provider = "gcp"
	ProviderAWS  Provider = "aws"
	ProviderKube Provider = "kube"
)

type PreflightID struct {
	Provider Provider        `json:"provider" yaml:"provider"`
	AWS      *IDProviderAWS  `json:"aws" yaml:"aws"`
	GCP      *IDProviderGCP  `json:"gcp" yaml:"gcp"`
	Kube     *IDProviderKube `json:"kube" yaml:"kube"`
}

func LoadConfig(filepath string) (*PreflightID, error) {
	l := log.WithFields(log.Fields{
		"fn": "LoadConfig",
	})
	l.Debug("loading config")
	var err error
	pf := &PreflightID{}
	bd, err := os.ReadFile(filepath)
	if err != nil {
		l.WithError(err).Error("error reading file")
		return pf, err
	}
	if err := yaml.Unmarshal(bd, pf); err != nil {
		// try with json
		if err := json.Unmarshal(bd, pf); err != nil {
			l.WithError(err).Error("error unmarshalling config")
			return pf, err
		}
	}
	return pf, err
}

func (p *PreflightID) Run() error {
	l := log.WithFields(log.Fields{
		"app": "preflight-id",
		"fn":  "p.Run",
	})
	l.Debug("running preflight-id")
	switch p.Provider {
	case ProviderAWS:
		return p.AWS.Run()
	case ProviderGCP:
		return p.GCP.Run()
	case ProviderKube:
		return p.Kube.Run()
	default:
		return errors.New("invalid provider")
	}
}
