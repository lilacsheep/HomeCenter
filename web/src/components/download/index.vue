<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-tabs tab-position="left" @tab-click="tabClick">
        <el-tab-pane label="下载中">
          <el-button-group>
            <el-button size="mini" type="primary" icon="el-icon-edit" @click="download.create.visit = true">创建下载</el-button>
            <el-button size="mini" type="primary" icon="el-icon-upload" @click="torrent.sync = true">上传种子</el-button>
          </el-button-group>
          <!-- <el-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
            <el-button slot="append" icon="el-icon-search"></el-button>
          </el-input> -->

          <el-table :data="download.tasks.not_finished" stripe size="mini" style="margin-top: 10px;">
            <el-table-column prop="gid" label="文件名" width="300">
              <template slot-scope="scope">
              {{getTaskName(scope.row)}}
              </template>
            </el-table-column>
            <el-table-column prop="totalLength" label="大小" width="100">
              <template slot-scope="scope">
                {{scope.row.totalLength | diskSize}}
              </template>
            </el-table-column>
            <el-table-column prop="completedLength" label="进度" width="80">
              <template slot-scope="scope">
                <span v-if="scope.row.completedLength == 0">0</span>
                <span v-else>{{(scope.row.completedLength / scope.row.totalLength * 100).toFixed(2)}} %</span>
              </template>
            </el-table-column>
            <el-table-column prop="downloadSpeed" label="速度">
              <template slot-scope="scope">
                {{scope.row.downloadSpeed | diskSize}}/秒
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="100">
              <template slot-scope="scope">
                <el-popconfirm v-if="scope.row.status == 1" title="是否继续该任务？" @onConfirm="start_task(scope.row)">
                  <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-caret-right"></el-button>
                </el-popconfirm>
                <el-popconfirm v-if="scope.row.status == 3" title="是否暂停该任务？" @onConfirm="cancel_task(scope.row)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-switch-button"></el-button>
                </el-popconfirm>
                <el-popconfirm title="是否删除该任务？" @onConfirm="remove_task(scope.row)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="已完成">
          <el-table :data="download.tasks.done" stripe size="mini" style="margin-top: 10px;">
            <el-table-column prop="gid" label="文件名" width="300">
              <template slot-scope="scope">
              {{getTaskName(scope.row)}}
              </template>
            </el-table-column>
            <el-table-column prop="totalLength" label="大小" width="100">
              <template slot-scope="scope">
                {{scope.row.totalLength | diskSize}}
              </template>
            </el-table-column>
            <el-table-column prop="md5" label="SH256"></el-table-column>
            <el-table-column label="操作" fixed="right" width="80">
              <template slot-scope="scope">
                <el-popconfirm title="是否删除该任务？" @onConfirm="remove_task(scope.row)">
                  <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="配置管理">
          <el-card>
            <el-form :model="settings.form" ref="ruleForm" label-width="100px" class="demo-ruleForm">
              <el-form-item label="下载路径" prop="path">
                <el-select v-model="settings.form.path" placeholder="请选择">
                  <el-option v-for="item in settings.nodes" :key="item.id" :label="item.name" :value="item.path"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="线程数量" prop="thread_num">
                <el-input-number size="small" v-model="settings.form.thread_num" :min="1" :step="2" :max="64"></el-input-number>
              </el-form-item>
              <el-form-item label="Aira2" prop="aria2_enable">
                <el-switch size="small" v-model="settings.form.aria2_enable"></el-switch>
              </el-form-item>
              <el-form-item label="地址" label-width="100px">
                <el-input size="small" v-model="settings.form.aria2_url" style="width: 220px"></el-input>
              </el-form-item>
              <el-form-item label="Token" label-width="100px">
                <el-input size="small" v-model="settings.form.aria2_token" style="width: 220px"></el-input>
              </el-form-item>
              <el-form-item>
                <el-button size="small" type="primary"  @click="submit_update_settings">立即更新</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </el-col>

    <el-dialog title="新建下载" :visible.sync="download.create.visit">
      <el-form :model="download.create.form" label-position="right">
        <el-form-item label="地址" label-width="100px">
          <el-input size="small" type="textarea" v-model="download.create.form.url" :rows="3"></el-input>
        </el-form-item>
        <el-form-item label="下载路径" prop="path" label-width="100px">
          <el-select size="small" v-model="download.create.form.path" placeholder="请选择">
            <el-option
              v-for="item in settings.nodes"
              :key="item.id"
              :label="item.name"
              :value="item.path">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="线程数量" prop="thread_num" label-width="100px">
          <el-input-number size="small" v-model="settings.form.thread_num" :min="1" :step="2" :max="64"></el-input-number>
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
  </el-row>
</template>

<script>

export default {
  data() {
    return {
      torrent: {
        sync: false
      },
      download: {
        tasks: {
          done: [],
          not_finished: []
        },
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
          path: "",
          thread_num: 0,
          aria2_enable: false,
          aria2_url: "",
          aria2_token: ""
        },
        nodes: []
      }
    }
  },
  methods: {
    submit_create_task () {
      let that = this
      this.$api.post('/download/create', this.download.create.form).then(function (response) {
        that.download.create.visit = false
        that.$message({message: '添加成功', type: 'success'})
        that.refresh_tasks()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    remove_task (item) {
      let that = this
      this.$api.post("/download/remove", {id: item.gid}).then(function (response) {
        that.$message({message: '删除成功', type: 'success'})
      })
    },
    cancel_task (item) {
      let that = this
      this.$api.post("/download/cancel", {id: item.id}).then(function (response) {
        that.$message({message: '暂停', type: 'success'})
      })
    },
    start_task (item) {
      let that = this
      this.$api.post("/download/start", {id: item.id}).then(function (response) {
        that.$message({message: '启动成功', type: 'success'})
      })
    },
    refresh_tasks () {
      let that = this
      this.$api.get("/download/tasks").then(function (response) {
        that.download.tasks.done = []
        that.download.tasks.not_finished = []
        if (response.detail) {
          response.detail["active"].forEach(function (item) {
            that.download.tasks.not_finished.push(item)          
          })
          response.detail["stopped"].forEach(function (item) {
            that.download.tasks.done.push(item)          
          })
        }
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
    filter_status(value, row) {
      return row.status === value;
    },
    refresh_nodes () {
      let that = this
      this.$api.get("/filesystem/nodes").then(function (response) {
        that.settings.nodes = response.detail
      })
    },
    tabClick: function (tab, event) {
      if (tab.index === '0') {
        this.timer = setInterval(this.refresh_tasks, 1000)
      } else if (tab.index === '1') {
        clearInterval(this.timer)
      } else if (tab.index === '2') {
        clearInterval(this.timer)
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
    this.refresh_nodes()
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

</style>