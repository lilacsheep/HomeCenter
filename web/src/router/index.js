import Router from 'vue-router'
import Herader from '../components/header.vue'
import Dashboard from '../components/index.vue'
import Roles from '../components/roles/index.vue'
import Monitor from '../components/monitor/index.vue'
import Download from '../components/download/index.vue'
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
            path: '/roles',
            name: 'roles',
            component: Roles,
            meta: {
              title: '规则配置'
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
            path: '/monitor',
            name: 'monitor',
            component: Monitor,
            meta: {
              title: '应用监控'
            }
          },
        ],
      }
    ]
  })

export default router


