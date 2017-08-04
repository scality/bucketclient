#!/usr/bin/env node

const logs = require('werelogs');
logs.configure({
    level: 'info',
    dump: 'error',
});

require('../lib/shell.js')();
