
/**
 * executable lambda runner script
 */
var spawn = require('child_process').spawn;
exports = module.exports = run;
exports.handler = run;

// wtf is this? it's something amazon recommends
// at least as of april 2015:
// https://aws.amazon.com/blogs/compute/running-executables-in-aws-lambda/
process.env.PATH = process.env.PATH + ':' + process.env.LAMBDA_TASK_ROOT;

var contexts = {};
var cproc;

cproc = spawn('./examples', [], {
  cwd: process.cwd(),
  env: process.env,
  stdio: [
    'pipe',
    'pipe',
    'pipe'
  ]
});


cproc.on('error', function (err) {
  console.error('[nodejs] process error')
  console.error(err.stack);
});

cproc.on('exit', function (c) {
  var time = Date.now()
  console.error('script exited w/status code %d', c);
  console.error('total time %d', time);
  process.exit(0);
});

cproc.stdout.on('data', function (chk) {
  console.log('stdoutdata', chk.toString('utf8'));
  var data = JSON.parse(chk);
  var ctx = contexts[data.requestId];

  if (ctx) {
    ctx.succeed(data);
    delete contexts[data.requestId];
  }
});

/**
 * Starts an arbitrary executable and then streams data to it
 * over stdout
 * @param  {Object} evt 
 * @param  {Object} ctx 
 */
function run(evt, ctx) {
  var start = Date.now();
  contexts[ctx.awsRequestId] = ctx;
  cproc.stdin.write(JSON.stringify({context: ctx, event: evt}) + '\n');
}