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
        <div id="xterm"></div>
      </a-col>
    </a-row>
  </a-layout-content>
  
</template>

<script>
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";
import "xterm/lib/xterm.js";

export default {
  data() {
    return {
    }
  },
  methods: {
    initXterm() {
      this.term = new Terminal({
        rendererType: "canvas", //渲染类型
        rows: 35, //行数
        convertEol: true, //启用时，光标将设置为下一行的开头
        scrollback: 10, //终端中的回滚量
        disableStdin: false, //是否应禁用输入
        cursorStyle: "underline", //光标样式
        cursorBlink: true, //光标闪烁
        theme: {
          foreground: "yellow", //字体
          background: "#060101", //背景色
          cursor: "help" //设置光标
        }
      });
      this.term.open(document.getElementById("xterm"));
      const fitAddon = new FitAddon();
      this.term.loadAddon(fitAddon);
      // 支持输入与粘贴方法
      let _this = this; //一定要重新定义一个this，不然this指向会出问题
      this.term.onData(function(key) {
        let order = ["stdin",key]; //这里key值是你输入的值，数据格式一定要找后端要！！！！
        _this.socket.onsend(JSON.stringify(order)); //转换为字符串
      });
    },
    init(url) {
      // 实例化socket
      this.socket = new WebSocket(url);
      // 监听socket连接
      this.socket.onopen = this.open;
      // 监听socket错误信息
      this.socket.onerror = this.error;
      // 监听socket消息
      this.socket.onmessage = this.getMessage;
      // 发送socket消息
      this.socket.onsend = this.send;
    },
    open: function() {
      console.log("socket连接成功");
      this.initXterm();
    },
    error: function() {
      console.log("连接错误");
    },
    close: function() {
      this.socket.close();
      console.log("socket已经关闭");
    },
    getMessage: function(msg) {
      this.term.write(JSON.parse(msg.data)[1]);
    },
    send: function(order) {
      this.socket.send(order);
    }
  },
  
  created: function () {},
  beforeDestroy () {},
  mounted: function () {}
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