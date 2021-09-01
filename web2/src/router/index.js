import Router from 'vue-router'
import Herader from '../components/header.vue'
import WebIndex from '../components/website/index.vue'
import AgentIndex from '../components/agent/index.vue'
import Login from '../components/login/login.vue'
import Vue from 'vue'


Vue.use(Router)

const router = new Router({
    mode: 'history',
    routes: [
      // 动态路径参数 以冒号开头
      { 
        path: '/',
        redirect: 'web',
        component: Herader,
        props: true,
        children:[
          {
            path: '/web',
            name: 'web',
            component: WebIndex,
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


