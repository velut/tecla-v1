import { Config } from '@/api/config';
import { OrganizerStatus } from '@/api/organizer';

/**
 * API represents the API for the server.
 */
export interface API extends ConfigValidatorAPI, OrganizerAPI {}

/**
 * ConfigValidatorAPI represents the API for the validation of configurations.
 */
export interface ConfigValidatorAPI {
    /**
     * validateConfig validates the given configuration,
     * returning an error of type ConfigValidationError if the configuration is not valid.
     */
    validateConfig: (config: Config) => Promise<void>;
}

/**
 * OrganizerAPI represents the API for the organizer.
 */
export interface OrganizerAPI {
    /**
     * loadConfig loads the given configuration, which must be valid, starting the organizer.
     */
    loadConfig: (config: Config) => Promise<OrganizerStatus>;

    /**
     * dropConfigWait removes the current configuration, if any, stopping the organizer.
     * All submitted operations, pending or in progress, are completed.
     */
    dropConfigWait: () => Promise<OrganizerStatus>;

    /**
     * dropConfig removes the current configuration, if any, stopping the organizer.
     * In progress operations are completed, pending operations are discarded.
     */
    dropConfig: () => Promise<OrganizerStatus>;

    /**
     * organizerStatus returns the organizer's status.
     */
    organizerStatus: () => Promise<OrganizerStatus>;

    /**
     * HandleHotkey handles the action corresponding to the pressed hotkey.
     * If the hotkey is the space character, the organizer advances to the next file.
     * If the hotkey is associated to a destination directory, the organizer creates
     * a file operation for the current file and then advances to the next file.
     * If the hotkey is unrecognized, the organizer does nothing.
     */
    handleHotkey: (hotkey: string) => Promise<OrganizerStatus>;
}

export const api = (window as unknown) as API;
export const configValidatorAPI = (window as unknown) as ConfigValidatorAPI;
export const organizerAPI = (window as unknown) as OrganizerAPI;
