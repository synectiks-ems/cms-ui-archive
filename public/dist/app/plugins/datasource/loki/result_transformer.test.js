import { LogLevel } from 'app/core/logs_model';
import { findCommonLabels, findUniqueLabels, formatLabels, getLogLevel, mergeStreamsToLogs, parseLabels, } from './result_transformer';
describe('getLoglevel()', function () {
    it('returns no log level on empty line', function () {
        expect(getLogLevel('')).toBe(LogLevel.unkown);
    });
    it('returns no log level on when level is part of a word', function () {
        expect(getLogLevel('this is information')).toBe(LogLevel.unkown);
    });
    it('returns same log level for long and short version', function () {
        expect(getLogLevel('[Warn]')).toBe(LogLevel.warning);
        expect(getLogLevel('[Warning]')).toBe(LogLevel.warning);
        expect(getLogLevel('[Warn]')).toBe('warning');
    });
    it('returns log level on line contains a log level', function () {
        expect(getLogLevel('warn: it is looking bad')).toBe(LogLevel.warn);
        expect(getLogLevel('2007-12-12 12:12:12 [WARN]: it is looking bad')).toBe(LogLevel.warn);
    });
    it('returns first log level found', function () {
        expect(getLogLevel('WARN this could be a debug message')).toBe(LogLevel.warn);
    });
});
describe('parseLabels()', function () {
    it('returns no labels on emtpy labels string', function () {
        expect(parseLabels('')).toEqual({});
        expect(parseLabels('{}')).toEqual({});
    });
    it('returns labels on labels string', function () {
        expect(parseLabels('{foo="bar", baz="42"}')).toEqual({ foo: 'bar', baz: '42' });
    });
});
describe('formatLabels()', function () {
    it('returns no labels on emtpy label set', function () {
        expect(formatLabels({})).toEqual('');
        expect(formatLabels({}, 'foo')).toEqual('foo');
    });
    it('returns label string on label set', function () {
        expect(formatLabels({ foo: 'bar', baz: '42' })).toEqual('{baz="42", foo="bar"}');
    });
});
describe('findCommonLabels()', function () {
    it('returns no common labels on empty sets', function () {
        expect(findCommonLabels([{}])).toEqual({});
        expect(findCommonLabels([{}, {}])).toEqual({});
    });
    it('returns no common labels on differing sets', function () {
        expect(findCommonLabels([{ foo: 'bar' }, {}])).toEqual({});
        expect(findCommonLabels([{}, { foo: 'bar' }])).toEqual({});
        expect(findCommonLabels([{ baz: '42' }, { foo: 'bar' }])).toEqual({});
        expect(findCommonLabels([{ foo: '42', baz: 'bar' }, { foo: 'bar' }])).toEqual({});
    });
    it('returns the single labels set as common labels', function () {
        expect(findCommonLabels([{ foo: 'bar' }])).toEqual({ foo: 'bar' });
    });
});
describe('findUniqueLabels()', function () {
    it('returns no uncommon labels on empty sets', function () {
        expect(findUniqueLabels({}, {})).toEqual({});
    });
    it('returns all labels given no common labels', function () {
        expect(findUniqueLabels({ foo: '"bar"' }, {})).toEqual({ foo: '"bar"' });
    });
    it('returns all labels except the common labels', function () {
        expect(findUniqueLabels({ foo: '"bar"', baz: '"42"' }, { foo: '"bar"' })).toEqual({ baz: '"42"' });
    });
});
describe('mergeStreamsToLogs()', function () {
    it('returns empty logs given no streams', function () {
        expect(mergeStreamsToLogs([]).rows).toEqual([]);
    });
    it('returns processed logs from single stream', function () {
        var stream1 = {
            labels: '{foo="bar"}',
            entries: [
                {
                    line: 'WARN boooo',
                    timestamp: '1970-01-01T00:00:00Z',
                },
            ],
        };
        expect(mergeStreamsToLogs([stream1]).rows).toMatchObject([
            {
                entry: 'WARN boooo',
                labels: { foo: 'bar' },
                key: 'EK1970-01-01T00:00:00Z{foo="bar"}',
                logLevel: 'warning',
                uniqueLabels: {},
            },
        ]);
    });
    it('returns merged logs from multiple streams sorted by time and with unique labels', function () {
        var stream1 = {
            labels: '{foo="bar", baz="1"}',
            entries: [
                {
                    line: 'WARN boooo',
                    timestamp: '1970-01-01T00:00:01Z',
                },
            ],
        };
        var stream2 = {
            labels: '{foo="bar", baz="2"}',
            entries: [
                {
                    line: 'INFO 1',
                    timestamp: '1970-01-01T00:00:00Z',
                },
                {
                    line: 'INFO 2',
                    timestamp: '1970-01-01T00:00:02Z',
                },
            ],
        };
        expect(mergeStreamsToLogs([stream1, stream2]).rows).toMatchObject([
            {
                entry: 'INFO 2',
                labels: { foo: 'bar', baz: '2' },
                logLevel: 'info',
                uniqueLabels: { baz: '2' },
            },
            {
                entry: 'WARN boooo',
                labels: { foo: 'bar', baz: '1' },
                logLevel: 'warning',
                uniqueLabels: { baz: '1' },
            },
            {
                entry: 'INFO 1',
                labels: { foo: 'bar', baz: '2' },
                logLevel: 'info',
                uniqueLabels: { baz: '2' },
            },
        ]);
    });
});
//# sourceMappingURL=result_transformer.test.js.map