import Router from 'vue-router'
import Herader from '../components/header.vue'
import Dashboard from '../components/index.vue'
import Monitor from '../components/monitor/index.vue'
import Download from '../components/download/index.vue'
import Filesystem from '../components/filesystem/index.vue'
import Other from '../components/other/index.vue'
import Login from '../components/login/index.vue'
import Vue from 'vue'


Vue.use(Router)

const router = new Router({
    mode: 'history',
    routes: [
      // 动态路径参数 以冒号开头
      { 
        path: '/',
        redirect: 'dashboard',
        component: Herader,
        children:[
          {
            path: '/dashboard',
            name: 'dashboard',
            component: Dashboard,
            meta: {
              title: '面板'
            }
          },
          {
            path: '/download',
            name: 'download',
            component: Download,
            meta: {
              title: '资源下载'
            }
          },
          {
            path: '/filesystem',
            name: 'filesystem',
            component: Filesystem,
            meta: {
              title: '资源管理'
            }
          },
          {
            path: '/monitor',
            name: 'monitor',
            component: Monitor,
            meta: {
              title: '应用监控'
            }
          },
          {
            path: '/other',
            name: 'other',
            component: Other,
            meta: {
              title: '其他功能'
            }
          },
          {
            path: '/login',
            name: 'login',
            component: Login,
            meta: {
              title: '登录'
            }
          },
        ],
      }
    ]
  })

export default router


