const RESTClient = require("./RESTClient");

class PromiseRESTClient extends RESTClient {
  async getRaftInformation(bucketName, reqUids, reqLogger) {
    return new Promise((resolve, reject) => {
      super.getRaftInformation(
        bucketName,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  async getBucketLeader(bucketName, reqUids, reqLogger) {
    return new Promise((resolve, reject) => {
      super.getBucketLeader(
        bucketName,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  // <BucketAndObjectOperations>

  async getBucketAttributes(bucketName, reqUids, reqLogger) {
    return new Promise((resolve, reject) => {
      super.getBucketAttributes(
        bucketName,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  async putBucketAttributes(bucketName, reqUids, attributes, reqLogger) {
    return new Promise((resolve, reject) => {
      super.putBucketAttributes(
        bucketName,
        reqUids,
        attributes,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  async createBucket(bucketName, reqUids, attributes, reqLogger) {
    return new Promise((resolve, reject) => {
      super.createBucket(
        bucketName,
        reqUids,
        attributes,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  async deleteBucket(bucketName, reqUids, reqLogger) {
    return new Promise((resolve, reject) => {
      super.deleteBucket(
        bucketName,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
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
   * @param {object} params - extra parameters for the case of versioning
   * @param {werelogs.Logger} [reqLogger] - Logger instance
   *
   * @return {undefined}
   */
  async putObject(bucketName, objName, objVal, reqUids, params, reqLogger) {
    return new Promise((resolve, reject) => {
      super.putObject(
        bucketName,
        objName,
        objVal,
        reqUids,
        params,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
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
   * @param {object} params - extra parameters for the case of versioning
   * @param {werelogs.Logger} [reqLogger] - Logger instance
   *
   * @return {undefined}
   */
  async getObject(bucketName, objName, reqUids, params, reqLogger) {
    return new Promise((resolve, reject) => {
      super.getObject(
        bucketName,
        objName,
        reqUids,
        params,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  /**
   * Get the attributes of a bucket and an object or a version of an object.
   *
   * @param {string} bucketName - bucket name
   * @param {string} objName - the name of the object
   * @param {string} reqUids - the identifier of the request
   * @param {object} params - extra parameters for the case of versioning
   * @param {werelogs.Logger} [reqLogger] - Logger instance
   *
   * @return {undefined}
   */
  async getBucketAndObject(bucketName, objName, reqUids, params, reqLogger) {
    return new Promise((resolve, reject) => {
      super.getBucketAndObject(
        bucketName,
        objName,
        reqUids,
        params,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
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
   * @param {object} params - extra parameters for the case of versioning
   * @param {werelogs.Logger} [reqLogger] - Logger instance
   *
   * @return {undefined}
   */
  async deleteObject(bucketName, objName, reqUids, params, reqLogger) {
    return new Promise((resolve, reject) => {
      super.deleteObject(
        bucketName,
        objName,
        reqUids,
        params,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  /**
   * List objects or versions of objects of a bucket.
   *
   * @param {string} bucketName - bucket name
   * @param {string} reqUids - the identifier of the request
   * @param {object} params - parameters for listing, now includes listing
   *                           all versions of all object of a bucket
   *                           Example: GET /foo?versions&delimiter=&prefix=
   * @param {werelogs.Logger} [reqLogger] - Logger instance
   *
   * @return {undefined}
   */
  async listObject(bucketName, reqUids, params, reqLogger) {
    return new Promise((resolve, reject) => {
      super.listObject(
        bucketName,
        reqUids,
        params,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  // </BucketAndObjectOperations>

  async healthcheck(log, callback) {
    return new Promise((resolve, reject) => {
      super.healthcheck(
        log,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        callback
      );
    });
  }

  async livecheck(log, callback) {
    return new Promise((resolve, reject) => {
      super.livecheck(
        log,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        callback
      );
    });
  }

  /**
   *   send a request to get all raft sessions
   *   get server's response and return it
   *
   *   @param {string} reqUids - the identifier of the request
   *   @param {werelogs.Logger} [reqLogger] - Logger instance
   *   @return {undefined}
   */
  async getAllRafts(reqUids, reqLogger) {
    return new Promise((resolve, reject) => {
      super.getAllRafts(
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        reqLogger
      );
    });
  }

  /**
   *   Get raft logs from bucketd
   *   get server's response and return it as a stream
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
   *   @param {werelogs.Logger} [logger] - Logger instance
   *   @return {undefined}
   */
  async getRaftLog(raftId, start, limit, targetLeader, reqUids, logger) {
    return new Promise((resolve, reject) => {
      super.getRaftLog(
        raftId,
        start,
        limit,
        targetLeader,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        logger
      );
    });
  }

  /**
   *   Get list of buckets associated with raft session
   *
   *   @param {string} raftId - raft session id
   *   @param {string} reqUids - the identifier of the request
   *   @param {werelogs.Logger} [logger] - Logger instance
   *   @return {undefined}
   */
  async getRaftBuckets(raftId, reqUids, logger) {
    return new Promise((resolve, reject) => {
      super.getRaftBuckets(
        raftId,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        logger
      );
    });
  }

  /**
   *   Get proper bucket information like raftSessionId or other statuses
   *
   *   @param {string} bucketName - bucketName
   *   @param {string} reqUids - the identifier of the request
   *   @param {werelogs.Logger} [logger] - Logger instance
   *   @return {undefined}
   */
  async getBucketInformation(bucketName, reqUids, logger) {
    return new Promise((resolve, reject) => {
      super.getBucketInformation(
        bucketName,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        logger
      );
    });
  }

  /**
   * @typedef {Object} BatchOperation
   * @property {string} key - key to operate on
   * @property {string} [type] - type of operation
   * @property {string} [value] - JSON stringified value
   */

  /**
   * Execute an atomic batch of operations on a bucket
   * @param {string} bucketName - bucketName
   * @param {Array<BatchOperation>} batch - batch of operations
   * @param {string} reqUids - the identifier of the request
   * @param {werelogs.Logger} [logger] - Logger instance
   * @return {undefined}
   */
  async execBatch(bucketName, batch, reqUids, logger) {
    return new Promise((resolve, reject) => {
      super.execBatch(
        bucketName,
        batch,
        reqUids,
        (err, res) => {
          if (err) {
            return reject(err);
          }
          return resolve(res);
        },
        logger
      );
    });
  }
}

module.exports = { RESTClient: PromiseRESTClient };
