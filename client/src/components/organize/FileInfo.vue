<template>
    <v-layout shrink wrap>
        <v-flex
            v-for="info in infos"
            :key="info.label"
            v-bind="{ [`xs${info.columns}`]: true }"
        >
            <v-text-field
                :label="info.label"
                :value="info.value"
                readonly
            ></v-text-field>
        </v-flex>
    </v-layout>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { File } from '@/api/file';
import { organizer } from '@/store/modules/organizer';

@Component
export default class FileInfo extends Vue {
    @organizer.State
    public currentFileIndex!: number;

    @organizer.State
    public currentFile!: File;

    @organizer.State
    public numFiles!: number;

    get infos(): Array<{ label: string; value: string; columns: string }> {
        return [
            {
                label: 'Position',
                value: `${this.currentFileIndex + 1}/${this.numFiles}`,
                columns: '2',
            },
            { label: 'File name', value: this.currentFile.name, columns: '4' },
            { label: 'Directory', value: this.currentFile.dir, columns: '5' },
            { label: 'Size', value: this.currentFileSize, columns: '1' },
        ];
    }

    get currentFileSize(): string {
        const units = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
        let i = 0;
        let size = this.currentFile.size;
        while (size >= 1024) {
            size /= 1024;
            i++;
        }
        return `${Math.round(size)} ${units[i]}`;
    }
}
</script>
