import * as tslib_1 from "tslib";
var _a;
import _ from 'lodash';
import { TimeSeries } from 'app/core/core';
import colors, { getThemeColor } from 'app/core/utils/colors';
export var LogLevel;
(function (LogLevel) {
    LogLevel["crit"] = "critical";
    LogLevel["critical"] = "critical";
    LogLevel["warn"] = "warning";
    LogLevel["warning"] = "warning";
    LogLevel["err"] = "error";
    LogLevel["error"] = "error";
    LogLevel["info"] = "info";
    LogLevel["debug"] = "debug";
    LogLevel["trace"] = "trace";
    LogLevel["unkown"] = "unkown";
})(LogLevel || (LogLevel = {}));
export var LogLevelColor = (_a = {},
    _a[LogLevel.critical] = colors[7],
    _a[LogLevel.warning] = colors[1],
    _a[LogLevel.error] = colors[4],
    _a[LogLevel.info] = colors[0],
    _a[LogLevel.debug] = colors[5],
    _a[LogLevel.trace] = colors[2],
    _a[LogLevel.unkown] = getThemeColor('#8e8e8e', '#dde4ed'),
    _a);
export var LogsMetaKind;
(function (LogsMetaKind) {
    LogsMetaKind[LogsMetaKind["Number"] = 0] = "Number";
    LogsMetaKind[LogsMetaKind["String"] = 1] = "String";
    LogsMetaKind[LogsMetaKind["LabelsMap"] = 2] = "LabelsMap";
})(LogsMetaKind || (LogsMetaKind = {}));
export var LogsDedupDescription;
(function (LogsDedupDescription) {
    LogsDedupDescription["none"] = "No de-duplication";
    LogsDedupDescription["exact"] = "De-duplication of successive lines that are identical, ignoring ISO datetimes.";
    LogsDedupDescription["numbers"] = "De-duplication of successive lines that are identical when ignoring numbers, e.g., IP addresses, latencies.";
    LogsDedupDescription["signature"] = "De-duplication of successive lines that have identical punctuation and whitespace.";
})(LogsDedupDescription || (LogsDedupDescription = {}));
export var LogsDedupStrategy;
(function (LogsDedupStrategy) {
    LogsDedupStrategy["none"] = "none";
    LogsDedupStrategy["exact"] = "exact";
    LogsDedupStrategy["numbers"] = "numbers";
    LogsDedupStrategy["signature"] = "signature";
})(LogsDedupStrategy || (LogsDedupStrategy = {}));
var LOGFMT_REGEXP = /(?:^|\s)(\w+)=("[^"]*"|\S+)/;
export var LogsParsers = {
    JSON: {
        buildMatcher: function (label) { return new RegExp("(?:{|,)\\s*\"" + label + "\"\\s*:\\s*\"?([\\d\\.]+|[^\"]*)\"?"); },
        getFields: function (line) {
            var fields = [];
            try {
                var parsed = JSON.parse(line);
                _.map(parsed, function (value, key) {
                    var fieldMatcher = new RegExp("\"" + key + "\"\\s*:\\s*\"?" + _.escapeRegExp(JSON.stringify(value)) + "\"?");
                    var match = line.match(fieldMatcher);
                    if (match) {
                        fields.push(match[0]);
                    }
                });
            }
            catch (_a) { }
            return fields;
        },
        getLabelFromField: function (field) { return (field.match(/^"(\w+)"\s*:/) || [])[1]; },
        getValueFromField: function (field) { return (field.match(/:\s*(.*)$/) || [])[1]; },
        test: function (line) {
            try {
                return JSON.parse(line);
            }
            catch (error) { }
        },
    },
    logfmt: {
        buildMatcher: function (label) { return new RegExp("(?:^|\\s)" + label + "=(\"[^\"]*\"|\\S+)"); },
        getFields: function (line) {
            var fields = [];
            line.replace(new RegExp(LOGFMT_REGEXP, 'g'), function (substring) {
                fields.push(substring.trim());
                return '';
            });
            return fields;
        },
        getLabelFromField: function (field) { return (field.match(LOGFMT_REGEXP) || [])[1]; },
        getValueFromField: function (field) { return (field.match(LOGFMT_REGEXP) || [])[2]; },
        test: function (line) { return LOGFMT_REGEXP.test(line); },
    },
};
export function calculateFieldStats(rows, extractor) {
    // Consider only rows that satisfy the matcher
    var rowsWithField = rows.filter(function (row) { return extractor.test(row.entry); });
    var rowCount = rowsWithField.length;
    // Get field value counts for eligible rows
    var countsByValue = _.countBy(rowsWithField, function (row) { return row.entry.match(extractor)[1]; });
    var sortedCounts = _.chain(countsByValue)
        .map(function (count, value) { return ({ count: count, value: value, proportion: count / rowCount }); })
        .sortBy('count')
        .reverse()
        .value();
    return sortedCounts;
}
export function calculateLogsLabelStats(rows, label) {
    // Consider only rows that have the given label
    var rowsWithLabel = rows.filter(function (row) { return row.labels[label] !== undefined; });
    var rowCount = rowsWithLabel.length;
    // Get label value counts for eligible rows
    var countsByValue = _.countBy(rowsWithLabel, function (row) { return row.labels[label]; });
    var sortedCounts = _.chain(countsByValue)
        .map(function (count, value) { return ({ count: count, value: value, proportion: count / rowCount }); })
        .sortBy('count')
        .reverse()
        .value();
    return sortedCounts;
}
var isoDateRegexp = /\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-6]\d[,\.]\d+([+-][0-2]\d:[0-5]\d|Z)/g;
function isDuplicateRow(row, other, strategy) {
    switch (strategy) {
        case LogsDedupStrategy.exact:
            // Exact still strips dates
            return row.entry.replace(isoDateRegexp, '') === other.entry.replace(isoDateRegexp, '');
        case LogsDedupStrategy.numbers:
            return row.entry.replace(/\d/g, '') === other.entry.replace(/\d/g, '');
        case LogsDedupStrategy.signature:
            return row.entry.replace(/\w/g, '') === other.entry.replace(/\w/g, '');
        default:
            return false;
    }
}
export function dedupLogRows(logs, strategy) {
    if (strategy === LogsDedupStrategy.none) {
        return logs;
    }
    var dedupedRows = logs.rows.reduce(function (result, row, index, list) {
        var previous = result[result.length - 1];
        if (index > 0 && isDuplicateRow(row, previous, strategy)) {
            previous.duplicates++;
        }
        else {
            row.duplicates = 0;
            result.push(row);
        }
        return result;
    }, []);
    return tslib_1.__assign({}, logs, { rows: dedupedRows });
}
export function getParser(line) {
    var parser;
    try {
        if (LogsParsers.JSON.test(line)) {
            parser = LogsParsers.JSON;
        }
    }
    catch (error) { }
    if (!parser && LogsParsers.logfmt.test(line)) {
        parser = LogsParsers.logfmt;
    }
    return parser;
}
export function filterLogLevels(logs, hiddenLogLevels) {
    if (hiddenLogLevels.size === 0) {
        return logs;
    }
    var filteredRows = logs.rows.reduce(function (result, row, index, list) {
        if (!hiddenLogLevels.has(row.logLevel)) {
            result.push(row);
        }
        return result;
    }, []);
    return tslib_1.__assign({}, logs, { rows: filteredRows });
}
export function makeSeriesForLogs(rows, intervalMs) {
    // currently interval is rangeMs / resolution, which is too low for showing series as bars.
    // need at least 10px per bucket, so we multiply interval by 10. Should be solved higher up the chain
    // when executing queries & interval calculated and not here but this is a temporary fix.
    // intervalMs = intervalMs * 10;
    // Graph time series by log level
    var seriesByLevel = {};
    var bucketSize = intervalMs * 10;
    var seriesList = [];
    for (var _i = 0, rows_1 = rows; _i < rows_1.length; _i++) {
        var row = rows_1[_i];
        var series = seriesByLevel[row.logLevel];
        if (!series) {
            seriesByLevel[row.logLevel] = series = {
                lastTs: null,
                datapoints: [],
                alias: row.logLevel,
                color: LogLevelColor[row.logLevel],
            };
            seriesList.push(series);
        }
        // align time to bucket size
        var time = Math.round(row.timeEpochMs / bucketSize) * bucketSize;
        // Entry for time
        if (time === series.lastTs) {
            series.datapoints[series.datapoints.length - 1][0]++;
        }
        else {
            series.datapoints.push([1, time]);
            series.lastTs = time;
        }
        // add zero to other levels to aid stacking so each level series has same number of points
        for (var _a = 0, seriesList_1 = seriesList; _a < seriesList_1.length; _a++) {
            var other = seriesList_1[_a];
            if (other !== series && other.lastTs !== time) {
                other.datapoints.push([0, time]);
                other.lastTs = time;
            }
        }
    }
    return seriesList.map(function (series) {
        series.datapoints.sort(function (a, b) {
            return a[1] - b[1];
        });
        return new TimeSeries(series);
    });
}
//# sourceMappingURL=logs_model.js.map