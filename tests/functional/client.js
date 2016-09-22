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
});
