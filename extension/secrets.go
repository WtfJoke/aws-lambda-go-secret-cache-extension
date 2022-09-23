package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	SecretId, SecretString string
}

func LoadSecrets() {
	// log current timestamp
	fmt.Println("Loading secrets...")
	fmt.Println("Current timestamp:", time.Now().UnixMilli())

	secretIds := readSecretIdsFromEnvironmentWhenStartsWithSecret()
	secrets := getSecretValuesFromListOfSecretIds(secretIds, getSecretValue)
	writeSecrets(secrets)

	fmt.Println("Finished loading secrets on", time.Now().UnixMilli())
}

func readSecretIdsFromEnvironmentWhenStartsWithSecret() []string {
	secretIds := []string{}
	envVars := os.Environ()
	for _, secret := range envVars {
		if strings.HasPrefix(secret, "SECRET_") {
			pair := strings.SplitN(secret, "=", 2)
			secretIds = append(secretIds, pair[1])
		}
	}
	return secretIds
}

func getSecretValuesFromListOfSecretIds(secretIds []string, getSecretValue SecretValueGetter) []Secret {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	checkError(err)
	client := secretsmanager.NewFromConfig(cfg)
	secretsCount := len(secretIds)
	secretsChannel := make(chan Secret, secretsCount)
	wg := sync.WaitGroup{}
	startTime := time.Now()
	for _, secretId := range secretIds {
		wg.Add(1)
		go func(secretId string) {
			fetchStartTime := time.Now()
			secretsChannel <- getSecretValue(client, secretId)
			log.Print("Fetched secret '", secretId, "' in ", time.Since(fetchStartTime))
			wg.Done()
		}(secretId)
	}
	wg.Wait()
	log.Print("All ", secretsCount, " secrets fetched in ", time.Since(startTime))
	close(secretsChannel)
	loadedSecrets := collectSecrets(secretsChannel)
	return loadedSecrets
}

type SecretsmanagerGetSecretValueApi interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SecretValueGetter func(client SecretsmanagerGetSecretValueApi, secretId string) Secret

func getSecretValue(client SecretsmanagerGetSecretValueApi, secretId string) Secret {
	log.Print("Fetch secret: ", secretId, " on ", time.Now().UnixMilli())
	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	checkError(err)

	return Secret{secretId, aws.ToString(output.SecretString)}
}

func collectSecrets(secretsChannel chan Secret) []Secret {
	var secrets []Secret
	for secret := range secretsChannel {
		secrets = append(secrets, secret)
	}
	return secrets
}

func writeSecrets(secrets []Secret) {
	marshaleldSecret, marshaledErr := json.Marshal(secrets)
	checkError(marshaledErr)
	secretFilePath := path.Join(os.TempDir(), "secrets.json")
	write_err := os.WriteFile(secretFilePath, marshaleldSecret, 0644)
	log.Print("Secrets written to ", secretFilePath)
	checkError(write_err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
