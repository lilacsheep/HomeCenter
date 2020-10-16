<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-button-group>
        <el-button size="mini" type="primary" icon="el-icon-edit" @click="roles.create.visit = true">新增规则</el-button>
      </el-button-group>
      <!-- <el-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
        <el-button slot="append" icon="el-icon-search"></el-button>
      </el-input> -->

      <el-table :data="rolesData" stripe border size="mini" style="margin-top: 10px;">
        <el-table-column prop="sub" label="域名" width="200">
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
            <el-popconfirm title="是否删除该规则？" @onConfirm="remove_role(scope.row)">
              <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>     
    </el-col>

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
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      rolesData: [],
      instanceData: [],
      proxy: {
        Status: false,
        Error: null
      },
      roles: {
        filter: undefined,
        create: {
          visit: false,
          style: "display: none",
          form: {
            url: undefined,
            status: undefined,
            instance_id: ""
          },
        }
      }
    }
  },
  methods: {
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
    remove_role (item) {
      let that = this
      this.$api.post('/proxy/role/remove', {id: item.id}).then(function (response) {
        that.$message({message: '删除成功', type: 'success'})
        that.refresh_roles()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    refresh_roles () {
      let that = this
      this.$api.get("/proxy/roles").then(function (response) {
        that.rolesData = response.detail
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
      this.instanceData.forEach(function (item) {
        if (item.id === id) {
          name = item.address
        }
      })
      return name
    }
  },
  created: function () {
    this.refresh_roles()
    this.refresh_instances()
  },
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
}

.el-card__body {
  padding: 0;
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