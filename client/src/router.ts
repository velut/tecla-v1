import Vue from 'vue';
import Router from 'vue-router';

import Home from '@/views/Home.vue';
import ConfigNew from '@/views/ConfigNew.vue';
import Organize from '@/views/Organize.vue';
import About from '@/views/About.vue';

Vue.use(Router);

export enum Routes {
    Home = 'home',
    ConfigNew = 'configNew',
    Organize = 'organize',
    About = 'about',
}

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/',
            name: Routes.Home,
            component: Home,
        },
        {
            path: '/' + Routes.ConfigNew,
            name: Routes.ConfigNew,
            component: ConfigNew,
        },
        {
            path: '/' + Routes.Organize,
            name: Routes.Organize,
            component: Organize,
        },
        {
            path: '/' + Routes.About,
            name: Routes.About,
            component: About,
        },
        { path: '*', redirect: { name: Routes.Home } },
    ],
});
