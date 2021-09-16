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
        <a-button type="primary" @click="form.group.add.visible=true" block>新增组</a-button>
        <a-tree style="background: #FFFFFF;height: 100%;margin: 0;" :tree-data="servergroup" :load-data="load_server" @select="server_select" :expandedKeys.sync="expandedKeys" show-icon>
          <template #title="{ key: treeKey, title }" >
            <a-dropdown :trigger="['contextmenu']">
              <span>{{ title }}</span>
              <template #overlay>
                <a-menu @click="({ key: menuKey }) => onContextMenuClick(treeKey, menuKey)">
                  <a-menu-item key="add_host" v-if="check_key(treeKey)"><a-icon type="edit" />新增主机</a-menu-item>
                  <a-menu-item key="edit_host" v-if="!check_key(treeKey)"><a-icon type="edit" />编辑</a-menu-item>
                  <a-menu-item key="delete_host" v-if="!check_key(treeKey)"><a-icon type="delete" />删除</a-menu-item>
                  <a-menu-item key="delete_node" v-if="check_key(treeKey)"><a-icon type="delete" />删除</a-menu-item>
                  <a-menu-item key="edit_node" v-if="check_key(treeKey)"><a-icon type="edit" />编辑</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </template>
          <a-icon slot="desktop" type="desktop" />
          <a-icon slot="folder" type="folder" />
        </a-tree>
      </a-col>
      <a-col :span="18" style="height: 100%;">
        <a-card :title="title" size="small" style="height: 100%;" :bodyStyle="{padding: 0}">
           <div id="xterm"></div>
        </a-card>
      </a-col>
      <a-modal title="新增主机" :visible="form.server.add.visible" @cancel="form.server.add.visible=false" @ok="add_server">
        <a-form-model :model="form.server.add.data" :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-model-item label="名称">
            <a-input v-model="form.server.add.data.name" />
          </a-form-model-item>
          <a-form-model-item label="分组">
            <a-select v-model="form.server.add.data.group" placeholder="please select your zone">
              <a-select-option :key="item.id" v-for="item in groups" :value="item.id">{{item.name}}</a-select-option>
            </a-select>
          </a-form-model-item>
          <a-form-model-item label="地址">
            <a-input-group compact>
              <a-input v-model="form.server.add.data.address" style="width: 200px" />
              <a-input style="width: 30px; pointer-events: none; backgroundColor: #fff" placeholder=":" disabled/>
              <a-input type="number" style="width: 100px; text-align: center;" v-model="form.server.add.data.port"/>
            </a-input-group>
          </a-form-model-item>
          <a-form-model-item label="用户">
            <a-input v-model="form.server.add.data.username" />
          </a-form-model-item>
          <a-form-model-item label="认证">
            <a-switch checked-children="密码" un-checked-children="密钥" @change="auth_switch"/>
          </a-form-model-item>
          <a-form-model-item :style="form.server.add.password_style" label="密码">
            <a-input v-model="form.server.add.data.password" />
          </a-form-model-item>
          <a-form-model-item :style="form.server.add.private_key_style" label="密钥">
            <a-input v-model="form.server.add.data.private_key" type="textarea" />
          </a-form-model-item>
          <a-form-model-item label="备注">
            <a-input v-model="form.server.add.data.remark" type="textarea" />
          </a-form-model-item>
        </a-form-model>
      </a-modal>
      <a-modal title="编辑主机" :visible="form.server.edit.visible" @cancel="form.server.edit.visible=false" @ok="edit_server">
        <a-form-model :model="form.server.edit.data" :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-model-item label="名称">
            <a-input v-model="form.server.edit.data.name" />
          </a-form-model-item>
          <a-form-model-item label="分组">
            <a-select v-model="form.server.edit.data.group" placeholder="please select your zone">
              <a-select-option :key="item.id" v-for="item in groups" :value="item.id">{{item.name}}</a-select-option>
            </a-select>
          </a-form-model-item>
          <a-form-model-item label="地址">
            <a-input-group compact>
              <a-input v-model="form.server.edit.data.address" style="width: 200px" />
              <a-input style="width: 30px; pointer-events: none; backgroundColor: #fff" placeholder=":" disabled/>
              <a-input type="number" style="width: 100px; text-align: center;" v-model="form.server.edit.data.port"/>
            </a-input-group>
          </a-form-model-item>
          <a-form-model-item label="用户">
            <a-input v-model="form.server.edit.data.username" />
          </a-form-model-item>
          <a-form-model-item label="认证">
            <a-switch v-model="form.server.edit.data.auth_mode" checked-children="密码" un-checked-children="密钥" @change="edit_auth_switch"/>
          </a-form-model-item>
          <a-form-model-item :style="form.server.edit.password_style" label="密码">
            <a-input v-model="form.server.edit.data.password" />
          </a-form-model-item>
          <a-form-model-item :style="form.server.edit.private_key_style" label="密钥">
            <a-input v-model="form.server.edit.data.private_key" type="textarea" />
          </a-form-model-item>
          <a-form-model-item label="备注">
            <a-input v-model="form.server.edit.data.remark" type="textarea" />
          </a-form-model-item>
        </a-form-model>
      </a-modal>
      <a-modal title="编辑组" :visible="form.group.edit.visible" @cancel="form.group.edit.visible=false" @ok="edit_group">
        <a-form-model :model="form.group.edit.data" :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-model-item label="组名">
            <a-input v-model="form.group.edit.data.name" />
          </a-form-model-item>
          <a-form-model-item label="备注">
            <a-input v-model="form.group.edit.data.remark" type="textarea" />
          </a-form-model-item>
        </a-form-model>
      </a-modal>
      <a-modal title="新增组" :visible="form.group.add.visible" @cancel="form.group.add.visible=false" @ok="add_group">
        <a-form-model :model="form.group.add.data" :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-model-item label="组名">
            <a-input v-model="form.group.add.data.name" />
          </a-form-model-item>
          <a-form-model-item label="备注">
            <a-input v-model="form.group.add.data.remark" type="textarea" />
          </a-form-model-item>
        </a-form-model>
      </a-modal>
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
      expandedKeys: [],
      labelCol: { span: 4 },
      wrapperCol: { span: 20 },
      groups: [],
      servers: {},
      form: {
        server: {
            add: {
              visible: false,
              password_style: "display: none",
              private_key_style: "",
              data: {
                name: "New",
                address: "",
                port: 22,
                group: 1,
                username: "",
                password: "",
                private_key: "",
                remark: ""
              }
            },
            edit: {
              visible: false,
              password_style: "display: none",
              private_key_style: "display: none",
              data: {
                id: 0,
                name: "",
                address: "",
                port: 22,
                group: 1,
                username: "",
                password: "",
                private_key: "",
                remark: "",
                auth_mode: true
              }
            }
        },
        group: {
          add: {
            visible: false,
            data: {name: "", remark: ""}
          },
          edit: {
            visible: false,
            data: {name: "", remark: ""}
          },
        }
      }
    }
  },
  methods: {
    check_key(key) {
      let t = typeof key
      if (t === 'string') {
        return false
      }
      return true
    },
    refresh_tree: function() {
      let data = [], that = this;
      this.$webssh.group.list(9999).then(function(response) {
        that.groups = response.detail
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
        let this_ = this
        this.$webssh.server.list(treeNode.dataRef.key).then(function(response) {
          treeNode.dataRef.children = []
          response.detail.forEach(function(item) {
            treeNode.dataRef.children.push({ title: item.name, key: `host-${item.id}`, slots: {icon: 'desktop'}, isLeaf: true})
          })
          this_.servers.set(treeNode.dataRef.key, response.detail)
        }).catch(function(response) {
            this_.$message.error("获取服务器信息失败："+response.message)
        })
        resolve();
      });
    },
    server_select(selectedKeys, event) {
      if (selectedKeys.length != 0) {
        let key = selectedKeys[0]
        if (typeof key == 'string') {
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
      }
    },
    onOpen() {
      this.term.fit()
      this.connection.send(JSON.stringify({type: "connect", cols: this.term.cols, rows: this.term.rows, host: this.host}))
    },
    onclose() {
      this.$message.error("连接中断")
    },
    onerror(error) {
       this.$message.error("连接中断: "+error)
    },
    onresize(e) {
      const msg = { type: "resize", ...e };
      this.term.fit()
    },
    init_term() {
        if (window.location.protocol === 'https:') {
            this.protocol = 'wss://'
        } else {
            this.protocol = 'ws://'
        }
        this.endpoint = `${this.protocol}127.0.0.1:8081/api/system/webssh`
        // obj.endpoint = `${obj.protocol}${window.location.host}/api/system/webssh`

        let initPtySize = this.termSize();
        let cols = initPtySize.cols;
        let rows = initPtySize.rows;
        const terminalContainer = document.getElementById("xterm")
        this.term = new Terminal({
            cursorBlink: true,
            cols: cols,
            rows: rows
        })
        this.term.open(terminalContainer, true)
        this.term.write('Connecting...')
        if (window.WebSocket) {
          // 如果支持websocket
          let ws = new WebSocket(this.endpoint)//后端接口位置
          this.connection = ws
          this.connection.onopen = this.onOpen
          this.connection.onclose = this.onclose
          this.connection.onerror = this.onerror
          this.term.on("resize", this.onresize);
          this.term.attach(this.connection)
        } else {
          // 否则报错
          console.log('WebSocket Not Supported' + this.endpoint)
        }
    },
    onContextMenuClick(treeKey, menuKey) {
      switch (menuKey) {
        case "add_host":
          this.open_add_server(treeKey)
          return
        case "edit_host":
          this.form.server.edit.visible = this.open_edit_server(treeKey)
          return
        case "delete_host":
          return
        case "edit_node":
          this.open_edit_group(treeKey)
          return
        case "delete_node":
          let this_ = this
          this.$confirm({
            title: '确认',
            content: '是否删除节点',
            okText: '确认',
            cancelText: '取消',
            onOk() {
              this_.$webssh.group.remove(treeKey).then(function(response) {
                this_.form.group.edit.visible = false
                this_.$message.success("删除成功")
                this_.refresh_tree()
              }).catch(function(response) {
                this_.$message.error("删除失败："+response.message)
                this_.form.group.edit.visible = false
              })
            }
          });
          return
      }
    },
    add_group() {
      let this_ = this
      this.$webssh.group.create(this.form.group.add.data).then(function(response) {
        this_.form.group.add.visible = false
        this_.$message.success("新增成功")
        this_.refresh_tree()
      }).catch(function(response) {
        this_.$message.error("添加失败："+response.message)
        this_.form.group.add.visible = false
      })
    },
    open_edit_group(key) {
      let this_ = this
      this.groups.forEach(function(item) {
        if (item.id == key) {
          this_.form.group.edit.data = item
          this_.form.group.edit.visible = true
        }
      })
    },
    edit_group() {
      let this_ = this
      this.$webssh.group.update(this.form.group.edit.data).then(function(response) {
        this_.form.group.edit.visible = false
        this_.$message.success("修改成功")
        this_.refresh_tree()
      }).catch(function(response) {
        this_.$message.error("修改失败："+response.message)
        this_.form.group.edit.visible = false
      })
    },
    add_server() {
      let this_ = this
      this.$webssh.server.create(this.form.server.add.data).then(function(response) {
        this_.form.server.add.visible = false
        this_.$message.success("添加成功")
        this_.refresh_tree()
      }).catch(function(response) {
        this_.$message.error("添加失败："+response.message)
        this_.form.server.add.visible = false
      })
    },
    open_add_server(key) {
      let this_ = this
      this.groups.forEach(function(item) {
        if (item.id == key) {
          this_.form.server.add.data.group = item.id
        }
      })
      this.form.server.add.visible = true
    },
    auth_switch(checked, event) {
      if (checked) {
        this.form.server.add.password_style = ""
        this.form.server.add.private_key_style = "display: none"
      } else {
        this.form.server.add.password_style = "display: none"
        this.form.server.add.private_key_style = ""
      }
    },
    edit_auth_switch(checked, event) {
      if (checked) {
        this.form.server.edit.password_style = ""
        this.form.server.edit.private_key_style = "display: none"
      } else {
        this.form.server.edit.password_style = "display: none"
        this.form.server.edit.private_key_style = ""
      }
    },
    open_edit_server(key) {
      let id = key.split("host-", -1)[1]
      let exist = false
      let this_ = this
      this.servers.forEach(function(v, k) {
        v.forEach(function(server) {
          if (server.id == id) {
            this_.form.server.edit.data = server
            exist = true
            if (this_.form.server.edit.data.password != "") {
              this_.form.server.edit.password_style = ""
              this_.form.server.edit.data.auth_mode = true
            }
            if (this_.form.server.edit.data.private_key != "") {
              this_.form.server.edit.data.auth_mode = false
              this_.form.server.edit.private_key_style = ""
            }
          }
        })
      })
      return exist
    },
    edit_server() {
      let this_ = this
      this.$webssh.server.update(this.form.server.edit.data).then(function (response) {
        Object.assign(this_.form.server.edit.data, this_.$options.data().form.server.edit.data)
        this_.form.server.edit.visible = false
        this_.$message.success("更新成功")
        this_.refresh_tree()
      }).catch(function(response) {
        this_.$message.error("更新失败："+response.message)
      }) 
    },
    termSize() {
        const init_width = 9;
        const init_height = 18;

        let windows_width = window.innerWidth;
        let windows_height = window.innerHeight - 200;

        return {
            cols: Math.floor(windows_width / init_width),
            rows: Math.floor(windows_height / init_height),
        }
    }
  },
  created: function () {
    this.refresh_tree()
    this.servers = new Map()
  },
  beforeDestroy() {
    this.connection.close()
    this.term.destroy()
  },
  mounted: function () {
    window.addEventListener("resize", this.onresize);
  }
};
</script>

<style>
.ant-tabs-bar {
  margin: 0 0 5px 0;
}

</style>