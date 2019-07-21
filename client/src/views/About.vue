<template>
    <v-container fluid>
        <div class="display-1">About Tecla</div>
        <v-container>
            <div class="subheading">
                <p>
                    {{ info.name }}
                    <br />
                    {{ info.description }}
                </p>
                <p>Version: {{ info.version }}</p>
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
                        v-for="credit in credits"
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

import info from '@/info';
import credits from '@/credits';

@Component
export default class About extends Vue {
    private info = info;

    private credits = credits;
}
</script>
