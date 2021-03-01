import api from "../api/api"
import { Message } from 'element-ui';


const aria2_uri = {
    "addUri": "/download/create",
    "addTorrent": "/download/torrent",
    "tasks": "/download/tasks",
    "remove": "/download/remove",
    "pause": "/download/task/pause",
    "unpause": "/download/task/unpause",
    "globalStat": "/download/global/stat",
    "taskStatus": "/download/task/status"
}

export function pause(gid, callback = function (response) { }) {
    api.post(aria2_uri['pause'], {id: gid}).then(function (response) {
        Message({ message: '已经暂停', type: 'success' })
        callback(response)
    }).catch(function (response) {
        Message({ message: '暂停失败:' + response.detail, type: 'error' })
    })
}

export function unpause(gid, callback = function (response) { }) {
    api.post(aria2_uri['unpause'], { id: gid }).then(function (response) {
        Message({ message: '启动成功', type: 'success' })
        callback(response)       
    }).catch(function (response) {
        Message({ message: '启动失败:' + response.detail, type: 'error' })
    })
}

export function globalStat(callback = function (response) { }) {
    api.get(aria2_uri['globalStat']).then(function (response) {
        callback(response)
    })
}

export function addUri(url, callback = function (response) { }) {
    api.post(aria2_uri['addUri'], {url: url}).then(function (response) {
        Message({ message: '创建成功', type: 'success' })
        callback(response)
    }).catch(function (response) {
        Message({ message: '创建失败:' + response.detail, type: 'error' })
    })
}

export function addTorrent(params={}, callback=function (response){}) {
    api.post(aria2_uri['addTorrent'], { url: url }).then(function (response) {
        Message({ message: '创建成功', type: 'success' })
        callback(response)       
    }).catch(function (response) {
        Message({ message: '创建失败:' + response.detail, type: 'error' })
    })
}

export function tasks(params = {}, callback = function (response) {}) {
    api.get(aria2_uri["tasks"], params).then(function (response) {    
        callback(response)
    }).catch(function (response) {
        console.log("获取任务列表失败: "+ response.detail)
    })
}

export function removeTask(gid, callback = function (response) {}) {
    api.get(aria2_uri["tasks"], {id: gid}).then(function (response) {
        Message({ message: '删除成功', type: 'success' })
        callback(response)
    }).catch(function (response) {
        Message({ message: '删除失败:' + response.detail, type: 'error' })
    })
}

export function taskStatus(gid, callback = function (response) {}) {
    api.post(aria2_uri["taskStatus"], { id: gid }).then(function (response) {
        callback(response)
    }).catch(function (response) {
        Message({ message: '获取任务状态失败:' + response.detail, type: 'error' })
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
    taskStatus: taskStatus
}

export {aria2Api}