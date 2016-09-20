'use strict';

const assert = require('assert');

/**
 * Turns the bootstrap list into an array of host/port objects.
 *
 * @param {String[]} list - the list to process
 *
 * @return {{ host: string, port: number }[]}
 */
function parseBootstrapList(list) {
    return list.map(elem => {
        const tmp = elem.split(':');
        const host = tmp[0];
        if (!tmp[1]) {
            return { host, port: 9000 };
        }
        const port = Number.parseInt(tmp[1], 10);
        assert.ok(!Number.isNaN(port), `port ${tmp[1]} is not a number`);
        return { host, port };
    });
}

module.exports = { parseBootstrapList };
