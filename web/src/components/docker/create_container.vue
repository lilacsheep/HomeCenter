<template>
    <a-layout style="padding: 0 12px 12px">
      <a-page-header
          style="padding: 0px 0px;margin: 8px;"
          title="新增容器"
          sub-title="自定义创建容器"
          @back="() => $router.go(-1)"
        />
      <a-layout-content
        :style="{ background: '#fff', padding: '12px', margin: 0, minHeight: '280px' }"
      >
        <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol">
          <a-form-model-item label="名称">
            <a-input v-model="form.name" />
          </a-form-model-item>
          <a-form-model-item label="镜像">
            <a-select v-model="form.region" placeholder="please select your zone">
              <a-select-option value="shanghai">
                Zone one
              </a-select-option>
              <a-select-option value="beijing">
                Zone two
              </a-select-option>
            </a-select>
          </a-form-model-item>
          <a-form-model-item label="端口随机">
            <a-switch v-model="form.delivery" />
          </a-form-model-item>
          <a-form-model-item label="重启策略">
            <a-radio-group v-model="form.restart" button-style="solid" size="small">
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
          <a-form-model-item label="高级选项">
            <a-checkbox-group v-model="form.type">
              <a-checkbox value="1" name="type">
                本地网络
              </a-checkbox>
            </a-checkbox-group>
          </a-form-model-item>
          <a-form-model-item label="Resources">
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
            <a-tab-pane key="1" tab="Tab Title 1">
              <p>Content of Tab Pane 1</p>
              <p>Content of Tab Pane 1</p>
              <p>Content of Tab Pane 1</p>
            </a-tab-pane>
            <a-tab-pane key="2" tab="Tab Title 2">
              <p>Content of Tab Pane 2</p>
              <p>Content of Tab Pane 2</p>
              <p>Content of Tab Pane 2</p>
            </a-tab-pane>
            <a-tab-pane key="3" tab="Tab Title 3">
              <p>Content of Tab Pane 3</p>
              <p>Content of Tab Pane 3</p>
              <p>Content of Tab Pane 3</p>
            </a-tab-pane>
          </a-tabs>
          <a-form-model-item
            v-for="(env, index) in form.envs"
            :key="env.key"
            v-bind="index === 0 ? formItemLayout : {}"
            :label="index === 0 ? 'Domains' : ''"
            :prop="'domains.' + index + '.value'"
            :rules="{
              required: true,
              message: 'domain can not be null',
              trigger: 'blur',
            }"
          >
            <a-input
              v-model="env.value"
              placeholder="please input domain"
              style="width: 60%; margin-right: 8px"
            />
            <a-icon
              v-if="form.envs.length > 1"
              class="dynamic-delete-button"
              type="minus-circle-o"
              :disabled="form.envs.length === 1"
              @click="removeDomain(env)"
            />
          </a-form-model-item>
          <a-form-model-item v-bind="formItemLayoutWithOutLabel">
            <a-button type="dashed" style="width: 60%" @click="addDomain">
              <a-icon type="plus" /> Add field
            </a-button>
          </a-form-model-item>
          <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
            <a-button type="primary" @click="onSubmit">
              Create
            </a-button>
            <a-button style="margin-left: 10px;">
              Cancel
            </a-button>
          </a-form-model-item>
        </a-form-model>
    </a-layout-content>
  </a-layout>
</template>

<script>
export default {
  data () {
    return {
      collapsed: false,
      formItemLayoutWithOutLabel: {
        wrapperCol: {
          xs: { span: 24, offset: 0 },
          sm: { span: 20, offset: 4 }
        }
      },
      labelCol: { span: 2 },
      wrapperCol: { span: 18 },
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
  methods: {
    onSubmit () {
      console.log('submit!', this.form)
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
  }
}
</script>
