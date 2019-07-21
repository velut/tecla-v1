import { Config } from '@/api/config';
import { File } from '@/api/file';

/**
 * OrganizerStatus represents the status of the organizer.
 */
export interface OrganizerStatus {
    config: Config;
    currentFile: File;
    currentFileIndex: number;
    numFiles: number;
}
