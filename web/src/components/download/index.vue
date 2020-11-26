<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-tabs tab-position="left" @tab-click="tabClick">
        <el-tab-pane label="下载中">
          <el-button-group>
            <el-button size="mini" type="primary" icon="el-icon-edit" @click="download.create.visit = true">创建下载</el-button>
          </el-button-group>
          <!-- <el-input size="mini" placeholder="请输入内容" v-model="roles.filter" style="width: 250px;float: right;">
            <el-button slot="append" icon="el-icon-search"></el-button>
          </el-input> -->

          <el-table :data="download.tasks.not_finished" stripe size="mini" style="margin-top: 10px;">
            <el-table-column prop="file_name" label="文件名" width="300"></el-table-column>
            <el-table-column prop="total_size" label="大小" width="100">
              <template slot-scope="scope">
                {{scope.row.total_size | diskSize}}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template slot-scope="scope">
                <el-tag v-if="scope.row.status == 0" size="mini">暂停</el-tag>
                <el-tag v-else-if="scope.row.status == 1" size="mini">等待中</el-tag>
                <el-tag v-else-if="scope.row.status == 2" size="mini">下载中</el-tag>
                <el-tag v-else-if="scope.row.status == 3" size="mini">已完成</el-tag>
                <el-tag v-else-if="scope.row.status == 99" size="mini">错误</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="progress" label="进度" width="80">
              <template slot-scope="scope">
                <span v-if="scope.row.progress == 0">0</span>
                <span v-else>{{scope.row.progress}} %</span>
              </template>
            </el-table-column>
            <el-table-column prop="speed" label="速度">
              <template slot-scope="scope">
                <span v-if="scope.row.status == 2">{{scope.row.speed | diskSize}}/秒</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="100">
              <template slot-scope="scope">
                <el-popconfirm v-if="scope.row.status == 0" title="是否继续该任务？" @onConfirm="start_task(scope.row)">
                  <el-button slot="reference" style="color: green" type="text" size="mini" icon="el-icon-caret-right"></el-button>
                </el-popconfirm>
                <el-popconfirm v-if="scope.row.status == 2" title="是否暂停该任务？" @onConfirm="cancel_task(scope.row)">
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
            <el-table-column prop="file_name" label="文件名" width="300"></el-table-column>
            <el-table-column prop="total_size" label="大小" width="100">
              <template slot-scope="scope">
                {{scope.row.total_size | diskSize}}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态"  width="80">
              <template slot-scope="scope">
                <el-tag v-if="scope.row.status == 0" size="mini">暂停</el-tag>
                <el-tag v-else-if="scope.row.status == 1" size="mini">等待中</el-tag>
                <el-tag v-else-if="scope.row.status == 2" size="mini">下载中</el-tag>
                <el-tag v-else-if="scope.row.status == 3" size="mini">已完成</el-tag>
                <el-tag v-else-if="scope.row.status == 99" size="mini">错误</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="md5" label="MD5"></el-table-column>

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
              <el-form-item label="通知开关" prop="notify_open">
                <el-switch size="small" v-model="settings.form.notify_open"></el-switch>
              </el-form-item>
              <el-form-item label="通知消息" prop="notify_message">
                <el-input size="small" type="textarea" v-model="settings.form.notify_message"></el-input>
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
  </el-row>
</template>

<script>

export default {
  data() {
    return {
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
          notify_open: true,
          notify_message: ""
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
      this.$api.post("/download/remove", {id: item.id}).then(function (response) {
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
          response.detail.forEach(function (item) {
            if (item.status == 3) {
              that.download.tasks.done.push(item)
            } else {
              that.download.tasks.not_finished.push(item)
            }
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

      } else if (tab.index === '2') {
        clearInterval(this.timer)
      }
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