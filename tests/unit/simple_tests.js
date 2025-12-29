'use strict';  

const assert = require('assert');
const async = require('async');
const { EventEmitter } = require('events');
const fs = require('fs');
const http = require('http');
const https = require('https');
const sinon = require('sinon');

const RESTClient = require('../../index').RESTClient;

const existBucket = {
    name: 'Zaphod',
    value: { status: 'alive' },
    raftInformation: {
        term: 1,
        cseq: 0,
        aseq: 5,
        prune: 0,
        ip: '127.0.0.1',
        port: 4242,
    },
    bucketInformation: {
        raftSessionId: 3,
        creating: false,
        deleting: false,
        version: 0,
        leader: {
            host: '127.0.0.4',
            port: 4500,
        },
    },
};
const existObject = {
    name: 'testObject',
    value: { key: 'value', metadata: 'test' },
};
const newBucketName = 'new-bucket';
const newBucketNameWithRSID = 'new-bucket-with-rsid';
const reqUids = 'REQ1';

let includeRaftSessionIdHeader = true;

function makeResponse(res, code, message) {
    /* eslint-disable no-param-reassign */
    res.statusCode = code;
    res.statusMessage = message;
    /* eslint-enable no-param-reassign */
}

const httpsOptions = {
    key: fs.readFileSync('./tests/utils/test.key', 'ascii'),
    cert: fs.readFileSync('./tests/utils/test.crt', 'ascii'),
    ca: [fs.readFileSync('./tests/utils/ca.crt', 'ascii')],
    requestCert: true,
};

const env = {
    http: {
        c: new RESTClient(['bucketclient.testing.local:9000']),
        s: handler => http.createServer(handler),
    },
    https: {
        s: handler => https.createServer(httpsOptions, handler),
        c: new RESTClient(['bucketclient.testing.local:9000'],
            undefined,
            true,
            httpsOptions.key,
            httpsOptions.cert,
            httpsOptions.ca[0]),
    },
};

function handler(req, res) {
    if (req.method === 'POST') {
        if (req.url === `/default/bucket/${existBucket.name}`) {
            makeResponse(res, 409, 'BucketAlreadyExists');
        } else if (req.url === `/default/bucket/${newBucketName}`) {
            makeResponse(res, 200, 'OK');
        } else if (req.url === `/default/bucket/${newBucketNameWithRSID}?raftsession=2`) {
            makeResponse(res, 200, 'OK');
        } else if (req.url === '/_/livecheck') {
            makeResponse(res, 200, 'OK');
        } else {
            assert.fail(`unexpected POST url: ${req.url}`);
        }
    } else if (req.method === 'GET') {
        if (req.url === `/default/attributes/${existBucket.name}`) {
            makeResponse(res, 200, 'OK');
            if (includeRaftSessionIdHeader) {
                res.setHeader('x-scal-raft-session-id', existBucket.bucketInformation.raftSessionId);
            }
            res.write(JSON.stringify(existBucket.value));
        } else if (req.url === `/default/parallel/${existBucket.name}/${existObject.name}`) {
            makeResponse(res, 200, 'OK');
            if (includeRaftSessionIdHeader) {
                res.setHeader('x-scal-raft-session-id', existBucket.bucketInformation.raftSessionId);
            }
            res.write(JSON.stringify(existObject.value));
        } else if (req.url === `/default/informations/${existBucket.name}`) {
            makeResponse(res, 200, 'OK');
            return res.end(JSON.stringify(existBucket.raftInformation));
        } else if (req.url === `/_/buckets/${existBucket.name}`) {
            makeResponse(res, 200, 'OK');
            return res.end(JSON.stringify(existBucket.bucketInformation));
        } else if (req.url === '/_/healthcheck') {
            makeResponse(res, 200, 'OK');
        } else if (req.url === '/_/healthcheck/simple') {
            makeResponse(res, 200, 'OK');
        } else {
            makeResponse(res, 404, 'NoSuchBucket');
        }
    } else if (req.method === 'DELETE') {
        if (req.url === `/default/bucket/${existBucket.name}`) {
            makeResponse(res, 200, 'OK');
        } else {
            makeResponse(res, 404, 'NoSuchBucket');
        }
    }
    return res.end();
}

Object.keys(env).forEach(key => {
    const e = env[key];
    describe(`Unit tests with mockup ${key} server`, () => {
        let server;
        let client;

        beforeEach('start server', done => {
            includeRaftSessionIdHeader = true;
            client = e.c;
            server = e.s(handler).on('error', done).listen(9000, done);
        });

        afterEach('stop server', done => {
            client.agent.destroy();
            server.close(done);
        });

        it('should create a new non-existing bucket', done => {
            client.createBucket(newBucketName, reqUids,
                '{}', done);
        });

        it('should create a new non-existing bucket on a given RAFT session', done => {
            client.createBucket(newBucketNameWithRSID, reqUids,
                '{}', done, null, { raftsession: 2 });
        });

        it('should try to create an already existing bucket and fail', done => {
            client.createBucket(existBucket.name, reqUids, '{}', err => {
                if (err) {
                    assert(err.is.BucketAlreadyExists);
                    assert.strictEqual(err.isExpected, true);
                    return done();
                }
                return done('Did not fail as expected');
            });
        });

        it('should get an existing bucket', done => {
            client.getBucketAttributes(existBucket.name, reqUids,
                (err, data) => {
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existBucket.value);
                    done(err);
                });
        });

        it('should get an existing bucket with raftSessionId', done => {
            includeRaftSessionIdHeader = true;
            client.getBucketAttributes(existBucket.name, reqUids,
                (err, data, raftSessionId) => {
                    assert.ifError(err);
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existBucket.value);
                    assert.strictEqual(raftSessionId, existBucket.bucketInformation.raftSessionId);
                    done();
                });
        });

        it('should get bucket attributes without raftSessionId when header absent', done => {
            includeRaftSessionIdHeader = false;
            client.getBucketAttributes(existBucket.name, reqUids,
                (err, data, raftSessionId) => {
                    assert.ifError(err);
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existBucket.value);
                    assert.strictEqual(raftSessionId, undefined);
                    done();
                });
        });

        it('should get bucket and object with raftSessionId', done => {
            includeRaftSessionIdHeader = true;
            client.getBucketAndObject(existBucket.name, existObject.name, reqUids,
                (err, data, raftSessionId) => {
                    assert.ifError(err);
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existObject.value);
                    assert.strictEqual(raftSessionId, existBucket.bucketInformation.raftSessionId);
                    done();
                });
        });

        it('should get bucket and object without raftSessionId when header absent', done => {
            includeRaftSessionIdHeader = false;
            client.getBucketAndObject(existBucket.name, existObject.name, reqUids,
                (err, data, raftSessionId) => {
                    assert.ifError(err);
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existObject.value);
                    assert.strictEqual(raftSessionId, undefined);
                    done();
                });
        });

        it('should get Raft informations on an existing bucket', done => {
            client.getRaftInformation(existBucket.name, reqUids,
                (err, data) => {
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existBucket.raftInformation);
                    done(err);
                });
        });

        it('should get Raft informations on an unexisting bucket', done => {
            client.getRaftInformation(newBucketName, reqUids,
                err => {
                    assert(err.is.NoSuchBucket);
                    assert.strictEqual(err.isExpected, true);

                    return done();
                });
        });

        it('should get Bucket informations on an existing bucket', done => {
            client.getBucketInformation(existBucket.name, reqUids,
                (err, data) => {
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existBucket.bucketInformation);
                    done(err);
                });
        });

        it('should get Bucket informations on an unexisting bucket', done => {
            client.getRaftInformation(newBucketName, reqUids,
                err => {
                    assert(err.is.NoSuchBucket);
                    assert.strictEqual(err.isExpected, true);

                    return done();
                });
        });

        it('should fetch non-existing bucket, sending back an error', done => {
            client.getBucketAttributes(newBucketName, reqUids, err => {
                if (err) {
                    assert(err.is.NoSuchBucket);
                    assert.strictEqual(err.isExpected, true);

                    return done();
                }
                return done(new Error('Did not fail as expected'));
            });
        });

        it('should delete an existing bucket', done => {
            client.deleteBucket(existBucket.name, reqUids, done);
        });

        it('should fetch non-existing bucket, sending back an error', done => {
            client.deleteBucket(newBucketName, reqUids, err => {
                if (err) {
                    assert(err.is.NoSuchBucket);
                    assert.strictEqual(err.isExpected, true);

                    return done();
                }
                return done(new Error('Did not fail as expected'));
            });
        });

        it('should return 200 on healthcheck request', done => {
            const log = e.c.createLogger();
            client.healthcheck(log, err => {
                assert.deepStrictEqual(err, null);
                return done();
            });
        });

        it('should return 200 on healthcheckSimple request', done => {
            const log = e.c.createLogger();
            client.healthcheckSimple(log, err => {
                assert.deepStrictEqual(err, null);
                return done();
            });
        });


        it('should return 200 on livecheck request', done => {
            const log = e.c.createLogger();
            client.livecheck(log, err => {
                assert.deepStrictEqual(err, null);
                return done();
            });
        });

        it('should be able to reuse a connection after an HTTP error status is received', done => {
            async.timesSeries(
                10,
                (i, next) => client.getRaftInformation(newBucketName, reqUids, err => {
                    assert(err.is.NoSuchBucket);
                    // trigger the node.js event loop after each iteration, to let the HTTP module
                    // a chance to cleanup the unique connection state and reuse it
                    setTimeout(next, 10);
                }),
                () => {
                    assert(client.agent.totalSocketCount === 1,
                        `expected total socket count to be 1, got ${client.agent.totalSocketCount}`);
                    done();
                });
        });
    });
});

describe('with a stubbed http.request', () => {
    let httpRequestStub;
    let client;

    beforeEach(() => {
        httpRequestStub = sinon.stub(http, 'request');
        client = new RESTClient(['bucketclient.testing.local:9000']);
    });

    afterEach(() => {
        httpRequestStub.restore();
        client.agent.destroy();
    });

    it('should call callback only once if HTTP request stream emits multiple errors', () => {
        const mockRequest = new EventEmitter();
        mockRequest.setNoDelay = () => {};
        mockRequest.end = () => {
            // emit two errors on the request stream
            mockRequest.emit('error', new Error('first error'));
            mockRequest.emit('error', new Error('second error'));
        };
        httpRequestStub.callsFake(() => mockRequest);

        const callbackSpy = sinon.spy();
        client.getObject('foobucket', 'fookey', '', callbackSpy);
        assert.strictEqual(callbackSpy.callCount, 1);
        assert.strictEqual(callbackSpy.getCall(0).args[0].code, 500);
    });

    it('should call callback only once if HTTP response stream emits multiple errors', () => {
        const mockResponse = new EventEmitter();
        mockResponse.statusCode = 200;

        const mockRequest = new EventEmitter();
        mockRequest.setNoDelay = () => {};
        mockRequest.end = () => {
            mockRequest.emit('response', mockResponse);

            // emit two errors on the response stream
            mockResponse.emit('error', new Error('first error'));
            mockResponse.emit('error', new Error('second error'));
        };
        httpRequestStub.callsFake(() => mockRequest);

        const callbackSpy = sinon.spy();
        client.getObject('foobucket', 'fookey', '', callbackSpy);
        assert.strictEqual(callbackSpy.callCount, 1);
        assert.strictEqual(callbackSpy.getCall(0).args[0].code, 500);
    });

    it('should return only body when returnHeaders is false', done => {
        const mockResponse = new EventEmitter();
        mockResponse.statusCode = 200;
        mockResponse.headers = { 'x-scal-raft-session-id': '42' };

        const mockRequest = new EventEmitter();
        mockRequest.setNoDelay = () => {};
        mockRequest.end = () => {
            mockRequest.emit('response', mockResponse);
            mockResponse.emit('data', Buffer.from('test body'));
            mockResponse.emit('end');
        };
        httpRequestStub.callsFake(() => mockRequest);

        const log = client.createLogger();
        client.request('GET', '/test', log, null, null, (err, body, headers) => {
            assert.ifError(err);
            assert.strictEqual(body, 'test body');
            assert.strictEqual(headers, undefined);
            done();
        }, false);
    });

    it('should return body and headers when returnHeaders is true', done => {
        const mockResponse = new EventEmitter();
        mockResponse.statusCode = 200;
        mockResponse.headers = { 'x-scal-raft-session-id': '42' };

        const mockRequest = new EventEmitter();
        mockRequest.setNoDelay = () => {};
        mockRequest.end = () => {
            mockRequest.emit('response', mockResponse);
            mockResponse.emit('data', Buffer.from('test body'));
            mockResponse.emit('end');
        };
        httpRequestStub.callsFake(() => mockRequest);

        const log = client.createLogger();
        client.request('GET', '/test', log, null, null, (err, body, headers) => {
            assert.ifError(err);
            assert.strictEqual(body, 'test body');
            assert.deepStrictEqual(headers, { 'x-scal-raft-session-id': '42' });
            done();
        }, true);
    });
});
