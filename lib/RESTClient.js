'use strict';

const errors = require('arsenal').errors;
const assert = require('assert');
const http = require('http');
const querystring = require('querystring');

const bunyanLogstash = require('bunyan-logstash');
const Logger = require('werelogs');
const shuffle = require('arsenal').shuffle;
const conf = require('../config.json');

function errorMap(mdError) {
    const map = {
        NoSuchBucket: 'NoSuchBucket',
        BucketAlreadyExists: 'BucketAlreadyExists',
        NoSuchKey: 'NoSuchKey',
        DBNotFound: 'NoSuchBucket',
        DBAlreadyExists: 'BucketAlreadyExists',
        ObjNotFound: 'NoSuchKey',
        NotImplemented: 'NotImplemented',
    };
    return map[mdError] ? map[mdError] : mdError;
}

class RESTClient {
    /**
     * Constructor for a REST client to bucketd
     *
     * @param {String[]} bootstrap bucketd hosts bootstrap list
     * @param {object} loggingConfig an object with these keys:
     * @param {String} loggingConfig.logLevel
     * @param {String} loggingConfig.dumpLevel
     * @param {String} loggingConfig.logstash.host
     * @param {Number} loggingConfig.logstash.port
     *
     * @return {undefined}
     */
    constructor(bootstrap, loggingConfig) {
        this.bootstrap = bootstrap === undefined ?
            conf.bootstrap : bootstrap;
        assert(this.bootstrap instanceof Array, 'bootstrap must be an Array');
        this.bootstrap = shuffle(this.bootstrap);
        this.setCurrentBootstrap(this.bootstrap[0]);
        this.port = 9000;
        this.httpAgent = new http.Agent({ keepAlive: true });

        this.setupLogging(loggingConfig);
    }

    setupLogging(config) {
        let options = undefined;
        if (config !== undefined) {
            options = {
                level: config.logLevel,
                dump: config.dumpLevel,
                streams: [
                    { stream: process.stdout},
                    {
                        type: 'raw',
                        stream: bunyanLogstash.createStream({
                            host: config.logstash.host,
                            port: config.logstash.port,
                        }),
                    }
                ],
            }
        }
        this.logging = new Logger('BucketClient', options);
    }

    createLogger(reqUids) {
        return reqUids ?
            this.logging.newRequestLoggerFromSerializedUids(reqUids) :
            this.logging.newRequestLogger();
    }

    _shiftCurrentBootstrapToEnd(log) {
        const previousEntry = this.bootstrap.shift();
        this.bootstrap.push(previousEntry);
        const newEntry = this.bootstrap[0];
        this.setCurrentBootstrap(newEntry);

        log.debug('bootstrap head moved',{ from: previousEntry, to: newEntry });
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

    getRaftInformation(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getRaftInformation', { bucketName });
        this._failover(0, 'GET', `/default/informations/${bucketName}`, log,
                        callback);
    }

    getBucketLeader(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getBucketLeader', { bucketName });
        this._failover(0, 'GET', `/default/leader/${bucketName}`, log,
                       callback);
    }

    getBucketAttributes(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getBucketAttributes', { bucketName });
        this._failover(0, 'GET', `/default/attributes/${bucketName}`, log,
                       callback);
    }

    putBucketAttributes(bucketName, reqUids, attributes, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof attributes === 'string', 'attributes must be a string');

        log.debug('putBucketAttributes', { bucketName, attr: attributes });
        this._failover(0, 'POST', `/default/attributes/${bucketName}`, log,
                attributes, callback);
    }

    createBucket(bucketName, reqUids, attributes, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof attributes === 'string', 'attributes must be a string');

        log.debug('createBucket', { bucketName, attr: attributes });
        this._failover(0, 'POST', `/default/bucket/${bucketName}`, log,
                attributes, callback);
    }

    deleteBucket(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('deleteBucket', { bucketName });
        this._failover(0, 'DELETE', `/default/bucket/${bucketName}`, log,
                       callback);
    }

    putObject(bucketName, objName, objVal, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        assert(typeof objVal === 'string', 'objVal must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);
        log.debug('putObject', { bucketName, objName, val: objVal, path });
        this._failover(0, 'POST', path, log,
                       JSON.stringify({ data: objVal }), callback);
    }

    getObject(bucketName, objName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);

        log.debug('getObject', { bucketName, objName, path });
        this._failover(0, 'GET', path, log, callback);
    }

    getBucketAndObject(bucketName, objName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/parallel/${bucketName}/`;
        path += encodeURIComponent(objName);
        log.debug('GetBucketAndObject', {
            bucketName,
            objName,
            path,
        });
        this._failover(0, 'GET', path, log, callback);
    }

    deleteObject(bucketName, objName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);

        log.debug('deleteObject', { bucketName, objName, path });
        this._failover(0, 'DELETE', path, log, callback);
    }

    listObject(bucketName, reqUids, params, cb) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        if (cb !== undefined) {
            assert(typeof params === 'object', 'params must be an object');
        }

        log.debug('listObject', { bucketName, params });
        this._failover(0, 'GET', `/default/bucket/${bucketName}`, log,
                       params, cb);
    }

    endRespond(res, ret, log, callback) {
        const code = res.statusCode;
        if (code <= 201) {
            log.debug('direct request to endpoint returned success',
                { httpCode: code });
            return callback(null, ret);
        }

        const error = errors[errorMap(res.statusMessage)];
        error.isExpected = true;
        log.debug('direct request to endpoint returned an expected error', {
            error });
        return callback(error, ret);
    }

    _failover(tries, method, beginPath, log, params, callback) {
        const argsStr = typeof params === 'function' ?
              '{}' : JSON.stringify(params);
        let count = tries;

        log.debug('sending request', { httpMethod: method, beginPath,
            params: argsStr, tries: count });

        this.request(method, beginPath, log, params, (err, data) => {
            if (err && !err.isExpected) {
                if (++count >= this.bootstrap.length) {
                    log.error('failover failed, giving up', { tries: count });
                    return callback(err);
                }
                return this._shiftCurrentBootstrapToEnd(log)
                    ._failover(count, method, beginPath, log,
                               params, callback);
            }
            log.debug('request received successful response');
            return callback(err, data);
        });
    }

    request(method, beginPath, log, params, callback) {
        let next;
        let path = beginPath;
        if (typeof params === 'function') {
            next = params;
        } else {
            next = callback;
        }
        assert(typeof next === 'function', 'callback must be a function');
        const headers = {
            'content-length': 0,
            'x-scal-request-uids': log.getSerializedUids(),
        };

        if (method === 'GET' && typeof params === 'object'
                && Object.keys(params).length !== 0) {
            path += `?${querystring.stringify(params)}`;
        }

        const option = {
            method,
            path,
            headers,
            host: this.getCurrentBootstrap(),
            port: this.getPort(),
            agent: this.httpAgent,
        };
        const ret = [];
        let retLen = 0;

        // somehow option can get cyclical so fields need to be chosen
        // for display
        log.debug('direct request to MD endpoint', { httpMethod: method,
            host: option.host, port: option.port });
        const req = http.request(option);

        if (method === 'POST') {
            const binData = new Buffer(params, 'binary');
            req.setHeader('content-type', 'application/octet-stream');
            req.setHeader('content-length', binData.length);
            req.write(binData, 'binary');
        }

        req.on('response', (res) => {
            res.on('data', (data) => {
                ret.push(data);
                retLen += data.length;
            }).on('error', next).on('end', () => {
                this.endRespond(res, Buffer.concat(ret, retLen).toString(), log, next);
            });
        }).on('error', next).end();
    }
}

module.exports = RESTClient;
