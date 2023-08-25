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
	kubeServiceAccount := preflightFlags.String("kube-service-account", "", "kube service account")
	awsArn := preflightFlags.String("aws-arn", "", "aws arn")
	gcpEmail := preflightFlags.String("gcp-email", "", "gcp email")
	configFile := preflightFlags.String("config", "", "config file to use")
	equiv := preflightFlags.Bool("equiv", false, "print equivalent command")
	preflightFlags.Parse(os.Args[1:])
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	preflightid.Logger = l.Logger
	l.Debugf("log level: %s", ll)
	pf := &preflightid.PreflightID{}
	if *configFile != "" {
		if pf, err = preflightid.LoadConfig(*configFile); err != nil {
			l.WithError(err).Error("error loading config")
			os.Exit(1)
		}
	}
	if *kubeServiceAccount != "" {
		pf.Kube = &preflightid.IDProviderKube{
			ServiceAccount: *kubeServiceAccount,
		}
	}
	if *awsArn != "" {
		pf.AWS = &preflightid.IDProviderAWS{
			ARN: *awsArn,
		}
	}
	if *gcpEmail != "" {
		pf.GCP = &preflightid.IDProviderGCP{
			Email: *gcpEmail,
		}
	}
	if *equiv {
		pf.Equivalent()
		os.Exit(0)
	}
	if err := pf.Run(); err != nil {
		l.Fatal(err)
	}
}
