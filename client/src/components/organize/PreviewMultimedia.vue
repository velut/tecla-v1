<template>
    <video controls autoplay loop muted :src="file.url">
        This browser does not support video/audio content.
    </video>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { File } from '@/api/file';

@Component
export default class PreviewMultimedia extends Vue {
    @Prop()
    public file!: File;
}
</script>

<style lang="scss" scoped>
video {
    max-width: 100%;
    max-height: 100%;
}
</style>
