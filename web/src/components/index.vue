<template>
  <el-row :gutter="20">
    <el-col :span="8">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-table :data="serverData" size="mini" :stripe="true" :show-header="false" style="width: 100%">
            <el-table-column prop="name" label="配置" width="60"></el-table-column>
            <el-table-column prop="value" label="值" width="80">
              <template slot-scope="scope">
                <span v-if="scope.row.key === 'status'">
                  <el-tag v-if="scope.row.value" type="success" size="mini" effect="plain">运行中</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">停止</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'balance'">
                  <el-tag v-if="!scope.row.value" type="success" size="mini" effect="plain">随机模式</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">循环模式</el-tag>
                </span>
                <span v-else-if="scope.row.key === 'auto_proxy'">
                  <el-tag v-if="scope.row.value" type="success" size="mini" effect="plain">自动</el-tag>
                  <el-tag v-else size="mini" type="danger" effect="plain">关闭</el-tag>
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
        </el-col>
        <!-- // proxy instances  -->
        <el-col :span="24" style="margin-top: 10px">
          <el-table :data="instanceData" :stripe="true" size="mini" style="width: 100%">
            <el-table-column prop="address">
              <template slot="header" slot-scope="">
                <el-button type="text" size="mini" icon="el-icon-circle-plus-outline" @click="add_instance">新增实例</el-button>
              </template>
              <template slot-scope="scope">
                <span>
                  <el-button type="text" size="mini" @click="edit_instance(scope.row)">
                    {{`ssh://${scope.row.username}@${scope.row.address}`}}</el-button>
                </span>
              </template>
            </el-table-column>
          </el-table>
        </el-col>
      </el-row>
    </el-col>

    <el-col :span="16">
      <el-card>
        <el-button-group>
          <el-button size="mini" type="primary" icon="el-icon-edit" @click="roles.create.visit = true">新增规则</el-button>
        </el-button-group>
        <el-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 300px;float: right;">
          <el-button slot="append" icon="el-icon-search"></el-button>
        </el-input>
      </el-card>
      <el-table :data="rolesData" stripe size="mini" style="width: 100%;margin-top: 10px">
        <el-table-column prop="sub" label="域名" width="200">
          <template slot-scope="scope">
            {{`${scope.row.sub}.${scope.row.domain}`}}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="规则" width="60" 
        :filters="[{ text: '代理', value: true }, { text: '封禁', value: false }]"
        :filter-method="filter_status"
        >
          <template slot-scope="scope">
            <el-tag v-if="scope.row.status" effect="plain" size="mini">转发</el-tag>
            <el-tag v-else size="mini" type="danger" effect="plain">封禁</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button type="text" size="mini" icon="el-icon-delete" @click="remove_role(scope.row)"></el-button>
          </template>
        </el-table-column>
      </el-table>      
    </el-col>
    <el-dialog title="修改服务" :visible.sync="server.edit.visit">
      <el-form :model="server.edit.form" label-position="right">
        <el-form-item label="名称" label-width="100px">
          <el-input v-model="server.edit.form.name" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="端口" label-width="100px">
          <el-input-number :min="81" :max="65534" controls-position="right" v-model="server.edit.form.port" size="small" autocomplete="off"></el-input-number>
        </el-form-item>
        <el-form-item label="负载" label-width="100px">
          <el-switch
            v-model="server.edit.form.status"
            active-text="轮训"
            inactive-text="随机">
          </el-switch>
        </el-form-item>
        <el-form-item label="自动" label-width="100px">
          <el-switch
            v-model="server.edit.form.auto_proxy"
            active-text="开启"
            inactive-text="关闭">
          </el-switch>
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

    <el-dialog title="修改实例" :visible.sync="instaces.edit.visit">
      <el-form :model="instaces.edit.form" label-position="right">
        <el-form-item style="display: none" label="ID" label-width="100px">
          <el-input v-model="instaces.edit.form.ID" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="地址" label-width="100px">
          <el-input v-model="instaces.edit.form.Address" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="用户" label-width="100px">
          <el-input v-model="instaces.edit.form.Username" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" label-width="100px">
          <el-input v-model="instaces.edit.form.Password" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="秘钥" label-width="100px">
          <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 6}" placeholder="请输入内容" v-model="instaces.edit.form.PrivateKey"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="instaces.edit.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_edit_server">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="新增规则" :visible.sync="roles.create.visit">
      <el-form :model="roles.create.form" label-position="right">
        <el-form-item label="域名" label-width="100px">
          <el-input v-model="roles.create.form.url" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="规则" label-width="100px">
          <el-switch v-model="roles.create.form.status" active-text="转发" inactive-text="封禁"></el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="roles.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_role">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      instanceData: [],
      serverData: [],
      rolesData: [],
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
            status: false,
            username: "",
            password: "",
            auto_proxy: false
          }
        }
      },
      instaces: {
        create: {
          visit: false,
          form: {
            address: undefined,
            username: undefined,
            password: undefined,
            privateKey: undefined
          },
        },
        edit: {
          visit: false,
          form: {
            id: undefined,
            address: undefined,
            username: undefined,
            password: undefined,
            private_key: undefined
          },
        }
      },
      roles: {
        filter: undefined,
        create: {
          visit: false,
          form: {
            url: undefined,
            status: undefined,
          },
        }
      }
    }
  },
  methods: {
    add_instance () {
      this.instaces.create.visit = true
    },
    edit_instance (item) {
      this.instaces.edit.visit = true
      this.instaces.edit.form = item
    },
    edit_server (item) {
      let that = this
      this.serverData.forEach(function (row) {
        switch (row.key) {
          case "name":
            that.server.edit.form.name = row.value
          case "port":
            that.server.edit.form.port = row.value
          case "balance":
            that.server.edit.form.status = row.value
          case "auto_proxy":
            that.server.edit.form.auto_proxy = row.value
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
    submit_create_role () {
      let that = this
      this.$api.post('/proxy/role/add', this.roles.create.form).then(function (response) {
        that.roles.create.visit = false
        that.$message({message: '修改成功', type: 'success'})
        that.refresh_roles()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    submit_edit () {
      console.log(this.editInstance)
    },
    remove_role (item) {
      let that = this
      this.$api.post('/proxy/role/remove', {id: item.id}).then(function (response) {
        that.$message({message: '删除成功', type: 'success'})
        that.refresh_roles()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    refresh_server () {
      let that = this
      this.$api.get("/proxy/server/info").then(function (response) {
        that.serverData = response.detail
      })
    },
    refresh_instances () {
      let that = this
      this.$api.get("/proxy/instances").then(function (response) {
        that.instanceData = response.detail
      })
    },
    refresh_roles () {
      let that = this
      this.$api.get("/proxy/roles").then(function (response) {
        that.rolesData = response.detail
      })
    },
    start_server () {
      let that = this
      this.$api.post("/proxy/server/start").then(function (response) {
        that.refresh_server()
      })
    },
    stop_server () {
      let that = this
      this.$api.post("/proxy/server/stop").then(function (response) {
        that.refresh_server()
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
    filter_status(value, row) {
      return row.status === value;
    },
  },
  created: function () {
    this.refresh_instances()
    this.refresh_roles()
    this.refresh_server()
  },
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
}
.el-card__body {
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
</style>