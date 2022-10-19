/**
 * A Lambda function that returns a static string
 */

const fs = require("fs");
const secrets = JSON.parse(fs.readFileSync("/tmp/secrets.json").toString());

exports.helloFromLambdaHandler = async () => {
  // If you change this message, you will need to change hello-from-lambda.test.js
  const message = "Hello from Lambda!";

  // All log statements are written to CloudWatch
  console.info(`${message}`);
  console.info(JSON.stringify({ secrets }, undefined, 2));

  return message;
};
