import Router from 'vue-router'
import Herader from '../components/header.vue'
import Dashboard from '../components/website/index.vue'
import AgentIndex from '../components/agent/index.vue'
import Login from '../components/login/login.vue'
import Download from '../components/download/index.vue'
import Users from '../components/user/index.vue'
import DDNS from '../components/ddns/index.vue'
import Containers from '../components/docker/containers.vue'
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
            path: '/ddns',
            name: 'ddns',
            component: DDNS,
            props: true
          },
          {
            path: '/users',
            name: 'users',
            component: Users,
            props: true
          },
          {
            path: '/containers',
            name: 'containers',
            component: Containers,
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


