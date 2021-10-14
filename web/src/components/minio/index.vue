<template>
    <a-layout-content :style="{ padding: '12px'}">
      <a-breadcrumb separator=">" style="margin: 12px 8px">
        <a-breadcrumb-item>首页</a-breadcrumb-item>
        <a-breadcrumb-item href="">
          对象存储
        </a-breadcrumb-item>
      </a-breadcrumb>
      <!-- // style="background: #fbfbfb;border: 1px solid #f4f4f4;height: 100%" -->
      <a-row :gutter="16" >
        <a-col :span="6">
          <a-button type="primary" @click="form.group.add.visible=true" block>新增组</a-button>
          <a-tree style="background: #FFFFFF;height: 100%;margin: 0;" :tree-data="servergroup" @select="server_select" :expandedKeys.sync="expandedKeys" show-icon>
            <template #title="{ key: treeKey, title }" >
              <a-dropdown :trigger="['contextmenu']">
                <span>{{ title }}</span>
                <template #overlay>
                  <a-menu @click="({ key: menuKey }) => onContextMenuClick(treeKey, menuKey)">
                    <!-- <a-menu-item key="add_host" v-if="check_key(treeKey)"><a-icon type="edit" />新增主机</a-menu-item>
                    <a-menu-item key="edit_host" v-if="!check_key(treeKey)"><a-icon type="edit" />编辑</a-menu-item>
                    <a-menu-item key="delete_host" v-if="!check_key(treeKey)"><a-icon type="delete" />删除</a-menu-item>
                    <a-menu-item key="delete_node" v-if="check_key(treeKey)"><a-icon type="delete" />删除</a-menu-item>
                    <a-menu-item key="edit_node" v-if="check_key(treeKey)"><a-icon type="edit" />编辑</a-menu-item> -->
                  </a-menu>
                </template>
              </a-dropdown>
            </template>
            <a-icon slot="desktop" type="desktop" />
            <a-icon slot="folder" type="folder" />
          </a-tree>
        </a-col>
        <a-col :span="18">
          <a-table :columns="columns" :data-source="objects" size="small">
            <a slot="name" slot-scope="text">{{ text }}</a>
          </a-table>
        </a-col>
      </a-row>
    </a-layout-content>
</template>

<script>
export default {
  data() {
    return {
      expandedKeys: [],
      objects: [],
      columns: [
        {
          title: 'Name',
          dataIndex: 'name',
          key: 'name',
          scopedSlots: { customRender: 'name' },
        },
        {
          title: 'Age',
          dataIndex: 'age',
          key: 'age',
          width: 80,
        },
        {
          title: 'Address',
          dataIndex: 'address',
          key: 'address 1',
          ellipsis: true,
        },
        {
          title: 'Long Column Long Column Long Column',
          dataIndex: 'address',
          key: 'address 2',
          ellipsis: true,
        },
        {
          title: 'Long Column Long Column',
          dataIndex: 'address',
          key: 'address 3',
          ellipsis: true,
        },
        {
          title: 'Long Column',
          dataIndex: 'address',
          key: 'address 4',
          ellipsis: true,
        },
      ]
    }
  },
  computed: {},
  created: function () {},
  methods: {},
};
</script>

<style>

</style>