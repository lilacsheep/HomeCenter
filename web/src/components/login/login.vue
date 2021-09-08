<template>
  <a-layout id="components-layout-demo-custom-trigger" :style="{backgroundImage: `url(${backgroundImage})`}">
    <a-layout-content :style="{ margin: '8px', padding: '12px', minHeight: '280px', 'text-align': 'center'}">
        <img :src="imgUrl" class="logo" style="height: 150px;"/>
        <a-row :gutter="16" style="margin: 0 auto; text-align:center;  background: #fff;width:400px; padding: 10px;border: 1px solid #f2f2f2">
            <a-form-model :model="formInline" @submit="handleSubmit" @submit.native.prevent style="">
              <a-form-model-item>
              <a-input v-model="formInline.username" placeholder="用户名">
                  <a-icon slot="prefix" type="user" style="color:rgba(0,0,0,.25)" />
              </a-input>
              </a-form-model-item>
              <a-form-model-item>
              <a-input v-model="formInline.password" type="password" placeholder="密码">
                  <a-icon slot="prefix" type="lock" style="color:rgba(0,0,0,.25)" />
              </a-input>
              </a-form-model-item>
              <a-form-model-item>
              <a-button style="width: 100px" type="primary" html-type="submit" :disabled="formInline.username === '' || formInline.password === ''">
                  登陆
              </a-button>
              </a-form-model-item>
          </a-form-model>
        </a-row>
    </a-layout-content>
  </a-layout>
</template>
<script>
export default {
  data() {
    return {
      imgUrl:require("../../assets/login.png"),
      backgroundImage: require("../../assets/beijing.jpg"),
      formInline: {
        username: '',
        password: '',
      },
    };
  },
  methods: {
    handleSubmit(e) {
      var that = this;
      this.$api.post("/login", this.formInline).then(function (response) {
          location.href = "/"
      }).catch(function (response) {
          that.$message.info(response.message)
      })
    },
  },
};
</script>
<style>
.ant-layout {
  height: 100%;
}

#background {
  background-image: url();
}
#components-layout-demo-custom-trigger .trigger {
  font-size: 18px;
  line-height: 32px;
  padding: 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

#components-layout-demo-custom-trigger .trigger:hover {
  color: #1890ff;
}

#components-layout-demo-custom-trigger .logo {
  height: 48px;
  margin: 5px 25px;
}
</style>
