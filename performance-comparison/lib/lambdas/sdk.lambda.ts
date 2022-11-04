import { SecretsManager } from "aws-sdk";
const secretsManager = new SecretsManager();
let coldStart = true;
export const handler = async () => {
  //load secret from secret manager
  const secretName = process.env.SECRET_NAME!;
  const secret = await secretsManager
    .getSecretValue({ SecretId: secretName })
    .promise();

  console.log("secret", secret);
  const resp = {
    coldStart,
    randomId: process.env["RANDOM_UUID"],
    secret: secret.SecretString,
  };
  coldStart = false;
  return resp;
};
