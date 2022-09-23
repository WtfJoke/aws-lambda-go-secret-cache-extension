package extension

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
)

func TestReadOneSecretIdFromEnv(t *testing.T) {
	os.Setenv("SECRET_1", "secret-id")
	envVars := readSecretIdsFromEnvironmentWhenStartsWithSecret()

	expected := []string{"secret-id"}
	assert.Equal(t, envVars, expected, "should be equal")

}
func TestReadMultipleSecretIdsFromEnv(t *testing.T) {
	os.Setenv("SECRET_1", "secret-id-1")
	os.Setenv("SECRET_2", "secret-id-2")
	envVars := readSecretIdsFromEnvironmentWhenStartsWithSecret()

	expected := []string{"secret-id-1", "secret-id-2"}
	assert.Equal(t, envVars, expected, "should be equal")
}

func TestNoSecretsInEnv(t *testing.T) {
	os.Clearenv()
	envVars := readSecretIdsFromEnvironmentWhenStartsWithSecret()

	// empty array
	assert.Equal(t, 0, len(envVars), "should not contain any envvars")
}

type mockGetSecretValueAPI func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)

func (m mockGetSecretValueAPI) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m(ctx, params, optFns...)
}

type mockClientGetSecretValueType func(t *testing.T, secretString string) SecretsmanagerGetSecretValueApi

func mockClientGetSecretValue(t *testing.T, secretString string) SecretsmanagerGetSecretValueApi {
	return mockGetSecretValueAPI(func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
		t.Helper()

		return &secretsmanager.GetSecretValueOutput{
			SecretString: aws.String(secretString),
		}, nil
	})
}

func TestGetSecretValue(t *testing.T) {
	testCases := []struct {
		client   mockClientGetSecretValueType
		secretId string
		expect   string
	}{
		{
			client:   mockClientGetSecretValue,
			secretId: "secret-id",
			expect:   "secret-string",
		},
	}

	for i, testCase := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			secret := getSecretValue(testCase.client(t, testCase.expect), testCase.secretId)

			assert.Equal(t, testCase.expect, secret.SecretString)
			assert.Equal(t, testCase.secretId, secret.SecretId)

		})
	}

}

func TestGetSecretValuesFromListOfSecretIds(t *testing.T) {

	mockSecretValueGetter := func(client SecretsmanagerGetSecretValueApi, secretId string) Secret {
		return Secret{secretId, secretId + "-value"}
	}

	secret := getSecretValuesFromListOfSecretIds([]string{"secret-id-2", "secret-id-1"}, mockSecretValueGetter)

	assert.Equal(t, []Secret{{"secret-id-1", "secret-id-1-value"}, {"secret-id-2", "secret-id-2-value"}}, secret)
}

func TestWriteSecrets(t *testing.T) {
	secret1 := Secret{"secret-id-1", "secret-value-1"}
	secret2 := Secret{"secret-id-2", "anyOtherSecretValue"}

	writeSecrets([]Secret{secret1, secret2})

	secretFilePath := path.Join(os.TempDir(), "secrets.json")
	content, err := os.ReadFile(secretFilePath)
	if err != nil {
		assert.Fail(t, "Error while reading secrets.json", err.Error())
	}
	payload := []Secret{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		assert.Fail(t, "Error while unmarshalling secrets.json", err.Error())
	}

	assert.Equal(t, 2, len(payload), "have written the wrong amount of secrets")
	assert.Equal(t, secret1, payload[0])
	assert.Equal(t, secret2, payload[1])
}
