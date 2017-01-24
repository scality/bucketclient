'use strict'; // eslint-disable-line strict

const errors = require('arsenal').errors;
const assert = require('assert');
const http = require('http');
const https = require('https');
const querystring = require('querystring');

const Logger = require('werelogs').Logger;
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
     * @param {String[]} bootstrap - bucketd hosts bootstrap list
     * @param {object} loggingConfig - an object with these keys:
     * @param {String} loggingConfig.logLevel - werelogs log level
     * @param {String} loggingConfig.dumpLevel - werelogs dump level
     * @param {Boolean} useHttps - use https
     * @param {string} key - Private certificate key content
     * @param {string} cert - Public certificate content
     * @param {string} ca - Certificate authority content
     *
     * @return {undefined}
     */
    constructor(bootstrap, loggingConfig, useHttps, key, cert, ca) {
        this.bootstrap = bootstrap === undefined ?
            conf.bootstrap : bootstrap;
        assert(this.bootstrap instanceof Array, 'bootstrap must be an Array');
        this.bootstrap = shuffle(this.bootstrap);
        this.setCurrentBootstrap(this.bootstrap[0]);
        this.port = 9000;
        this._useHttps = useHttps;
        this.transport = useHttps ? https : http;
        this._cert = cert;
        this._key = key;
        if (this._useHttps) {
            this.agent = new https.Agent({
                keepAlive: true,
                ca: ca ? [ca] : undefined,
                requestCert: true,
            });
        } else {
            this.agent = new http.Agent({ keepAlive: true });
        }
        this.setupLogging(loggingConfig);
    }

    setupLogging(config) {
        let options = undefined;
        if (config !== undefined) {
            options = {
                level: config.logLevel,
                dump: config.dumpLevel,
            };
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

        log.debug('bootstrap head moved',
            { from: previousEntry, to: newEntry });
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
            null, null, callback);
    }

    getBucketLeader(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getBucketLeader', { bucketName });
        this._failover(0, 'GET', `/default/leader/${bucketName}`, log,
            null, null, callback);
    }

    // <BucketAndObjectOperations>

    getBucketAttributes(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getBucketAttributes', { bucketName });
        this._failover(0, 'GET', `/default/attributes/${bucketName}`, log,
            null, null, callback);
    }

    putBucketAttributes(bucketName, reqUids, attributes, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof attributes === 'string', 'attributes must be a string');

        log.debug('putBucketAttributes', { bucketName, attr: attributes });
        this._failover(0, 'POST', `/default/attributes/${bucketName}`, log,
            null, attributes, callback);
    }

    createBucket(bucketName, reqUids, attributes, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof attributes === 'string', 'attributes must be a string');

        log.debug('createBucket', { bucketName, attr: attributes });
        this._failover(0, 'POST', `/default/bucket/${bucketName}`, log,
            null, attributes, callback);
    }

    deleteBucket(bucketName, reqUids, callback) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('deleteBucket', { bucketName });
        this._failover(0, 'DELETE', `/default/bucket/${bucketName}`, log,
            null, null, callback);
    }

    /**
     * Create or udpate an object or a version of an object.
     * Examples:
     * - creating an object: PUT /foo/bar
     * - updating a version: PUT /foo/bar?versionId=1234567890
     *
     * @param {string} bucketName - bucket name
     * @param {string} objName - the name of the object
     * @param {string} objVal - the value of the object
     * @param {string} reqUids - the identifier of the request
     * @param {function} callback - callback
     * @param {object} params - extra parameters for the case of versioning
     *
     * @return {undefined}
     */
    putObject(bucketName, objName, objVal, reqUids, callback, params) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        assert(typeof objVal === 'string', 'objVal must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);
        log.debug('putObject', { bucketName, objName, val: objVal, path });
        this._failover(0, 'POST', path, log, params, objVal, callback);
    }

    /**
     * Get an object or a version of an object.
     * Examples:
     * - getting an object: GET /foo/bar
     * - getting a version: GET /foo/bar?versionId=1234567890
     *
     * @param {string} bucketName - bucket name
     * @param {string} objName - the name of the object
     * @param {string} reqUids - the identifier of the request
     * @param {function} callback - callback
     * @param {object} params - extra parameters for the case of versioning
     *
     * @return {undefined}
     */
    getObject(bucketName, objName, reqUids, callback, params) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);

        log.debug('getObject', { bucketName, objName, path });
        this._failover(0, 'GET', path, log, params, null, callback);
    }

    /**
     * Get the attributes of a bucket and an object or a version of an object.
     *
     * @param {string} bucketName - bucket name
     * @param {string} objName - the name of the object
     * @param {string} reqUids - the identifier of the request
     * @param {function} callback - callback
     * @param {object} params - extra parameters for the case of versioning
     *
     * @return {undefined}
     */
    getBucketAndObject(bucketName, objName, reqUids, callback, params) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/parallel/${bucketName}/`;
        path += encodeURIComponent(objName);
        log.debug('GetBucketAndObject', {
            bucketName,
            objName,
            path,
            params,
        });
        this._failover(0, 'GET', path, log, params, null, callback);
    }

    /**
     * Delete an object or a version of an object.
     * Examples:
     * - deleting an object: DELETE /foo/bar
     * - deleting a version: DELETE /foo/bar?versionId=1234567890
     *
     * @param {string} bucketName - bucket name
     * @param {string} objName - the name of the object
     * @param {string} reqUids - the identifier of the request
     * @param {function} callback - callback
     * @param {object} params - extra parameters for the case of versioning
     *
     * @return {undefined}
     */
    deleteObject(bucketName, objName, reqUids, callback, params) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);

        log.debug('deleteObject', { bucketName, objName, path, params });
        this._failover(0, 'DELETE', path, log, params, null, callback);
    }

    /**
     * List objects or versions of objects of a bucket.
     *
     * @param {string} bucketName - bucket name
     * @param {string} reqUids - the identifier of the request
     * @param {object} params - parameters for listing, now includes listing
     *                           all versions of all object of a bucket
     *                           Example: GET /foo?versions&delimiter=&prefix=
     * @param {function} cb - callback
     *
     * @return {undefined}
     */
    listObject(bucketName, reqUids, params, cb) {
        const log = this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        if (cb !== undefined) {
            assert(typeof params === 'object', 'params must be an object');
        }

        log.debug('listObject', { bucketName, params });
        this._failover(0, 'GET', `/default/bucket/${bucketName}`, log, params,
            null, cb);
    }

    // </BucketAndObjectOperations>

    healthcheck(log, callback) {
        this.request('GET', '/_/healthcheck', log, null, null,
            callback);
    }

    endRespond(res, ret, log, callback) {
        const code = res.statusCode;
        if (code <= 201) {
            log.debug('direct request to endpoint returned success',
                { httpCode: code });
            return callback(null, ret);
        }

        let error = errors[errorMap(res.statusMessage)];
        if (error) {
            error.isExpected = true;
            log.debug('direct request to endpoint returned an expected error', {
                error });
            return callback(error, ret);
        }
        error = errors.InternalError.customizeDescription('unexpected error');
        error.isExpected = false;
        log.debug('direct request to endpoint returned an unexpected error',
            { error });
        return callback(error, ret);
    }

    /**
     * Send a request with its associated parameters and data.
     *
     * @param {number} tries - times trying to send the request
     * @param {string} method - the HTTP method of the request
     * @param {string} beginPath - formated path without parameters
     * @param {object} log - logger
     * @param {object} params - parameters of the request
     * @param {string} data - data of the request
     * @param {function} callback - callback
     *
     * @return {object} - return value of the callback
     */
    _failover(tries, method, beginPath, log, params, data, callback) {
        log.debug('sending request', { httpMethod: method, beginPath, tries });

        this.request(method, beginPath, log, params, data, (err, _data) => {
            if (err && !err.isExpected) {
                const count = tries + 1;
                if (count >= this.bootstrap.length) {
                    log.error('failover failed, giving up', { tries: count });
                    return callback(err);
                }
                return this._shiftCurrentBootstrapToEnd(log)
                    ._failover(count, method, beginPath, log,
                               params, data, callback);
            }
            log.debug('request received successful response');
            return callback(err, _data);
        });
    }

    /**
     * For all requests, there might always be both extra parameters and extra
     * data, such as PUT /object?acl&versionId=1234567890 for changing the acl
     * of a specific version and POST /foo?delete for batch deleting a list of
     * objects and/or versions which are attached in the data of the request.
     *
     * By forwarding these data items, we can make bucketclient general enough
     * for upcoming features without putting the burden of sub-categorizing
     * requests on S3 and without having to make new specific bucketclient
     * functions for new specific features.
     *
     * @param {string} method - the HTTP method of the request
     * @param {string} beginPath - formated path without parameters
     * @param {object} log - logger
     * @param {object} params - parameters of the request
     * @param {string} data - data of the request
     * @param {function} callback - callback
     *
     * @return {undefined}
     */
    request(method, beginPath, log, params, data, callback) {
        let path = beginPath;
        assert(typeof callback === 'function', 'callback must be a function');
        const headers = {
            'content-length': 0,
            'x-scal-request-uids': log.getSerializedUids(),
        };

        if (params) {
            path += `?${querystring.stringify(params)}`;
        }

        const options = {
            method,
            path,
            headers,
            host: this.getCurrentBootstrap(),
            hostname: this.getCurrentBootstrap(),
            port: this.getPort(),
            agent: this.agent,
        };
        if (this._cert && this._key) {
            options.key = this._key;
            options.cert = this._cert;
        }
        const ret = [];
        let retLen = 0;

        // somehow options can get cyclical so fields need to be chosen
        // for display
        log.debug('direct request to MD endpoint', { httpMethod: method,
            host: options.host, port: options.port });
        const req = this.transport.request(options);
        req.setNoDelay();

        if (data) {
            /*
            * Encoding data to binary provides a hot path to write data
            * directly to the socket, without node.js trying to encode the data
            * over and over again.
            */
            const binData = Buffer.from(data, 'utf8');
            req.setHeader('content-type', 'application/octet-stream');
            /*
            * Using Buffer.bytelength is not required here because data is
            * binary encoded, data.length would give us the exact byte length
            */
            req.setHeader('content-length', binData.length);
            req.write(binData);
        }

        req.on('response', res => {
            res.on('data', data => {
                ret.push(data);
                retLen += data.length;
            }).on('error', callback).on('end', () => {
                this.endRespond(res, Buffer.concat(ret, retLen).toString(),
                    log, callback);
            });
        }).on('error', error => {
            // covers system errors like ECONNREFUSED, ECONNRESET etc.
            log.error('error sending request to metadata', { error });
            callback(errors.InternalError);
        }).end();
    }
}

module.exports = RESTClient;
