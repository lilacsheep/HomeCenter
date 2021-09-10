<template>
  <a-layout-content style="padding: 12px;">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          用户管理
        </a-breadcrumb-item>
      </a-breadcrumb>
    <a-row :gutter="20">
      <a-col :span="24">
        <div id="terminal" class="xterm" />
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
    }
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
    this.endpoint = `${this.protocol}${window.location.host}${window.location.port}/api/system/webssh`
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