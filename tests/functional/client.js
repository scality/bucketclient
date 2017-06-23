'use strict'; // eslint-disable-line strict

const assert = require('assert');
const errors = require('arsenal').errors;

const BucketClient = require('../../index').RESTClient;

const bucketName = 'HanSolo';
const bucketAttributes = { status: 'dead' };
const reqUids = 'REQ1';

describe('Bucket Client tests', function testClient() {
    this.timeout(0);
    let client;

    before('Create the client', () => {
        client = new BucketClient();
    });

    it('should create a new bucket', done => {
        client.createBucket(bucketName, reqUids,
                            JSON.stringify(bucketAttributes), done);
    });

    it('should try to create the same bucket and fail', done => {
        client.createBucket(
            bucketName,
            reqUids,
            JSON.stringify(bucketAttributes),
            err => {
                if (err) {
                    errors.BucketAlreadyExists.isExpected = true;
                    assert.deepStrictEqual(err, errors.BucketAlreadyExists);
                    return done();
                }
                return done('Did not fail as expected');
            });
    });
    it('should get the created bucket', done => {
        client.getBucketAttributes(bucketName, reqUids, (err, data) => {
            const ret = JSON.parse(data);
            if (ret.status !== 'dead') {
                return done(new Error('Did not fetch the data correctly'));
            }
            return done(err);
        });
    });

    it('should delete the created bucket', done => {
        client.deleteBucket(bucketName, reqUids, err => done(err));
    });

    it('should fetch non-existing bucket, sending back an error', done => {
        client.getBucketAttributes(bucketName, reqUids, err => {
            if (err) {
                errors.NoSuchBucket.isExpected = true;
                assert.deepStrictEqual(err, errors.NoSuchBucket);
                return done();
            }
            return done(new Error('Did not fail as expected'));
        });
    });

    it('get all raftSessions', done => {
        client.getAllRafts(null, (err, msg) => {
            if (err) {
                return done(err);
            }
            const rafts = JSON.parse(msg);
            assert.strictEqual(Array.isArray(rafts), true);
            assert.strictEqual(rafts.length >= 1, true);
            assert.strictEqual(rafts.every(raft =>
                typeof raft.id === 'number' &&
                Array.isArray(raft.raftMembers)
            ), true);
            return done();
        });
    });


    it('should get logs from bucketd', done => {
        const start = 1;
        const end = 10;
        client.getRaftLog(0, 1, 10, true, null, (err, data) => {
            if (err) {
                return done(err);
            }
            const obj = JSON.parse(data);
            assert.strictEqual(obj.hasOwnProperty('info'), true);
            const info = obj.info;
            assert.strictEqual(['start', 'end', 'cseq', 'prune']
                .every(key => obj.info.hasOwnProperty(key)), true);
            assert.strictEqual(info.start >= start, true);
            assert.strictEqual(info.end <= end, true);
            // NOTE: this check will be removed when pruned logs are
            // retrieved also
            assert.strictEqual(info.prune <= info.start, true);
            assert.strictEqual(info.cseq >= info.end, true);
            assert.strictEqual(obj.hasOwnProperty('log'), true);
            const logs = obj.log;
            assert.strictEqual(Array.isArray(logs), true);
            assert.strictEqual(logs.length >= 1, true);
            assert.strictEqual(logs.every(log =>
                (typeof log === 'object') && (Object.keys(log).length > 0)
            ), true);
            return done();
        });
    });
});
