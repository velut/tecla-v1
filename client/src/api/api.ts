import { Config } from '@/api/config';
import { Credits, Info } from '@/api/info';
import { OrganizerStatus } from '@/api/organizer';

/**
 * API represents the API for the server.
 */
export interface API
    extends ConfigValidatorAPI,
        OrganizerAPI,
        AppInfoAPI,
        DialogAPI {}

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
     * restoreConfig TODO:
     */
    restoreConfig: () => Promise<Config>;

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

/**
 * AppInfoAPI represents the API for retrieving information about the application.
 */
export interface AppInfoAPI {
    /**
     * appInfo returns information about the application.
     */
    appInfo: () => Promise<Info>;

    /**
     * appCredits returns information about the libraries included in the application.
     */
    appCredits: () => Promise<Credits>;
}

/**
 * DialogAPI represents the API for interacting with dialogs.
 */
export interface DialogAPI {
    /**
     * selectDirectory opens a directory selection dialog
     * and returns the path of the selected directory.
     */
    selectDirectory: () => Promise<string>;
}

export const api = (window as unknown) as API;
export const configValidatorAPI = (window as unknown) as ConfigValidatorAPI;
export const organizerAPI = (window as unknown) as OrganizerAPI;
export const appInfoAPI = (window as unknown) as AppInfoAPI;
export const dialogAPI = (window as unknown) as DialogAPI;
