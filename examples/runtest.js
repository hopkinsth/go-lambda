var m = require('./');

var evt = {
  "idk": "somestuff"
};
var ctx = {};
ctx.succeed = function () {
  console.log(arguments, "yay");
};

ctx.fail = function (e) {
  console.log(e);
};

ctx.awsRequestId = "123405959";

m.handler(evt, ctx);