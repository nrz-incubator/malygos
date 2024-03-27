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

	m.logger.WithValues("kubeconfig", m.kubeconfig).Info("configuration set")

	return nil
}
