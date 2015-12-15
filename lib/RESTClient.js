'use strict';

const assert = require('assert');
const http = require('http');
const querystring = require('querystring');

class RESTClient {
    constructor(ip, port) {
        if (ip) {
            assert(typeof this.ip === 'string', 'ip must be a string');
        }
        if (port) {
            assert(typeof this.port === 'number', 'port must be a number');
        }
        this.ip = ip || '127.0.0.1';
        this.port = port || 9000;
    }

    getIp() {
        return this.ip;
    }

    getPort() {
        return this.port;
    }

    getBucketLeader(bucketName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this.request('GET', `/default/leader/${bucketName}`, callback);
    }

    getBucketAttributes(bucketName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this.request('GET', `/default/attributes/${bucketName}`, callback);
    }

    putBucketAttributes(bucketName, attributes, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this.request('POST', `/default/attributes/${bucketName}`,
                attributes, callback);
    }

    createBucket(bucketName, attributes, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this.request('POST', `/default/bucket/${bucketName}`,
                attributes, callback);
    }

    deleteBucket(bucketName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this.request('DELETE', `/default/bucket/${bucketName}`, callback);
    }

    putObject(bucketName, objName, objVal, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        assert(typeof objVal === 'string', 'objVal must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += objName.replace(/\//g, '%2F');
        this.request('POST', path, { data: objVal }, callback);
    }

    getObject(bucketName, objName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += objName.replace(/\//g, '%2F');
        this.request('GET', path, callback);
    }

    deleteObject(bucketName, objName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += objName.replace(/\//g, '%2F');
        this.request('DELETE', path, callback);
    }

    listObject(bucketName, params, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        if (callback !== undefined) {
            assert(typeof params === 'object', 'params must be an object');
        }
        this.request('GET', `/default/bucket/${bucketName}`, params, callback);
    }

    endRespond(res, ret, callback) {
        const code = res.statusCode;
        if (code <= 201) {
            callback(null, JSON.parse(ret));
        } else {
            try {
                callback(JSON.parse(ret));
            } catch (e) {
                callback(e);
            }
        }
    }

    request(method, beginPath, params, callback) {
        let next;
        let path = beginPath;
        if (typeof params === 'function') {
            next = params;
        } else {
            next = callback;
        }
        assert(typeof next === 'function', 'callback must be a function');
        if (method === 'GET' && typeof params === 'object'
                && Object.keys(params).length !== 0) {
            path += `?${querystring.stringify(params)}`;
        }
        const option = {
            method,
            path,
            host: this.getIp(),
            port: this.getPort(),
        };
        let ret = '';
        const req = http.request(option);

        if (method === 'POST' && typeof params === 'object') {
            req.write(JSON.stringify(params));
        }

        req.on('response', (res) => {
            res.on('data', (data) => {
                ret += data.toString();
            }).on('error', next).on('end', () => {
                this.endRespond(res, ret, next);
            });
        }).on('error', next).end();
    }
}

module.exports = RESTClient;
