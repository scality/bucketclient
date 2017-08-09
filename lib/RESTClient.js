'use strict'; // eslint-disable-line strict

const assert = require('assert');
const http = require('http');
const https = require('https');
const querystring = require('querystring');
const werelogs = require('werelogs');

const errors = require('arsenal').errors;

function errorMap(mdError) {
    const map = {
        NoSuchBucket: 'NoSuchBucket',
        BucketAlreadyExists: 'BucketAlreadyExists',
        NoSuchKey: 'NoSuchKey',
        DBNotFound: 'NoSuchBucket',
        DBAlreadyExists: 'BucketAlreadyExists',
        ObjNotFound: 'NoSuchKey',
        NotImplemented: 'NotImplemented',
        InvalidRange: 'InvalidRange',
        BadRequest: 'BadRequest',
    };
    return map[mdError] ? map[mdError] : mdError;
}

class RESTClient {
    /**
     * Constructor for a REST client to bucketd
     *
     * @param {String|String[]} host - bucketd host (or deprecated bucketd
     *                                 hosts bootstrap list)
     * @param {werelogs.API} [logApi] - object providing a constructor function
     *                                for the Logger object
     * @param {Boolean} useHttps - use https
     * @param {string} key - Private certificate key content
     * @param {string} cert - Public certificate content
     * @param {string} ca - Certificate authority content
     *
     * @return {undefined}
     */
    constructor(host, logApi, useHttps, key, cert, ca) {
        // keep for now, compat reason
        let remoteHost;
        if (host instanceof Array) {
            remoteHost = host[0];
        } else {
            remoteHost = host;
        }
        const hostPort = remoteHost.split(':');
        this.remoteHost = hostPort[0];
        const port = hostPort.length > 1 ?
            Number.parseInt(hostPort[1], 10) : 9000;
        assert.ok(!Number.isNaN(port), `port ${port} is not a number`);
        this.port = port;
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
        this.setupLogging(logApi);
    }

    /*
     * Create a dedicated logger for Bucketclient, from the provided werelogs
     * API instance.
     *
     * @param {werelogs.API} [logApi] - object providing a constructor function
     *                                for the Logger object
     * @return {undefined}
     */
    setupLogging(logApi) {
        this.logging = new (logApi || werelogs).Logger('BucketClient');
    }

    createLogger(reqUids) {
        return reqUids ?
            this.logging.newRequestLoggerFromSerializedUids(reqUids) :
            this.logging.newRequestLogger();
    }

    getPort() {
        return this.port;
    }

    getRaftInformation(bucketName, reqUids, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getRaftInformation', { bucketName });
        this.request('GET', `/default/informations/${bucketName}`, log,
            null, null, callback);
    }

    getBucketLeader(bucketName, reqUids, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getBucketLeader', { bucketName });
        this.request('GET', `/default/leader/${bucketName}`, log,
            null, null, callback);
    }

    // <BucketAndObjectOperations>

    getBucketAttributes(bucketName, reqUids, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('getBucketAttributes', { bucketName });
        this.request('GET', `/default/attributes/${bucketName}`, log,
            null, null, callback);
    }

    putBucketAttributes(bucketName, reqUids, attributes, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof attributes === 'string', 'attributes must be a string');

        log.debug('putBucketAttributes', { bucketName, attr: attributes });
        this.request('POST', `/default/attributes/${bucketName}`, log,
            null, attributes, callback);
    }

    createBucket(bucketName, reqUids, attributes, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof attributes === 'string', 'attributes must be a string');

        log.debug('createBucket', { bucketName, attr: attributes });
        this.request('POST', `/default/bucket/${bucketName}`, log,
            null, attributes, callback);
    }

    deleteBucket(bucketName, reqUids, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');

        log.debug('deleteBucket', { bucketName });
        this.request('DELETE', `/default/bucket/${bucketName}`, log,
            null, null, callback);
    }

    /**
     * Create or udpate an object or a version of an object.
     * Examples:
     * - creating an object: PUT /foo/bar
     * - updating a version: PUT /foo/bar?versionId=1234567890
     * @param {string} bucketName - bucket name
     * @param {string} objName - the name of the object
     * @param {string} objVal - the value of the object
     * @param {string} reqUids - the identifier of the request
     * @param {function} callback - callback
     * @param {object} params - extra parameters for the case of versioning
     * @param {werelogs.Logger} [reqLogger] - Logger instance
     *
     * @return {undefined}
     */
    putObject(bucketName, objName, objVal, reqUids, callback, params,
        reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        assert(typeof objVal === 'string', 'objVal must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);
        log.debug('putObject', { bucketName, objName, path });
        this.request('POST', path, log, params, objVal, callback);
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
     * @param {werelogs.Logger} [reqLogger] - Logger instance
     *
     * @return {undefined}
     */
    getObject(bucketName, objName, reqUids, callback, params, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);

        log.debug('getObject', { bucketName, objName, path });
        this.request('GET', path, log, params, null, callback);
    }

    /**
     * Get the attributes of a bucket and an object or a version of an object.
     *
     * @param {string} bucketName - bucket name
     * @param {string} objName - the name of the object
     * @param {string} reqUids - the identifier of the request
     * @param {function} callback - callback
     * @param {object} params - extra parameters for the case of versioning
     * @param {werelogs.Logger} [reqLogger] - Logger instance
     *
     * @return {undefined}
     */
    getBucketAndObject(bucketName, objName, reqUids, callback, params,
        reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
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
        this.request('GET', path, log, params, null, callback);
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
     * @param {werelogs.Logger} [reqLogger] - Logger instance
     *
     * @return {undefined}
     */
    deleteObject(bucketName, objName, reqUids, callback, params, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        assert(typeof objName === 'string', 'objName must be a string');
        let path = `/default/bucket/${bucketName}/`;
        path += encodeURIComponent(objName);

        log.debug('deleteObject', { bucketName, objName, path, params });
        this.request('DELETE', path, log, params, null, callback);
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
     * @param {werelogs.Logger} [reqLogger] - Logger instance
     *
     * @return {undefined}
     */
    listObject(bucketName, reqUids, params, cb, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        assert(typeof bucketName === 'string', 'bucketName must be a string');
        if (cb !== undefined) {
            assert(typeof params === 'object', 'params must be an object');
        }

        log.debug('listObject', { bucketName, params });
        this.request('GET', `/default/bucket/${bucketName}`, log, params,
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
        log.debug('sending request', { httpMethod: method, beginPath });
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
            host: this.remoteHost,
            hostname: this.remoteHost,
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

    /**
    *   send a request to get all raft sessions
    *   get server's response and return it
    *
    *   @param {string} reqUids - the identifier of the request
    *   @param {callback} callback - callback(err, info)
    *   @param {werelogs.Logger} [reqLogger] - Logger instance
    *   @return {undefined}
    */
    getAllRafts(reqUids, callback, reqLogger) {
        const log = reqLogger || this.createLogger(reqUids);
        const path = '/_/raft_sessions/';
        log.debug('getting all raftSessions');
        this.request('GET', path, log, null, null, callback);
    }

    /**
    *   Get raft logs from bucketd
    *   get server's response and return it
    *
    *   @param {string} raftId - raft session id
    *   @param {number} [start=undefined] - starting sequence number. If it is
    *       not given, its value will be 1
    *   @param {number} [limit=undefined] - maximum number of log records
    *       to return. It is at most of 10K. If it is not given, max 10K logs
    *       would be return.
    *   @param {boolean} [targetLeader=undefined] - true: from leader instead of
    *       follower
    *   @param {string} reqUids - the identifier of the request
    *   @param {callback} callback - callback(err, info)
    *   @param {werelogs.Logger} [logger] - Logger instance
    *   @return {undefined}
    */
    getRaftLog(raftId, start, limit, targetLeader, reqUids, callback, logger) {
        const log = logger || this.createLogger(reqUids);
        const path = `/_/raft_sessions/${raftId}/log`;
        const query = {};
        if (start !== undefined && start !== null) {
            query.begin = start;
        }
        if (limit !== undefined && limit !== null) {
            query.limit = limit;
        }
        if (targetLeader !== undefined && targetLeader !== null) {
            query.targetLeader = targetLeader;
        }
        log.debug('getting raft log', { raftId, path, query });
        this.request('GET', path, log, query, null, callback);
    }
}

module.exports = RESTClient;
