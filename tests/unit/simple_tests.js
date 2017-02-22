'use strict'; // eslint-disable-line strict

const assert = require('assert');
const fs = require('fs');
const http = require('http');
const https = require('https');

const errors = require('arsenal').errors;

const RESTClient = require('../../index.js').RESTClient;

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
};
const nonExistBucket = { name: 'Ford' };
const reqUids = 'REQ1';

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
        } else if (req.url === `/default/bucket/${nonExistBucket.name}`) {
            makeResponse(res, 200, 'OK');
        }
    } else if (req.method === 'GET') {
        if (req.url === `/default/attributes/${existBucket.name}`) {
            makeResponse(res, 200, 'OK');
            res.write(JSON.stringify(existBucket.value));
        } else if (req.url === `/default/informations/${existBucket.name}`) {
            makeResponse(res, 200, 'OK');
            return res.end(JSON.stringify(existBucket.raftInformation));
        } else if (req.url === '/_/healthcheck') {
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
            server = e.s(handler).listen(9000, done).on('error', done);
            client = e.c;
        });

        afterEach('stop server', () => { server.close(); });

        it('should create a new non-existing bucket', done => {
            client.createBucket(nonExistBucket.name, reqUids,
                                '{ status: "dead" }', done);
        });

        it('should try to create an already existing bucket and fail', done => {
            client.createBucket(existBucket.name, reqUids, '{}', err => {
                if (err) {
                    const error = errors.BucketAlreadyExists;
                    error.isExpected = true;
                    assert.deepStrictEqual(err, error);
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

        it('should get Raft informations on an existing bucket', done => {
            client.getRaftInformation(existBucket.name, reqUids,
                (err, data) => {
                    const ret = JSON.parse(data);
                    assert.deepStrictEqual(ret, existBucket.raftInformation);
                    done(err);
                });
        });

        it('should get Raft informations on an unexisting bucket', done => {
            client.getRaftInformation(nonExistBucket.name, reqUids,
                err => {
                    const error = errors.NoSuchBucket;
                    error.isExpected = true;
                    assert.deepStrictEqual(err, error);
                    return done();
                });
        });

        it('should fetch non-existing bucket, sending back an error', done => {
            client.getBucketAttributes(nonExistBucket.name, reqUids, err => {
                if (err) {
                    const error = errors.NoSuchBucket;
                    error.isExpected = true;
                    assert.deepStrictEqual(err, error);
                    return done();
                }
                return done(new Error('Did not fail as expected'));
            });
        });

        it('should delete an existing bucket', done => {
            client.deleteBucket(existBucket.name, reqUids, done);
        });

        it('should fetch non-existing bucket, sending back an error', done => {
            client.deleteBucket(nonExistBucket.name, reqUids, err => {
                if (err) {
                    const error = errors.NoSuchBucket;
                    error.isExpected = true;
                    assert.deepStrictEqual(err, error);
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
    });
});
