<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-breadcrumb v-if="node.files.length > 0" separator-class="el-icon-arrow-right">
        <el-breadcrumb-item><a @click="select_index">首页</a></el-breadcrumb-item>
        <el-breadcrumb-item v-for="(dir, index) in node.dirs" :key="index"><a @click="select_node(node)">{{dir.name}}</a></el-breadcrumb-item>
      </el-breadcrumb>
    </el-col>
    
    <el-col :span="24" style="margin-top: 15px;">
      <el-button size="small" type="primary" @click.prevent="node.create.visit = true">新增节点</el-button>
      <el-button size="small" type="success" @click.prevent="node.file.visit = true" plain>上传文件</el-button>
      <el-button size="small" type="warning" @click="open_create_dir_dialog" plain>新增目录</el-button>
    </el-col>
    
    <el-col :span="4" v-for="(node, index) in nodes" :key="index">
      <el-card style="margin-top: 15px;" :body-style="{ padding: '0' }">
        <img :src="icon.folder" style="margin: 0 5px" @click="select_node(node)">
        <div style="padding: 0 5px;">
          <span @click="select_node(node)">{{node.name}}</span>
          <el-dropdown style="float: right;">
            <el-button type="text" class="icon"  icon="el-icon-more-outline"></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>下载</el-dropdown-item>
              <el-dropdown-item>删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-card>
    </el-col>

    <el-col v-if="node.files !== null && node.files.length > 0" :span="24" style="margin-top: 15px;">
      <el-table :data="node.files" :highlight-current-row="true" stripe row-key="path" size="mini" :tree-props="{children: 'children', hasChildren: 'hasChildren'}" @row-click="table_row_click"	:default-sort="{prop: 'is_dir', order: 'descending'}">
        <el-table-column prop="name" label="名称" width="300" :show-overflow-tooltip="true">
          <template slot-scope="scope">
            <span v-if="scope.row.is_dir">
              <i class="el-icon-folder"> {{scope.row.name}}</i>
            </span>
            <el-button size="mini" v-else icon="el-icon-document" @click="view_file(scope.row)" type="text">{{scope.row.name}}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100">
          <template slot-scope="scope">
            <span v-if="!scope.row.is_dir">{{scope.row.size | diskSize}}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="类型" prop="is_dir" sortable></el-table-column>
        <el-table-column label="日期" prop="create_at"></el-table-column>
        <el-table-column label="操作" fixed="right" width="80">
          <template slot-scope="scope">
            <el-popconfirm v-if="!scope.row.has_children" title="是否删除该文件"  @onConfirm="remove_file(scope.row)">
              <el-button slot="reference" style="color: red" type="text" size="mini" icon="el-icon-delete"></el-button>
            </el-popconfirm>
            <el-button v-if="!scope.row.is_dir" type="text" size="mini" icon="el-icon-download" @click="download_file(scope.row)"></el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-col>
    <el-dialog title="新建节点" :visible.sync="node.create.visit"  width="28%">
      <el-form :model="node.create.form" label-position="right">
        <el-form-item label="名称" label-width="100px">
          <el-input v-model="node.create.form.name" size="small" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="路径" label-width="100px">
          <el-input v-model="node.create.form.path" size="small" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="node.create.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="submit_create_node">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="新建目录" :visible.sync="node.dir.visit"  width="28%">
      <el-form :model="node.dir.form" label-position="right">
        <el-form-item label="路径" label-width="100px">
          <el-input v-model="node.dir.form.path" size="small" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="node.dir.visit = false">取 消</el-button>
        <el-button size="small" type="primary" @click="creat_node_dir">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="上传文件" :visible.sync="node.file.visit"  width="380px">
      <el-upload :data='{"path": node.file.upload_path, "node_id": node.info.id}' class="upload-demo" drag action="/api/filesystem/file/upload" :on-success="upload_success" multiple>
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
      </el-upload>
    </el-dialog>

    <el-dialog :visible.sync="node.view.visit" :width="node.view.width" @closed="close_view">
      <span v-if="node.view.type === 'img'">
        <img :src="node.view.src" fit="fill" :style="node.view.style" class="image">
      </span>
      <span v-if="node.view.type === 'video'">
        <video-player id="videobox" class="video-player vjs-custom-skin" controls :width="node.view.width" :playsinline="true" :options="node.playerOptions" ref='videoRef'>
        </video-player>
      </span>
    </el-dialog>

  </el-row>
</template>

<script>

export default {
  data() {
    return {
      icon: {
          folder: require("../../assets/icon/folder.png"),
          file: require("../../assets/icon/file.png"),
          avi: require("../../assets/icon/avi.png"),
          file_style: "margin: 0 30px",
          folder_style: "margin: 0 5px",
      },
      nodes: [],
      node: {
        info: {},
        files: [],
        dirs: [],
        create: {
          visit: false,
          form: {
            name: "",
            path: ""
          }
        },
        dir: {
          visit: false,
          form: {
            path: "",
            node: ""
          }
        },
        file: {
          visit: false,
          upload_path: ""
        },
        view: {
          width: "500px",
          visit: false,
          style: "",
          src: ""
        },
        playerOptions: {
          playbackRates: [0.5, 1.0, 1.5, 2.0], // 可选的播放速度
            autoplay: false, // 如果为true,浏览器准备好时开始回放。
            muted: false, // 默认情况下将会消除任何音频。
            loop: false, // 是否视频一结束就重新开始。
            preload: 'auto', // 建议浏览器在<video>加载元素后是否应该开始下载视频数据。auto浏览器选择最佳行为,立即开始加载视频（如果浏览器支持）
            language: 'zh-CN',
            aspectRatio: '16:9', // 将播放器置于流畅模式，并在计算播放器的动态大小时使用该值。值应该代表一个比例 - 用冒号分隔的两个数字（例如"16:9"或"4:3"）
            fluid: true, // 当true时，Video.js player将拥有流体大小。换句话说，它将按比例缩放以适应其容器。
            sources: [{
              type: "video/mp4", // 类型
              src: "" // url地址
            }],
            poster: '', // 封面地址
            notSupportedMessage: '此视频暂无法播放，请稍后再试', // 允许覆盖Video.js无法播放媒体源时显示的默认信息。
            controlBar: {
              timeDivider: true, // 当前时间和持续时间的分隔符
              durationDisplay: true, // 显示持续时间
              remainingTimeDisplay: false, // 是否显示剩余时间功能
              fullscreenToggle: true // 是否显示全屏按钮
            }
        }
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
    submit_create_node () {
      let that = this
      this.$api.post('/filesystem/node/create', this.node.create.form).then(function (response) {
        that.node.create.visit = false
        that.$message({message: '添加成功', type: 'success'})
        that.refresh_nodes()
      }).catch(function (response) {
        that.$message.error({message: response.message, type: 'error'})
      })
    },
    refresh_nodes () {
      let that = this
      this.$api.get("/filesystem/nodes").then(function (response) {
        that.nodes = response.detail
      })
    },
    filter_status(value, row) {
      return row.status === value;
    },
    select_node(node) {
      let that = this
      this.$api.post("/filesystem/files", {id: node.id}).then(function (response) {
        that.nodes = []
        that.node.info = node
        that.node.files = response.detail.files ? response.detail.files : []
        that.node.dirs = response.detail.dirs  
      })
    },
    select_index() {
      this.node.info = {}
      this.node.files = []
      this.node.dirs = []
      this.refresh_nodes()
    },
    remove_file(row) {
      let that = this
      let params = {node_id: this.node.info.id, path: row.path}
      this.$api.post("/filesystem/file/remove", params).then(function (response) {
        that.select_node(that.node.info)
        that.$message({type: "success", "message": "删除成功"})
      }).catch(function (response) {
        that.$message({type: "warning", "message": response.message})
      })
    },
    open_create_dir_dialog () {
      if (Object.keys(this.node.info).length === 0) {
        this.$message({type: "warning", "message": "请进入目录后操作"})
      } else {
        this.node.dir.visit = true
      }
    },
    creat_node_dir() {
      let that = this
      this.$api.post('/filesystem/dir/create', {'node_id': this.node.info.id, 'path': this.node.dir.form.path}).then(function (response) {
        that.$message({message: '添加成功', type: 'success'})
        that.select_node(that.node.info)
        that.node.dir.visit = false
      }).catch(function (response) {
        that.$message({type: "error", "message": response.message})
        that.node.dir.visit = false
      })
    },
    table_row_click (row, column, event) {
      if (row.is_dir) {
        this.node.file.upload_path = row.path
      } else {
        var index= row.path.lastIndexOf(row.name);
        this.node.file.upload_path = row.path.substr(0, index-1)//截取路径字符串
      }
    },
    upload_success (response, file, fileList) {
      this.select_node(this.node.info)
    },
    GetPercent(num, total) {
        num = parseFloat(num);
        total = parseFloat(total);
        if (isNaN(num) || isNaN(total)) {
            return "-";
        }
        return total <= 0 ? "0%" : (Math.round(num / total * 10000) / 100.00)+"%";
    },
    view_file (row) {
      let that = this
      var index= row.path.lastIndexOf(".");
      var ext = row.path.substr(index+1);
      this.node.view.src = "/api/filesystem/download?path="+row.path
      this.$api.post("/filesystem/file/info", {path: row.path}).then(function (response) {
        if (response.detail.type == "img") {
          that.node.view.width = `${response.detail.width + 20}px`
          if (response.detail.width >= 480) {
            that.node.view.width = "500px"
            let height = that.GetPercent(480, response.detail.width)
            that.node.view.style = `width: 480px; height: ${height}`
            that.node.view.type = response.detail.type
            that.node.view.visit = true
          }
        } else if (response.detail.type == "video") {
          that.node.view.width = "50%"
          that.node.view.type = response.detail.type
          that.node.view.visit = true
          that.node.playerOptions.sources[0].src = that.node.view.src
        } else {
          that.$message({type: "error", "message": "不支持该文件预览"})
        }
      })
    },
    close_view: function () {
      this.node.view = {
        width: "500px",
        visit: false,
        style: "",
        src: ""
      }
      this.node.playerOptions.sources[0].src = ""
    },
    download_file(row) {
      window.open("/api/filesystem/download?path="+row.path)
    }
  },
  created: function () {
    this.refresh_nodes()
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
  /* border-bottom: 1px solid whitesmoke; */
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

.font-size {
  font-size: 64px;
}

.icon {
  transform:rotate(90deg);
  -ms-transform:rotate(90deg); /* Internet Explorer */
  -moz-transform:rotate(90deg); /* Firefox */
  -webkit-transform:rotate(90deg); /* Safari 和 Chrome */
  -o-transform:rotate(90deg); /* Opera */
  margin-top: -8px;
}
</style>