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
    server: {
        buckets: function () {
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
        objects: function (bucket_name, prefix="") {
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
    }

}


export default {
    minio
}