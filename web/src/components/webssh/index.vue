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
        <a-tree style="background: #FFFFFF;height: 100%;margin: 0;" :tree-data="servergroup" @select="server_select" :expandedKeys.sync="expandedKeys" show-icon>
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
        <a-tabs v-model="tab_connection.activeKey" hide-add type="editable-card" @change="tabChange"  @edit="onEdit">
          <a-button-group slot="tabBarExtraContent">
            <a-button icon="sync" @click="on_reconnect">重连</a-button>
            <a-button icon="thunderbolt" @click="on_clean">清屏</a-button>
            <a-button icon="api" @click="on_connection_close">断开</a-button>
          </a-button-group>
          <a-tab-pane forceRender v-for="pane in tab_connection.panes" :key="pane.key" :tab="pane.title" :closable="pane.closable">
            <a-spin tip="连接中" :spinning="pane.spinning" :delay="500">
              <div :id="`xterm-${pane.id}`"></div>
            </a-spin>
          </a-tab-pane>
        </a-tabs>
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
          <a-form-model-item label="代理">
            <a-input-group>
              <a-row :gutter="8">
                <a-col :span="3">
                  <a-switch checked-children="开" un-checked-children="关" v-model="form.server.add.data.use_proxy" />
                </a-col>
                <a-col :span="5">
                  <a-select v-model="form.server.add.data.proxy_id" :disabled="!form.server.add.data.use_proxy" style="width: 200px">
                    <a-select-option :key="item.id" v-for="item in all_servers" :value="item.id">{{item.name}}</a-select-option>
                  </a-select>
                </a-col>
              </a-row>
            </a-input-group>
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
          <a-form-model-item label="代理">
            <a-input-group>
              <a-row :gutter="8">
                <a-col :span="3">
                  <a-switch checked-children="开" un-checked-children="关" v-model="form.server.edit.data.use_proxy" />
                </a-col>
                <a-col :span="5">
                  <a-select v-model="form.server.edit.data.proxy_id" :disabled="!form.server.edit.data.use_proxy" style="width: 200px">
                    <a-select-option :key="item.id" v-for="item in all_servers" :value="item.id">{{item.name}}</a-select-option>
                  </a-select>
                </a-col>
              </a-row>
            </a-input-group>
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
      tab_connection: {
        activeKey: 0,
        panes : [
          {
            id: 1,
            key: 0, 
            title: "默认连接", 
            spinning: false, 
            connection: undefined, 
            term: undefined,
            cols: 0,
            rows: 0,
            host: "",
            closable: false
          }
        ]
      },
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
      all_servers: [],
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
                remark: "",
                use_proxy: false,
                proxy_id: 0
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
    onEdit(targetKey, action) {
      switch (action) {
        case "remove":
          let activeKey = this.tab_connection.activeKey;
          let lastIndex;
          this.tab_connection.panes.forEach((pane, i) => {
            if (pane.key === targetKey) {
              lastIndex = i - 1;
              pane.connection.close()
              pane.term.dispose()
            }
          });
          const panes = this.tab_connection.panes.filter(pane => pane.key !== targetKey);
          
          if (panes.length && activeKey === targetKey) {
            if (lastIndex >= 0) {
              activeKey = panes[lastIndex].key;
            } else {
              activeKey = panes[0].key;
            }
          } 
          
          if (lastIndex == -1) {
            this.tab_connection.panes = [{
              id: 1,
              key: 0, 
              title: "默认连接", 
              spinning: false, 
              connection: undefined, 
              term: undefined,
              cols: 0,
              rows: 0,
              host: "",
              closable: false
            }]
            this.tab_connection.activeKey = 0
            this.$nextTick(() => {
              let initPtySize = this.termSize();
              let cols = initPtySize.cols;
              let rows = initPtySize.rows;
              let terminalContainer = document.getElementById(`xterm-${this.tab_connection.panes[0].id}`)
              this.tab_connection.panes[0].cols = cols
              this.tab_connection.panes[0].rows = rows
              this.tab_connection.panes[0].term = new Terminal({
                  cursorBlink: false,
                  cols: cols,
                  rows: rows
              })
              this.tab_connection.panes[0].term.open(terminalContainer, true)
            })
          } else {
            this.tab_connection.panes = panes;
            this.tab_connection.activeKey = activeKey;
            this.tabChange(activeKey)
          }
          
      }
    },
    on_clean() {
      this.tab_connection.panes[this.tab_connection.activeKey].connection.send("clear\r")
    },
    on_connection_close() {
      this.tab_connection.panes[this.tab_connection.activeKey].connection.close()
    },
    on_reconnect() {
      this.tab_connection.panes[this.tab_connection.activeKey].connection.close()
      this.tab_connection.panes[this.tab_connection.activeKey].term.dispose()
      this.openTab(this.tab_connection.activeKey+1)
    },
    refresh_tree: function() {
      let data = [], that = this;
      this.all_server(function () {
        that.$webssh.group.list(9999).then(function(response) {
          that.groups = response.detail
          response.detail.forEach(function(item) {
            let children = []
            that.all_servers.forEach(function (server) {
              if (server.group == item.id) {
                children.push({ title: server.name, key: `host-${server.id}`, slots: {icon: 'desktop'}, isLeaf: true})
              }
            })
            data.push({title: item.name,key: item.id,slots: {icon: 'folder'}, children: children})
          })
          that.servergroup = data
        }).catch(function(response) {
          that.$message.error(response.message)
        })
      })
    },
    server_select(selectedKeys, event) {
      if (selectedKeys.length != 0) {
        let key = selectedKeys[0]
        if (typeof key == 'string') {
          let node = event.selectedNodes[0]
          this.host = node.key
          this.title = node.data.props.title
          let index = this.tab_connection.panes.length

          let id = key.split("host-", -1)[1]
          let selected_server = undefined

          this.all_servers.forEach((server) => {
            if (server.id == id) {
              selected_server = server
            }
          })
          
          if (!selected_server) {
            this.$message.error('未发现该服务器信息...')
            return
          }

          if ((index == 1) && (this.tab_connection.panes[0].title == "默认连接")) {
            this.tab_connection.panes[0].title = selected_server.name
            this.tab_connection.panes[0].spinning = true
            this.tab_connection.panes[0].host = node.key
            this.tab_connection.panes[0].closable = true
            this.tab_connection.panes[0].term.setOption('cursorBlink', true)
            this.openTab(index)
          } else {
            const panes = this.tab_connection.panes;
            const activeKey = index;
            let params = {
              key: index, 
              title: selected_server.name, 
              id: index+1, 
              spinning: true, 
              term: undefined,
              connection: undefined,
              cols: 0,
              rows: 0,
              host: node.key,
            }
            panes.push(params)
            this.panes = panes;
            this.tab_connection.activeKey = activeKey;
            this.openTab(index+1)
          }
        }
      }
    },
    openTab(index=0) {
      if (!window.WebSocket) {
        this.$message.error("此游览器不支持WebSocket")
        return
      }

      let ws = new WebSocket(this.endpoint)
      let host = this.tab_connection.panes[index-1].host
      let rows = 0
      let cols = 0
      this.tab_connection.panes[index-1].connection = ws
      this.tab_connection.panes[index-1].connection.onopen = function() {
        ws.send(JSON.stringify({type: "connect", cols: cols, rows: rows, host: host}))
      }
      if (index == 1) {
        cols = this.tab_connection.panes[index-1].cols
        rows = this.tab_connection.panes[index-1].rows
        this.tab_connection.panes[index-1].connection.onmessage = (evt) => {
          let data = JSON.parse(evt.data)
          switch (data.type) {
            case "error":
              this.$message.error("连接错误: "+ data.message)
              this.tab_connection.panes[index-1].spinning = false
            case "success":
              this.tab_connection.panes[index-1].term.attach(this.tab_connection.panes[index-1].connection)
              this.tab_connection.panes[index-1].connection.onmessage = function(evt) {}
              this.tab_connection.panes[index-1].spinning = false
              this.tab_connection.panes[index-1].term.fit()
          }
        }
      } else {
        let initPtySize = this.termSize()
        cols = initPtySize.cols;
        rows = initPtySize.rows;
        setTimeout(() => {
          let terminalContainer = document.getElementById(`xterm-${this.tab_connection.panes[index-1].id}`)
          this.tab_connection.panes[index-1].cols = cols
          this.tab_connection.panes[index-1].rows = rows
          this.tab_connection.panes[index-1].term = new Terminal({
              cursorBlink: true,
              cols: cols,
              rows: rows
          })
          this.tab_connection.panes[index-1].term.open(terminalContainer, true)
          this.tab_connection.panes[index-1].connection.onmessage = (evt) => {
            let data = JSON.parse(evt.data)
            switch (data.type) {
              case "error":
                this.$message.error("连接错误: "+ data.message)
                this.tab_connection.panes[index-1].spinning = false
              case "success":
                this.tab_connection.panes[index-1].term.attach(this.tab_connection.panes[index-1].connection)
                this.tab_connection.panes[index-1].connection.onmessage = function(evt) {}
                this.tab_connection.panes[index-1].spinning = false
                this.tab_connection.panes[index-1].term.fit()
                this.tab_connection.panes[index-1].term.focus()
            }
          }
        }, 50)
      }
      // ws.connection.onerror = () => {
      //   this.$message.error("连接中断: " + error)
      //   this.tab_connection.panes[index-1].spinning = false
      // }
      // this.connection.onclose = () => {
      //   this.$message.error("连接中断")
      // }
      // this.connection.onerror = (err) => {
      //   this.$message.error("连接错误: "+err)
      // }
    },
    tabChange(activeKey) {
      this.tab_connection.panes[activeKey].term.selectAll()
      let w = this.tab_connection.panes[activeKey].term.getSelection()
      this.tab_connection.panes[activeKey].term.dispose()
      
      this.tab_connection.panes[activeKey].term = new Terminal({
          cursorBlink: true,
          cols: this.tab_connection.panes[activeKey].cols,
          rows: this.tab_connection.panes[activeKey].rows
      })
      let terminalContainer = document.getElementById(`xterm-${this.tab_connection.panes[activeKey].id}`)
      
      this.tab_connection.panes[activeKey].term.write(w.trim('', 'right'))
      this.tab_connection.panes[activeKey].term.open(terminalContainer)
      this.tab_connection.panes[activeKey].term.attach(this.tab_connection.panes[activeKey].connection)
      this.tab_connection.panes[activeKey].term.focus()
    },
    onresize(e) {
      const msg = { type: "resize", ...e };
      this.tab_connection.panes.forEach((item) => {
        item.connection.send(msg)
        item.term.fit()
      })
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
      this.all_servers.forEach(function(server) {
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
    },
    all_server(fn = function() {}) {
      let this_ = this
      this.$webssh.server.list().then(function (response) {
        this_.all_servers = response.detail ? response.detail : []
        if (fn) {
          fn()
        }
      }).catch(function(response) {
        this_.$message.error("获取服务器信息失败："+response.message)
      })
    }
  },
  created: function () {
    this.refresh_tree()
    this.servers = new Map()
    if (window.location.protocol === 'https:') {
      this.protocol = 'wss://'
    } else {
      this.protocol = 'ws://'
    }
    let host = this.$apihost == ""? window.location.host: this.$apihost
    this.endpoint = `${this.protocol}${host}/api/system/webssh`
    this.$nextTick(() => {
      let initPtySize = this.termSize();
      let cols = initPtySize.cols;
      let rows = initPtySize.rows;
      let terminalContainer = document.getElementById(`xterm-${this.tab_connection.panes[0].id}`)
      this.tab_connection.panes[0].cols = cols
      this.tab_connection.panes[0].rows = rows
      this.tab_connection.panes[0].term = new Terminal({
          cursorBlink: false,
          cols: cols,
          rows: rows
      })
      this.tab_connection.panes[0].term.open(terminalContainer, true)
    })
  },
  mounted: function () {    
    
  },
  beforeDestroy() {
    this.connection.close()
    this.term.dispose()
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