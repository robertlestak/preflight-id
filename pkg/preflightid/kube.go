package preflightid

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type IDProviderKube struct {
	ServiceAccount string `json:"serviceAccount"`
}

func (k *IDProviderKube) Run() error {
	l := log.WithFields(log.Fields{
		"app":      "preflight-id",
		"provider": "kube",
		"fn":       "k.Run",
	})
	l.Debug("running preflight-id")
	if k.ServiceAccount == "" {
		return errors.New("service account name not configured")
	}
	tokenFile := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return err
	}
	token := strings.TrimSpace(string(tokenBytes))
	l.Debugf("token: %s", token)

	// Decode the JWT token to extract claims
	type serviceAccountClaims struct {
		Name string `json:"name"`
		Uid  string `json:"uid"`
	}
	type kubeClaims struct {
		ServiceAccount serviceAccountClaims `json:"serviceaccount"`
	}
	claims := struct {
		Kubernetes kubeClaims `json:"kubernetes.io"`
	}{}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid JWT token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return err
	}
	serviceAccountName := claims.Kubernetes.ServiceAccount.Name
	l.Debugf("service account name: %s", serviceAccountName)
	if serviceAccountName == "" {
		return errors.New("invalid JWT token claims")
	}
	if k.ServiceAccount != "" && k.ServiceAccount != serviceAccountName {
		l.WithError(err).Errorf("service account name mismatch: %s != %s", k.ServiceAccount, serviceAccountName)
		return errors.New("service account name mismatch")
	}
	l.Info("service account name match")
	return nil
}
