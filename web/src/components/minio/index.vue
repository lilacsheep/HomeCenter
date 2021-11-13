<template>
  <a-layout-content :style="{ padding: '12px' }">
    <a-breadcrumb separator=">" style="margin: 12px 8px">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item href=""> 对象存储 </a-breadcrumb-item>
    </a-breadcrumb>
    <!-- // style="background: #fbfbfb;border: 1px solid #f4f4f4;height: 100%" -->
    <a-tabs default-active-key="1" tab-position="left">
      <a-tab-pane key="1" tab="文件管理">
        <a-row :gutter="20">
          <a-col :span="6">
            <a-button type="primary" @click="bucket.create.visible=true" block>新增Bucket</a-button>
            <a-list id="bucketlist" bordered :data-source="buckets"  style="background: #FFF">
              <a-list-item slot="renderItem" slot-scope="item">
                <a-button type="link" style="padding:0" @click="select_bucket(item.name)">{{ item.name }}</a-button>
              <a-icon type="setting" slot="actions"/>
              </a-list-item>
            </a-list>
          </a-col>
          <a-col :span="18" style="background: #fff;padding: 0;">
            <a-row :gutter="20">
              <a-col :span="20">
                <a-breadcrumb id="breadcrumb" separator="/" style="margin-bottom: 5px">
                  <a-breadcrumb-item >
                    <a style="font-size: 16px;padding: 2px;margin-right:0" @click="select_breadcrumb(-1)">{{now.bucket}}</a>
                  </a-breadcrumb-item>
                  <a-breadcrumb-item :key="p" v-for="(p, i) in now.path">
                    <a style="font-size: 16px;padding: 2px;margin-right:0" @click="select_breadcrumb(i)">{{ p }}</a>
                  </a-breadcrumb-item>
                </a-breadcrumb>
              </a-col>
              <a-col :span="4" style="float: right" >
                <a-upload name="file" :multiple="true" :action="upload_action" :data="{name: now.bucket}" @change="upload_change">
                  <a-button size="small"> <a-icon type="upload" /> 上 传 </a-button>
                </a-upload>
              </a-col>
            </a-row>
            
            <a-table :columns="columns" :data-source="objects" size="small" style="background: white" :scroll="{ x: 1300 }">
              <template slot="name" slot-scope="text, row">
                <a-button icon="file" size="small" type="link" v-if="row.etag !== ''">{{ getFileName(row.name) }}</a-button>
                <a-button icon="folder" @click="select_dir(row)" size="small" v-else type="link">
                  {{ get_dirname(row.name) }}</a-button>
              </template>
              <span slot="size" slot-scope="text">{{ text | diskSize }}</span>
              <span slot="lastModified" slot-scope="text">{{text | dateformat}}</span>
              <template slot="action" slot-scope="text, row">
                <a-button-group size="small">
                  <a-button @click="download_file(row)"><a-icon type="download" /></a-button>
                  <a-button @click="share_file(row)"><a-icon type="share-alt" /></a-button>
                  <a-button @click="object_info_open(row)"><a-icon type="setting" /></a-button>
                </a-button-group>
              </template>
              
            </a-table>
          </a-col>
        </a-row>
      </a-tab-pane>
      <a-tab-pane key="2" tab="本地配置">
        <a-row :gutter="20">
          <a-col :span="24">
            <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol" style="background: #ffffff; padding: 10px">
              <a-form-model-item label="启动">
                <a-switch
                  checked-children="开"
                  un-checked-children="关"
                  v-model="form.auto_start"
                />
              </a-form-model-item>
              <a-form-model-item label="认证密钥">
                <a-input v-model="form.access_key" />
              </a-form-model-item>
              <a-form-model-item label="认证密码">
                <a-input v-model="form.secret_key" />
              </a-form-model-item>
              <a-form-model-item label="存储目录">
                <a-input v-model="form.save_path" />
              </a-form-model-item>
              <a-form-model-item label="配置目录">
                <a-input v-model="form.config_dir" />
              </a-form-model-item>
              <a-form-model-item label="对外端口">
                <a-input-number :max="65534" :min="80" v-model="form.port" />
              </a-form-model-item>
              <a-form-model-item label="控制台">
                <a-switch
                  checked-children="开"
                  un-checked-children="关"
                  v-model="form.webui"
                />
              </a-form-model-item>
              <a-form-model-item label="控制台端口">
                <a-input-number
                  :max="65534"
                  :min="80"
                  v-model="form.webui_port"
                />
              </a-form-model-item>
              <a-form-model-item label="域名">
                <a-input v-model="form.minio_domain" />
              </a-form-model-item>
              <a-form-model-item label="服务地址">
                <a-input v-model="form.server_url" />
              </a-form-model-item>
              <a-form-model-item :wrapper-col="{ span: 14, offset: 3 }">
                <a-button
                  type="primary"
                  :loading="update_loading"
                  @click="update_setting"
                >
                  应用配置
                </a-button>
              </a-form-model-item>
            </a-form-model>
          </a-col>
        </a-row>
      </a-tab-pane>
    </a-tabs>
    <a-drawer title="创建" placement="right" :closable="false">
      <a-card title="Card Title">
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 100%; text-align: center" :hoverable="false">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 100%; text-align: center" :hoverable="false">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
      </a-card>
    </a-drawer>
    <a-modal v-model="bucket.create.visible" title="创建Bucket" @ok="create_bucket">
      <a-input placeholder="Bucket Name" v-model="bucket.create.name"/>
    </a-modal>
    <a-drawer title="创建" placement="right" :visible="object.info.visible" width="60%" :closable="false" @close="object_info_close">
      <a-card title="Card Title">
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 50%; text-align: center">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 100%; text-align: center" :hoverable="false">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
        <a-card-grid style="width: 100%; text-align: center" :hoverable="false">
          <a-statistic
            title="Feedback"
            :value="11.28"
            :precision="2"
            suffix="%"
            :value-style="{ color: '#3f8600' }"
            style="margin-right: 50px"
          >
            <template #prefix>
              <a-icon type="arrow-up" />
            </template>
          </a-statistic>
        </a-card-grid>
      </a-card>
    </a-drawer>
    <a-modal v-model="object.share.visible" title="分享连接" :footer="null">
      <a-result
          status="success"
          title="已经成功复制到剪切板中"
          :sub-title="object.share.url"
        >
        </a-result>
    </a-modal>
  </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      expandedKeys: [],
      objects: [],
      minios: [],
      labelCol: { span: 3 },
      wrapperCol: { span: 16 },
      update_loading: false,
      buckets: [],
      now: {
        bucket: "",
        path: [],
      },
      object: {
        info: {
          visible: false,
          data: {}
        },
        share: {
          visible: false,
          url: ""
        }
      },
      bucket: {
        create: {
          visible: false,
          name: ""
        }
      },
      upload_action: "",
      form: {
        auto_start: false,
        access_key: "",
        secret_key: undefined,
        config_dir: undefined,
        delivery: false,
        port: 0,
        webui: false,
        webui_port: 0,
        minio_domain: "",
        server_url: "",
      },
      columns: [
        {
          title: "文件名",
          dataIndex: "name",
          key: "name",
          scopedSlots: { customRender: "name" },
        },
        {
          title: "Size",
          dataIndex: "size",
          key: "size",
          fixed: 'right',
          scopedSlots: { customRender: "size" },
        },
        {
          title: "修改日期",
          dataIndex: "lastModified",
          key: "lastModified",
          fixed: 'right',
          scopedSlots: { customRender: "lastModified" },
        },
        {
          title: '操作',
          key: 'operation',
          fixed: 'right',
          width: 120,
          scopedSlots: { customRender: 'action' },
        },
      ],
    };
  },
  created: function () {
    let host = this.$apihost == ""? window.location.host: this.$apihost
    this.upload_action = `${window.location.protocol}//${host}/api/minio/object/upload`
    this.reload_setting();
    this.reload_buckets(true);
  },
  methods: {
    reload_buckets: function (reload_object=false) {
      this.$minio.buckets.list().then((response) => {
        this.buckets = response.detail ? response.detail : [];
        if (reload_object) {
          this.now.bucket = this.buckets[0].name
          this.load_object()
        }
      }).catch((response) => {
        this.$message.error(response.message);
      });
    },
    upload_change: function (event) {
      if (event.file.status === "done") {
        this.load_object()
      }
    },
    create_bucket() {
      this.$minio.buckets.create(this.bucket.create.name).then((response) => {
        this.$message.info("创建成功");
        this.bucket.create.visible = false
        this.reload_buckets(false)
      }).catch((response) => {
        this.$message.error(response.message);
      });
    },
    select_dir(row) {
      let dirs = row.name.substring(0, row.name.length-1).split("/")
      this.now.path.push(dirs.pop());
      this.load_object();
    },
    select_file(row) {},
    object_info_open(row) {
      this.$minio.objects.info(this.now.bucket, row.name).then((response) => {
        this.object.data = response.detail
        this.object.info.visible = true
      }).catch((response) => {
        this.$message.error("获取信息失败: "+response.message)
      })
    },
    object_info_close(e) {
      this.object.info.visible = false
    },
    download_file(row) {
      let host = this.$apihost == ""? window.location.host: this.$apihost
      window.open(`${window.location.protocol}//${host}/api/minio/object/download?bucket_name=${this.now.bucket}&object_name=${row.name}`, "_blank")
    },
    share_file(row) {
      this.$minio.objects.share(this.now.bucket, row.name).then((response) => {
        this.object.share.visible = true
        this.object.share.url = response.detail
        console.log(response.detail)
        this.$copyText(this.object.share.url)
      }).catch((response) => {
        this.$message.error("获取信息失败: "+response.message)
      })
    },
    select_bucket(bucket) {
      this.now.path = []
      this.now.bucket = bucket
      this.load_object()
    },
    select_breadcrumb(index) {
      this.now.path = this.now.path.slice(0, index + 1);
      this.load_object();
    },
    getFileName: function (path) {
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

      return fileName;
    },
    get_dirname(path) {
      let dirs = path.substring(0, path.length-1).split("/")
      return dirs.pop()
    },
    load_object: function () {
      this.$minio.objects.list(this.now.bucket, this.now.path.join('/')+'/').then((response) => {
        this.objects = response.detail;
      }).catch((response) => {
        this.$message.error(response.message);
      });
    },
    update_setting: function () {
      this.update_loading = true;
      this.$minio.settings.update(this.form).then((response) => {
        this.$message.info("更新成功");
        this.update_loading = false;
      }).catch((response) => {
        this.$message.error(response.message);
        this.update_loading = false;
      });
    },
    reload_setting: function () {
      this.$minio.settings.query().then((response) => {
        this.form = response.detail;
      }).catch((response) => {
        this.$message.error(response.message);
      });
    },
  },
};
</script>

<style>
#breadcrumb .ant-breadcrumb-separator {
  margin: 0 5px;
  font-size: 16px;
}
#bucketlist .ant-list-item {
  padding-right: 5px;
  padding-left: 5px;
  padding: 5px;
}
</style>