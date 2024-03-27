package clusterregistrar

import "fmt"

type secretTypeError struct{}

func newSecretTypeError() *secretTypeError {
	return &secretTypeError{}
}

func (e *secretTypeError) Error() string {
	return fmt.Sprintf("secret type is not %s", managementClusterSecretType)
}
