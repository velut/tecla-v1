import Vue from 'vue';
import Vuex from 'vuex';

import organizer from '@/store/modules/organizer';

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        organizer,
    },
    state: {},
    mutations: {},
    actions: {},
});
