<template>
    <v-navigation-drawer app permanent>
        <v-list>
            <template v-for="(group, groupID) in itemGroups">
                <v-list-tile
                    v-for="item in group"
                    :key="`${groupID}-${item.text}`"
                    :to="item.to"
                    exact
                    active-class="accent--text"
                >
                    <v-list-tile-action>
                        <v-icon>
                            {{ item.icon }}
                        </v-icon>
                    </v-list-tile-action>
                    <v-list-tile-content>
                        <v-list-tile-title>
                            {{ item.text }}
                        </v-list-tile-title>
                    </v-list-tile-content>
                </v-list-tile>
                <v-divider
                    :key="`${groupID}-divider`"
                    v-if="groupID !== 'about'"
                ></v-divider>
            </template>
        </v-list>
    </v-navigation-drawer>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Location } from 'vue-router';

export interface NavigationDrawerItemsByID {
    [id: string]: NavigationDrawerItem[];
}

export interface NavigationDrawerItem {
    icon: string;
    text: string;
    to: Location;
}

@Component
export default class NavigationDrawer extends Vue {
    public itemGroups: NavigationDrawerItemsByID = {
        home: [
            {
                icon: 'home',
                text: 'Home',
                to: { name: 'home' },
            },
        ],
        config: [
            {
                icon: 'note_add',
                text: 'New configuration',
                to: { name: 'configNew' },
            },
        ],
        organize: [
            {
                icon: 'view_carousel',
                text: 'Organize',
                to: { name: 'organize' },
            },
        ],
        about: [
            {
                icon: 'info',
                text: 'About Tecla',
                to: { name: 'about' },
            },
        ],
    };
}
</script>
