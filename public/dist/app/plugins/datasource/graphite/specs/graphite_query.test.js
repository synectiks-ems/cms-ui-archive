import gfunc from '../gfunc';
import GraphiteQuery from '../graphite_query';
describe('Graphite query model', function () {
    var ctx = {
        datasource: {
            getFuncDef: gfunc.getFuncDef,
            getFuncDefs: jest.fn().mockReturnValue(Promise.resolve(gfunc.getFuncDefs('1.0'))),
            waitForFuncDefsLoaded: jest.fn().mockReturnValue(Promise.resolve(null)),
            createFuncInstance: gfunc.createFuncInstance,
        },
        templateSrv: {},
        targets: [],
    };
    beforeEach(function () {
        ctx.target = { refId: 'A', target: 'scaleToSeconds(#A, 60)' };
        ctx.queryModel = new GraphiteQuery(ctx.datasource, ctx.target, ctx.templateSrv);
    });
    describe('when updating targets with nested queries', function () {
        beforeEach(function () {
            ctx.target = { refId: 'D', target: 'asPercent(#A, #C)' };
            ctx.targets = [
                { refId: 'A', target: 'first.query.count' },
                { refId: 'B', target: 'second.query.count' },
                { refId: 'C', target: 'diffSeries(#A, #B)' },
                { refId: 'D', target: 'asPercent(#A, #C)' },
            ];
            ctx.queryModel = new GraphiteQuery(ctx.datasource, ctx.target, ctx.templateSrv);
        });
        it('targetFull should include nested queries', function () {
            ctx.queryModel.updateRenderedTarget(ctx.target, ctx.targets);
            var targetFullExpected = 'asPercent(first.query.count, diffSeries(first.query.count, second.query.count))';
            expect(ctx.queryModel.target.targetFull).toBe(targetFullExpected);
        });
        it('should not hang on circular references', function () {
            ctx.target.target = 'asPercent(#A, #B)';
            ctx.targets = [{ refId: 'A', target: 'asPercent(#B, #C)' }, { refId: 'B', target: 'asPercent(#A, #C)' }];
            ctx.queryModel.updateRenderedTarget(ctx.target, ctx.targets);
            // Just ensure updateRenderedTarget() is completed and doesn't hang
            expect(ctx.queryModel.target.targetFull).toBeDefined();
        });
    });
});
//# sourceMappingURL=graphite_query.test.js.map