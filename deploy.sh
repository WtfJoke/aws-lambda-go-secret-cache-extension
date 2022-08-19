# shellcheck disable=SC2164
GOOS=linux GOARCH=amd64 go build -o bin/extensions/go-example-extension main.go
chmod +x bin/extensions/go-example-extension
cd bin
zip -r extension.zip extensions/

aws lambda publish-layer-version \
 --layer-name "Secrets-Lambda-Extension-Layer" \
 --zip-file  "fileb://extension.zip"

aws lambda update-function-configuration \
 --function-name foobar \
 --layers $(aws lambda list-layer-versions --layer-name Secrets-Lambda-Extension-Layer  \
--max-items 1 --no-paginate --query 'LayerVersions[0].LayerVersionArn' \
--output text)

# aws logs describe-log-groups --query 'logGroups[?starts_with(logGroupName,`/aws/lambda/Test`)].logGroupName' \
# --output table | awk '{print $2}' | grep -v ^$ | while read x; do aws logs delete-log-group --log-group-name $x; done
