import * as cdk from "aws-cdk-lib";
import { Duration } from "aws-cdk-lib";
import { AssetCode, LayerVersion } from "aws-cdk-lib/aws-lambda";
import { NodejsFunction } from "aws-cdk-lib/aws-lambda-nodejs";
import { Secret } from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";
import { resolve } from "path";
import { TestExecuter } from "./TestExecuter";

export class PerformanceComparisonStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const chunks: number = 5;
    const chunkSize: number = 100;
    const delayBetweenChunks = Duration.seconds(30);

    const secret = new Secret(this, "Secret");

    const extensionLambdaFunction = new NodejsFunction(this, "lambda", {
      entry: resolve(__dirname, "lambdas", "extension.lambda.ts"),
      environment: {
        SECRET_NAME: secret.secretName,
      },

      layers: [
        new LayerVersion(this, "extensionLayer", {
          code: new AssetCode(resolve(__dirname, "../../bin", "extension.zip")),
        }),
      ],
    });
    secret.grantRead(extensionLambdaFunction);

    const sdkLambdaFunction = new NodejsFunction(this, "sdkLambda", {
      entry: resolve(__dirname, "lambdas", "sdk.lambda.ts"),

      environment: {
        SECRET_NAME: secret.secretName,
      },
    });
    secret.grantRead(sdkLambdaFunction);

    new TestExecuter(this, "ExtensionTestExecuter", {
      chunks,
      chunkSize,
      delayBetweenChunks,
      lambdaFunction: extensionLambdaFunction,
    });

    new TestExecuter(this, "SdkTestExecuter", {
      chunks,
      chunkSize,
      delayBetweenChunks,
      lambdaFunction: sdkLambdaFunction,
    });
  }
}
