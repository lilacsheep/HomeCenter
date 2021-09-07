<template>
    <a-layout style="padding: 0 12px 12px">
        <a-breadcrumb style="margin: 16px 0">
          <a-breadcrumb-item>主页</a-breadcrumb-item>
          <a-breadcrumb-item>容器相关</a-breadcrumb-item>
          <a-breadcrumb-item>卷组列表</a-breadcrumb-item>
        </a-breadcrumb>
        <a-layout-content :style="{ background: '#fff', padding: '12px', margin: 0, minHeight: '280px' }">
          <a-row :gutter="16">
            <a-col :span="18">
              <a-button-group>
                <a-button type="primary" @click="showModal">创建卷组</a-button>
              </a-button-group>
            </a-col>
            <a-col :span="6">
              <!--  @search="onSearch" -->
                <a-input-search style="width: 200px;float: right;" placeholder="input search text" enter-button />
            </a-col>
            <a-col class="table" :span="24">
              <a-table :columns="columns" :data-source="data" :loading="tableloading" size="small">
                <span slot="name" slot-scope="text">{{ text.slice(1,19) }}</span>
                <span slot="action" slot-scope="text, record">
                  <a>Invite</a>
                  <a-divider type="vertical" />
                  <a-popconfirm v-if="data.length" title="确认删除卷组么" @confirm="() => onDelete(record)">
                    <a href="javascript:;">Delete</a>
                  </a-popconfirm>
                </span>
              </a-table>
            </a-col>
          </a-row>
        </a-layout-content>
        <a-modal v-model="visible" :okText="okText" :cancelText="cancelText" :maskClosable='false' title="新增镜像" @ok="handleOk" :bodyStyle='{"padding": "10px"}'>
          <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol">
            <a-form-model-item label="名称">
              <a-input v-model="form.Name" placeholder="e.g. MyVolume" />
            </a-form-model-item>
            <!-- <a-form-model-item label="挂载">
              <a-input v-model="form.Mountpoint" />
            </a-form-model-item> -->
            <a-form-model-item label="驱动">
                <a-radio-group default-value="local" button-style="solid" @change="driverChange">
                  <a-radio-button value="local">
                      本地
                  </a-radio-button>
                  <a-radio-button value="nfs">
                      NFS
                  </a-radio-button>
                  <a-radio-button value="nfs4">
                      NFSv4
                  </a-radio-button>
                </a-radio-group>
            </a-form-model-item>
            <a-form-model-item :style="hideStyle" label="地址">
              <a-input v-model="NfsAddress" @change="addressChange" placeholder="e.g. my.nfs-server.com or xxx.xxx.xxx.xxx"/>
            </a-form-model-item>
            <a-form-model-item :style="hideStyle" label="目录">
              <a-input v-model="form.DriverOpts.device" placeholder="e.g. /share OR :/share" />
            </a-form-model-item>
            <a-form-model-item :style="hideStyle" label="参数">
              <a-input v-model="DriverOptions" />
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
    key: 'Name',
    scopedSlots: { customRender: 'name' }
  },
  {
    title: '参数',
    key: 'Options',
    dataIndex: 'Options'
  },
  {
    title: '操作',
    key: 'action',
    scopedSlots: { customRender: 'action' },
    width: 200,
    fixed: 'right'
  }
]

export default {
  data () {
    return {
      columns,
      data: [],
      tableloading: true,
      visible: false,
      hideStyle: 'display: none',
      DriverOptions: 'rw,noatime,rsize=8192,wsize=8192,tcp,timeo=14',
      NfsAddress: '',
      labelCol: { span: 3 },
      wrapperCol: { span: 20 },
      okText: '拉取',
      cancelText: '取消',
      form: {
        Name: '',
        Driver: 'local',
        Mountpoint: '',
        DriverOpts: {
          type: '',
          o: '',
          device: ''
        }
      }
    }
  },
  methods: {
    showModal () {
      this.visible = true
    },
    driverChange (e) {
      if (e.target.value !== 'local') {
        this.$data.hideStyle = ''
        this.$data.form.DriverOpts.type = e.target.value
      } else {
        this.$data.hideStyle = 'display: none'
        this.$data.NfsAddress = ''
        this.$data.form.DriverOpts.type = ''
        this.$data.form.DriverOpts.o = ''
        this.$data.form.DriverOpts.device = ''
      }
    },
    addressChange (e) {
      this.$data.form.DriverOpts.o = `addr=${this.$data.NfsAddress},${this.$data.DriverOptions}`
    },
    handleOk (e) {
      this.visible = false
      let that = this
      this.$data.tableloading = true
      this.$api.post('/volumes/create', that.$data.form).then(function (response) {
        that.$message.info('ok')
        that.$data.tableloading = false
        that.$api.post('/volumes/list').then(function (response) {
          that.$data.data = response.data
        })
      }).catch(function (resp) {
        that.$message.info(resp.message)
      })
    },
    onDelete (row) {
      let that = this
      that.$data.tableloading = true
      this.$api.post('/volumes/remove', {'id': row.Name}).then(function (response) {
        that.$message.info('删除成功')
        that.$data.tableloading = false
        that.$api.post('/volumes/list').then(function (response) {
          that.$data.data = response.data
        })
      }).catch(function (resp) {
        that.$message.error(resp.message)
      })
    }
  },
  created: function () {
    let that = this
    this.$api.post('/volumes/list').then(function (response) {
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
