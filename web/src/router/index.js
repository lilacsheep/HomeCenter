import Router from 'vue-router'
import Herader from '../components/header.vue'
import Dashboard from '../components/index.vue'
import Roles from '../components/roles/index.vue'
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
        ],
      }
    ]
  })

export default router


