const assert = require('assert');

const parse = require('../../lib/bootstraplist').parseBootstrapList;

describe('Bootstrap list tests', () => {
    const list = new Array(1024).fill(0)
                                .map((item, index) => `localhost:${index}`);

    it('should give back the right list for each member', () => {
        parse(list).forEach((item, index) => {
            assert.strictEqual(item.host, 'localhost');
            assert.strictEqual(item.port, index);
        });
    });

    it('should give back the right port', () => {
        list.forEach((item, index) => {
            assert.strictEqual(parse([item])[0].port, index);
        });
    });

    it('should default to port 9000', () => {
        parse(new Array(200).fill(0)
                            .map((item, index) => `${index}`)
        ).forEach((item, index) => {
            assert.strictEqual(item.host, index.toString());
            assert.strictEqual(item.port, 9000);
        });
    });
});
