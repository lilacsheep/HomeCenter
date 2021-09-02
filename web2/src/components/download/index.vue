<template>
  <a-layout-content style="padding: 12px;">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item href="">
        文件下载
      </a-breadcrumb-item>
    </a-breadcrumb>
    <a-row :gutter="16" style="background: #fbfbfb;border: 1px solid #f4f4f4;height: 100%">
      <a-col>
        <a-tabs default-active-key="1" @change="tabClick">
          <a-tab-pane key="1" tab="下载列表">
            <a-button-group>
              <a-button size="mini" type="primary" icon="edit" @click="download.create.visit = true">创建下载</a-button>
              <a-button size="mini" type="primary" icon="upload" @click="torrent.sync = true">上传种子</a-button>
            </a-button-group>
          
          <!-- <a-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
            <el-button slot="append" icon="el-icon-search"></el-button>
          </a-input> -->
          <span style="float: right;">
            <a-radio-group v-model="task.query.status" size="mini">
              <a-radio-button size="small" value="全部">全部</a-radio-button>
              <a-radio-button size="small" value="下载中">下载中</a-radio-button>
              <a-radio-button size="small" value="已完成">已完成</a-radio-button>
              <a-radio-button size="small" value="其他">其他</a-radio-button>
            </a-radio-group>
            <a-tag color="green" size="small"><i class="arrow-up"></i>{{global.upload | diskSize}}/秒</a-tag>
            <a-tag color="orange" size="small"><i class="arrow-down"></i>{{global.download | diskSize}}/秒</a-tag>
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
            <a-card style="padding: 5px">
              <a-form-model :model="settings.form" :label-col="labelCol" :wrapper-col="wrapperCol">
                <a-form-model-item label="地址" label-width="100px">
                  <a-input size="small" v-model="settings.form.aria2_url" style="width: 220px"></a-input>
                </a-form-model-item>
                <a-form-model-item label="Token" label-width="100px">
                  <a-input size="small" v-model="settings.form.aria2_token" style="width: 220px"></a-input>
                </a-form-model-item>
                <a-form-model-item label="清理" label-width="100px">
                  <a-input v-model="settings.form.auto_clean" style="width: 220px" size="small">
                    <template slot="append">MB</template>
                  </a-input>
                  <span style="color: #909399;font-size: 11px">自动清理BT文件夹内不满足文件大小的文件, 0为关闭</span>
                </a-form-model-item>
                <a-form-model-item label="自动同步">
                  <el-select size="small" v-model="settings.form.auto_update_bt_tracker" placeholder="请选择">
                    <a-select-option label="关闭" value=""></a-select-option>
                    <a-select-option label="每小时" value="@hourly"></a-select-option>
                    <a-select-option label="每天" value="@every 24h"></a-select-option>
                  </el-select>
                  <span style="color: #909399;font-size: 11px">自动同步最新的BT服务器</span>
                </a-form-model-item>
                <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
                  <a-button size="small" type="primary"  @click="submit_update_settings">立即更新</a-button>
                </a-form-model-item>
              </a-form-model>
            </a-card>
          </a-tab-pane>
          <a-tab-pane key="3" tab="Aria2配置">
            <a-table :columns="aria2.columns"  :data="aria2.options" size="mini" style="margin-top: 10px;">
            </a-table>
          </a-tab-pane>
        </a-tabs>
        <el-dialog title="新建下载" :visible.sync="download.create.visit">
          <el-form :model="download.create.form" label-position="right">
            <a-form-model-item label="地址" label-width="100px">
              <a-input size="small" type="textarea" v-model="download.create.form.url" :rows="3"></a-input>
            </a-form-model-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button size="small" @click="download.create.visit = false">取 消</el-button>
            <el-button size="small" type="primary" @click="submit_create_task">确 定</el-button>
          </div>
        </el-dialog>
        <el-dialog custom-class="upload_diglog" :title="种子上传" :visible.sync="torrent.sync"  width="380px">
          <el-upload class="upload-demo" drag action="/api/download/torrent">
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          </el-upload>
        </el-dialog>
        <el-drawer :withHeader="false" :visible.sync="task.visible" :before-close="taskInfoClose" size="60%">
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
        </el-drawer>
      </a-col>
    </a-row>
  </a-layout-content>
</template>

<script>
import { aria2Api } from "../../scripts/aria2"
export default {
  data() {
    return {
      torrent: {
        sync: false
      },
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
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
      aria2Api.taskStatus(info.gid, function (response) {
        that.task.info.status = response.detail
        that.task.visible = true
      })
    },
    submit_create_task () {
      let that = this
      aria2Api.addUri(this.download.create.form.url, function (){
        that.download.create.visit = false
        that.refresh_tasks()
      })
    },
    remove_task (item) {
      aria2Api.removeTask(item.gid)
    },
    cancel_task (item) {
      aria2Api.pause(item.gid)
    },
    start_task (item) {
      aria2Api.unpause(item.gid)
    },
    refresh_tasks () {
      let that = this, tasks = []
      aria2Api.tasks({}, function (response) {
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
      })
      aria2Api.globalStat(function (response) {
        that.global.upload = response.detail.uploadSpeed
        that.global.download = response.detail.downloadSpeed
      })
    },
    refresh_settings () {
      let that = this
      this.$api.get("/download/settings").then(function (response) {
        that.settings.form = response.detail
        that.download.create.form.thread_num = response.detail.thread_num
        that.download.create.form.path = response.detail.path
      })
    },
    submit_update_settings () {
      let that = this
      this.$api.post("/download/settings/update", this.settings.form).then(function (response) {
        that.$message({message: '更新成功', type: 'success'})
      }).catch(function (resp) {
        that.$message({message: '更新失败', type: 'error'})
      })
    },
    refresh_global_options() {
      let that = this, options = []
      aria2Api.globalOptions(function (response) {
        that.aria2.options = response.detail
      })
    },
    tabClick: function (tab, event) {
      if (tab.index === '0') {
        if (!this.timer) {
          this.timer = setInterval(this.refresh_tasks, 1000)
        }
      } else if (tab.index === '1') {
        if (this.timer) {
          clearInterval(this.timer)
          this.timer = undefined
        } 
      } else if (tab.index === '2') {
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
    this.refresh_settings()
    this.refresh_tasks()
    this.timer = setInterval(this.refresh_tasks, 1000)
  },
  beforeDestroy () {
    clearInterval(this.timer)
  },
  mounted: function () {}
};
</script>

<style>
.el-card__header {
  padding: 5px;
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

.el-drawer__header {
  margin-bottom: 0px;
  padding: 10px 10px 0;
}

.el-drawer__body {
  padding: 5px;
}

.el-drawer__body table tr td {
  border: 1px solid #f2f2f2;
}
</style>