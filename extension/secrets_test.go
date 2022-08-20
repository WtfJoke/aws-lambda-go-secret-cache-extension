package extension

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
)

func TestReadOneSecretIdFromEnv(t *testing.T) {
  os.Setenv("SECRET_1", "secret-id")
  envVars:= readSecretIdsFromEnvironmentWhenStartsWithSecret()

  expected := []string{"secret-id"}
  assert.Equal(t, envVars,expected, "should be equal")

}
func TestReadMultipleSecretIdsFromEnv(t *testing.T) {
  os.Setenv("SECRET_1", "secret-id-1")
  os.Setenv("SECRET_2", "secret-id-2")
  envVars:= readSecretIdsFromEnvironmentWhenStartsWithSecret()

  expected := []string{"secret-id-1", "secret-id-2"}
  assert.Equal(t, envVars,expected, "should be equal")
}

func TestNoSecretsInEnv(t *testing.T) {
 
  envVars:= readSecretIdsFromEnvironmentWhenStartsWithSecret()

  // empty array

  expected := []string{}
  assert.Equal(t, envVars,expected, "should be equal")
}


type mockGetSecretValueAPI func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)

func (m mockGetSecretValueAPI) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m(ctx, params, optFns...)
}
func TestGetSecretValue(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) SecretsmanagerGetSecretValueApi
		secretId string
		expect string
	}{
		{
			client: func(t *testing.T) SecretsmanagerGetSecretValueApi {
				return mockGetSecretValueAPI(func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
					t.Helper()
		
					return &secretsmanager.GetSecretValueOutput{
						SecretString: aws.String("secret-string"),
					}, nil
				})
			},
			secretId: "secret-id",
			expect: "secret-string",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			secret := getSecretValue(tt.client(t), tt.secretId)

			assert.Equal(t, tt.expect, secret.SecretString)
			assert.Equal(t, tt.secretId, secret.SecertId)
			
		})
	}
}
func TestGetSecretValuesFromListOfSecretIds(t *testing.T) {

	mockSecretValueGetter := func (client SecretsmanagerGetSecretValueApi, secretId string) Secret {
		return Secret{secretId, secretId + "-value"}
	}


	secret := getSecretValuesFromListOfSecretIds( []string{"secret-id-2", "secret-id-1"}, mockSecretValueGetter)

	assert.Equal(t,  []Secret{{"secret-id-1", "secret-id-1-value"},{"secret-id-2", "secret-id-2-value"}} , secret)
}