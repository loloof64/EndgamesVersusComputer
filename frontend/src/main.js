import 'core-js/stable';
import 'regenerator-runtime/runtime';
import '@mdi/font/css/materialdesignicons.css';
import Vue from 'vue';
import Vuetify from 'vuetify';
import 'vuetify/dist/vuetify.min.css';

import VueI18n from 'vue-i18n'

Vue.use(VueI18n);
Vue.use(Vuetify);

import App from './App.vue';

Vue.config.productionTip = false;
Vue.config.devtools = true;

import Wails from '@wailsapp/runtime';
import i18n from './i18n';

require('../node_modules/@loloof64/chessboard-component/dist/index');

Wails.Init(() => {
	new Vue({
		vuetify: new Vuetify({
			icons: {
				iconfont: 'mdi'
			},
			theme: {
				dark: true
			}
		}),
		i18n,
		render: h => h(App)
	}).$mount('#app');
});
