<template>
    <a-layout-content :style="{ padding: '12px'}">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          容器模板
        </a-breadcrumb-item>
      </a-breadcrumb>
      <a-row :gutter="16" >
        <a-button @click="open_template">添加模板</a-button>
        <a-drawer :visible="template.visible" :closable="false" @close="close_template" width="60%" :bodyStyle='{background: "#fbfbfb"}'>
          <a-form-model :model="form" v-bind="formItemLayoutWithOutLabel" style="height: 100%;">
            <a-form-model-item label="名称" v-bind="formItemLayout">
              <a-input  style="width: 60%; margin-right: 8px"/>
            </a-form-model-item>
            <a-form-model-item v-bind="formItemLayout" label="镜像">
              <a-input v-model="form.name" />
            </a-form-model-item>
            <a-form-model-item v-bind="formItemLayout" label="端口随机">
              <a-switch v-model="form.delivery" />
            </a-form-model-item>
            <a-form-model-item v-bind="formItemLayout" label="重启策略">
              <a-radio-group v-model="form.restart" button-style="solid">
                <a-radio-button value="Unless stopped">
                  除非停止
                </a-radio-button>
                <a-radio-button value="always">
                  总是
                </a-radio-button>
                <a-radio-button value="onfailure">
                  失败
                </a-radio-button>
                <a-radio-button value="never">
                  从不
                </a-radio-button>
              </a-radio-group>
            </a-form-model-item>
            <a-form-model-item v-bind="formItemLayout" label="高级选项">
              <a-checkbox-group v-model="form.type">
                <a-checkbox value="1" name="type">
                  本地网络
                </a-checkbox>
              </a-checkbox-group>
            </a-form-model-item>
            <a-form-model-item v-bind="formItemLayout" label="Resources">
              <a-radio-group v-model="form.resource">
                <a-radio value="1">
                  Sponsor
                </a-radio>
                <a-radio value="2">
                  Venue
                </a-radio>
              </a-radio-group>
            </a-form-model-item>
            <a-tabs type="card">
              <a-tab-pane key="1" tab="环境变量">
                <a-form-model-item v-for="(env, index) in form.envs" :key="env.key" v-bind="index === 0 ? formItemLayout : {}"
                  :label="index === 0 ? 'Domains' : ''"
                  :prop="'domains.' + index + '.value'">
                  <a-input v-model="env.value" placeholder="please input domain" style="width: 60%; margin-right: 8px"/>
                  <a-icon v-if="form.envs.length > 1" class="dynamic-delete-button" type="minus-circle-o"  :disabled="form.envs.length === 1" @click="removeDomain(env)"/>
                </a-form-model-item>
                <a-form-model-item v-bind="formItemLayoutWithOutLabel">
                  <a-button type="dashed" style="width: 60%" @click="addDomain">
                    <a-icon type="plus" /> Add field
                  </a-button>
                </a-form-model-item>
              </a-tab-pane>
              <a-tab-pane key="2" tab="卷组">
                <p>Content of Tab Pane 2</p>
                <p>Content of Tab Pane 2</p>
                <p>Content of Tab Pane 2</p>
              </a-tab-pane>
              <a-tab-pane key="3" tab="端口绑定">
                <p>Content of Tab Pane 3</p>
                <p>Content of Tab Pane 3</p>
                <p>Content of Tab Pane 3</p>
              </a-tab-pane>
            </a-tabs>
            
            <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
              <a-button type="primary" @click="onSubmit">
                Create
              </a-button>
              <a-button style="margin-left: 10px;">
                Cancel
              </a-button>
            </a-form-model-item>
          </a-form-model>
        </a-drawer>
        

      </a-row>
    </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      collapsed: false,
      formItemLayoutWithOutLabel: {
        wrapperCol: {
          xs: { span: 24, offset: 0 },
          sm: { span: 20, offset: 4 }
        }
      },
      formItemLayout: {
        labelCol: {
          xs: { span: 24 },
          sm: { span: 4 },
        },
        wrapperCol: {
          xs: { span: 24 },
          sm: { span: 20 },
        },
      },
      template: {
        visible: false
      },
      form: {
        name: '',
        region: undefined,
        envs: [],
        delivery: false,
        type: [],
        resource: '',
        desc: ''
      }
    }
  },

  created: function () {
    
  },
  methods: {
    open_template() {
      this.template.visible = true
    },
    close_template() {
      this.template.visible = false
    },
    removeDomain (item) {
      let index = this.form.envs.indexOf(item)
      if (index !== -1) {
        this.form.envs.splice(index, 1)
      }
    },
    addDomain () {
      this.form.envs.push({
        value: '',
        key: Date.now()
      })
    }
  },
};
</script>

<style>
.ant-list-bordered .ant-list-header {
  padding-left: 0;
}
.ant-list-header, .ant-list-footer {
  padding: 0;
}

.ant-list-item {
  padding: 5px;
}
.card-container {
  background: #fbfbfb;
  overflow: hidden;
  /* padding: 24px; */
}
.card-container > .ant-tabs-card > .ant-tabs-content {
  margin-top: -16px;
}

.card-container > .ant-tabs-card > .ant-tabs-content > .ant-tabs-tabpane {
  background: #fff;
  padding: 0;
}

.card-container > .ant-tabs-card > .ant-tabs-bar {
  border-color: #fff;
}

.card-container > .ant-tabs-card > .ant-tabs-bar .ant-tabs-tab {
  border-color: transparent;
  background: transparent;
}

.card-container > .ant-tabs-card > .ant-tabs-bar .ant-tabs-tab-active {
  /* border-color: #fff; */
  background: #fff;
}

.ant-collapse > .ant-collapse-item > .ant-collapse-header {
  padding: 4px 8px;
  padding-left: 40px;
}
</style>