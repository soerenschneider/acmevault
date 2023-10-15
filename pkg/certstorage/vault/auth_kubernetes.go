package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"golang.org/x/net/context"
)

const (
	KubernetesAuthName             = "kubernetes"
	defaultServiceAccountTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token" // #nosec G101
	defaultMount                   = "kubernetes"
)

type KubernetesAuth struct {
	role                    string
	serviceAccountTokenFile string
	mount                   string
}

func NewVaultKubernetesAuth(role string, mountPath string) (*KubernetesAuth, error) {
	return &KubernetesAuth{
		role:                    role,
		mount:                   mountPath,
		serviceAccountTokenFile: defaultServiceAccountTokenFile,
	}, nil
}

func (t *KubernetesAuth) Logout(ctx context.Context, client *api.Client) error {
	path := "auth/token/revoke-self"
	_, err := client.Logical().Write(path, map[string]interface{}{})
	return err
}

func (t *KubernetesAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	k8sAuth, err := kubernetes.NewKubernetesAuth(
		t.role,
		kubernetes.WithServiceAccountTokenPath(t.serviceAccountTokenFile),
		kubernetes.WithMountPath(t.mount))

	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes kubernetes method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes kubernetes: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no kubernetes info was returned after login")
	}

	return authInfo, nil
}
