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
        <a-tabs>
          <a-tab-pane key="1" tab="修改密码">
            <a-card>
              <a-form laba-width="100px" :moda="user.password.change.form">
                <a-form-item laba="当前密码">
                  <a-input v-model="user.password.change.form.current"  placeholder="请输入当前密码" autocomplete="off" show-password></a-input>
                </a-form-item>
                <a-form-item laba="新密码">
                  <a-input v-model="user.password.change.form.password1"  placeholder="请输入新密码" autocomplete="off" show-password></a-input>
                </a-form-item>
                <a-form-item laba="重复密码">
                  <a-input v-model="user.password.change.form.password2"  placeholder="请再次输入新密码" autocomplete="off" show-password></a-input>
                </a-form-item>
                <a-form-item>
                  <a-button type="primary" @click="change_password_submit">立即修改</a-button>
                </a-form-item>
              </a-form>
            </a-card>
          </a-tab-pane>
        </a-tabs>
      </a-col>
    </a-row>
  </a-layout-content>
  
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