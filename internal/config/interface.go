package config

import "github.com/aliases/internal/domain"

// Manager describes configuration persistence behavior.
type Manager interface {
	Load() (domain.Config, error)
	Save(config domain.Config) error
	Exists() (bool, error)
}
