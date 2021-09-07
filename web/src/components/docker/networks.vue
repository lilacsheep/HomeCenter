<template>
    <a-layout style="padding: 0 12px 12px">
        <a-breadcrumb style="margin: 16px 0">
          <a-breadcrumb-item>主页</a-breadcrumb-item>
          <a-breadcrumb-item>容器相关</a-breadcrumb-item>
          <a-breadcrumb-item>网络列表</a-breadcrumb-item>
        </a-breadcrumb>
        <a-layout-content :style="{ background: '#fff', padding: '12px', margin: 0, minHeight: '280px' }">
          <a-row :gutter="16">
            <a-col :span="18">
              <a-button-group>
                <a-button type="primary" @click="showModal">新增镜像</a-button>
              </a-button-group>
            </a-col>
            <a-col :span="6">
              <!--  @search="onSearch" -->
              <a-input-search style="width: 200px;float: right;" placeholder="input search text" enter-button />
            </a-col>
            <a-col class="table" :span="24">
              <a-table :columns="columns" :data-source="data" :loading="tableloading" size="small">
                <span slot="Id" slot-scope="text">{{ text.slice(7,19) }}</span>
                <span slot="Created" slot-scope="text">{{ text | dateformat }}</span>
                <span slot="Used" slot-scope="text, record">{{ JSON.stringify(record.Containers) == "{}" ? '未使用' : "使用"}}</span>
                <span slot="Subnet" slot-scope="text, record">{{ record.IPAM.Config.length > 0 ? record.IPAM.Config[0].Subnet :  '无'}}</span>
                <span slot="action" slot-scope="text, record">
                  <a>Invite</a>
                  <a-divider type="vertical" />
                  <a-popconfirm v-if="data.length" title="确认删除镜像么" @confirm="() => onDelete(record)">
                    <a href="javascript:;">Delete</a>
                  </a-popconfirm>
                </span>
              </a-table>
            </a-col>
          </a-row>
        </a-layout-content>
        <a-modal v-model="visible" :okText="okText" :cancelText="cancelText" :maskClosable='false' title="新增镜像" @ok="handleOk" :bodyStyle='{"padding": "10px"}'>
          <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol">
            <a-form-model-item label="镜像名">
              <a-input v-model="form.ref" />
            </a-form-model-item>
            <a-form-model-item label="用户名">
              <a-input v-model="form.username" />
            </a-form-model-item>
            <a-form-model-item label="密码">
              <a-input v-model="form.password" />
            </a-form-model-item>
          </a-form-model>
        </a-modal>
      </a-layout>
</template>

<script>
const columns = [
  {
    title: 'Name',
    dataIndex: 'Name',
    key: 'Name'
  },
  {
    title: '子网',
    dataIndex: 'IPAM',
    key: 'IPAM',
    scopedSlots: { customRender: 'Subnet' }
  },
  {
    title: '类型',
    key: 'Driver',
    dataIndex: 'Driver'
  },
  {
    title: '使用',
    key: 'Containers',
    dataIndex: 'Containers',
    scopedSlots: { customRender: 'Used' }
  },
  {
    title: '创建时间',
    key: 'Created',
    dataIndex: 'Created',
    scopedSlots: { customRender: 'Created' }
  },
  {
    title: '操作',
    key: 'action',
    scopedSlots: { customRender: 'action' }
  }
]

export default {
  data () {
    return {
      columns,
      data: [],
      tableloading: true,
      visible: false,
      labelCol: { span: 3 },
      wrapperCol: { span: 20 },
      okText: '拉取',
      cancelText: '取消',
      form: {
        ref: '',
        username: '',
        password: ''
      }
    }
  },
  methods: {
    showModal () {
      this.visible = true
    },
    handleOk (e) {
      this.visible = false
      let that = this
      this.$data.tableloading = true
      this.$api.post('/networks/create', that.$data.form).then(function (response) {
        that.$message.info('ok')
        that.$data.tableloading = false
        that.$api.post('/networks/list').then(function (response) {
          that.$data.data = response.data
        })
      }).catch(function (resp) {
        that.$message.info(resp.message)
      })
    },
    onDelete (row) {
      let that = this
      if (row.Name === 'host') {
        that.$message.warning('该网络不能删除')
        return
      }
      if (row.Name === 'bridge') {
        that.$message.warning('该网络不能删除')
        return
      }
      if (row.Name === 'none') {
        that.$message.warning('该网络不能删除')
        return
      }
      that.$data.tableloading = true
      this.$api.post('/networks/remove', {'id': row.Id}).then(function (response) {
        that.$message.info('删除成功')
        that.$data.tableloading = false
        that.$api.post('/networks/list').then(function (response) {
          that.$data.data = response.data
        })
      }).catch(function (resp) {
        that.$message.error(resp.message)
      })
    }
  },
  created: function () {
    let that = this
    this.$api.post('/networks/list').then(function (response) {
      that.$data.data = response.data
      that.$data.tableloading = false
    })
  }
}
</script>

<style>
.table {
  margin-top: 10px;
}

.ant-table-thead > tr >th{
  color: white;
  background: #5069d6 !important;
}

.ant-table-tbody > tr >th{
  color: white;
  background: #5069d6 !important;
}

.ant-table-small > .ant-table-content > .ant-table-body {
    margin: 0 0px;
}

.ant-modal-header {
  padding: 8px 12px;
}
</style>
