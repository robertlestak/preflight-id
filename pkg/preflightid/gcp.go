package preflightid

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/compute/metadata"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

type IDProviderGCP struct {
	Email string `json:"email" yaml:"email"`
}

func (p *IDProviderGCP) Run() error {
	l := log.WithFields(log.Fields{
		"app":      "preflight-id",
		"provider": "gcp",
		"fn":       "p.Run",
	})
	l.Debug("running preflight-id")
	if p.Email == "" {
		return errors.New("email not configured")
	}
	// Initialize a GCP client with the appropriate credentials
	ctx := context.Background()
	client, err := iam.NewService(ctx, option.WithScopes(iam.CloudPlatformScope))
	if err != nil {
		l.WithError(err).Error("Failed to initialize GCP client")
		return err
	}

	// Get the list of authorized accounts using the service client.
	response, err := client.Projects.ServiceAccounts.List("projects/-").Do()
	if err != nil {
		l.WithError(err).Error("Failed to retrieve authorized accounts")
		return err
	}
	for _, account := range response.Accounts {
		if strings.EqualFold(account.Email, p.Email) {
			l.Debugf("Service Account match: %s", account.Email)
			return nil
		}
	}
	if metadata.OnGCE() {
		vmIdentity, err := metadata.Email("default")
		if err != nil {
			l.WithError(err).Error("Failed to retrieve VM Identity")
			return err
		}
		if strings.EqualFold(vmIdentity, p.Email) {
			l.Debugf("VM Identity match: %s", vmIdentity)
			return nil
		}
	}
	l.Errorf("Service Account not found: %s", p.Email)
	return errors.New(fmt.Sprintf("Service Account not found"))
}
