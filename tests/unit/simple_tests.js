'use strict';

const assert = require('assert');
const http = require('http');
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
            port: 4242
    }
};
const nonExistBucket = { name: 'Ford' };
const reqUids = 'REQ1';

function makeResponse(res, code, message) {
    res.statusCode = code;
    res.statusMessage = message;
}

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
    res.end();
}

describe('Unit tests with mockup server', function tests() {
    let server;
    let client;

    beforeEach('start server', done => {
        server = http.createServer(handler).listen(9000);
        client = new RESTClient([ 'localhost', ]);
        done();
    });

    afterEach('stop server', () => { server.close(); });

    it('should create a new non-existing bucket', done => {
        client.createBucket(nonExistBucket.name, reqUids,
                            '{ status: "dead" }', done);
    });

    it('should try to create an already existing bucket and fail', done => {
        client.createBucket(existBucket.name, reqUids, '{}', (err) => {
            if (err) {
                const error = new Error('BucketAlreadyExists');
                error.isExpected = true;
                assert.deepStrictEqual(err, error);
                return done();
            }
            done('Did not fail as expected');
        });
    });

    it('should get an existing bucket', done => {
        client.getBucketAttributes(existBucket.name, reqUids, (err, data) => {
            const ret = JSON.parse(data);
            assert.deepStrictEqual(ret, existBucket.value);
            done(err);
        });
    });

    it('should get Raft informations on an existing bucket', done => {
        client.getRaftInformation(existBucket.name, reqUids, (err, data) => {
            const ret = JSON.parse(data);
            assert.deepStrictEqual(ret, existBucket.raftInformation);
            done(err);
        });
    });

    it('should get Raft informations on an existing bucket', done => {
        client.getRaftInformation(nonExistBucket.name, reqUids, (err, data) => {
            const error = new Error('NoSuchBucket');
            error.isExpected = true;
            assert.deepStrictEqual(err, error);
            return done();
            done(err);
        });
    });

    it('should fetch non-existing bucket, sending back an error', done => {
        client.getBucketAttributes(nonExistBucket.name, reqUids, (err) => {
            if (err) {
                const error = new Error('NoSuchBucket');
                error.isExpected = true;
                assert.deepStrictEqual(err, error);
                return done();
            }
            done(new Error('Did not fail as expected'));
        });
    });

    it('should delete an existing bucket', done => {
        client.deleteBucket(existBucket.name, reqUids, done);
    });

    it('should fetch non-existing bucket, sending back an error', done => {
        client.deleteBucket(nonExistBucket.name, reqUids, err => {
            if (err) {
                const error = new Error('NoSuchBucket');
                error.isExpected = true;
                assert.deepStrictEqual(err, error);
                return done();
            }
            done(new Error('Did not fail as expected'));
        });
    });
});
