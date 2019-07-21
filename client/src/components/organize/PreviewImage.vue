<template>
    <img :src="file.url" />
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { File } from '@/api/file';

@Component
export default class PreviewImage extends Vue {
    @Prop()
    public file!: File;
}
</script>

<style lang="scss" scoped>
img {
    max-width: 100%;
    max-height: 100%;
}
</style>
