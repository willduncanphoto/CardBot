package main

import "github.com/illwill/cardbot/internal/config"

// runSetup executes first-time/--setup prompts and persists config.
func runSetup(
	cfg *config.Config,
	cfgPath string,
	promptDestinationFn func(string) string,
	promptNamingFn func(string) string,
) error {
	cfg.Destination.Path = config.ContractPath(promptDestinationFn(cfg.Destination.Path))
	cfg.Naming.Mode = config.NormalizeNamingMode(promptNamingFn(cfg.Naming.Mode))

	if cfgPath == "" {
		return nil
	}
	return config.Save(cfg, cfgPath)
}
