import { Module, VuexModule, Mutation, Action } from 'vuex-module-decorators';
import { namespace } from 'vuex-class';
import to from 'await-to-js';

import { Config, DstDir } from '@/api/config';
import { File } from '@/api/file';
import { organizerAPI } from '@/api/api';
import { OrganizerStatus } from '@/api/organizer';

export const organizer = namespace('organizer');

@Module({ namespaced: true })
export default class Organizer extends VuexModule {
    /**
     * config represents the configuration currently in use by the organizer.
     */
    public config: Config | null = null;

    /**
     * currentFile represents the file currently displayed by the client.
     */
    public currentFile: File | null = null;

    /**
     * currentFileIndex represents the index of the file
     * currently displayed by the client.
     */
    public currentFileIndex: number = 0;

    /**
     * numFiles represents the number of files to be organized.
     */
    public numFiles: number = 0;

    /**
     * hasCurrentFile returns true if the organizer is active
     * and has a file to display.
     */
    public get hasCurrentFile(): boolean {
        return this.isActive && this.currentFile !== null;
    }

    /**
     * isActive returns true if the organizer is active.
     */
    public get isActive(): boolean {
        return this.config !== null;
    }

    /**
     * dstDirs returns the list of destination directories.
     */
    public get dstDirs(): DstDir[] {
        return this.config!.dst.dirs;
    }

    @Action({ commit: 'setStatus' })
    public async loadConfig(config: Config) {
        const [err, status] = await to<OrganizerStatus, string>(
            organizerAPI.loadConfig(config),
        );
        if (err) {
            console.error(err);
            return null;
        }

        return status;
    }

    @Action({ commit: 'setStatus' })
    public async dropConfigWait() {
        const [err, status] = await to<OrganizerStatus, string>(
            organizerAPI.dropConfigWait(),
        );
        if (err) {
            console.error(err);
            return null;
        }

        return status;
    }

    @Action({ commit: 'setStatus' })
    public async dropConfig() {
        const [err, status] = await to<OrganizerStatus, string>(
            organizerAPI.dropConfig(),
        );
        if (err) {
            console.error(err);
            return null;
        }

        return status;
    }

    @Action({ commit: 'setStatus' })
    public async updateStatus() {
        const [err, status] = await to<OrganizerStatus, string>(
            organizerAPI.organizerStatus(),
        );
        if (err) {
            console.error(err);
            return null;
        }

        return status;
    }

    @Action({ commit: 'setStatus' })
    public async handleHotkey(hotkey: string) {
        const [err, status] = await to<OrganizerStatus, string>(
            organizerAPI.handleHotkey(hotkey),
        );
        if (err) {
            console.error(err);
            return null;
        }

        return status;
    }

    @Mutation
    private setStatus(status: OrganizerStatus | null) {
        if (status) {
            this.config = status.config;
            this.currentFile = status.currentFile;
            this.currentFileIndex = status.currentFileIndex;
            this.numFiles = status.numFiles;
        }
    }
}
