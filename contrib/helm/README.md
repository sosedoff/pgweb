## Helm Chart

Helm in the package manager for Kubernetes.

Usage:

```
# Install helm-git plugin.
helm plugin install https://github.com/aslafy-z/helm-git --version master

# Add helm repository
helm repo add pgweb git+https://github.com/sosedoff/pgweb@contrib/helm

# Install pgweb on your cluster.
helm install pgweb pgweb/pgweb --set databaseUrl="postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"

```