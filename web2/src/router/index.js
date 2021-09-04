import Router from 'vue-router'
import Herader from '../components/header.vue'
import Dashboard from '../components/website/index.vue'
import AgentIndex from '../components/agent/index.vue'
import Login from '../components/login/login.vue'
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
        props: true,
        children:[
          {
            path: '/dashboard',
            name: 'dashboard',
            component: Dashboard,
            props: true
          },
          {
            path: '/download',
            name: 'download',
            component: Download,
            props: true
          },
          {
            path: '/agent',
            name: 'agent',
            component: AgentIndex,
            props: true
          },
        ],
      },
      {
        path: '/login',
        component: Login,
      }
    ]
  })

export default router


