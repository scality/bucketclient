export type Callback = (error: Error, res?: any) => void;

export type BatchOperation = {
    type?: string,
    key: string,
    value?: string,
};

export type RecreateBackupsParams = {
    formatVersion?: number,
    compression?: string,
    recreateMinBseq?: number,
    recreateMaxBseq?: number,
};

export class RESTClient {
    constructor(
        host: string | string[],
        logApi?: object | null,
        useHttps?: boolean,
        key?: string,
        cert?: string,
        ca?: string,
    );

    setupLogging(logApi?: object | null): void;
    createLogger(reqUids?: string[]): any;
    getPort(): number;
    getRaftInformation(bucketName: string, reqUids: string[], callback: Callback, reqLogger?: any): void;
    getBucketLeader(bucketName: string, reqUids: string[], callback: Callback, reqLogger?: any): void;
    getBucketAttributes(bucketName: string, reqUids: string[], callback: Callback, reqLogger?: any): void;
    putBucketAttributes(bucketName: string, reqUids: string[], attributes: string, callback: Callback, reqLogger: any): void;
    createBucket(bucketName: string, reqUids: string[], attributes: string, callback: Callback, reqLogger: any): void;
    deleteBucket(bucketName: string, reqUids: string[], callback: Callback, reqLogger: any): void;
    putObject(bucketName: string, objName: string, objVal: string, reqUids: string[], callback: Callback, params: any, reqLogger: any): void;
    getObject(bucketName: string, objName: string, reqUids: string[], callback: Callback, params: any, reqLogger: any): void;
    getBucketAndObject(bucketName: string, objName: string, reqUids: string[], callback: Callback, params: any, reqLogger: any): void;
    deleteObject(bucketName: string, objName: string, reqUids: string[], callback: Callback, params: any, reqLogger: any): void;
    listObject(bucketName: string, reqUids: string[], params: any, cb: Callback, reqLogger: any): void;
    getAllRafts(reqUids: string[], callback: Callback, reqLogger: any): void;
    getRaftLog(raftId: string, start: number, limit: number, targetLeader: boolean, reqUids: string[], callback: Callback, reqLogger: any): void;
    getRaftBuckets(raftId: string, reqUids: string[], callback: Callback, reqLogger: any): void;
    getBucketInformation(bucketName: string, reqUids: string[], callback: Callback, reqLogger: any): void;
    execBatch(bucketName: string, batch: Array<BatchOperation>, reqUids: string[], callback: Callback, reqLogger: any): void;
    appendToLog(bucketName: string, batch: Array<BatchOperation>, raftsession: string | null | undefined, reqUids: string[], callback: Callback, reqLogger: any): void;
    reindexBackups(raftId: string, recreateBackupsParams: RecreateBackupsParams | null | undefined, reqUids: string[], callback: Callback, logger?: any): void;
    getBackupReindexJobStatus(raftId: string, reqUids: string[], callback: Callback, logger?: any): void;
    abortBackupReindexJob(raftId: string, reqUids: string[], callback: Callback, logger?: any): void;
    healthcheck(log: any, callback: Callback): void;
    healthcheckSimple(log: any, callback: Callback): void;
    livecheck(log: any, callback: Callback): void;
    private endRespond(res: any, log: any, callback: Callback): void;
    private request(method: string, beginPath: string, log: any, params: any, data: any, callback: Callback): void;
    private requestStreamed(method: string, beginPath: string, log: any, params: any, data: any, callback: Callback): void;
}

export function shell(): void;
