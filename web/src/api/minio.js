const base = require("./base")
const httpService = base.httpService 


export const minio = {
    settings: {
        query: function () {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/settings',
                    method: 'post',
                    data: {},
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
        update: function (data = {}) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/settings/update',
                    method: 'post',
                    data: data,
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
    },
    buckets: {
        list: function() {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/buckets',
                    method: 'post',
                    data: {},
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
        create: function(name="") {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/buckets/create',
                    method: 'post',
                    data: {name: name},
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
    },
    objects: {
        list: function (bucket_name, prefix="") {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/objects',
                    method: 'post',
                    data: {bucket_name: bucket_name, prefix: prefix},
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
        info: function(bucket_name, object_name) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/object/info',
                    method: 'post',
                    data: {bucket_name: bucket_name, object_name: object_name},
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
        share: function(bucket_name, object_name) {
            return new Promise((resolve, reject) => {
                httpService({
                    url: '/minio/object/share',
                    method: 'post',
                    data: {bucket_name: bucket_name, object_name: object_name},
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
    }
}


export default {
    minio
}