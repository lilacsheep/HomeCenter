import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import router from './router'
import api from "./api/api"
import VueClipboard from 'vue-clipboard2'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';


Vue.config.productionTip = false;

axios.defaults.withCredentials = true;
Vue.prototype.$axios = axios;
Vue.prototype.$api = api;
Vue.use(VueClipboard)
Vue.use(ElementUI)

Vue.filter('hideStr', (src, endOf) => {
  if (src.length > endOf) {
    return `${src.substr(0, endOf)}...`
  }
  return src
})

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
