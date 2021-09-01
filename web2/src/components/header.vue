<template>
  <a-layout id="components-layout-demo-custom-trigger">
    <a-layout-sider v-model="collapsed" :trigger="null" collapsible>
      <img :src="imgUrl" class="logo" :style="logo_style"/>
      <a-menu theme="dark" mode="inline" :default-selected-keys="['1']">
        <a-menu-item key="1">
          <router-link :to="{path:'/web'}">
            <a-icon type="dashboard" />
            <span>控制面板</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="2">
          <router-link :to="{path:'/agent'}">
            <a-icon type="bars" />
            <span>节点操作</span>
          </router-link>
        </a-menu-item>
        <!-- <a-sub-menu key="3">
          <span slot="title"><a-icon type="solution" /><span>权限管理</span></span>
          <a-menu-item-group key="31">
            <template slot="title"> <a-icon type="team" /><span>用户</span> </template>
            <a-menu-item key="32">
              用户管理
            </a-menu-item>
            <a-menu-item key="33">
              权限管理
            </a-menu-item>
          </a-menu-item-group>
        </a-sub-menu> -->
      </a-menu>
    </a-layout-sider>
    <a-layout>
      <a-layout-header style="background: #fff; padding: 0">
        <a-icon
          class="trigger"
          :type="collapsed ? 'menu-unfold' : 'menu-fold'"
          @click="hideMenuClick"
        />
      </a-layout-header>
        <router-view></router-view>
    </a-layout>
  </a-layout>
</template>
<script>
export default {
  data() {
    return {
      collapsed: false,
      imgUrl:require("../assets/logo.png"),
      logo_style: "",
    };
  },
  methods: {
    hideMenuClick: function () {
      this.collapsed = !this.collapsed
      this.collapsed ? this.logo_style = "display: none" : this.logo_style = ""
      this.collapsed ? this.imgUrl = "" : this.imgUrl = require("../assets/logo.png")
    }
  },
};
</script>
<style>
.ant-layout {
  height: 100%;
}

.ant-layout-header {
  height: 48px!important;
  line-height: 48px!important;
}

#components-layout-demo-custom-trigger {
  min-width: 1080px;
}

#components-layout-demo-custom-trigger .trigger {
  font-size: 18px;
  line-height: 32px;
  padding: 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

#components-layout-demo-custom-trigger .trigger:hover {
  color: #1890ff;
}

#components-layout-demo-custom-trigger .logo {
  height: 48px;
  margin: 5px 25px;
}
</style>
