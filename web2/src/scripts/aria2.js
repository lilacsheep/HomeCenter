import api from "../api/api"
import { message } from 'ant-design-vue';


const aria2_uri = {
    "addUri": "/download/create",
    "addTorrent": "/download/torrent",
    "tasks": "/download/tasks",
    "remove": "/download/remove",
    "pause": "/download/task/pause",
    "unpause": "/download/task/unpause",
    "globalStat": "/download/global/stat",
    "taskStatus": "/download/task/status",
    "globalOptions": "/download/global/options"
}

export function pause(gid, callback = function (response) { }) {
    api.post(aria2_uri['pause'], {id: gid}).then(function (response) {
        message.success('已经暂停')
        callback(response)
    }).catch(function (response) {
        message.error('暂停失败:' + response.message)
    })
}

export function unpause(gid, callback = function (response) { }) {
    api.post(aria2_uri['unpause'], { id: gid }).then(function (response) {
        message.success('启动成功')
        callback(response)       
    }).catch(function (response) {
        message.error('启动失败:' + response.message)
    })
}

export function globalStat(callback = function (response) { }) {
    api.get(aria2_uri['globalStat']).then(function (response) {
        callback(response)
    })
}

export function addUri(url, callback = function (response) { }) {
    api.post(aria2_uri['addUri'], {url: url}).then(function (response) {
        message.success('创建成功')
        callback(response)
    }).catch(function (response) {
        Message({ message: '创建失败:' + response.message, type: 'error' })
    })
}

export function addTorrent(params={}, callback=function (response){}) {
    api.post(aria2_uri['addTorrent'], { url: url }).then(function (response) {
        message.success('创建成功')
        callback(response)       
    }).catch(function (response) {
        Message({ message: '创建失败:' + response.message, type: 'error' })
    })
}

export function tasks(params = {}, callback = function (response) {}) {
    api.get(aria2_uri["tasks"]).then(function (response) {    
        callback(response)
    }).catch(function (response) {
        console.log("获取任务列表失败: " + response.message)
    })
}

export function removeTask(gid, callback = function (response) {}) {
    api.post(aria2_uri["remove"], {id: gid}).then(function (response) {
        message.success('删除成功')
        callback(response)
    }).catch(function (response) {
        message.error('删除失败:' + response.message)
    })
}

export function taskStatus(gid, callback = function (response) {}) {
    api.post(aria2_uri["taskStatus"], { id: gid }).then(function (response) {
        callback(response)
    }).catch(function (response) {
        message.error('获取任务状态失败:' + response.message)
    })
}

export function globalOptions(callback = function (response) { }) {
    api.get(aria2_uri["globalOptions"]).then(function (response) {
        callback(response)
    }).catch(function (response) {
        message.error('获取全局配置失败:' + response.message)
    })
}

var aria2Api = {
    addUri: addUri,
    globalStat: globalStat,
    pause: pause,
    unpause: unpause,
    addTorrent: addTorrent,
    tasks: tasks,
    removeTask: removeTask,
    taskStatus: taskStatus,
    globalOptions: globalOptions,
}

export {aria2Api}