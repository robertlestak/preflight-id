# preflight-id

A preflight check to validate the expected identity is bound in the environment.

## Build

```bash
make
```

## Usage

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