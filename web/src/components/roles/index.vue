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
      rolesData: [],
      proxy: {
        Status: false,
        Error: null
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
  },
  created: function () {
    this.refresh_roles()
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