package main

var jsloader = `
/**
 * lambda-phage entrypoint 
 * at time of writing, AWS Lambda
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
cproc = spawn('./' + require('./cfg').bin), [], {
  cwd: process.cwd(),
  env: process.env,
  stdio: [
    'pipe',
    'pipe',
    'pipe'
  ]
});


cproc.on('error', function (err) {
  console.error('parent process error', err.stack);
});

cproc.on('exit', function (c) {
  console.error('go program exited w/status code %d', c);
  process.exit(c);
});

cproc.stderr.pipe(process.stderr);
cproc.stdout.on('data', function (chk) {

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

function decodeMaybe(s) {
  try {
    return JSON.parse(s);
  } catch (e) {
    console.error('invalid json from child', e);
    return {};
  }
}
`
