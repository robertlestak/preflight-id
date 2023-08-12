# preflight-id

A preflight check to validate the expected identity is bound in the environment.

## Build

```bash
make
```

## Install

NOTE: you will need `curl`, `bash`, and `jq` installed for the install script to work. It will attempt to install the binary in `/usr/local/bin` and will require `sudo` access. You can override the install directory by setting the `INSTALL_DIR` environment variable.

```bash
curl -sSL https://raw.githubusercontent.com/robertlestak/preflight-id/main/scripts/install.sh | bash
```

## Usage

```bash
Usage of preflight-id:
  -aws-arn string
        aws arn
  -gcp-email string
        gcp email
  -kube-service-account string
        kube service account
  -log-level string
        log level (default "info")
```

### AWS

```bash
preflight-id \
    -aws-arn arn:aws:iam::123456789012:role/role-name
```

### GCP

```bash
preflight-id \
    -gcp-email my-example@my-project.google.com
```

### Kubernetes

```bash
preflight-id \
    -kube-service-account my-service-account
```