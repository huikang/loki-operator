package manifests

import (
	"fmt"
	"strings"

	"github.com/ViaQ/loki-operator/internal/manifests/internal/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LokiConfigMap creates the single configmap containing the loki configuration for the whole cluster
func LokiConfigMap(stackName, namespace string) (*corev1.ConfigMap, error) {
	b, err := config.Build(config.Options{
		FrontendWorker: config.Address{
			FQDN: "",
			Port: 0,
		},
		GossipRing: config.Address{
			FQDN: fqdn(LokiGossipRingService(stackName).GetName(), namespace),
			Port: gossipPort,
		},
		Querier: config.Address{
			FQDN: serviceNameQuerierHTTP(stackName),
			Port: httpPort,
		},
		StorageDirectory: strings.TrimRight(dataDirectory, "/"),
		Namespace:        namespace,
	})
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   lokiConfigMapName(stackName),
			Labels: commonLabels(stackName),
		},
		BinaryData: map[string][]byte{
			config.LokiConfigFileName: b,
		},
	}, nil
}

func lokiConfigMapName(stackName string) string {
	return fmt.Sprintf("loki-config-%s", stackName)
}
