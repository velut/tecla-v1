<template>
    <object :data="`${file.url}#toolbar=0&view=Fit`" type="application/pdf" />
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { File } from '@/api/file';

@Component
export default class PreviewPdf extends Vue {
    @Prop()
    public file!: File;
}
</script>

<style lang="scss" scoped>
object {
    width: 100%;
    height: 100%;
}
</style>
