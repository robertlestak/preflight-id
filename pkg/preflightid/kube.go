package preflightid

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type IDProviderKube struct {
	ServiceAccount string `json:"serviceAccount" yaml:"serviceAccount"`
}

func (k *IDProviderKube) Run() error {
	l := log.WithFields(log.Fields{
		"preflight": "id",
		"provider":  "kube",
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
		failStr := "failed - no service account name in JWT token"
		l.Error(failStr)
		return errors.New(failStr)
	}
	if k.ServiceAccount != "" && k.ServiceAccount != serviceAccountName {
		failStr := fmt.Sprintf("failed - expected %s, got %s", k.ServiceAccount, serviceAccountName)
		l.Error(failStr)
		return errors.New(failStr)
	}
	l.Info("passed")
	return nil
}
