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
	Equiv          bool   `json:"equiv" yaml:"equiv"`
}

func (k *IDProviderKube) RunEquiv() bool {
	return k.Equiv
}

func (k *IDProviderKube) Equivalent() {
	l := Logger
	l.Debug("printing equivalent command")
	cmd := `sh -c 'EXPECTED="%s"; TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token); PAYLOAD=$(echo "$TOKEN" | cut -d. -f2); DECODED_PAYLOAD=$(echo "$PAYLOAD" | base64 -d 2>/dev/null); SERVICE_ACCOUNT=$(echo "$DECODED_PAYLOAD" | jq -r '.sub'); SERVICE_ACCOUNT_NAME=$(echo "$SERVICE_ACCOUNT" | cut -d: -f4); if [ "$SERVICE_ACCOUNT_NAME" != "$EXPECTED" ]; then echo "Service account $SERVICE_ACCOUNT_NAME does not match expected $EXPECTED"; exit 1; fi'`
	cmd = fmt.Sprintf(cmd, k.ServiceAccount)
	fmt.Println(cmd)
}

func (k *IDProviderKube) Run() error {
	l := Logger.WithFields(log.Fields{
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
