import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import router from './router'
import moment from 'moment'
import api from "./api/api"
import webssh from "./api/server"
import VueClipboard from 'vue-clipboard2'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';
import VCharts from 'v-charts'

Vue.config.productionTip = false;

axios.defaults.withCredentials = true;
Vue.prototype.$axios = axios;
Vue.prototype.$moment = moment;
Vue.prototype.$api = api;
Vue.prototype.$webssh = webssh
Vue.use(VueClipboard)
Vue.use(Antd)
Vue.use(VCharts)

Vue.filter('dateformat', (dataStr, pattern = 'YYYY-MM-DD HH:mm') =>{
  return moment(dataStr).format(pattern)
})

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

Vue.filter("formatSeconds", (value) => { 
    var theTime = parseInt(value);// 需要转换的时间秒 
    var theTime1 = 0;// 分 
    var theTime2 = 0;// 小时 
    var theTime3 = 0;// 天
    if(theTime > 60) { 
    theTime1 = parseInt(theTime/60); 
    theTime = parseInt(theTime%60); 
    if(theTime1 > 60) { 
      theTime2 = parseInt(theTime1/60); 
      theTime1 = parseInt(theTime1%60); 
      if(theTime2 > 24){
      //大于24小时
      theTime3 = parseInt(theTime2/24);
      theTime2 = parseInt(theTime2%24);
      }
    } 
    } 
    var result = '';
    if(theTime > 0){
      result = ""+parseInt(theTime)+"秒";
    }
    if(theTime1 > 0) { 
      result = ""+parseInt(theTime1)+"分"+result; 
    } 
    if(theTime2 > 0) { 
      result = ""+parseInt(theTime2)+"小时"+result; 
    } 
    if(theTime3 > 0) { 
      result = ""+parseInt(theTime3)+"天"+result; 
    }
    return result; 
 })
new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
