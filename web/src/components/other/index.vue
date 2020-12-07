<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-tabs tab-position="left" @tab-click="tabClick">
        <el-tab-pane label="DDNS">
          <el-button-group>
            <el-button size="mini" type="primary" icon="el-icon-edit" @click="ddns.settings.create.visit = true">新增同步</el-button>
          </el-button-group>

          <el-table :data="ddns.roles" stripe size="mini" style="margin-top: 10px;">
            <el-table-column prop="domain" label="域名" width="150">
              <template slot-scope="scope">
                <el-button type="text" @click="get_setting(scope.row.id)" >{{`${scope.row.sub_domain}.${scope.row.domain}`}}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="provider" label="渠道" width="100"></el-table-column>
            <el-table-column prop="net_card" label="网卡" width="80">
              <template slot-scope="scope">
                {{search_net_card(scope.row.net_card)}}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template slot-scope="scope">
                <i v-if="scope.row.status" style="color: green" class="el-icon-success"></i>
                <i v-else style="color: red" class="el-icon-error"></i>
              </template>
            </el-table-column>
            <el-table-column prop="updated_on" label="更新时间"></el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="内网穿透"></el-tab-pane>
      </el-tabs>
    </el-col>
    <el-dialog title="新建同步" :visible.sync="ddns.settings.create.visit" width="500px">
      <el-form :model="ddns.settings.create.form" label-position="right" label-width="100px">
        <el-form-item label="域名">
          <el-input size="small" placeholder="域名" v-model="ddns.settings.create.form.domain" style="width: 300px"></el-input>
        </el-form-item>
        <el-form-item label="模式">
          <el-switch v-model="ddns.settings.create.form.mode" active-text="本地网卡" inactive-text="公网IP" @change="ddns_mode_change"></el-switch>
        </el-form-item>
        <el-form-item label="网卡" :style="ddns.cards_style">
          <el-select size="small" v-loading="ddns.cards_style_loading" v-model="ddns.settings.create.form.net_card" placeholder="请选择">
            <el-option v-for="item in ddns.net_cards" :key="item.index" :label="item.name" :value="item.index" :disabled="item.flags && item.flags[0] !== 'up'"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="同步间隔">
          <el-select size="small" v-model="ddns.settings.create.form.time_interval" placeholder="请选择">
            <el-option label="每小时" value="@hourly"></el-option>
            <el-option label="半小时" value="@every 30m"></el-option>
            <el-option label="每天" value="@every 24h"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="供应商">
          <el-radio-group v-model="ddns.settings.create.form.provider" size="mini" @change="provider_change">
            <el-radio-button label="dnspod"></el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="DNSPod" :style="ddns.dnspod_style">
          <el-input size="small" placeholder="ID" v-model="ddns.settings.create.form.dnspod_id" style="width: 120px"></el-input>
          <el-input size="small" placeholder="Token" v-model="ddns.settings.create.form.dnspod_token" style="width: 200px"></el-input>
        </el-form-item>
        <el-form-item>
          <el-switch v-model="ddns.add_mode" inactive-text="新增记录" active-text="原有记录"  @change="ddns_add_change"></el-switch>
        </el-form-item>
        <el-form-item label="子域名" :style="ddns.sub_domain_style">
          <el-input size="small" placeholder="子域名" v-model="ddns.settings.create.form.sub_domain" style="width: 120px" @blur="check_subdomain"></el-input>
        </el-form-item>
        <el-form-item label="记录" :style="ddns.records_style">
          <el-select size="small" :loading="ddns.records_loading" @visible-change="refresh_record" @change="record_change" v-model="ddns.settings.create.form.record_id" placeholder="请选择">
            <el-option v-for="item in ddns.records" :key="item.id" :label="`${item.name}.${ddns.settings.create.form.domain}(${item.type})`" :value="item.id" :disabled="item.type == 'NS'"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="ddns.settings.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_settings">确 定</el-button>
      </div>
    </el-dialog>

    <el-drawer title="我是标题" :visible.sync="ddns.setting.visit" :with-header="false" size="40%">
      <el-button-group style="margin-bottom: 10px;">
        <el-button type="warning" v-if="ddns.setting.info.setting.status" size="small" icon="el-icon-circle-close" @click="disable_role">关闭同步</el-button>
        <el-button type="success" v-else size="small" icon="el-icon-circle-close" @click="start_role">开启同步</el-button>
        <el-button type="danger" size="small" icon="el-icon-delete" @click="delete_role">删除记录</el-button>
      </el-button-group>
      <table class="descriptions">
        <tbody>
          <tr>
            <th class="title">域名</th>
            <td class="details">{{ddns.setting.info.setting.sub_domain}}.{{ddns.setting.info.setting.domain}}</td>
          </tr>
          <tr>
            <th class="title">运营商</th>
            <td class="details">{{ddns.setting.info.setting.provider}}</td>
          </tr>
          <tr>
            <th class="title">状态</th>
            <td class="details">
              <i v-if="ddns.setting.info.setting.status" style="color: green" class="el-icon-success"></i>
              <i v-else style="color: red" class="el-icon-error"></i>
            </td>
          </tr>
          <tr>
            <th class="title">同步周期</th>
            <td class="details">
              <el-select size="small" v-model="ddns.setting.info.setting.time_interval" @change="setting_change" placeholder="请选择">
                <el-option label="每小时" value="@hourly"></el-option>
                <el-option label="半小时" value="@every 30m"></el-option>
                <el-option label="每天" value="@every 1d"></el-option>
              </el-select>
            </td>
          </tr>
          <tr>
            <th class="title">同步模式</th>
            <td class="details">
              <el-tag size="small" type="success" v-if="ddns.setting.info.setting.use_public_ip">公网地址</el-tag>
              <el-tag size="small" type="warning" v-else>网卡模式({{search_net_card(ddns.setting.info.setting.net_card)}})</el-tag>
            </td>
          </tr>
          <tr>
            <th class="title">任务状态</th>
            <td class="details">
              <el-tag size="small" effect="plain">{{task_status_filter(ddns.setting.info.entry.status)}}</el-tag>
            </td>
          </tr>
          <tr>
            <th class="title">更新记录</th>
            <td class="details">
              <el-timeline>
                <el-timeline-item
                  v-for="(info, index) in ddns.setting.info.setting.history"
                  :key="index"
                  :type="info.status ? 'danger': 'success'"
                  :timestamp="info.date" placement="top">
                  <span v-if="info.status == 0">同步: {{info.value}}</span>
                  <span v-else>同步失败: {{info.error}}</span>
                  {{info.error}}
                </el-timeline-item>
              </el-timeline>
            </td>
          </tr>
        </tbody>
      </table>
    </el-drawer>
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      ddns: {
        records: [],
        roles: [],
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
    tabClick: function (tab, event) {
      if (tab.index === '0') {
      } else if (tab.index === '1') {

      } else if (tab.index === '2') {
      }
    },
    submit_create_settings: function () {
      let that = this
      this.$api.post("/ddns/setting/create", this.ddns.settings.create.form).then(function (response) {
        that.$message({message: '添加成功', type: 'success'})
        that.refresh_settings()
      }).catch(function (response) {
        that.$message({message: response.message, type: 'error'})
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