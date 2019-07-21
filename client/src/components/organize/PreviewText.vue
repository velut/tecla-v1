<template>
    <object :data="file.url" type="text/plain" />
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { File } from '@/api/file';

@Component
export default class PreviewText extends Vue {
    @Prop()
    public file!: File;
}
</script>

<style lang="scss" scoped>
object {
    width: 80%;
    height: 100%;
    background-color: white;
}
</style>
