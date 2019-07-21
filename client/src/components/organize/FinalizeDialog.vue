<template>
    <v-dialog :value="show" persistent max-width="450">
        <v-card>
            <v-card-title class="headline" primary-title>
                Finalizing file operations
            </v-card-title>

            <v-card-text>
                <div class="subheading">
                    Please wait for file operations to complete.
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import { Location } from 'vue-router';

import { File } from '@/api/file';
import { Routes } from '@/router';
import { organizer } from '@/store/modules/organizer';

@Component
export default class FinalizeDialog extends Vue {
    @Prop()
    public show!: boolean;

    @Prop()
    public action!: () => Promise<void>;

    @Prop()
    public location!: Location;

    @organizer.Getter
    public isActive!: boolean;

    public mounted() {
        if (!this.isActive) {
            return;
        }

        this.action().then(() => {
            // Display dialog for some time.
            setTimeout(() => {
                this.$router.push(this.location);
            }, 750);
        });
    }
}
</script>
