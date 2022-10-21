# AWS Lambda Extension secret cache

This AWS lambda extension fetches secrets from the aws secretsmanger during init phase of the lambda cold start and then exposes them on a local webserver, so that the lambda can access the secret value by making a http call to the local webserver.

**Note** This is a proof-of-concept implementation. Use this in production on your own risks.

## Compile package and dependencies

You can compile the lambda extension to both currently supported AWS lambda runtimes. The defualt build runtime is `amd64`.

To compile for the ARM runtime set an env variable `export ARCH=arm64`.

Building and saving package into a `bin/extensions` directory:

```bash
./scripts/build.sh
```

This creates a go binary for the target architecture.

## Layer Setup Process

The extensions .zip file should contain a root directory called `extensions/`, where the extension executables are located. In this project we must include the `secret-cache-extension` binary.

Creating zip package for the extension:

```bash
./scripts/bundle.sh
```

Ensure that you have aws-cli v2 for the commands below.
Publish a new layer using the `extension.zip`. The output of the following command should provide you a layer arn.

```bash
aws lambda publish-layer-version \
 --layer-name "secret-cache-extension" \
 --region <use your region> \
 --zip-file  "fileb://extension.zip"
```

Note the LayerVersionArn that is produced in the output.
eg. `"LayerVersionArn": "arn:aws:lambda:<region>:123456789012:layer:<layerName>:1"`

Add the newly created layer version to a Lambda function.

## Function Invocation and Extension Execution

When invoking the function, you should now see log messages from the example extension similar to the following:

```
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    EXTENSION Name: go-example-extension State: Ready Events: [INVOKE,SHUTDOWN]
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    START RequestId: 9ca08945-de9b-46ec-adc6-3fe9ef0d2e8d Version: $LATEST
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    [go-example-extension]  Registering...
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    [go-example-extension]  Register response: {
                "functionName": "my-function",
                "functionVersion": "$LATEST",
                "handler": "function.handler"
        }
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    [go-example-extension]  Waiting for event...
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    [go-example-extension]  Received event: {
                "eventType": "INVOKE",
                "deadlineMs": 1234567890123,
                "requestId": "9ca08945-de9b-46ec-adc6-3fe9ef0d2e8d",
                "invokedFunctionArn": "arn:aws:lambda:<region>:123456789012:function:my-function",
                "tracing": {
                        "type": "X-Amzn-Trace-Id",
                        "value": "XXXXXXXXXX"
                }
        }
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    [go-example-extension]  Waiting for event...
    ...
    ...
    Function logs...
    ...
    ...
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    END RequestId: 9ca08945-de9b-46ec-adc6-3fe9ef0d2e8d
    XXXX-XX-XXTXX:XX:XX.XXX-XX:XX    REPORT RequestId: 9ca08945-de9b-46ec-adc6-3fe9ef0d2e8d Duration: 3.78 ms	Billed Duration: 100 ms	Memory Size: 128 MB	Max Memory Used: 59 MB	Init Duration: 264.75 ms
```
