<template>
  <a-layout-content style="padding: 12px;">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
            动态域名
        </a-breadcrumb-item>
    </a-breadcrumb>
    <a-row :span="24">
        <a-col :span="24">
            <a-button-group>
                <a-button type="primary" icon="edit" @click="ddns.settings.create.visit = true">新增同步</a-button>
            </a-button-group>
        </a-col>
        <a-col :span="24" style="margin-top: 10px">
            <a-table style="background: #FFFFFF" :data-source="ddns.roles" size="small" :columns="ddns.columns" >
            </a-table>
        </a-col>
    </a-row>
    <a-modal title="新建同步" :visible.sync="ddns.settings.create.visit" width="500px">
      <a-form-model :model="ddns.settings.create.form" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-model-item label="域名">
          <a-input placeholder="域名" v-model="ddns.settings.create.form.domain"></a-input>
        </a-form-model-item>
        <a-form-model-item label="模式">
          <a-switch v-model="ddns.settings.create.form.mode" checked-children="本地网卡" un-checked-children="公网IP" @change="ddns_mode_change"></a-switch>
        </a-form-model-item >
        <a-form-model-item  label="网卡" :style="ddns.cards_style">
          <a-select v-loading="ddns.cards_style_loading" v-model="ddns.settings.create.form.net_card" placeholder="请选择">
            <a-select-option v-for="item in ddns.net_cards" :key="item.index" :value="item.index" :disabled="item.flags && item.flags[0] !== 'up'">{{item.name}}</a-select-option>
          </a-select>
        </a-form-model-item>
        <a-form-model-item label="同步间隔">
          <a-select v-model="ddns.settings.create.form.time_interval" placeholder="请选择">
            <a-select-option value="@hourly">每小时</a-select-option>
            <a-select-option value="@every 30m">半小时</a-select-option>
            <a-select-option value="@every 24h">每天</a-select-option>
          </a-select>
        </a-form-model-item>
        <a-form-model-item label="供应商">
          <a-radio-group v-model="ddns.settings.create.form.provider" @change="provider_change">
            <a-radio-button value="dnspod">DNSPod</a-radio-button>
          </a-radio-group>
        </a-form-model-item>
        <a-form-model-item label="DNSPod" :style="ddns.dnspod_style">
          <a-input placeholder="ID" v-model="ddns.settings.create.form.dnspod_id" style="width: 120px"></a-input>
          <a-input placeholder="Token" v-model="ddns.settings.create.form.dnspod_token" style="width: 200px"></a-input>
        </a-form-model-item>
        <a-form-model-item>
          <a-switch v-model="ddns.add_mode" un-checked-children="新增记录" checked-children="原有记录"  @change="ddns_add_change"></a-switch>
        </a-form-model-item>
        <a-form-model-item label="子域名" :style="ddns.sub_domain_style">
          <a-input placeholder="子域名" v-model="ddns.settings.create.form.sub_domain" style="width: 120px" @blur="check_subdomain"></a-input>
        </a-form-model-item>
        <a-form-model-item label="记录" :style="ddns.records_style">
          <a-select :loading="ddns.records_loading" @visible-change="refresh_record" @change="record_change" v-model="ddns.settings.create.form.record_id" placeholder="请选择">
            <a-select-option v-for="item in ddns.records" :key="item.id" :label="`${item.name}.${ddns.settings.create.form.domain}(${item.type})`" :value="item.id" :disabled="item.type == 'NS'"></a-select-option>
          </a-select>
        </a-form-model-item>
      </a-form-model>
      <div slot="footer">
        <a-button @click="ddns.settings.create.visit = false">取 消</a-button>
        <a-button type="primary" @click="submit_create_settings">确 定</a-button>
      </div>
    </a-modal>
  </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      ddns: {
        records: [],
        roles: [],
        columns: [
          {title: '域名', dataIndex: 'domain', key: 'domain'},
          {title: '渠道', dataIndex: 'provider', key: 'provider'},
          {title: '网卡', dataIndex: 'net_card', key: 'net_card'},
          {title: '状态', dataIndex: 'status', key: 'status'},
          {title: '更新时间', dataIndex: 'updated_on', key: 'updated_on'},
        ],
        net_cards: [],
        cards_style: "display: none",
        cards_style_loading: true,
        records_loading: false,
        dnspod_style: "",
        add_mode: false,
        sub_domain_style: "",
        records_style: "display: none",
        settings: {
          loading: false,
          create: {
            visit: false,
            form: {
              domain: "",
              sub_domain: "",
              provider: "dnspod",
              net_card: "",
              time_interval: "@hourly",
              mode: false,
              dnspod_id: "",
              dnspod_token: "",
              record_id: ""
            }
          },
        },
        setting: {
          visit: false,
          info: {
            entry: {
              name: "",
              time: "",
              status: -1
            },
            setting: {
              id: "",
              sub_domain: "",
              domain: "",
              provider: "",
              time_interval: "",
              use_public_ip: "",
              net_card: "",
              status: false,
              history: []
            },
          }
        }
      }
    }
  },
  methods: {
    submit_create_settings: function () {
      let that = this
      this.$api.post("/ddns/setting/create", this.ddns.settings.create.form).then(function (response) {
        that.$message.success('添加成功')
        that.refresh_settings()
      }).catch(function (response) {
        that.$message.error(response.message)
      })
      this.ddns.settings.create.visit = false
    },
    ddns_mode_change: function (value) {
      if (value) {
        this.ddns.cards_style = ""
        this.refresh_net_cards()
      } else{
        this.ddns.cards_style = "display: none"
      }
    },
    ddns_add_change: function (value) {
      if (value) {
        this.ddns.records_style = ""
        this.ddns.sub_domain_style = "display: none"
      } else {
        this.ddns.records_style = "display: none"
        this.ddns.sub_domain_style = ""
      }
    },
    provider_change: function (value) {
      if (value === "dnspod") {
        this.ddns.dnspod_style = ""
      } else {
        this.ddns.dnspod_style = "display: none"
      }
    },
    search_net_card: function (index) {
      let card = "未知"
      this.ddns.net_cards.forEach(function (item) {
        if (item.index == index) {
          card = item.name
        }
      })
      return card
    },
    refresh_net_cards: function () {
      let that = this
      this.ddns.cards_style_loading = true
      this.$api.get("/ddns/netcards").then(function (response) {
        that.ddns.net_cards = response.detail
        that.ddns.cards_style_loading = false
      }).catch(function (response) {
        that.ddns.cards_style_loading = false
      })
    },
    check_subdomain: function (event) {
      if (this.ddns.settings.create.form.sub_domain == "") {
        this.ddns.settings.create.form.sub_domain = "@"
      }
    },
    refresh_settings: function () {
      let that = this
      this.$api.get("/ddns/settings").then(function (response) {
        that.ddns.roles = response.detail
      })
    },
    record_change: function (value) {
      let that = this
      this.ddns.records.forEach(function (item) {
        if (item.id == value) {
          that.ddns.settings.create.form.sub_domain = item.name
        }
      })
    },
    refresh_record: function (status) {
      if (status) {
        let that = this
        this.ddns.records_loading = true
        let params = {
          provider: this.ddns.settings.create.form.provider,
          domain: this.ddns.settings.create.form.domain,
          sub_domain: "",
          dnspod_id: this.ddns.settings.create.form.dnspod_id,
          dnspod_token: this.ddns.settings.create.form.dnspod_token,
        }
        if (this.ddns.settings.create.form.dnspod_id == "") {
          this.$message.error({message: "请先填写ID", type: 'error'})
          this.ddns.records_loading = false
          return
        } else if (this.ddns.settings.create.form.dnspod_token == "") {
          this.$message.error({message: "请先填写Token", type: 'error'})
          this.ddns.records_loading = false
          return
        }
        this.$api.post("/ddns/records", params).then(function (response) {
          that.ddns.records = response.detail
          that.ddns.records_loading = false
        }).catch(function (response) {
          that.ddns.records_loading = false
          that.$message.error({message: response.message, type: 'error'})
        })
      } else {
        this.records = []
      }
    },
    get_setting: function (id) {
      let that = this
      this.$api.post("/ddns/setting", {id: id}).then(function (response) {
        that.ddns.setting.visit = true
        that.ddns.setting.info = response.detail
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    },
    setting_change: function (value) {
      let that = this
      this.$api.post('/ddns/setting/refresh', {id: this.ddns.setting.info.setting.id, time_interval: value}).then(function (response) {
        that.get_setting(that.ddns.setting.info.setting.id)
        that.refresh_settings()
        that.$message({message: "修改成功", type: 'success'})
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    },
    disable_role: function () {
      let that = this
      this.$api.post('/ddns/setting/stop', {id: this.ddns.setting.info.setting.id}).then(function (response) {
        that.get_setting(that.ddns.setting.info.setting.id)
        that.refresh_settings()
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    },
    delete_role: function () {
      let that = this
      this.$api.post('/ddns/setting/remove', {id: this.ddns.setting.info.setting.id}).then(function (response) {
        that.$message({message: "删除成功", type: 'success'})
        that.refresh_settings()
        that.ddns.setting.visit = false
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
      })
    },
    start_role: function () {
      let that = this
      this.$api.post('/ddns/setting/start', {id: this.ddns.setting.info.setting.id}).then(function (response) {
        that.get_setting(that.ddns.setting.info.setting.id)
        that.refresh_settings()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    task_status_filter: function (status) {
      let s = ""
      switch (status) {
        case 0:
          s = "准备"
          break
        case 1:
          s = "运行"
          break
        case 2:
          s = "停止"
          break
        case 3:
          s = "重置"
          break
        case -1:
          s = "关闭"
          break  
        default:
          break
      }
      return s
    }
    
  },
  
  created: function () {
    this.refresh_settings()
    this.refresh_net_cards()
  },
  beforeDestroy () {},
  mounted: function () {}
};
</script>