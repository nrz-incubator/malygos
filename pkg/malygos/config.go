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

	return nil
}
