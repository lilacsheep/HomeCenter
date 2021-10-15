<template>
    <a-layout-content :style="{ padding: '12px'}">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          对象存储
        </a-breadcrumb-item>
      </a-breadcrumb>
      <!-- // style="background: #fbfbfb;border: 1px solid #f4f4f4;height: 100%" -->
      <a-tabs default-active-key="1" tab-position="left">
        <a-tab-pane key="1" tab="文件管理">
          <a-row :gutter="20">
            <a-col :span='6'>
              <a-button type="primary" @click="form.group.add.visible=true" block>新增组</a-button>
                <a-tree style="background: #FFFFFF;height: 100%;margin: 0;" :tree-data="minios" @select="server_select" :expandedKeys.sync="expandedKeys" show-icon>
                  <template #title="{ key: treeKey, title }" >
                    <a-dropdown :trigger="['contextmenu']">
                      <span>{{ title }}</span>
                      <template #overlay>
                        <a-menu @click="({ key: menuKey }) => onContextMenuClick(treeKey, menuKey)">
                          <a-menu-item key="add_host" v-if="check_key(treeKey)"><a-icon type="edit" />新增主机</a-menu-item>
                          <a-menu-item key="edit_host" v-if="!check_key(treeKey)"><a-icon type="edit" />编辑</a-menu-item>
                          <a-menu-item key="delete_host" v-if="!check_key(treeKey)"><a-icon type="delete" />删除</a-menu-item>
                          <a-menu-item key="delete_node" v-if="check_key(treeKey)"><a-icon type="delete" />删除</a-menu-item>
                          <a-menu-item key="edit_node" v-if="check_key(treeKey)"><a-icon type="edit" />编辑</a-menu-item>
                        </a-menu>
                      </template>
                    </a-dropdown>
                  </template>
                  <a-icon slot="desktop" type="desktop" />
                  <a-icon slot="folder" type="folder" />
                </a-tree>
            </a-col>
            <a-col :span='18'>
              <a-table :columns="columns" :data-source="objects" size="small">
                <a slot="name" slot-scope="text">{{ text }}</a>
              </a-table>
            </a-col>
          </a-row>
        </a-tab-pane>
        <a-tab-pane key="2" tab="配置管理">
          <a-row :gutter="20">
            <a-col :span="10">
              <a-card title="Card Title">
                <a-card-grid style="width:50%;text-align:center">
                  <a-statistic
                    title="Feedback"
                    :value="11.28"
                    :precision="2"
                    suffix="%"
                    :value-style="{ color: '#3f8600' }"
                    style="margin-right: 50px"
                  >
                    <template #prefix>
                      <a-icon type="arrow-up" />
                    </template>
                  </a-statistic>
                </a-card-grid>
                <a-card-grid style="width:50%;text-align:center">
                  <a-statistic
                    title="Feedback"
                    :value="11.28"
                    :precision="2"
                    suffix="%"
                    :value-style="{ color: '#3f8600' }"
                    style="margin-right: 50px"
                  >
                    <template #prefix>
                      <a-icon type="arrow-up" />
                    </template>
                  </a-statistic>
                </a-card-grid>
                <a-card-grid style="width:50%;text-align:center">
                  <a-statistic
                    title="Feedback"
                    :value="11.28"
                    :precision="2"
                    suffix="%"
                    :value-style="{ color: '#3f8600' }"
                    style="margin-right: 50px"
                  >
                    <template #prefix>
                      <a-icon type="arrow-up" />
                    </template>
                  </a-statistic>
                </a-card-grid>
                <a-card-grid style="width:50%;text-align:center">
                  <a-statistic
                    title="Feedback"
                    :value="11.28"
                    :precision="2"
                    suffix="%"
                    :value-style="{ color: '#3f8600' }"
                    style="margin-right: 50px"
                  >
                    <template #prefix>
                      <a-icon type="arrow-up" />
                    </template>
                  </a-statistic>
                </a-card-grid>
                <a-card-grid style="width:100%;text-align:center" :hoverable="false">
                  <a-statistic
                    title="Feedback"
                    :value="11.28"
                    :precision="2"
                    suffix="%"
                    :value-style="{ color: '#3f8600' }"
                    style="margin-right: 50px"
                  >
                    <template #prefix>
                      <a-icon type="arrow-up" />
                    </template>
                  </a-statistic>
                </a-card-grid>
                <a-card-grid style="width:100%;text-align:center" :hoverable="false">
                  <a-statistic
                    title="Feedback"
                    :value="11.28"
                    :precision="2"
                    suffix="%"
                    :value-style="{ color: '#3f8600' }"
                    style="margin-right: 50px"
                  >
                    <template #prefix>
                      <a-icon type="arrow-up" />
                    </template>
                  </a-statistic>
                </a-card-grid>
              </a-card>
            </a-col>
            <a-col :span="14">
              <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol" style="background: #FFFFFF; padding: 10px;">
                <a-form-model-item label="认证密钥">
                  <a-input v-model="form.access_key" />
                </a-form-model-item>
                <a-form-model-item label="认证密码">
                  <a-input v-model="form.secret_key" />
                </a-form-model-item>
                <a-form-model-item label="存储目录">
                  <a-input v-model="form.dir" />
                </a-form-model-item>
                <a-form-model-item label="对外端口">
                  <a-input-number :max="65534" :min="80" v-model="form.port" />
                </a-form-model-item>
                <a-form-model-item label="控制台">
                  <a-input-number :max="65534" :min="80" v-model="form.port" />
                </a-form-model-item>
                <a-form-model-item label="域名">
                  <a-input v-model="form.domain" />
                </a-form-model-item>
                <a-form-model-item label="服务地址">
                  <a-input v-model="form.server_url" />
                </a-form-model-item>
                
                <a-form-model-item :wrapper-col="{ span: 14, offset: 3 }">
                  <a-button type="primary" @click="onSubmit">
                    更新并启动
                  </a-button>
                </a-form-model-item>
              </a-form-model>
            </a-col>
          </a-row>
          
        </a-tab-pane>
      </a-tabs>
    </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      expandedKeys: [],
      objects: [],
      minios: [],
      labelCol: { span: 4 },
      wrapperCol: { span: 20 },
      form: {
        name: '',
        region: undefined,
        date1: undefined,
        delivery: false,
        type: [],
        resource: '',
        desc: '',
      },
      columns: [
        {
          title: 'Name',
          dataIndex: 'name',
          key: 'name',
          scopedSlots: { customRender: 'name' },
        },
        {
          title: 'Age',
          dataIndex: 'age',
          key: 'age',
          width: 80,
        },
        {
          title: 'Address',
          dataIndex: 'address',
          key: 'address 1',
          ellipsis: true,
        },
        {
          title: 'Long Column Long Column Long Column',
          dataIndex: 'address',
          key: 'address 2',
          ellipsis: true,
        },
        {
          title: 'Long Column Long Column',
          dataIndex: 'address',
          key: 'address 3',
          ellipsis: true,
        },
        {
          title: 'Long Column',
          dataIndex: 'address',
          key: 'address 4',
          ellipsis: true,
        },
      ]
    }
  },
  computed: {},
  created: function () {},
  methods: {
    server_select: function () {}
  },
};
</script>

<style>

</style>