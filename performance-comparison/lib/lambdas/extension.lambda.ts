import { request } from "http";
let coldStart = true;
export const handler = async () => {
  const randomId = process.env["RANDOM_UUID"];
  const secretName = process.env.SECRET_NAME;

  if (!randomId) {
    console.log(process.env);
  }

  // make a http request as promise

  const options = {
    method: "GET",
    hostname: "localhost",
    port: "4000",
    path: `/secrets?name=${secretName}`,
  };

  const req = await httpRequest(options, null);
  console.log("req", req);
  const resp = { coldStart, randomId: process.env["RANDOM_UUID"], req };

  coldStart = false;
  return resp;
};

function httpRequest(params: any, postData: any) {
  return new Promise(function (resolve, reject) {
    var req = request(params, function (res) {
      // reject on bad status
      // @ts-ignore
      if (res.statusCode < 200 || res.statusCode >= 300) {
        return reject(new Error("statusCode=" + res.statusCode));
      }
      // cumulate data
      var body: any = [];
      res.on("data", function (chunk) {
        body.push(chunk);
      });
      // resolve on end
      res.on("end", function () {
        try {
          body = Buffer.concat(body).toString();
        } catch (e) {
          reject(e);
        }
        resolve(body);
      });
    });
    // reject on request error
    req.on("error", function (err) {
      // This is not a "Second reject", just a different sort of failure
      reject(err);
    });
    if (postData) {
      req.write(postData);
    }
    // IMPORTANT
    req.end();
  });
}
