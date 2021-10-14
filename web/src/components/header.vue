<template>
  <a-layout id="components-layout-demo-custom-trigger">
    <a-layout-sider v-model="collapsed" :trigger="null" collapsible>
      <img :src="imgUrl" class="logo" :style="logo_style"/>
      <a-menu theme="dark" :open-keys="openKeys" mode="inline" :default-selected-keys="['1']" :selectedKeys="selectedKeys" @select="select" @openChange="onChange">
        <a-menu-item key="1">
          <router-link :to="{path:'/dashboard'}">
            <a-icon type="dashboard" />
            <span>控制面板</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="2">
          <router-link :to="{path:'/ddns'}">
            <a-icon type="bars" />
            <span>动态域名</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="3">
          <router-link :to="{path:'/download'}">
            <a-icon type="download" />
            <span>离线下载</span>
          </router-link>
        </a-menu-item>
        <a-sub-menu key="4">
          <span slot="title"><a-icon type="appstore" /><span>Docker</span></span>
            <a-menu-item key="41">
              <router-link :to="{path:'/containers'}">
                <a-icon type="build" />
                <span>容器管理</span>
              </router-link>
            </a-menu-item>
            <a-menu-item key="42">
              <router-link :to="{path:'/container/template'}">
                <a-icon type="profile" />
                <span>容器模板</span>
              </router-link>
            </a-menu-item>
            <a-menu-item key="43">
              <router-link :to="{path:'/docker/images'}">
                <a-icon type="profile" />
                <span>镜像管理</span>
              </router-link>
            </a-menu-item>
        </a-sub-menu>
        <a-menu-item key="WebSSH">
          <router-link :to="{path:'/webssh'}">
            <a-icon type="code" />
            <span>网页SSH</span>
          </router-link>
        </a-menu-item>
        <a-menu-item key="minio">
          <router-link :to="{path:'/minio'}">
            <a-icon type="file" />
            <span>对象存储</span>
          </router-link>
        </a-menu-item>
        <a-sub-menu key="5">
          <span slot="title"><a-icon type="solution" /><span>权限管理</span></span>
          <a-menu-item-group key="50">
            <template slot="title"> <a-icon type="team" /><span>用户</span> </template>
            <a-menu-item key="51">
              <router-link :to="{path:'/users'}">
                <span>用户管理</span>
              </router-link>
            </a-menu-item>
          </a-menu-item-group>
        </a-sub-menu>
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
      selectedKeys: ["1"],
      rootSubmenuKeys: ['4', '5'],
      openKeys: []
    };
  },
  methods: {
    hideMenuClick: function () {
      this.collapsed = !this.collapsed
      this.collapsed ? this.logo_style = "display: none" : this.logo_style = ""
      this.collapsed ? this.imgUrl = "" : this.imgUrl = require("../assets/logo.png")
    },
    select: function(item) {
      this.selectedKeys = item.selectedKeys
    },
    onChange(openKeys) {
      const latestOpenKey = openKeys.find(key => this.openKeys.indexOf(key) === -1);
      if (this.rootSubmenuKeys.indexOf(latestOpenKey) === -1) {
        this.openKeys = openKeys;
      } else {
        this.openKeys = latestOpenKey ? [latestOpenKey] : [];
      }
    }
  },
  created: function() {
    var params = {
      "/dashboard": "1", "/ddns":"2", 
      "/download":"3", "/users":"51", 
      "/webssh": "WebSSH", 
      "/containers": "41",
      "/container/template": "42",
      "/docker/images": "43",
      "/users": "51",
      "/minio": "minio"
    }
    
    let key = params[this.$route.path]
    if (key !== "") {
      this.selectedKeys = [key]
      if (this.selectedKeys == "41") {
        this.openKeys = ["4"]
      } else if (this.selectedKeys == "42") {
        this.openKeys = ["4"]
      } else if (this.selectedKeys == "43") {
        this.openKeys = ["4"]
      } else if (this.selectedKeys == "51"){
        this.openKeys = ["5"]
      } else {
        this.openKeys = []
      }
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
}

.ant-table-thead > tr >th{
  color: white;
  background: #5069d6 !important;
}

.ant-table-tbody > tr >th{
  color: white;
  background: #5069d6 !important;
}

.ant-table-small > .ant-table-content > .ant-table-body {
    margin: 0 0px;
}
</style>
