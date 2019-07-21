<template>
    <v-dialog :value="show" persistent max-width="850">
        <v-card>
            <v-card-title class="headline" primary-title>
                Stop organizing files?
            </v-card-title>

            <v-card-text>
                <div class="subheading">
                    File operations already in progress will be completed before
                    leaving.
                    <br />
                    Pending file operations can be either discarded or
                    completed.
                </div>
            </v-card-text>
            <v-divider></v-divider>

            <v-card-actions>
                <v-spacer></v-spacer>

                <v-btn color="accent" flat="flat" @click="cancel">
                    Cancel
                </v-btn>

                <v-btn color="accent" flat="flat" @click="discard">
                    Discard pending operations
                </v-btn>

                <v-btn color="accent" flat="flat" @click="complete">
                    Complete all operations
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from 'vue-property-decorator';

@Component
export default class LeaveDialog extends Vue {
    @Prop()
    public show!: boolean;

    @Emit('cancel')
    private cancel() {
        return;
    }

    @Emit('discard')
    private discard() {
        return;
    }

    @Emit('complete')
    private complete() {
        return;
    }
}
</script>
