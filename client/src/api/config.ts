import { OpType } from '@/api/operation';

export interface Config {
    id: number;
    name: string;
    src: ConfigSrc;
    dst: ConfigDst;
    ops: ConfigOps;
}

export interface ConfigSrc {
    dir: string;
    includeSubdirs: boolean;
    defaultOpType: OpType;
}

export interface ConfigDst {
    dirs: DstDir[];
}

export interface DstDir {
    hotkey: string;
    dir: string;
}

export interface ConfigOps {
    numWorkers: number;
    maxTries: number;
}

export interface ConfigValidationError {
    errors: ConfigErrorsByKey;
}

export interface ConfigErrorsByKey {
    [key: string]: string;
}
