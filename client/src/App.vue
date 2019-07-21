<template>
    <v-app dark>
        <NavigationDrawer></NavigationDrawer>
        <v-content>
            <router-view></router-view>
        </v-content>
    </v-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import NavigationDrawer from '@/components/navigation/NavigationDrawer.vue';

@Component({
    components: {
        NavigationDrawer,
    },
})
export default class App extends Vue {}
</script>

<style lang="scss">
html,
body {
    width: 100%;
    height: 100%;
}

html {
    overflow-y: auto !important;
}
</style>
