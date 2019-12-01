package config

import (
	"errors"

	cron "github.com/robfig/cron/v3"
)

func (cfg *Config) Validate() []error {
	errs := make([]error, 0)

	if len(cfg.Restic.Password) <= 0 {
		errs = append(errs, errors.New("missing restic password"))
	}

	if cfg.Schedule != "" {
		specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err := specParser.Parse(cfg.Schedule); err != nil {
			errs = append(errs, err)
		}
	}

	// TODO much more validation rules to fail immediately

	return errs
}
