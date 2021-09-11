<template>
  <a-layout-content style="padding: 12px;">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          用户管理
        </a-breadcrumb-item>
      </a-breadcrumb>
    <a-row :gutter="20">
      <a-col :span="6">
        <a-tree :tree-data="treeData" show-icon default-expand-all :default-selected-keys="['0-0-0']" style="background: #FFFFFF">
          <a-icon slot="switcherIcon" type="down" />
          <a-icon slot="smile" type="smile-o" />
          <a-icon slot="meh" type="smile-o" />
          <template slot="custom" slot-scope="{ selected }">
            <a-icon :type="selected ? 'frown' : 'frown-o'" />
          </template>
        </a-tree>
      </a-col>
      <a-col :span="16">
        <a-tabs v-model="activeKey" hide-add type="editable-card" @edit="onEdit">
          <a-tab-pane v-for="pane in panes" :key="pane.key" :tab="pane.title" :closable="pane.closable">
            <div id="terminal" class="xterm" />
          </a-tab-pane>
        </a-tabs>
        
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
    const panes = []
    return {
      treeData: [
        {
          title: 'parent 1',
          key: '0-0',
          slots: {
            icon: 'smile',
          },
          children: [
            { title: 'leaf', key: '0-0-0', slots: { icon: 'meh' } },
            { title: 'leaf', key: '0-0-1', scopedSlots: { icon: 'custom' } },
          ],
        },
      ],
      panes,
      newTabIndex: 0,
      term: null,
        endpoint: null,
        connection: null,
        content: '',
        protocol: null,
        option: {
          operate: 'connect',
          host: '192.168.*.*',//你要连接的终端的ip
          port: '22',
          username: '*',//你要连接的终端的用户名和密码
          password: '*'
        }
    }
  },
  methods: {
    send(data) {
      this.connection.send(JSON.stringify(data))
    },
    onError(error) {
      // 连接失败回调
      this.term.write('Error: ' + error + '\r\n')
      console.log('Error: ' + error + '\r\n')
    },
    onConnect() {
    //   this.connection.send(JSON.stringify(this.option))
    },
    onClose() {
      // 连接关闭回调
      this.term.write('\rconnection closed')
      console.log('\rconnection closed')
    },
    onData(data) {
      // 收到数据时回调
      this.term.write(data)
      console.log(data)
    },
    callback(key) {
      console.log(key);
    },
    onEdit(targetKey, action) {
      this[action](targetKey);
    },
    add() {
      const panes = this.panes;
      const activeKey = `newTab${this.newTabIndex++}`;
      panes.push({
        title: `New Tab ${activeKey}`,
        content: `Content of new Tab ${activeKey}`,
        key: activeKey,
      });
      this.panes = panes;
      this.activeKey = activeKey;
    },
    remove(targetKey) {
      let activeKey = this.activeKey;
      let lastIndex;
      this.panes.forEach((pane, i) => {
        if (pane.key === targetKey) {
          lastIndex = i - 1;
        }
      });
      const panes = this.panes.filter(pane => pane.key !== targetKey);
      if (panes.length && activeKey === targetKey) {
        if (lastIndex >= 0) {
          activeKey = panes[lastIndex].key;
        } else {
          activeKey = panes[0].key;
        }
      }
      this.panes = panes;
      this.activeKey = activeKey;
    },
  },
  created: function () {},
  beforeDestroy() {
    this.connection.close()
    this.term.destroy()
  },
  mounted: function () {
    if (window.location.protocol === 'https:') {
      this.protocol = 'wss://'
    } else {
      this.protocol = 'ws://'
    }
    this.endpoint = `${this.protocol}${window.location.host}/api/system/webssh`
    const terminalContainer = document.getElementById('terminal')
    this.term = new Terminal({
      // 光标闪烁
      cursorBlink: true
    })
    // this.term.on('data', (data) => {
    //   // 键盘输入时的回调函数
    //   this.connection.send(data)
    // })
    this.term.open(terminalContainer, true)
    this.term.write('Connecting...')
    if (window.WebSocket) {
      // 如果支持websocket
      console.log(this.endpoint)
      this.connection = new WebSocket(this.endpoint)//后端接口位置
    } else {
      // 否则报错
      this.onError('WebSocket Not Supported')
    }
    this.connection.onopen = this.onConnect
    this.connection.onclose = this.onClose
    this.connection.onerror = this.onError
    this.term.attach(this.connection)
  }
};
</script>

<style>
.a-card__header {
  padding: 5px;
}

.a-card__body {
  padding: 20px;
}

.a-dialog__header {
  padding: 10px 10px 5px;
  border-bottom: 1px solid whitesmoke;
}

.a-dialog__headerbtn {
  top: 12px;
}
.a-dialog__body {
  padding: 15px 10px;
}
.a-dialog__footer {
  border-top: 1px solid whitesmoke;
  padding: 5px 10px 10px;
}

.a-drawer__body {
  padding: 10px;
}

.descriptions  {
  width: 100%;
  margin-bottom: 10px;
}

.descriptions .title {
  background: #fafafa;
  border: 1px solid #e8e8e8;
  padding: 5px;
  font-size: 14px;
  font-weight: 400;
  line-height: 1.5;
  text-align: left;
}

.descriptions .details {
  border: 1px solid #e8e8e8;
  padding: 5px;
}
</style>