<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-tabs tab-position="left" @tab-click="tabClick">
        <el-tab-pane label="下载任务">
          <el-button-group>
            <el-button size="mini" type="primary" icon="el-icon-edit" @click="download.create.visit = true">创建下载</el-button>
            <el-button size="mini" type="primary" icon="el-icon-upload" @click="torrent.sync = true">上传种子</el-button>
          </el-button-group>
          
          <!-- <el-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
            <el-button slot="append" icon="el-icon-search"></el-button>
          </el-input> -->
          <span style="float: right;">
            <el-radio-group v-model="task.query.status" size="mini">
              <el-radio-button label="全部"></el-radio-button>
              <el-radio-button label="下载中"></el-radio-button>
              <el-radio-button label="已完成"></el-radio-button>
            </el-radio-group>
            <el-tag type="success" size="small"><i class="el-icon-top"></i>{{global.upload | diskSize}}/秒</el-tag>
            <el-tag type="danger" size="small"><i class="el-icon-bottom"></i>{{global.download | diskSize}}/秒</el-tag>
          </span>
          <el-table :data="download.tasks" stripe size="mini" style="margin-top: 10px;">
            <el-table-column prop="gid" label="文件名">
              <template slot-scope="scope">
                <el-button size="mini" type="text" @click="taskInfoOpen(scope.row)">{{getTaskName(scope.row)}}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="totalLength" label="大小" width="100">
              <template slot-scope="scope">
                {{scope.row.totalLength | diskSize}}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template slot-scope="scope">
                <span v-if="scope.row.status == 'active'">{{(scope.row.completedLength / scope.row.totalLength * 100).toFixed(2)}}%</span>
                <span v-else-if="scope.row.status == 'stopped'">已停止</span>
                <span v-else-if="scope.row.status == 'paused'">已暂停</span>
                <span v-else-if="scope.row.status == 'complete'">已完成</span>
                <span v-else>{{scope.row.status}}</span>
              </template>
            </el-table-column>

            <el-table-column prop="downloadSpeed" label="下载" width="100">
              <template slot-scope="scope">
                {{scope.row.downloadSpeed | diskSize}}/秒
              </template>
            </el-table-column>
            <el-table-column prop="uploadSpeed" label="上传" width="100">
              <template slot-scope="scope">
                {{scope.row.uploadSpeed | diskSize}}/秒
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="100">
              <template slot-scope="scope">
                <el-popconfirm v-if="scope.row.status == 'paused'" title="是否继续该任务？" @onConfirm="start_task(scope.row)">
                  <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-caret-right"></el-button>
                </el-popconfirm>
                <el-popconfirm v-if="scope.row.status == 'active'" title="是否暂停该任务？" @onConfirm="cancel_task(scope.row)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-switch-button"></el-button>
                </el-popconfirm>
                <el-popconfirm title="是否删除该任务？" @onConfirm="remove_task(scope.row)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="配置管理">
          <el-card body-style="padding: 5px">
            <el-form :model="settings.form" ref="ruleForm" label-width="100px" class="demo-ruleForm">
              <el-form-item label="地址" label-width="100px">
                <el-input size="small" v-model="settings.form.aria2_url" style="width: 220px"></el-input>
              </el-form-item>
              <el-form-item label="Token" label-width="100px">
                <el-input size="small" v-model="settings.form.aria2_token" style="width: 220px"></el-input>
              </el-form-item>
              <el-form-item label="清理" label-width="100px">
                <el-input v-model="settings.form.auto_clean" style="width: 220px" size="small">
                  <template slot="append">MB</template>
                </el-input>
                <span style="color: #909399;font-size: 11px">自动清理BT文件夹内不满足文件大小的文件, 0为关闭</span>
              </el-form-item>
              <el-form-item label="自动同步">
                <el-select size="small" v-model="settings.form.auto_update_bt_tracker" placeholder="请选择">
                  <el-option label="关闭" value=""></el-option>
                  <el-option label="每分钟" value="@every 1m"></el-option>
                  <el-option label="每小时" value="@hourly"></el-option>
                  <el-option label="半小时" value="@every 30m"></el-option>
                  <el-option label="每天" value="@every 24h"></el-option>
                </el-select>
                <span style="color: #909399;font-size: 11px">自动同步最新的BT服务器</span>
              </el-form-item>
              <el-form-item>
                <el-button size="small" type="primary"  @click="submit_update_settings">立即更新</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
        <el-tab-pane label="Aria2配置">
          <el-table :data="aria2.options" stripe size="mini" style="margin-top: 10px;">
            <el-table-column prop="key" label="配置"></el-table-column>
            <el-table-column prop="value" label="值"></el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-col>

    <el-dialog title="新建下载" :visible.sync="download.create.visit">
      <el-form :model="download.create.form" label-position="right">
        <el-form-item label="地址" label-width="100px">
          <el-input size="small" type="textarea" v-model="download.create.form.url" :rows="3"></el-input>
        </el-form-item>
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
      <el-table :data="task.info.status.files" stripe size="mini" style="margin-top: 10px;" max-height="500">
         <el-table-column prop="path" label="文件">
           <template slot-scope="scope">
             {{scope.row.path.split("/").slice(-1)[0]}}
           </template>
         </el-table-column>
         <el-table-column prop="completedLength" label="进度" width="100">
           <template slot-scope="scope">{{(scope.row.completedLength / scope.row.length * 100).toFixed(2)}}%</template>
          </el-table-column>
         <el-table-column prop="length" label="大小" width="120">
           <template slot-scope="scope">
             {{scope.row.length | diskSize}}
           </template>
         </el-table-column>
      </el-table>
    </el-drawer>
  </el-row>
</template>

<script>
import { aria2Api } from "../../scripts/aria2"
export default {
  data() {
    return {
      torrent: {
        sync: false
      },
      aria2: {
        options: []
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
        } else {
          response.detail.forEach(function (item) {
            if ((item.status != "active") && (item.status != "error") && (!item.followedBy)) {
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