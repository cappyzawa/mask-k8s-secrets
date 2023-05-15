# mask-k8s-secrets

This command masks your kustomize secrets.

## How to install

This command can be install from [GitHub Releases](https://github.com/cappyzawa/mask-k8s-secrets/releases).

## How to use

```
$ mask-k8s-secrets

# use kustomize
$ kustomize build . | mask-k8s-secrets

# use a raw manifest
$ cat ./manifest.yaml | mask-k8s-secrets
```

Output is as follows.

```yaml
apiVersion: v1
data:
    password: '***'
    username: '***'
kind: Secret
metadata:
    name: secret-basic-auth
type: kubernetes.io/basic-auth
```
