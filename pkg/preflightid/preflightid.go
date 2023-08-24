package preflightid

import (
	"encoding/json"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	Logger *log.Logger
)

func init() {
	if Logger == nil {
		Logger = log.New()
		Logger.SetOutput(os.Stdout)
		Logger.SetLevel(log.InfoLevel)
	}
}

type IDProvider interface {
	Run() error
	Equivalent()
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
	l := Logger.WithFields(log.Fields{
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

func NewPreflighter(provider Provider, config *PreflightID) (IDProvider, error) {
	l := Logger.WithFields(log.Fields{
		"fn":       "NewPreflighter",
		"provider": provider,
	})
	l.Debug("creating preflighter")
	switch provider {
	case ProviderAWS:
		return config.AWS, nil
	case ProviderGCP:
		return config.GCP, nil
	case ProviderKube:
		return config.Kube, nil
	default:
		return nil, errors.New("invalid provider")
	}
}

func (p *PreflightID) InferProvider() error {
	l := Logger.WithFields(log.Fields{
		"fn": "InferProvider",
	})
	l.Debug("inferring provider")
	var inferredProvider Provider
	if p.Provider != "" {
		l.Debug("provider already configured")
		return nil
	}
	if p.AWS != nil {
		inferredProvider = ProviderAWS
	} else if p.GCP != nil {
		inferredProvider = ProviderGCP
	} else if p.Kube != nil {
		inferredProvider = ProviderKube
	}
	if inferredProvider == "" {
		return errors.New("unable to infer provider")
	}
	p.Provider = inferredProvider
	return nil
}

func (p *PreflightID) Run() error {
	l := Logger.WithFields(log.Fields{
		"preflight": "id",
	})
	l.Debug("running preflight-id")
	if err := p.InferProvider(); err != nil {
		l.WithError(err).Error("error inferring provider")
		return err
	}
	preflighter, err := NewPreflighter(p.Provider, p)
	if err != nil {
		l.WithError(err).Error("error creating preflighter")
		return err
	}
	if err := preflighter.Run(); err != nil {
		l.WithError(err).Error("error running preflighter")
		return err
	}
	return nil
}
