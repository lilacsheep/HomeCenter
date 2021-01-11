<template>
  <el-row :gutter="20" v-loading="loading">
    <el-tabs tab-position="left" @tab-click="tabClick">
      <el-tab-pane label="代理配置">
        <el-col :span="8">
          <el-table :data="server.info.data" size="mini" :stripe="true" :show-header="false" style="width: 100%">
            <el-table-column prop="name" label="配置" width="60"></el-table-column>
            <el-table-column prop="value" label="值" width="130">
              <template slot-scope="scope">
                <span v-if="scope.row.key === 'status'">
                  <el-tag v-if="scope.row.value" type="success" size="mini" effect="plain">运行中</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">停止</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'balance'">
                  <el-tag v-if="!scope.row.value" type="success" size="mini" effect="plain">随机模式</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">循环模式</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'proxy_mode'">
                  <el-tag v-if="scope.row.value == 1" type="success" size="mini" effect="plain">全部转发</el-tag>
                  <el-tag v-else-if="scope.row.value == 2" size="mini" type="danger" effect="plain">规则转发</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">DNS转发</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'instances'">
                  <el-tag v-for="(item, i) in scope.row.value" size="mini" type="danger" effect="plain" :key="i">{{item.address}}</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'auto_start'">
                  <el-tag v-if="scope.row.value" type="success" size="mini" effect="plain">开机自启</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">手动启动</el-tag>
                </span>
                <span v-else>{{scope.row.value}}</span>
              </template>
            </el-table-column>
            <el-table-column prop="operation" align="right">
              <template slot-scope="scope">
                <span v-if="scope.row.key === 'status'">
                  <el-button v-if="scope.row.value" @click="stop_server" type="text" size="small">停止</el-button>
                  <el-button v-else @click="start_server" type="text" size="small">启动</el-button>
                </span>
                <span v-else-if="scope.row.key === 'port'">
                  <el-button @click.native="edit_server(scope.row)" type="text" size="small">编辑</el-button>
                </span>
              </template>
            </el-table-column>
          </el-table>

          <el-card v-if="server.info.instances.length > 0" style="margin-top: 10px;" :body-style="{padding: '5px'}">
            <div slot="header" class="clearfix">
              <span>代理池</span>
            </div>
            <div v-for="(item, i) in server.info.instances" :key="i">
              <el-tag  size="mini" effect="plain" >{{item.address}}</el-tag>
              <el-popconfirm title="是否将实例从代理池中移除？" @onConfirm="instance_remove_from_pool(item.id)">
                <el-button slot="reference" style="color: red;float: right;" type="text" size="mini" icon="el-icon-remove"></el-button>
              </el-popconfirm>
              <el-divider></el-divider>
            </div>
          </el-card>
        </el-col>
        <!-- // proxy instances  -->
        <el-col :span="16">
          <el-table :data="instanceData" :stripe="true" size="mini" style="width: 100%">
            <el-table-column prop="address">
              <template slot="header" slot-scope="scope">
                <el-button type="text" size="mini" icon="el-icon-circle-plus-outline" @click="instaces.create.visit = true">新增实例</el-button>
              </template>
              <template slot-scope="scope">
                  <el-button type="text" size="mini" @click="edit_instance(scope.row)">
                    {{`ssh://${scope.row.username}@${scope.row.address}`}}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="country" label="国家" width="180">
              <template slot-scope="scope">
                  {{scope.row.country}}
              </template>
            </el-table-column>
            <el-table-column prop="delay" label="延迟" width="100">
              <template slot-scope="scope" >
                <el-tag v-if="scope.row.delay < 100" size="mini" type="success" effect="plain">{{scope.row.delay}} ms</el-tag>
                <el-tag v-else-if="scope.row.delay < 200" size="mini" type="warning" effect="plain">{{scope.row.delay}} ms</el-tag>
                <el-tag v-else size="mini" type="danger" effect="plain">{{scope.row.delay}} ms</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template slot-scope="scope" >
                <el-popconfirm title="是否将实例删除？" @onConfirm="instance_remove(scope.row.id)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-error"></el-button>
                </el-popconfirm>
                <el-popconfirm v-if="!scope.row.status" title="是否将实例加入代理池？" @onConfirm="instance_into_pool(scope.row.id)">
                  <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-circle-plus"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-col>
      </el-tab-pane>
      <el-tab-pane label="访问记录">
        <el-input size="small" placeholder="请输入内容" v-model="roles.logs.filter" style="width: 300px;float: right;margin-bottom: 10px;">
          <el-button @click="submit_roles_logs_filter" slot="append" icon="el-icon-search"></el-button>
        </el-input>
        <el-table :data="roleLogs" stripe border size="mini">
          <el-table-column prop="domain" label="站点"></el-table-column>
          <el-table-column prop="error" :show-overflow-tooltip="true" label="错误"></el-table-column>
          <el-table-column prop="times" label="次数" width="100"></el-table-column>
          <el-table-column label="操作" fixed="right" width="100">
            <template slot-scope="scope">
              <el-popconfirm title="是否将站点加入负载策略？" @onConfirm="add_log_to_role(scope.row)">
                <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-circle-plus"></el-button>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination class="pagination" background layout="prev, pager, next" :total="roles.logs.pagination.total" @current-change="role_logs_pagination_change"></el-pagination>
      </el-tab-pane>
      <el-tab-pane label="规则配置">
        <el-button-group>
          <el-button size="small" type="primary" icon="el-icon-edit" @click="roles.create.visit = true">新增规则</el-button>
        </el-button-group>
        <el-input size="small" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
          <el-button slot="append" icon="el-icon-search" @click="submit_roles_filter"></el-button>
        </el-input>

        <el-table :data="rolesData" stripe border size="mini" style="margin-top: 10px;">
          <el-table-column prop="sub" label="域名" width="300">
            <template slot-scope="scope">
              {{`${scope.row.sub}.${scope.row.domain}`}}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="规则"
          :filters="[{ text: '代理', value: true }, { text: '封禁', value: false }]"
          :filter-method="filter_status"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.status">
                <el-tag effect="plain" size="mini">转发 -> {{get_instance(scope.row.instance_id)}}</el-tag>
              </span>
              <el-tag v-else size="mini" type="danger" effect="plain">封禁</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right" width="100">
            <template slot-scope="scope">
              <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-edit" @click="change_role(scope.row)"></el-button>
              <el-popconfirm title="是否删除该规则？" @onConfirm="remove_role(scope.row)">
                <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination class="pagination" background layout="prev, pager, next" :total="roles.pagination.total" @current-change="roles_pagination_change"></el-pagination>

      </el-tab-pane>
    </el-tabs>

    <el-dialog title="修改服务" :visible.sync="server.edit.visit" width="450px">
      <el-form :model="server.edit.form" label-position="right">
        <el-form-item label="端口" label-width="80px">
          <el-input-number :min="81" :max="65534" controls-position="right" v-model="server.edit.form.port" size="small" autocomplete="off"></el-input-number>
        </el-form-item>
        <el-form-item label="DNS" label-width="80px">
          <el-input v-model="server.edit.form.dns_addr" size="small" autocomplete="off" style="width: 300px"></el-input>
        </el-form-item>
        <el-form-item label="负载" label-width="80px">
          <el-switch v-model="server.edit.form.status" active-text="轮训" inactive-text="随机"></el-switch>
        </el-form-item>
        <el-form-item label="模式" label-width="80px">
          <el-select v-model="server.edit.form.proxy_mode" size="small" placeholder="请选择">
            <el-option label="全代理模式" :value="1"></el-option>
            <el-option label="规则代理" :value="2"></el-option>
            <el-option label="DNS代理" :value="3"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="server.edit.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_edit_server">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="新增实例" :visible.sync="instaces.create.visit">
      <el-form :model="instaces.create.form" label-position="right">
        <el-form-item label="地址" label-width="100px">
          <el-input v-model="instaces.create.form.address" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="用户" label-width="100px">
          <el-input v-model="instaces.create.form.username" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" label-width="100px">
          <el-input v-model="instaces.create.form.password" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="秘钥" label-width="100px">
          <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 6}" placeholder="请输入内容" v-model="instaces.create.form.private_key"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="instaces.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_instance">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="修改实例" :visible.sync="instaces.edit.visit" width="450px">
      <el-form :model="instaces.edit.form" label-position="right">
        <el-form-item style="display: none" label="ID" label-width="80px">
          <el-input v-model="instaces.edit.form.ID" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="地址" label-width="80px">
          <el-input v-model="instaces.edit.form.address" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="归属" label-width="80px">
          <el-select size="small" v-model="instaces.edit.form.country_code" filterable placeholder="请选择">
            <el-option v-for="item in countrys" :key="item.code" :label="item.cn" :value="item.code"></el-option>
          </el-select>
          <el-tooltip content="选择后将强制使用你选择的国家不在自动判断" placement="top">
            <el-checkbox v-model="instaces.edit.form.force_country">强制</el-checkbox>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="用户" label-width="80px">
          <el-input v-model="instaces.edit.form.username" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" label-width="80px">
          <el-input v-model="instaces.edit.form.password" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="秘钥" label-width="80px">
          <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 6}" placeholder="请输入内容" v-model="instaces.edit.form.private_key"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="instaces.edit.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_edit_instance">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="新增规则" :visible.sync="roles.create.visit">
      <el-form :model="roles.create.form" label-position="right">
        <el-form-item label="域名" label-width="100px">
          <el-input v-model="roles.create.form.url" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="规则" label-width="100px">
          <el-switch v-model="roles.create.form.status" active-text="转发" inactive-text="封禁" @change="role_change"></el-switch>
        </el-form-item>
        <el-form-item label="实例" label-width="100px" :style="roles.create.style">
          <el-select v-model="roles.create.form.instance_id" size="small" placeholder="请选择转发">
            <el-option label="默认转发" value=""></el-option>
            <el-option v-for="(item, i) in instanceData" :key="i" :label="item.address" :value="item.id"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="roles.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_role">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="修改规则" :visible.sync="roles.change.visit">
      <el-form :model="roles.change.form" label-position="right">
        <el-form-item label="规则" label-width="100px">
          <el-switch v-model="roles.change.form.status" active-text="转发" inactive-text="封禁" @change="role_change"></el-switch>
        </el-form-item>
        <el-form-item label="实例" label-width="100px" :style="roles.create.style">
          <el-select v-model="roles.change.form.instance_id" size="small" placeholder="请选择转发">
            <el-option label="默认转发" value=""></el-option>
            <el-option v-for="(item, i) in instanceData" :key="i" :label="item.address" :value="item.id"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="roles.change.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_change_role">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      loading: false,
      instanceData: [],
      serverData: [],
      countrys: [],
      proxy: {
        Status: false,
        Error: null
      },
      server: {
        edit: {
          visit: false,
          form: {
            name: "",
            port: 0,
            dns_addr: "",
            status: false,
            username: "",
            password: "",
            proxy_mode: "",
          }
        },
        info: {
          data: [],
          instances: []
        }
      },
      instaces: {
        create: {
          visit: false,
          form: {
            address: undefined,
            username: undefined,
            password: undefined,
            privateKey: undefined,
          },
        },
        edit: {
          visit: false,
          form: {
            id: undefined,
            address: undefined,
            username: undefined,
            password: undefined,
            private_key: undefined,
            force_country: false,
            country_code: undefined,
          },
        }
      },
      rolesData: [],
      instanceData: [],
      roleLogs: [],
      roles: {
        pagination: {
          total: 1000,
          page: 1,
          limit: 10
        },
        filter: undefined,
        create: {
          visit: false,
          style: "display: none",
          form: {
            url: undefined,
            status: undefined,
            instance_id: ""
          },
        },
        change: {
          visit: false,
          style: "display: none",
          form: {
            id: undefined,
            status: undefined,
            instance_id: ""
          },
        },
        logs:{
          filter: undefined,
          pagination: {
            total: 1000,
            page: 1,
            limit: 10
          }
        },
      }
    }
  },
  methods: {
    tabClick: function (tab, event) {
      if (tab.index === '0') {
        this.refresh_instances()
        this.refresh_server()
      } else if (tab.index === '1') {
        this.refresh_visit_logs(this.roles.logs.pagination.page, this.roles.logs.pagination.limit)
      } else if (tab.index === '2') {
        this.refresh_roles(this.roles.pagination.page, this.roles.pagination.limit)
      }
    },
    edit_instance (item) {
      this.refresh_countrys()
      this.instaces.edit.visit = true
      this.instaces.edit.form = item
    },
    edit_server (item) {
      let that = this
      this.server.info.data.forEach(function (row) {
        switch (row.key) {
          case "name":
            that.server.edit.form.name = row.value
            break
          case "port":
            that.server.edit.form.port = row.value
            break
          case "dns_addr":
            that.server.edit.form.dns_addr = row.value
            break
          case "balance":
            that.server.edit.form.status = row.value
            break
          case "proxy_mode":
            that.server.edit.form.proxy_mode = row.value
            break
        }
      })
      this.server.edit.visit = true
    },
    submit_create_instance: function () {
      let that = this
      this.instaces.create.visit = false
      this.$api.post("/proxy/instance/create", this.instaces.create.form).then(function (response) {
        that.$message({type: "success", message: '创建成功'})
        that.refresh_instances()
      }).catch(function (response) {
        that.$message({type: "error", message: response.message})
      })
    },
    submit_edit_server () {
      let that = this
      this.$api.post('/proxy/server/update', this.server.edit.form).then(function (response) {
        that.server.edit.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    submit_edit_instance () {
      let that = this
      this.$api.post('/proxy/instance/update', this.instaces.edit.form).then(function (response) {
        that.instaces.edit.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    instance_remove_from_pool (id) {
      let that = this
      this.$api.post('/proxy/instance/pool/remove', {'id': id}).then(function (response) {
        that.$message({message: '移除成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    instance_remove (id) {
      let that = this
      this.$api.post('/proxy/instance/remove', {'id': id}).then(function (response) {
        that.$message({message: '移除成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    instance_into_pool (id) {
      let that = this
      this.$api.post('/proxy/instance/pool/add', {'id': id}).then(function (response) {
        that.$message({message: '加入成功', type: 'success'})
        that.refresh_instances()
        that.refresh_server()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    refresh_server () {
      let that = this
      this.$api.get("/proxy/server/info").then(function (response) {
        that.server.info.data = response.detail.data
        that.server.info.instances = response.detail.instances
      })
    },
    refresh_instances () {
      let that = this
      this.$api.get("/proxy/instances").then(function (response) {
        that.instanceData = response.detail
      })
    },
    refresh_visit_logs (page, limit) {
      let that = this
      let params = {limit: limit, page: page}
      if (this.roles.logs.filter != "") {
        params["filter"] = this.roles.logs.filter
      }
      this.$api.get("/proxy/logs", params).then(function (response) {
        that.roleLogs = response.detail
        that.roles.logs.pagination.total = response.hasOwnProperty("count") ? response.count : 0
      })
    },
    add_log_to_role(row) {
      let that = this
      this.$api.post("/proxy/log/add", {id: row.id}).then(function (response) {
        that.$notify({title: '添加成功', message: '添加规则成功', type: 'success'});
        that.refresh_visit_logs(this.roles.logs.pagination.page, this.roles.logs.pagination.limit)
      }).catch(function (response) {
        that.$notify({title: '添加失败', message: response.message, type: 'error'});
      })
    },
    start_server () {
      this.loading = true
      let that = this
      this.$api.post("/proxy/server/start").then(function (response) {
        that.refresh_server()
        that.loading = false
        that.$notify({title: '启动成功', message: '服务器启动成功', type: 'success'});
      }).then(function (response) {
        that.$notify({title: '启动失败', message: response.message, type: 'warning'});
        that.loading = false
      })
    },
    stop_server () {
      this.loading = true
      let that = this
      this.$api.post("/proxy/server/stop").then(function (response) {
        that.refresh_server()
        that.loading = false
        that.$notify({title: '停止成功', message: '服务器启动成功', type: 'success'});
      }).then(function (response) {
        that.$notify({title: '停止失败', message: response.message, type: 'warning'});
        that.loading = false
      })
    },
    start () {
      if (this.proxy.Status) {
        this.$message.warning('服务已经启动')
      } else{
        this.start_server()
      }
    },
    stop () {
      if (!this.proxy.Status) {
        this.$message.warning('服务已经停止')
      } else{
        this.stop_server()
      }
    },
    submit_create_role () {
      let that = this
      this.$api.post('/proxy/role/add', this.roles.create.form).then(function (response) {
        that.roles.create.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_roles(that.roles.logs.pagination.page, that.roles.logs.pagination.limit)
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    remove_role (item) {
      let that = this
      this.$api.post('/proxy/role/remove', {id: item.id}).then(function (response) {
        that.$message({message: '删除成功', type: 'success'})
        that.refresh_roles(that.roles.logs.pagination.page, that.roles.logs.pagination.limit)
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    change_role (role) {
      this.role_change(role.status)
      this.roles.change.form = role
      this.roles.change.visit = true
    },
    submit_change_role () {
      let that = this
      this.$api.post("/proxy/role/change", this.roles.change.form).then(function (response) {
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_roles(that.roles.logs.pagination.page, that.roles.logs.pagination.limit)
      }).catch(function (response) {
        that.$message({message: '修改失败', type: 'error'})
      })
      this.roles.change.visit = false
    },
    refresh_roles (page, limit) {
      let that = this
      let params = {page: page, limit: limit}
      if (this.roles.filter != "") {
        params["filter"] = this.roles.filter
      }
      this.$api.get("/proxy/roles", params).then(function (response) {
        that.rolesData = response.detail
        that.roles.pagination.total = response.hasOwnProperty("count") ? response.count : 0
      })
    },
    filter_status(value, row) {
      return row.status === value;
    },
    refresh_instances () {
      let that = this
      this.$api.get("/proxy/instances").then(function (response) {
        that.instanceData = response.detail
      })
    },
    role_change: function (status) {
      if (status) {
        this.roles.create.style = ""
      } else {
        this.roles.create.style = "display: none"
      }
    },
    get_instance: function (id) {
      name = "默认转发"
      if (this.instanceData != null) {
        this.instanceData.forEach(function (item) {
          if (item.id === id) {
            name = item.address
          }
        })
      }
      return name
    },
    role_logs_pagination_change(value) {
      this.roles.logs.pagination.page = value
      this.refresh_visit_logs(this.roles.logs.pagination.page, this.roles.logs.pagination.limit)
    },
    submit_roles_logs_filter() {
      this.refresh_visit_logs(this.roles.logs.pagination.page, this.roles.logs.pagination.limit)
    },
    submit_roles_filter() {
      this.refresh_roles(this.roles.pagination.page, this.roles.pagination.limit)
    },
    roles_pagination_change(value) {
      this.roles.pagination.page = value
      this.refresh_roles(this.roles.pagination.page, this.roles.pagination.limit)
    },
    refresh_countrys() {
      let that = this
      this.$api.get("/common/countrys").then(function (response) {
        that.countrys = response.detail
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    }
  },
  created: function () {
    this.refresh_instances()
    this.refresh_server()
    this.refresh_roles(this.roles.logs.pagination.page, this.roles.logs.pagination.limit)
    this.refresh_visit_logs(this.roles.logs.pagination.page, this.roles.logs.pagination.limit)
  },
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
}

.el-dialog__header {
  padding: 10px 10px 5px;
  border-bottom: 1px solid whitesmoke;
}

.el-dialog__headerbtn {
  top: 12px;
}
.el-dialog__body {
  padding: 15px 10px;
}
.el-dialog__footer {
  border-top: 1px solid whitesmoke;
  padding: 5px 10px 10px;
}

.el-divider--horizontal {
  margin: 3px 0;
}

.pagination {
  margin-top: 10px;
  float: right;
}
</style>