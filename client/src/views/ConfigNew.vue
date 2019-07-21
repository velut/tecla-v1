<template>
    <v-container fluid>
        <div class="display-1">
            New configuration
        </div>
        <ConfigForm></ConfigForm>
    </v-container>
</template>


<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import ConfigForm from '@/components/config/ConfigForm.vue';

@Component({
    components: {
        ConfigForm,
    },
})
export default class ConfigNew extends Vue {}
</script>
