export interface Info {
    name: string;
    description: string;
    homepage: string;
    repository: string;
    version: string;
    commit: string;
    copyright: string;
    warranty: string;
    license: string;
}

export interface Credits {
    credits: Library[];
}

export interface Library {
    name: string;
    homepage: string;
    license: string;
}
