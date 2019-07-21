<template>
    <div class="container" v-if="showContainer">
        <FileInfo></FileInfo>
        <FilePreview></FilePreview>
        <HotkeyButtons></HotkeyButtons>
        <LeaveDialog
            :show="showLeaveDialog"
            @cancel="cancelLeaveDialog"
            @discard="leaveAndDiscard"
            @complete="leaveAndComplete"
        ></LeaveDialog>
    </div>
    <v-container v-else>
        <FinalizeDialog
            :show="true"
            :action="finalizeAction"
            :location="finalizeLocation"
        ></FinalizeDialog>
    </v-container>
</template>


<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import store from '@/store/store';
import { Route, Location } from 'vue-router';

import { organizer } from '@/store/modules/organizer';
import { Routes } from '@/router';
import FileInfo from '@/components/organize/FileInfo.vue';
import FilePreview from '@/components/organize/FilePreview.vue';
import HotkeyButtons from '@/components/organize/HotkeyButtons.vue';
import FinalizeDialog from '@/components/organize/FinalizeDialog.vue';
import LeaveDialog from '@/components/organize/LeaveDialog.vue';

@Component({
    components: {
        FileInfo,
        FilePreview,
        HotkeyButtons,
        FinalizeDialog,
        LeaveDialog,
    },
})
export default class Organize extends Vue {
    @organizer.Getter
    public isActive!: boolean;

    @organizer.Getter
    public hasCurrentFile!: boolean;

    @organizer.Action
    public handleHotkey!: (hotkey: string) => Promise<void>;

    @organizer.Action
    public dropConfigWait!: () => Promise<void>;

    @organizer.Action
    public dropConfig!: () => Promise<void>;

    public showLeaveDialog: boolean = false;

    public isLeaving: boolean = false;

    public finalizeAction: () => Promise<void> = this.dropConfigWait;

    public finalizeLocation: Location = { name: Routes.Home };

    public beforeRouteEnter(to: Route, from: Route, next: any) {
        const isOrganizerActive: boolean = store.getters['organizer/isActive'];
        if (isOrganizerActive) {
            next();
        } else {
            next({ name: Routes.ConfigNew });
        }
    }

    public beforeRouteLeave(to: Route, from: Route, next: any) {
        if (!this.isActive) {
            next();
            return;
        }

        this.finalizeLocation = { path: to.path };
        this.showLeaveDialog = true;
    }

    public mounted() {
        window.addEventListener('keypress', this.handleKeypress);
    }

    public beforeDestroy() {
        window.removeEventListener('keypress', this.handleKeypress);
    }

    public handleKeypress(e: KeyboardEvent) {
        const hotkey = String.fromCharCode(e.keyCode);
        this.handleHotkey(hotkey);
    }

    public get showContainer(): boolean {
        return this.hasCurrentFile && !this.isLeaving;
    }

    public cancelLeaveDialog() {
        this.showLeaveDialog = false;
        this.finalizeLocation = { name: Routes.Home };
    }

    public leaveAndDiscard() {
        this.showLeaveDialog = false;
        this.finalizeAction = this.dropConfig;
        this.isLeaving = true;
    }

    public leaveAndComplete() {
        this.showLeaveDialog = false;
        this.finalizeAction = this.dropConfigWait;
        this.isLeaving = true;
    }
}
</script>

<style lang="scss" scoped>
.container {
    // Fill height
    height: 100vh;
    // Fill width
    flex: 1;
    // Contain flex items in a column
    display: flex;
    flex-direction: column;
    // Emulate vuetify's `$grid-gutters.xl` container padding
    padding: 24px;
}
</style>
