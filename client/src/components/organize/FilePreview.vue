<template>
    <div class="preview-container">
        <PreviewImage
            v-if="fileKind === 'image'"
            :file="currentFile"
        ></PreviewImage>
        <PreviewMultimedia
            v-else-if="fileKind === 'multimedia'"
            :file="currentFile"
        ></PreviewMultimedia>
        <PreviewPdf
            v-else-if="fileKind === 'pdf'"
            :file="currentFile"
        ></PreviewPdf>
        <PreviewText
            v-else-if="fileKind === 'text' && isSmallFile"
            :file="currentFile"
        ></PreviewText>
        <PreviewDefault v-else :file="currentFile"></PreviewDefault>
    </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import PreviewImage from '@/components/organize/PreviewImage.vue';
import PreviewMultimedia from '@/components/organize/PreviewMultimedia.vue';
import PreviewPdf from '@/components/organize/PreviewPdf.vue';
import PreviewText from '@/components/organize/PreviewText.vue';
import PreviewDefault from '@/components/organize/PreviewDefault.vue';
import { File } from '@/api/file';
import { organizer } from '@/store/modules/organizer';

@Component({
    components: {
        PreviewImage,
        PreviewMultimedia,
        PreviewPdf,
        PreviewText,
        PreviewDefault,
    },
})
export default class FilePreview extends Vue {
    @organizer.State
    public currentFile!: File;

    get fileKind(): string {
        switch (this.currentFile.ext) {
            // Images
            case '.jpeg':
            case '.jpg':
            case '.png':
            case '.apng':
            case '.gif':
            case '.bmp':
            case '.ico':
            case '.webp':
            case '.svg':
                return 'image';
            // Multimedia, see https://www.chromium.org/audio-video
            case '.flac':
            case '.mp4':
            case '.m4a':
            case '.mp3':
            case '.ogv':
            case '.ogm':
            case '.ogg':
            case '.oga':
            case '.opus':
            case '.webm':
            case '.wav':
                // Video tag can be used for both audio and video.
                return 'multimedia';
            // Pdf
            case '.pdf':
                return 'pdf';
            // Plaintext
            case '.txt':
                return 'text';
            // Unknown
            default:
                return 'other';
        }
    }

    get isSmallFile(): boolean {
        // 1MB or smaller.
        return this.currentFile.size <= 1048576;
    }
}
</script>

<style lang="scss" scoped>
.preview-container {
    // Fill remaining parent's space
    flex: 1;
    // Resize contained items correctly
    display: flex;
    justify-content: center;
    align-items: center;
    // Shrink vertically
    min-height: 0;
}
</style>
