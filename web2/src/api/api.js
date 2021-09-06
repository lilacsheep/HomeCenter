import axios from "axios"

// let apiUrl = "http://127.0.0.1:8200/api/control"
// let apiUrl = "http://192.168.2.23:8200/api/control"
let apiUrl = "/api"


const httpService = axios.create({
    baseURL: apiUrl, // 需自定义
    // 请求超时时间
    timeout: 30000 // 需自定义
});


httpService.interceptors.request.use(
    config => {
        // 根据条件加入token-安全携带
        if (true) { // 需自定义
            // 让每个请求携带token
            // config.headers['User-Token'] = window.sessionStorage.getItem("token");
        }
        return config;
    }, 
    error => {
        // 请求错误处理
        Promise.reject(error);
    }
)

// respone拦截器
httpService.interceptors.response.use(
    response => {
        // 统一处理状态
        const res = response.data;
        if (res.code !== 200) { // 需自定义
            // 返回异常
            if (res.code === 401 || res.code === 403) {
                window.location.href = "/login";
            } else {
                return Promise.reject({
                    status: res.code,
                    message: res.detail
                });
            }
        } else {
            return response.data;
        }
    },
    // 处理处理
    error => {
         if (error && error.response) {
            switch (error.response.status) {
                case 400:
                    error.message = '错误请求';
                    break;
                case 401:
                    error.message = '未授权，请重新登录';
                    window.location.href = "/login";
                    break;
                case 403:
                    error.message = '未授权，请重新登录';
                    window.location.href = "/login";
                    break;
                case 404:
                    error.message = '请求错误,未找到该资源';
                    break;
                case 405:
                    error.message = '请求方法未允许';
                    break;
                case 408:
                    error.message = '请求超时';
                    break;
                case 500:
                    error.message = '服务器端出错';
                    break;
                case 501:
                    error.message = '网络未实现';
                    break;
                case 502:
                    error.message = '网络错误';
                    break;
                case 503:
                    error.message = '服务不可用';
                    break;
                case 504:
                    error.message = '网络超时';
                    break;
                case 505:
                    error.message = 'http版本不支持该请求';
                    break;
                default:
                    error.message = `未知错误${error.response.status}`;
            }
        } else {
            error.message = "连接到服务器失败";
        }
        return Promise.reject(error);
    }
)

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

export function systemInfo() {
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
}

export function systemProcesses() {
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

export function updateServer(params) {
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
}

export function createInstance(params) {
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
}

export function updateInstance(params) {
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
}

export function  instance_remove(id) {
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
}

export function  refresh_server() {
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
}

export function refresh_instances() {
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

export function start_server() {
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
}

export function stop_server() {
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

export function tasks_list(params = {}) {
    return new Promise((resolve, reject) => {
        httpService({
            url: '/download/tasks',
            method: 'post',
            data: params,
        }).then(response => {
            resolve(response);
        }).catch(error => {
            reject(error);
        });
    });
}

export function refresh_settings(params={}) {
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
}

export function aria2_settings_update(params = {}) {
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
}

export function aria2_global_stats() {
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
}

export function aria2_global_options() {
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

export default {
    get,
    post,
    head,
    fileUpload,
    systemInfo,
    systemProcesses,
    updateServer,
    createInstance,
    updateInstance,
    instance_remove,
    refresh_server,
    refresh_instances,
    start_server,
    stop_server,
    tasks_list,
    refresh_settings,
    aria2_settings_update,
    aria2_global_stats,
    aria2_global_options,
}