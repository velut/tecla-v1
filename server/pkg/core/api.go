package core

// API represents the API for the core package.
type API interface {
	ConfigValidatorAPI
	OrganizerAPI
}

// ConfigValidatorAPI represents the API for the validation of configurations.
type ConfigValidatorAPI interface {
	// ValidateConfig validates the given configuration,
	// returning an error of type ConfigValidationError if the configuration is not valid.
	ValidateConfig(config *Config) error
}

// OrganizerAPI represents the API for the organizer.
type OrganizerAPI interface {
	// LoadConfig loads the given configuration, which must be valid, starting the organizer.
	LoadConfig(config *Config) (*OrganizerStatus, error)

	// DropConfigWait removes the current configuration, if any, stopping the organizer.
	// All submitted operations, pending or in progress, are completed.
	DropConfigWait() (*OrganizerStatus, error)

	// DropConfig removes the current configuration, if any, stopping the organizer.
	// In progress operations are completed, pending operations are discarded.
	DropConfig() (*OrganizerStatus, error)

	// OrganizerStatus returns the organizer's status.
	OrganizerStatus() (*OrganizerStatus, error)

	// HandleHotkey handles the action corresponding to the pressed hotkey.
	// If the hotkey is the space character, the organizer advances to the next file.
	// If the hotkey is associated to a destination directory, the organizer creates
	// a file operation for the current file and then advances to the next file.
	// If the hotkey is unrecognized, the organizer does nothing.
	HandleHotkey(hotkey string) (*OrganizerStatus, error)
}
