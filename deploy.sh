# shellcheck disable=SC2164
./scripts/build.sh
./scripts/bundle.sh

aws lambda publish-layer-version \
 --layer-name "Secrets-Lambda-Extension-Layer" \
 --zip-file  "fileb://bin/extension.zip"

aws lambda update-function-configuration \
 --function-name foobar \
 --layers $(aws lambda list-layer-versions --layer-name Secrets-Lambda-Extension-Layer  \
--max-items 1 --no-paginate --query 'LayerVersions[0].LayerVersionArn' \
--output text)

# aws logs describe-log-groups --query 'logGroups[?starts_with(logGroupName,`/aws/lambda/Test`)].logGroupName' \
# --output table | awk '{print $2}' | grep -v ^$ | while read x; do aws logs delete-log-group --log-group-name $x; done
