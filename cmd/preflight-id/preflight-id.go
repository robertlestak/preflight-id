package main

import (
	"flag"
	"os"

	"github.com/robertlestak/preflight-id/pkg/preflightid"
	log "github.com/sirupsen/logrus"
)

func init() {
	ll, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func main() {
	l := log.WithFields(log.Fields{
		"app": "preflight-id",
	})
	l.Debug("starting preflight-id")
	preflightFlags := flag.NewFlagSet("preflight-id", flag.ExitOnError)
	logLevel := preflightFlags.String("log-level", log.GetLevel().String(), "log level")
	provider := preflightFlags.String("provider", "", "provider. one of: aws, gcp, kube")
	kubeServiceAccount := preflightFlags.String("kube-service-account", "", "kube service account")
	awsArn := preflightFlags.String("aws-arn", "", "aws arn")
	gcpEmail := preflightFlags.String("gcp-email", "", "gcp email")
	preflightFlags.Parse(os.Args[1:])
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	if *provider == "" {
		// infer provider from flags
		if *kubeServiceAccount != "" {
			*provider = "kube"
		} else if *awsArn != "" {
			*provider = "aws"
		} else if *gcpEmail != "" {
			*provider = "gcp"
		}
	}
	l.Debugf("log level: %s", ll)
	l.Debugf("provider: %s", *provider)
	pf := &preflightid.PreflightID{
		Provider: preflightid.Provider(*provider),
	}
	switch pf.Provider {
	case preflightid.ProviderAWS:
		pf.AWS = &preflightid.IDProviderAWS{
			ARN: *awsArn,
		}
	case preflightid.ProviderGCP:
		pf.GCP = &preflightid.IDProviderGCP{
			Email: *gcpEmail,
		}
	case preflightid.ProviderKube:
		pf.Kube = &preflightid.IDProviderKube{
			ServiceAccount: *kubeServiceAccount,
		}
	}
	if err := pf.Run(); err != nil {
		l.Fatal(err)
	}
}
