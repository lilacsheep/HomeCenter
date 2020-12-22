<template>
  <el-card class="login">
    <div slot="header" class="clearfix">
      <span>用户登录</span>
    </div>
    <el-form :model="login.form">
      <el-form-item label="账号" label-width="80px">
        <el-input v-model="login.form.username" size="small" placeholder="请输入账号" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="密码" label-width="80px">
        <el-input v-model="login.form.password" size="small" placeholder="请输入密码" autocomplete="off" show-password></el-input>
      </el-form-item>
      <el-form-item label-width="200px">
        <el-button size="small" @click="login_submit" type="primary">登 录</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>

export default {
  data() {
    return {
      login: {
        form: {
          username: "",
          password: "",
        }
      }
    }
  },
  methods: {
    login_submit: function() {
      let that = this
      this.$api.post("/login", this.login.form).then(function (response) {
        window.location.href = "/dashboard"
      }).catch(function (response) {
        that.$message({"message": "用户名或密码错误", "type": "warning"})
      })
    }
  },
  created: function () {},
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 10px;
  border-bottom: 0;
  font-size: 32px;
  text-align: center;
}

.el-card__body {
  padding: 5px;
}

.login {
  height: 100%;
  width: 500px;
  margin: 0 auto;
}

</style>