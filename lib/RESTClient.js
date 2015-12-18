'use strict';

const assert = require('assert');
const http = require('http');
const querystring = require('querystring');

const shuffle = require('arsenal').shuffle;
const conf = require('../config.json');

class RESTClient {
    constructor(bootstrap) {
        this.bootstrap = bootstrap === undefined ?
            conf.bootstrap : bootstrap;
        assert(this.bootstrap instanceof Array, 'bootstrap must be an Array');
        this.bootstrap = shuffle(this.bootstrap);
        this.setCurrentBootstrap(this.bootstrap[0]);
        this.port = 9000;
    }

    _shiftCurrentBootstrapToEnd() {
        this.bootstrap.push(this.bootstrap.shift());
        this.setCurrentBootstrap(this.bootstrap[0]);
        return this;
    }

    setCurrentBootstrap(host) {
        this.current = host;
        return this;
    }

    getCurrentBootstrap() {
        return this.current;
    }

    getPort() {
        return this.port;
    }

    getBucketLeader(bucketName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this._failover(0, 'GET', `/default/leader/${bucketName}`, callback);
    }

    getBucketAttributes(bucketName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this._failover(0, 'GET', `/default/attributes/${bucketName}`, callback);
    }

    putBucketAttributes(bucketName, attributes, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this._failover(0, 'POST', `/default/attributes/${bucketName}`,
                attributes, callback);
    }

    createBucket(bucketName, attributes, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this._failover(0, 'POST', `/default/bucket/${bucketName}`,
                attributes, callback);
    }

    deleteBucket(bucketName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        this._failover(0, 'DELETE', `/default/bucket/${bucketName}`, callback);
    }

    putObject(bucketName, objName, objVal, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        assert(typeof objVal === 'string', 'objVal must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += objName.replace(/\//g, '%2F');
        this._failover(0, 'POST', path, { data: objVal }, callback);
    }

    getObject(bucketName, objName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += objName.replace(/\//g, '%2F');
        this._failover(0, 'GET', path, callback);
    }

    deleteObject(bucketName, objName, callback) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += objName.replace(/\//g, '%2F');
        this._failover(0, 'DELETE', path, callback);
    }

    listObject(bucketName, params, cb) {
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        if (cb !== undefined) {
            assert(typeof params === 'object', 'params must be an object');
        }
        this._failover(0, 'GET', `/default/bucket/${bucketName}`, params, cb);
    }

    endRespond(res, ret, callback) {
        const code = res.statusCode;
        if (code <= 201) {
            const value = ret ? JSON.parse(ret) : null;
            return callback(null, value);
        }
        const error = new Error(res.statusCode);
        error.isExpected = true;
        error.message = res.statusMessage;
        return callback(error, ret);
    }

    _failover(tries, method, beginPath, params, callback) {
        let count = tries;
        this.request(method, beginPath, params, (err, data) => {
            if (err && !err.isExpected) {
                if (++count >= this.bootstrap.length) {
                    return callback(err);
                }
                this._shiftCurrentBootstrapToEnd()
                    ._failover(count, method, beginPath, params, callback);
            }
            return callback(err, data);
        });
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
            host: this.getCurrentBootstrap(),
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
