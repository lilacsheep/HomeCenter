<template>
  <a-layout-content style="padding: 12px;">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          用户管理
        </a-breadcrumb-item>
      </a-breadcrumb>
    <a-row :gutter="20" >
      <a-col :span="6">
        <a-button-group size="small">
          <a-button>新增组</a-button>
          <a-button>新增主机</a-button>
        </a-button-group>
        <a-tree  style="background: #FFFFFF;height: 100%;margin: 5px;" :tree-data="servergroup" :load-data="load_server" @select="server_select" show-icon>
          <a-icon slot="desktop" type="desktop" />
          <a-icon slot="folder" type="folder" />
        </a-tree>
      </a-col>
      <a-col :span="18" style="height: 100%;">
        <a-card :title="title" size="small" style="height: 100%;" :bodyStyle="{padding: 0}">
           <div id="xterm"></div>
        </a-card>
      </a-col>
    </a-row>
  </a-layout-content>
</template>

<script>
import {Terminal } from 'xterm'
import * as fit from 'xterm/lib/addons/fit/fit'
import * as attach from 'xterm/lib/addons/attach/attach'
import 'xterm/dist/xterm.css'
Terminal.applyAddon(attach)
Terminal.applyAddon(fit)

export default {
  data() {
    return {
      servergroup: [],
      term: null,
      endpoint: null,
      connection: null,
      content: '',
      protocol: null,
      host: 1,
      title: "默认",
    }
  },
  methods: {
    refresh_tree: function() {
      let data = [], that = this;
      this.$webssh.group.list(9999).then(function(response) {
        response.detail.forEach(function(item) {
          data.push({title: item.name,key: item.id,slots: {icon: 'folder'}, children: []})
        })
        that.servergroup = data
      }).catch(function(response) {
        that.$message.error(response.message)
      })
    },
    load_server: function(treeNode) {
      return new Promise(resolve => {
        this.$webssh.server.list(treeNode.dataRef.key).then(function(response) {
          treeNode.dataRef.children = []
          response.detail.forEach(function(item) {
            treeNode.dataRef.children.push({ title: item.name, key: `host-${item.id}`, slots: {icon: 'desktop'}, isLeaf: true})
          })
        }).catch(function(response) {
            that.$message.error("获取服务器信息失败："+response.message)
        })
        resolve();
      });
    },
    server_select(selectedKeys, event) {
      if ((selectedKeys.length != 0) && (selectedKeys[0].startsWith("host"))) {
        if ((this.term) && (this.connection)) {
          let that = this
          
          this.$confirm({
            content: '检测到当前已经存在连接是否继续',
            onOk() {
              return new Promise((resolve, reject) => {
                let node = event.selectedNodes[0]
                console.log(node)
                that.connection.close()
                that.term.destroy()
                that.host = node.key
                that.title = node.data.props.title
                that.init_term()
                resolve()
              }).catch(() => console.log('Oops errors!'));
            },
            cancelText: '取消',
            onCancel() {
              that.$destroyAll()
              return
            },
          });
        } else {
          let node = event.selectedNodes[0]
          this.host = node.key
          this.title = node.data.props.title
          this.init_term()
        }
      }
    },
    onOpen() {
      this.connection.send(JSON.stringify({type: "connect", cols: this.term.cols, rows: this.term.rows, host: this.host}))
    },
    onclose() {
      this.$message.error("连接中断")
    },
    onerror(error) {
       this.$message.error("连接中断: "+error)
    },
    init_term() {
        if (window.location.protocol === 'https:') {
            this.protocol = 'wss://'
        } else {
            this.protocol = 'ws://'
        }
        this.endpoint = `${this.protocol}127.0.0.1:8081/api/system/webssh`
        // obj.endpoint = `${obj.protocol}${window.location.host}/api/system/webssh`

        const terminalContainer = document.getElementById("xterm")
        this.term = new Terminal({
            cursorBlink: true,
        })
        this.term.open(terminalContainer, true)
        this.term.write('Connecting...')
        if (window.WebSocket) {
            // 如果支持websocket
            let ws = new WebSocket(this.endpoint)//后端接口位置
            this.connection = ws
            } else {
            // 否则报错
            console.log('WebSocket Not Supported' + this.endpoint)
        }

        this.connection.onopen = this.onOpen
        this.connection.onclose = this.onclose
        this.connection.onerror = this.onerror
        this.term.attach(this.connection)
    }
  },
  created: function () {
    this.refresh_tree()
  },
  beforeDestroy() {
    this.connection.close()
    this.term.destroy()
  },
  mounted: function () {}
};
</script>

<style>
.ant-tabs-bar {
  margin: 0 0 5px 0;
}

</style>