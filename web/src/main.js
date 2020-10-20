import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import router from './router'
import api from "./api/api"
import VueClipboard from 'vue-clipboard2'
import ElementUI from 'element-ui';
import VeLine from 'v-charts/lib/line.common'
import 'element-ui/lib/theme-chalk/index.css';
import numerify from 'numerify'
import numerifyBytes from 'numerify/lib/plugins/bytes.umd.js'


Vue.config.productionTip = false;

axios.defaults.withCredentials = true;
Vue.prototype.$axios = axios;
Vue.prototype.$api = api;
Vue.use(VueClipboard)
Vue.use(ElementUI)
Vue.component(VeLine.name, VeLine)
numerify.register('bytes', numerifyBytes)

Vue.filter('hideStr', (src, endOf) => {
  if (src.length > endOf) {
    return `${src.substr(0, endOf)}...`
  }
  return src
})

Vue.filter('diskSize', (num) => {
  if (num === 0) return '0 B';
  let k = 1024; //设定基础容量大小
  let sizeStr = ['B','KB','MB','GB','TB','PB','EB','ZB','YB']; //容量单位
  let i = 0; //单位下标和次幂
  for(let l=0;l<8;l++){   //因为只有8个单位所以循环八次
      if(num / Math.pow(k, l) < 1){ //判断传入数值 除以 基础大小的次幂 是否小于1，这里小于1 就代表已经当前下标的单位已经不合适了所以跳出循环
          break; //小于1跳出循环
      }
      i = l; //不小于1的话这个单位就合适或者还要大于这个单位 接着循环
  }
  return (num / Math.pow(k, i)).toFixed(2) + ' ' + sizeStr[i];  //循环结束 或 条件成立 返回字符
})

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
