const base = require("./base")
const httpService = base.httpService 
/*网络请求部分*/

/*
 *  get请求
 *  url:请求地址
 *  params:参数
 */
export function get(url, params = {}) {
    return new Promise((resolve, reject) => {
        httpService({
            url: url,
            method: 'get',
            params: params
        }).then(response => {
            resolve(response);
        }).catch(error => {
            reject(error);
        });
    });
}

/*
 *  post请求
 *  url:请求地址
 *  params:参数
 * */
export function post(url, params = {}) {
    return new Promise((resolve, reject) => {
        httpService({
            url: url,
            method: 'post',
            data: params,
            headers: {
                'Content-type': 'application/json;charset=UTF-8'
            },

        }).then(response => {
            resolve(response);
        }).catch(error => {
            reject(error);
        });
    });
}

/*
 *  head请求
 *  url:请求地址
 *  params:参数
 * */
export function head(url, params = {}) {
    return new Promise((resolve, reject) => {
        httpService({
            url: url,
            method: 'head',
            params: params
        }).then(response => {
            resolve(response);
        }).catch(error => {
            reject(error);
        });
    });
}
/*
 *  文件上传
 *  url:请求地址
 *  params:参数
 * */
export function fileUpload(url, params = {}) {
    return new Promise((resolve, reject) => {
        httpService({
            url: url,
            method: 'post',
            data: params,
            headers: { 'Content-Type': 'multipart/form-data' }
        }).then(response => {
            resolve(response);
        }).catch(error => {
            reject(error);
        });
    });
}

export const system = {
    info: function() {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/system/info',
                method: 'get',
                params: {}
            }).then(response => {
                resolve(response);
            }).catch(error => {
                reject(error);
            });
        });
    },
    process: function(pid) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/system/process',
                method: 'post',
                params: {pid: pid}
            }).then(response => {
                resolve(response);
            }).catch(error => {
                reject(error);
            });
        });
    },
    processes: function() {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/system/processes',
                method: 'get',
                params: {}
            }).then(response => {
                resolve(response);
            }).catch(error => {
                reject(error);
            });
        });
    }
}

export const proxy = {
    server: {
        update: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/server/update',
                    method: 'post',
                    data: params,
                    headers: {
                        'Content-type': 'application/json;charset=UTF-8'
                    },
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        info: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/server/info',
                    method: 'get',
                    data: {},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        start: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/server/start',
                    method: 'post',
                    data: {},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        stop: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/server/stop',
                    method: 'post',
                    data: {},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        }
    },
    instance: {
        create: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/instance/create',
                    method: 'post',
                    data: params,
                    headers: {
                        'Content-type': 'application/json;charset=UTF-8'
                    },
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        update: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/instance/update',
                    method: 'post',
                    data: params,
                    headers: {
                        'Content-type': 'application/json;charset=UTF-8'
                    },
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        remove: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/instance/remove',
                    method: 'post',
                    data:  { id: id },
                    headers: {
                        'Content-type': 'application/json;charset=UTF-8'
                    },
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        list: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/proxy/instances',
                    method: 'get',
                    data: {},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        }
    }
}

export const aria2 = {
    task: {
        list: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/tasks',
                    method: 'get',
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        pause: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/task/pause',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        unpause: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/task/unpause',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        add_uri: function(url) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/create',
                    method: 'post',
                    data: {url: url},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        remove: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/remove',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        status: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/task/status',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        }
    },
    global: {
        stats: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/global/stat',
                    method: 'get',
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        options: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/global/options',
                    method: 'get',
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        }
    },
    settings: {
        info: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/settings',
                    method: 'get',
                    data: params,
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        update: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/download/settings/update',
                    method: 'post',
                    data: params,
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
    }
}

export const docker = {
    container: {
        list: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/list',
                    method: 'post',
                    data: params,
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        start: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/start',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        stop: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/stop',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        restart: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/restart',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        pause: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/pause',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        unpause: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/unpause',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        info: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/info',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        stats: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/stats',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        rename: function(id, name) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/stats',
                    method: 'post',
                    data: {id: id, name: name},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        create: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/container/create',
                    method: 'post',
                    data: params,
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        }
    },
    images: {
        list: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/image/list',
                    method: 'post',
                    data: params,
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        pull: function(params={}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/image/pull',
                    method: 'post',
                    data: params,
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        },
        remove: function(id) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/docker/image/remove',
                    method: 'post',
                    data: {id: id},
                }).then(response => {
                    resolve(response);
                }).catch(error => {
                    reject(error);
                });
            });
        }
    }
}


export default {
    get,
    post,
    head,
    fileUpload,
    system,
    proxy,
    aria2,
    docker,
}