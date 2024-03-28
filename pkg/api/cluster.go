package api

import (
	"github.com/nrz-incubator/malygos/pkg/errors"
	"golang.org/x/mod/semver"
)

func (c *Cluster) ValidateInputs() error {
	if c.Id != nil {
		return errors.NewInvalidArgumentError("id field is not allowed")
	}

	if c.Region == "" {
		return errors.NewInvalidArgumentError("region field is required")
	}

	if c.Name == "" {
		return errors.NewInvalidArgumentError("name field is required")
	}

	if c.Version == "" {
		return errors.NewInvalidArgumentError("version field is required")
	}

	if !semver.IsValid(c.Version) {
		return errors.NewInvalidArgumentError("version field is not a valid semver")
	}

	if semver.Compare(c.Version, "v1.28.0") < 0 {
		return errors.NewInvalidArgumentError("version field must be >= v1.28.0")
	}

	return nil
}
