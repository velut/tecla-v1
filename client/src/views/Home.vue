<template>
    <v-container fluid>
        <div class="display-1">Home</div>
        <v-container v-if="info">
            <v-container pt-0>
                <v-alert :value="true" type="warning" outline>
                    <div class="subheading">
                        Tecla can move, copy, overwrite, and delete files.
                        <br />
                        <b>Always have a backup copy of your files!</b>
                    </div>
                </v-alert>
            </v-container>
            <div class="subheading">
                <p>
                    Welcome to Tecla, the interactive file organizer.
                </p>
                <p>
                    To start, click on
                    <i>New configuration</i>.
                </p>
            </div>

            <div class="subheading mt-5">
                <p>
                    {{ info.copyright }}
                </p>
                <p>
                    {{ info.warranty }}; for details see the license in
                    <i>About Tecla</i>.
                </p>
                <p>
                    By using this program, you confirm that you accept its
                    license.
                </p>
            </div>
        </v-container>
    </v-container>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import { Info } from '@/api/info';
import { appInfoAPI } from '@/api/api';
import appInfo from '@/info';

@Component
export default class Home extends Vue {
    private info: Info | null = null;

    public beforeCreate() {
        appInfoAPI.appInfo().then((res) => (this.info = res));
    }
}
</script>
