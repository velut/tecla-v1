<template>
    <v-container fluid>
        <div class="display-1">About Tecla</div>
        <v-container v-if="info && credits">
            <div class="subheading">
                <p>
                    {{ info.name }}
                    <br />
                    {{ info.description }}
                </p>
                <p>
                    Version: {{ info.version }}
                    <br />
                    Commit: {{ info.commit }}
                </p>
                <p>
                    Homepage:
                    <a
                        :href="info.homepage"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        {{ info.homepage }}
                    </a>
                </p>
                <p>
                    {{ info.copyright }}
                </p>
                <p>{{ info.warranty }}; for details see the license.</p>
                <v-expansion-panel>
                    <v-expansion-panel-content lazy>
                        <template v-slot:header>
                            <div>Show license</div>
                        </template>
                        <v-card>
                            <v-card-text>
                                <v-textarea
                                    :value="info.license"
                                    auto-grow
                                    disabled
                                    box
                                >
                                </v-textarea>
                            </v-card-text>
                        </v-card>
                    </v-expansion-panel-content>
                </v-expansion-panel>

                <p class="mt-5">
                    Tecla includes the following software
                </p>
                <v-expansion-panel>
                    <v-expansion-panel-content
                        v-for="credit in credits.credits"
                        :key="credit.name"
                        lazy
                    >
                        <template v-slot:header>
                            <div>{{ credit.name }}</div>
                        </template>
                        <v-card>
                            <v-card-text>
                                <p>
                                    Homepage:
                                    <a
                                        :href="credit.homepage"
                                        target="_blank"
                                        rel="noopener noreferrer"
                                    >
                                        {{ credit.homepage }}
                                    </a>
                                </p>
                                <p>
                                    License:
                                </p>
                                <v-textarea
                                    :value="credit.license"
                                    auto-grow
                                    disabled
                                    box
                                >
                                </v-textarea>
                            </v-card-text>
                        </v-card>
                    </v-expansion-panel-content>
                </v-expansion-panel>
            </div>
        </v-container>
    </v-container>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import { Info, Credits } from '@/api/info';
import { appInfoAPI } from '@/api/api';

@Component
export default class About extends Vue {
    private info: Info | null = null;
    private credits: Credits | null = null;

    public beforeCreate() {
        appInfoAPI.appInfo().then((res) => (this.info = res));
        appInfoAPI.appCredits().then((res) => (this.credits = res));
    }
}
</script>
