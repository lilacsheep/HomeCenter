<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-tabs tab-position="left" @tab-click="tabClick">
        <el-tab-pane label="修改密码">
          <el-card>
            <el-form label-width="100px" :model="user.password.change.form">
              <el-form-item label="当前密码">
                <el-input v-model="user.password.change.form.current"  placeholder="请输入当前密码" autocomplete="off" show-password></el-input>
              </el-form-item>
              <el-form-item label="新密码">
                <el-input v-model="user.password.change.form.password1"  placeholder="请输入新密码" autocomplete="off" show-password></el-input>
              </el-form-item>
              <el-form-item label="重复密码">
                <el-input v-model="user.password.change.form.password2"  placeholder="请再次输入新密码" autocomplete="off" show-password></el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="change_password_submit">立即修改</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
        <el-tab-pane label="新增用户">
          <el-card>
            <el-form label-width="40px" :model="user.create.form">
              <el-form-item label="账号">
                <el-input v-model="user.create.form.username"></el-input>
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="user.create.form.password1" placeholder="请输入新密码" autocomplete="off" show-password></el-input>
              </el-form-item>
              <el-form-item label="确认">
                <el-input v-model="user.create.form.password2" placeholder="请再次输入新密码" autocomplete="off" show-password></el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="create_user_submit">立即创建</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
        <el-tab-pane label="用户列表"></el-tab-pane>
      </el-tabs>
    </el-col>
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      user: {
        create: {
          form: {
            username: "",
            password1: "",
            password2: "",
          }
        },
        password: {
          change: {
            form: {
              current: "",
              password1: "",
              password2: "",
            }
          }
        }
      }
    }
  },
  methods: {
    tabClick: function (tab, event) {
      if (tab.index === '0') {
      } else if (tab.index === '1') {

      } else if (tab.index === '2') {
      }
    },
    change_password_submit: function () {
      let that = this
      this.$api.post("/auth/change/self/password", this.user.password.change.form).then(function (response) {
        that.$message({message: "修改成功", type: 'success'})
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    },
    create_user_submit: function() {
      let that = this
      this.$api.post("/auth/create/user", this.user.create.form).then(function (response) {
        that.$message({message: "创建成功", type: 'success'})
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    }
  },
  
  created: function () {
    this.$api.post("/auth/self").then(function (response) {}).catch(function (response) {})
  },
  beforeDestroy () {},
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
}

.el-card__body {
  padding: 20px;
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

.el-drawer__body {
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