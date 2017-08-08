'use strict'; // eslint-disable-line strict

const readline = require('readline');

const RESTClient = require('./RESTClient');
const host = process.argv[2] || 'localhost';
const client = new RESTClient([host]);
const werelogs = require('werelogs');

const logging = new werelogs.Logger('BucketClientShell');
const bucket = {};
let rl;

const commands = [
    'createbucket ',
    'deletebucket ',
    'listobject ',
    'deleteobject ',
    'putobject ',
    'getobject ',
    'getbucketattr ',
    'putbucketattr ',
    'exit ',
];

const listParams = [
    'prefix=',
    'marker=',
    'maxKeys=',
    'delimiter=',
];

function find(tab, string) {
    const hits = tab.filter(c => c.indexOf(string) === 0);
    if (hits.length !== 0) {
        return [hits, string];
    } else if (tab.length <= 1) {
        return [[string], string];
    }
    return [tab, string];
}

function completer(line) {
    const tab = line.split(' ');
    if (tab.length === 1) {
        return find(commands, line);
    } else if (tab.length === 2) {
        return find(Object.keys(bucket), tab[1]);
    } else if (tab[0] === 'listobject') {
        return find(listParams, tab[tab.length - 1]);
    } else if (tab.length === 3) {
        if (bucket[`${tab[1]} `] !== undefined) {
            return find(bucket[`${tab[1]} `], tab[2]);
        }
    }
    return line;
}

function makeJSON(str) {
    const tab = str.split(';');
    const obj = {};
    tab.forEach(elem => {
        const tmp = elem.split('=');
        const key = tmp[0];
        const val = tmp[1];
        obj[key] = val;
    });
    return obj;
}

function createBucket(tab, callback) {
    if (tab.length !== 3) {
        callback('createbucket : usage : <bucket_name> <attributes>');
    } else {
        const begin = process.hrtime();
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        client.createBucket(tab[1], reqUids, tab[2], err => {
            const end = process.hrtime(begin);
            rl.prompt();
            process.stdout.write(`time : ${end[1] / 1e6} .ms\n`);
            if (bucket[`${tab[1]} `] === undefined) {
                bucket[`${tab[1]} `] = [];
            }
            callback(err, `Database ${tab[1]} successfully created`);
        });
    }
}

function deleteBucket(tab, callback) {
    if (tab.length !== 2) {
        callback('deletebucket : usage : <bucket_name>');
    } else {
        if (bucket[`${tab[1]} `] === undefined) {
            bucket[`${tab[1]} `] = [];
        }
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        const begin = process.hrtime();
        client.deleteBucket(tab[1], reqUids, err => {
            const end = process.hrtime(begin);
            rl.prompt();
            process.stdout.write(`time : ${end[1] / 1e6} .ms\n`);
            if (err === null || err === undefined) {
                if (bucket[`${tab[1]} `] !== undefined) {
                    delete bucket[`${tab[1]} `];
                }
            }
            callback(err, `Database ${tab[1]} succesfully deleted`);
        });
    }
}

function putObject(tab, callback) {
    if (tab.length !== 4) {
        callback('putObject : usage : <bucket_name> <obj_name> <obj_value>');
    } else {
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        const begin = process.hrtime();
        client.putObject(tab[1], tab[2], tab[3], reqUids, err => {
            const end = process.hrtime(begin);
            rl.prompt();
            process.stdout.write(`time : ${end[1] / 1e6} .ms\n`);
            if (err === null || err === undefined) {
                if (bucket[`${tab[1]} `] === undefined) {
                    bucket[`${tab[1]} `] = [];
                }
                bucket[`${tab[1]} `].push(`${tab[2]} `);
            }
            callback(err, `Successfully put ${tab[2]}`);
        });
    }
}

function getBucketAttributes(tab, callback) {
    if (tab.length !== 2) {
        callback('getBucketAttributes: usage: <bucketName>');
    } else {
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        client.getBucketAttributes(tab[1], reqUids, callback);
    }
}

function putBucketAttributes(tab, callback) {
    if (tab.length !== 3) {
        callback('putBucketAttributes: usage: '
                + '<bucketName> <bucketAttributes>');
    } else {
        const attributes = JSON.stringify(makeJSON(tab[2]));
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        client.putBucketAttributes(tab[1], reqUids, attributes, callback);
    }
}

function getObject(tab, callback) {
    if (tab.length !== 3) {
        callback('getObject : usage : <bucket_name> <obj_name>');
    } else {
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        const begin = process.hrtime();
        client.getObject(tab[1], tab[2], reqUids, (err, value) => {
            const end = process.hrtime(begin);
            rl.prompt();
            process.stdout.write(`time : ${end[1] / 1e6} .ms\n`);
            if (err === null || err === undefined) {
                if (bucket[`${tab[1]} `] === undefined) {
                    bucket[`${tab[1]} `] = [];
                } else if (bucket[`${tab[1]} `].indexOf(`${tab[2]} `) >= 0) {
                    bucket[`${tab[1]} `].push(`${tab[2]} `);
                }
            }
            callback(err, `value ${value}`);
        });
    }
}

function deleteObject(tab, callback) {
    if (tab.length !== 3) {
        callback('getObject : usage : <bucket_name> <obj_name>');
    } else {
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        const begin = process.hrtime();
        client.deleteObject(tab[1], tab[2], reqUids, (err, value) => {
            const end = process.hrtime(begin);
            rl.prompt();
            process.stdout.write(`time : ${end[1] / 1e6} .ms\n`);
            if (err === null || err === undefined) {
                if (bucket[`${tab[1]} `] !== undefined) {
                    const tmp = bucket[`${tab[1]} `];
                    tmp.splice(bucket[`${tab[1]} `].indexOf(`${tab[2]} `), 1);
                }
            }
            callback(err, value);
        });
    }
}

function listObject(tab, callback) {
    let error = 'listObject : usage : ';
    error += '<bucket_name> [ params=value || ex: prefix=test ]';
    if (tab.length < 2) {
        callback(error);
    } else {
        const params = {};
        tab.splice(2).forEach(curVal => {
            const split = curVal.split('=');
            const index = listParams.indexOf(`${split[0]}=`);
            if (index < 0) {
                return callback(error);
            }
            params[split[0]] = curVal.substring(listParams[index].length);
            return {};
        });
        const logger = logging.newRequestLogger();
        const reqUids = logger.getSerializedUids();
        const begin = process.hrtime();
        client.listObject(tab[1], reqUids, params, (err, value) => {
            const end = process.hrtime(begin);
            rl.prompt();
            process.stdout.write(`time : ${end[1] / 1e6} .ms\n`);
            if (err !== null || err !== undefined) {
                if (bucket[`${tab[1]} `] === undefined) {
                    bucket[tab[1]] = [];
                }
            }
            callback(err, value);
        });
    }
}

function error(tab) {
    rl.prompt();
    if (tab[0] !== '') {
        process.stdout.write(`command not found : ${tab[0]}\n`);
    }
}

function print(err, value) {
    rl.prompt();
    if (err) {
        if (typeof err === 'object') {
            process.stdout.write(` err :  ${err.code} : ${err.message}\n`);
        } else {
            process.stdout.write(` err : ${err}\n`);
        }
    } else {
        rl.prompt();
        process.stdout.write(value);
        process.stdout.write('\n');
    }
    rl.prompt();
}

function parseCommand(tab) {
    const index = commands.indexOf(`${tab[0]} `);
    if (index > -1) {
        switch (index) {
        case 0 :
            createBucket(tab, print);
            break;
        case 1 :
            deleteBucket(tab, print);
            break;
        case 2 :
            listObject(tab, print);
            break;
        case 3 :
            deleteObject(tab, print);
            break;
        case 4 :
            putObject(tab, print);
            break;
        case 5 :
            getObject(tab, print);
            break;
        case 6 :
            getBucketAttributes(tab, print);
            break;
        case 7 :
            putBucketAttributes(tab, print);
            break;
        default :
            rl.close();
            break;
        }
    } else {
        error(tab);
    }
}

function main() {
    rl = readline.createInterface(process.stdin, process.stdout, completer);
    rl.setPrompt('client> ');
    rl.prompt();

    rl.on('line', l => {
        const line = l.trim();
        const tab = line.split(' ');
        parseCommand(tab);
        rl.prompt();
    }).on('close', () => {
        process.stdout.write('client > Exiting process.\n');
        process.exit(0);
    });
}

module.exports = main;
