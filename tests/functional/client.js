'use strict'; // eslint-disable-line strict
const assert = require('assert');

const BucketClient = require('../../index').RESTClient;

const bucketName = 'HanSolo';
const bucketAttributes = { status: 'dead' };

describe('Bucket Client tests', function testClient() {
    this.timeout(0);
    let client;

    before('Create the client', () => {
        client = new BucketClient();
    });

    it('should create a new bucket', done => {
        client.createBucket(bucketName, JSON.stringify(bucketAttributes),
                            done);
    });

    it('should try to create the same bucket and fail', done => {
        client.createBucket(bucketName, JSON.stringify(bucketAttributes),
            err => {
                if (err) {
                    const error = new Error(409);
                    error.isExpected = true;
                    assert.deepStrictEqual(err, error);
                    return done();
                }
                done('Did not fail as expected');
            });
    });

    it('should get the created bucket', done => {
        client.getBucketAttributes(bucketName, (err, data) => {
            const ret = JSON.parse(data);
            if (ret.status !== 'dead') {
                return done(new Error('Did not fetch the data correctly'));
            }
            done(err);
        });
    });

    it('should delete the created bucket', done => {
        client.deleteBucket(bucketName, (err) => {
            done(err);
        });
    });

    it('should fetch non-existing bucket, sending back an error', done => {
        client.getBucketAttributes(bucketName, (err) => {
            if (err) {
                const error = new Error(404);
                error.isExpected = true;
                assert.deepStrictEqual(err, error);
                return done();
            }
            done(new Error('Did not fail as expected'));
        });
    });
});
