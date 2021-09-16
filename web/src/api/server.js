const base = require("./base")
const httpService = base.httpService 

export const server = {
    list: function (group) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/server/list',
                method: 'post',
                data: {group: group, limit: 9999},
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
                url: '/server/remove',
                method: 'post',
                data: {id: id},
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
    update: function(params = {}) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/server/update',
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
    create: function(params={}) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/server/create',
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
}

export const group = {
    list: function (limit=10, page=1) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/server/group/list',
                method: 'post',
                data: {page: page, limit: limit},
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
    create: function (params={}) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/server/group/create',
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
    remove: function (id) {
        return new Promise((resolve, reject) => {
            httpService({
                url: '/server/group/remove',
                method: 'post',
                data: {id: id},
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
}
export default {
    server,
    group,
}