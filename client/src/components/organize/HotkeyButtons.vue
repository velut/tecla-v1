<template>
    <v-layout shrink wrap justify-center class="layout">
        <v-flex md2 lg1>
            <v-tooltip top>
                <template v-slot:activator="{ on }">
                    <v-btn
                        v-on="on"
                        @click="nextFile"
                        class="text-none"
                        color="#616161"
                    >
                        Skip
                    </v-btn>
                </template>
                <span>Press spacebar to skip this file</span>
            </v-tooltip>
        </v-flex>
        <v-flex v-for="dstDir in dstDirs" :key="dstDir.hotkey" md2 lg1>
            <v-tooltip top>
                <template v-slot:activator="{ on }">
                    <v-btn
                        v-on="on"
                        dark
                        @click="handleHotkey(dstDir.hotkey)"
                        class="text-none"
                        color="#616161"
                    >
                        {{ dstDir.hotkey }}
                    </v-btn>
                </template>
                <span>Send to {{ dstDir.dir }}</span>
            </v-tooltip>
        </v-flex>
    </v-layout>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { DstDir } from '@/api/config';
import { organizer } from '@/store/modules/organizer';

@Component
export default class HotkeyButtons extends Vue {
    @organizer.Action
    public handleHotkey!: (hotkey: string) => Promise<void>;

    @organizer.Getter
    public dstDirs!: DstDir[];

    private nextFile() {
        this.handleHotkey(' ');
    }
}
</script>

<style lang="scss" scoped>
.layout {
    margin-top: 12px;
}

.v-btn {
    min-width: 60px;
}
</style>
