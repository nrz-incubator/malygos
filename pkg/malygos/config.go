package malygos

import (
	"fmt"
	"os"
)

func (m *Malygos) readConfiguration() error {
	m.kubeconfig = os.Getenv("KUBECONFIG")
	if m.kubeconfig == "" {
		return fmt.Errorf("KUBECONFIG variable not set")
	}

	m.managementNamespace = os.Getenv("MANAGEMENT_NAMESPACE")
	if m.managementNamespace == "" {
		return fmt.Errorf("MANAGEMENT_NAMESPACE variable not set")
	}

	m.logger.WithValues("kubeconfig", m.kubeconfig,
		"managementNamespace", m.managementNamespace,
	).Info("configuration set")

	return nil
}
