# Welcome to your CDK TypeScript project

This is a blank project for CDK development with TypeScript.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

- `npm run build` compile typescript to js
- `npm run watch` watch for changes and compile
- `npm run test` perform the jest unit tests
- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template

```
fields @timestamp, @message
| filter @type="REPORT" and ispresent(@initDuration)
| stats count() as coldStarts, avg(@initDuration), min(@initDuration), max(@initDuration) ,avg(@billedDuration), min(@billedDuration), max(@billedDuration) by  bin (10m) as foo
| sort foo desc
| limit 2000



```

**CloudWatch Logs Insights**  
region: eu-central-1  
log-group-names: /aws/lambda/PerformanceComparisonStack-lambda8B5974B5-i1KtsBfGai8w  
start-time: -3600s  
end-time: 0s  
query-string:

```
fields @timestamp, @message
| filter @type="REPORT" and ispresent(@initDuration)
| stats count() as coldStarts, avg(@initDuration), min(@initDuration), max(@initDuration) ,avg(@billedDuration), min(@billedDuration), max(@billedDuration) by  bin (10m) as foo
| sort foo desc
| limit 2000
```

### Lambda extension

---

| foo                     | coldStarts | avg(@initDuration) | min(@initDuration) | max(@initDuration) | avg(@billedDuration) | min(@billedDuration) | max(@billedDuration) |
| ----------------------- | ---------- | ------------------ | ------------------ | ------------------ | -------------------- | -------------------- | -------------------- |
| 2022-11-04 16:10:00.000 | 500        | 275.908            | 225.85             | 753.63             | 271.044              | 212                  | 356                  |

---

### Lambda SDK

---

| foo                     | coldStarts | avg(@initDuration) | min(@initDuration) | max(@initDuration) | avg(@billedDuration) | min(@billedDuration) | max(@billedDuration) |
| ----------------------- | ---------- | ------------------ | ------------------ | ------------------ | -------------------- | -------------------- | -------------------- |
| 2022-11-04 16:30:00.000 | 500        | 455.2082           | 363.83             | 992.14             | 623.664              | 562                  | 772                  |

---
