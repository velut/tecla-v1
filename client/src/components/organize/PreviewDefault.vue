<template>
    <v-sheet dark class="pa-5" elevation="6">
        <div class="title text-xs-center pa-2">
            Preview not available for
        </div>
        <div class="headline text-xs-center pa-2">
            {{ file.name }}
        </div>
    </v-sheet>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { File } from '@/api/file';

@Component
export default class PreviewDefault extends Vue {
    @Prop()
    public file!: File;
}
</script>
