package preflightid

import (
	"errors"

	log "github.com/sirupsen/logrus"
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
	Provider Provider        `json:"provider"`
	AWS      *IDProviderAWS  `json:"aws"`
	GCP      *IDProviderGCP  `json:"gcp"`
	Kube     *IDProviderKube `json:"kube"`
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
