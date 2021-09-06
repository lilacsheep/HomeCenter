<template>
  <a-layout-content style="padding: 12px;">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item href="">
        文件下载
      </a-breadcrumb-item>
    </a-breadcrumb>
    <a-row :gutter="16" style="border: 1px solid #f4f4f4;height: 100%">
      <a-col>
        <a-tabs default-active-key="1" @change="tabClick" tab-position="left">
          <a-tab-pane key="1" tab="下载列表">
            <a-button-group>
              <a-button type="primary" icon="edit" @click="download.create.visit = true">创建下载</a-button>
              <a-upload name="file" :fileList="false" :loading="download.uploadloading" :multiple="false" @change="uploadChange" accept=".torrent" action="/api/download/torrent">
                <a-button  type="primary" icon="upload">上传种子</a-button>
              </a-upload>
            </a-button-group>
            <span style="float: right;">
              <a-radio-group v-model="task.query.status" size="small" button-style="solid">
                <a-radio-button value="全部">全部</a-radio-button>
                <a-radio-button value="下载中">下载中</a-radio-button>
                <a-radio-button value="已完成">已完成</a-radio-button>
                <a-radio-button value="其他">其他</a-radio-button>
              </a-radio-group>
              <a-tag color="green" size="small"><a-icon type="arrow-up"></a-icon>{{global.upload | diskSize}}/秒</a-tag>
              <a-tag color="orange" size="small"><a-icon type="arrow-down"></a-icon>{{global.download | diskSize}}/秒</a-tag>
            </span>
            <a-table :data="download.tasks" stripe size="mini" style="margin-top: 10px;">
              <a-table-column prop="gid" label="文件名">
                <template slot-scope="scope">
                  <el-button size="mini" type="text" @click="taskInfoOpen(scope.row)">{{getTaskName(scope.row)}}</el-button>
                </template>
              </a-table-column>
              <a-table-column prop="totalLength" label="大小" width="100">
                <template slot-scope="scope">
                  {{scope.row.totalLength | diskSize}}
                </template>
              </a-table-column>
              <a-table-column prop="status" label="状态" width="100">
                <template slot-scope="scope">
                  <span v-if="scope.row.status == 'active'">{{(scope.row.completedLength / scope.row.totalLength * 100).toFixed(2)}}%</span>
                  <span v-else-if="scope.row.status == 'stopped'">已停止</span>
                  <span v-else-if="scope.row.status == 'paused'">已暂停</span>
                  <span v-else-if="scope.row.status == 'complete'">已完成</span>
                  <span v-else>{{scope.row.status}}</span>
                </template>
              </a-table-column>

              <a-table-column prop="downloadSpeed" label="下载" width="100">
                <template slot-scope="scope">
                  {{scope.row.downloadSpeed | diskSize}}/秒
                </template>
              </a-table-column>
              <a-table-column prop="uploadSpeed" label="上传" width="100">
                <template slot-scope="scope">
                  {{scope.row.uploadSpeed | diskSize}}/秒
                </template>
              </a-table-column>
              <a-table-column label="操作" fixed="right" width="100">
                <template slot-scope="scope">
                  <el-popconfirm v-if="scope.row.status == 'paused'" title="是否继续该任务？" @onConfirm="start_task(scope.row)">
                    <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-caret-right"></el-button>
                  </el-popconfirm>
                  <el-popconfirm v-if="scope.row.status == 'error'" title="是否继续该任务？" @onConfirm="start_task(scope.row)">
                    <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-caret-right"></el-button>
                  </el-popconfirm>
                  <el-popconfirm v-if="scope.row.status == 'active'" title="是否暂停该任务？" @onConfirm="cancel_task(scope.row)">
                    <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-switch-button"></el-button>
                  </el-popconfirm>
                  <el-popconfirm title="是否删除该任务？" @onConfirm="remove_task(scope.row)">
                    <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
                  </el-popconfirm>
                </template>
              </a-table-column>
            </a-table>
          </a-tab-pane>
          <a-tab-pane key="2" tab="配置管理">
            <a-card style="padding: 0">
              <a-form-model :model="settings.form" :label-col="labelCol" :wrapper-col="wrapperCol">
                <a-form-model-item label="地址">
                  <a-input v-model="settings.form.aria2_url" style="width: 300px"></a-input>
                </a-form-model-item>
                <a-form-model-item label="Token">
                  <a-input v-model="settings.form.aria2_token" style="width: 300px"></a-input>
                </a-form-model-item>
                <a-form-model-item label="清理">
                  <a-input v-model="settings.form.auto_clean" style="width: 300px" addon-after="MB"></a-input>
                  <span style="color: #909399;font-size: 11px">自动清理BT文件夹内不满足文件大小的文件, 0为关闭</span>
                </a-form-model-item>
                <a-form-model-item label="自动同步">
                  <a-select v-model="settings.form.auto_update_bt_tracker" placeholder="请选择" style="width: 300px">
                    <a-select-option value="">关闭</a-select-option>
                    <a-select-option value="@hourly">每小时</a-select-option>
                    <a-select-option value="@every 24h">每天</a-select-option>
                  </a-select>
                  <span style="color: #909399;font-size: 11px">自动同步最新的BT服务器</span>
                </a-form-model-item>
                <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
                  <a-button type="primary"  @click="submit_update_settings">立即更新</a-button>
                </a-form-model-item>
              </a-form-model>
            </a-card>
          </a-tab-pane>
          <a-tab-pane key="3" tab="Aria2配置">
            <a-table :columns="aria2.columns"  :data="aria2.options" size="small">
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </a-col>
    </a-row>
    <a-modal title="新建下载" v-model="download.create.visit">
      <template slot="footer">
        <a-button @click="download.create.visit = false">取 消</a-button>
        <a-button type="primary" @click="submit_create_task">确 定</a-button>
      </template>
      <a-form-model :model="download.create.form" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-model-item label="地址">
          <a-input type="textarea" v-model="download.create.form.url" :rows="3"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
    <a-drawer :withHeader="false" :visible.sync="task.visible" :before-close="taskInfoClose" size="60%">
      <table style="border: 1px solid #f2f2f2" width="100%">
        <tr>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">文件名</td>
          <td colspan="3">{{task.info.filename}}</td>
        </tr>
        <tr>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">GID</td>
          <td colspan="3">{{task.info.status.gid}}</td>
        </tr>
        <tr>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">错误</td>
          <td colspan="3">{{task.info.status.errorMessage}}</td>
        </tr>
        <tr>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">大小</td>
          <td>{{task.info.status.totalLength | diskSize}}</td>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">状态</td>
          <td>{{task.info.status.status}}</td>
        </tr>
        <tr>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">已上传</td>
          <td>{{task.info.status.uploadLength | diskSize}}</td>
          <td width="100px" style="background-color: #f2f2f2;padding: 0">已完成</td>
          <td>{{task.info.status.completedLength | diskSize}}</td>
        </tr>
      </table>
      <a-table :data="task.info.status.files" stripe size="mini" style="margin-top: 10px;" max-height="500">
        <a-table-column prop="path" label="文件">
          <template slot-scope="scope">
            {{scope.row.path.split("/").slice(-1)[0]}}
          </template>
        </a-table-column>
        <a-table-column prop="completedLength" label="进度" width="100">
          <template slot-scope="scope">{{(scope.row.completedLength / scope.row.length * 100).toFixed(2)}}%</template>
          </a-table-column>
        <a-table-column prop="length" label="大小" width="120">
          <template slot-scope="scope">
            {{scope.row.length | diskSize}}
          </template>
        </a-table-column>
      </a-table>
    </a-drawer>
  </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      torrent: {
        sync: false
      },
      labelCol: { span: 2 },
      wrapperCol: { span: 16 },
      aria2: {
        options: [],
        columns: [
          {title: '键', dataIndex: 'key', key: 'key'},
          {title: '值', dataIndex: 'value', key: 'value'},
          ]
      },
      task: {
        visible: false,
        query: {
          status: "全部"
        },
        info: {
          filename: "",
          status: {}
        }
      },
      global: {
        upload: 0,
        download: 0
      },
      download: {
        tasks: [],
        uploadloading: false,
        create: {
          visit: false,
          form: {
            url: undefined,
            path: "",
            thread_num: 0,
          },
        }
      },
      settings: {
        form: {
          aria2_url: "",
          aria2_token: "",
          auto_clean: 0,
          auto_update_bt_tracker: ""
        },
        nodes: []
      }
    }
  },
  methods: {
    taskInfoClose() {
      this.task.visible = false
    },
    taskInfoOpen(info) {
      let that = this
      this.task.info.filename = this.getTaskName(info)

      this.$api.aria2_task_status(info.gid).then(function (response) {
        that.task.info.status = response.detail
        that.task.visible = true
      }).catch(function(response) {
        that.$message.error('获取任务信息失败：'+response.message)
      })
    },
    submit_create_task () {
      let that = this
      this.$api.aria2_add_uri(this.download.create.form.url).then(function (response) {
        that.download.create.visit = false
        that.refresh_tasks()
      })
    },
    remove_task (item) {
      this.$api.aria2_remove_task(item.gid)
    },
    cancel_task (item) {
      this.$api.aria2_task_pause(item.gid)
    },
    start_task (item) {
      this.$api.aria2_task_unpause(item.gid)
    },
    refresh_tasks () {
      if (this.settings.form.aria2_url == "") {
        return
      }
      let that = this, tasks = []
      this.$api.aria2_tasks().then(function (response) {
        if (that.task.query.status === "全部") {
          response.detail.forEach(function (item) {
            if ((item.status != "error") && (!item.followedBy)){
              tasks.push(item)
            }
          })
          that.download.tasks = tasks
        } else if (that.task.query.status === '下载中') {
          response.detail.forEach(function (item) {
            if (item.status == "active") {
              tasks.push(item)
            }
          }) 
          that.download.tasks = tasks
        } else if (that.task.query.status === '已完成') {
          response.detail.forEach(function (item) {
            if ((item.status == "complete") && (item.status != "error") && (!item.followedBy)) {
              tasks.push(item)
            }
          })
          that.download.tasks = tasks
        } else {
          response.detail.forEach(function (item) {
            if ((item.status != "active") && (item.status != "complete") && (!item.followedBy)) {
              tasks.push(item)
            }
          })
          that.download.tasks = tasks
        }
      }).catch(function (response) {
        that.$message.error("刷新任务列表失败: "+response.message)
      })
      this.$api.aria2_global_stats().then(function (response) {
        that.global.upload = response.detail.uploadSpeed
        that.global.download = response.detail.downloadSpeed
      })
    },
    refresh_settings () {
      let that = this
      this.$api.aria2_refresh_settings().then(function (response) {
        that.settings.form = response.detail
        that.download.create.form.thread_num = response.detail.thread_num
        that.download.create.form.path = response.detail.path
      })
    },
    submit_update_settings () {
      let that = this
      this.$api.aria2_settings_update(this.settings.form).then(function (response) {
        that.$message.success('更新成功')
      }).catch(function (resp) {
        that.$message.error('更新失败')
      })
    },
    refresh_global_options() {
      let that = this, options = []
      this.$api.aria2_global_options().then(function (response) {
        that.aria2.options = response.detail
      }).catch(function (response) {
        that.$message.error('获取全局配置失败:' + response.message)
      })
    },
    tabClick: function (key) {
      if (key === '1') {
        if (!this.timer) {
          this.timer = setInterval(this.refresh_tasks, 1000)
        }
      } else if (key === '2') {
        if (this.timer) {
          clearInterval(this.timer)
          this.timer = undefined
        } 
      } else if (key === '3') {
        if (this.timer) {
          clearInterval(this.timer)
          this.timer = undefined
        }
        this.refresh_global_options()
      }
    },
    getTaskName: function(info) {
      let taskName = "Unknown"
      if (info.bittorrent && info.bittorrent.info) {
        taskName = info.bittorrent.info.name;
      }
      if (!taskName && info.files && info.files.length > 0) {
        taskName = this.getFileName(info.files[0]);
      }
      return taskName
    },
    uploadChange(info) {
      if (info.file.status !== 'uploading') {
        this.download.uploadloading = true
      }
      if (info.file.status === 'done') {
        this.download.uploadloading = false
      }
    },
    getFileName: function (file) {
      if (!file) {
          return '';
      }

      var path = file.path;
      var needUrlDecode = false;

      if (!path && file.uris && file.uris.length > 0) {
          path = file.uris[0].uri;
          needUrlDecode = true;
      }

      var index = path.lastIndexOf('/');

      if (index <= 0 || index === path.length) {
          return path;
      }

      var fileNameAndQueryString = path.substring(index + 1);
      var queryStringStartPos = fileNameAndQueryString.indexOf('?');
      var fileName = fileNameAndQueryString;

      if (queryStringStartPos > 0) {
          fileName = fileNameAndQueryString.substring(0, queryStringStartPos);
      }

      if (needUrlDecode) {
        fileName = decodeURI(fileName);
      }

      return fileName;
    }
  },
  created: function () {
    let that = this
    this.$api.aria2_refresh_settings().then(function (response) {
      that.settings.form = response.detail
      that.download.create.form.thread_num = response.detail.thread_num
      that.download.create.form.path = response.detail.path
    }).catch(function(response) {
      that.$message.error("加载配置失败: "+response.message)
    })
    this.$api.aria2_tasks().then(function (response) {
      that.refresh_tasks()
      that.timer = setInterval(this.refresh_tasks, 1000)
    }).catch(function(response) {
      that.$message.error("任务列表获取失败: "+response.message)
    })
  },
  beforeDestroy () {
    clearInterval(this.timer)
  },
  mounted: function () {}
};
</script>

<style>

</style>