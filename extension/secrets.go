package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	SecertId, SecretString string
}

func LoadSecrets() {
	// log current timestamp
	fmt.Println("Loading secrets...")
	fmt.Println("Current timestamp:", time.Now().UnixMilli())

	secretIds := readSecretIdsFromEnvironmentWhenStartsWithSecret()
	secretsChannel := make(chan Secret, len(secretIds))
	go getSecretValuesFromListOfSecretIds(secretIds, secretsChannel)
	secrets := collectSecrets(secretsChannel)
	writeSecrets(secrets)

	fmt.Println("finished timestamp:", time.Now().UnixMilli())
	log.Print("finished loading secrets")
}

func readSecretIdsFromEnvironmentWhenStartsWithSecret() []string {
	var secretIds []string
	for _, secret := range os.Environ() {
		if strings.HasPrefix(secret, "SECRET_") {
			pair := strings.SplitN(secret, "=", 2)
			secretIds = append(secretIds, pair[1])
		}
	}
	return secretIds
}

func getSecretValuesFromListOfSecretIds(secretIds []string, secretsChannel chan Secret) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	checkError(err)
	client := secretsmanager.NewFromConfig(cfg)

	for _, secretId := range secretIds {
		log.Print("Fetch secret: "+secretId+" on ", time.Now().UnixMilli())
		secret := getSecretValue(client, secretId)
		secretsChannel <- secret
	}
	close(secretsChannel)
}

func getSecretValue(client *secretsmanager.Client, secretId string) Secret {
	log.Print("Fetch secret: ", secretId, " on ", time.Now().UnixMilli())
	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	log.Print("After fetching secret ", time.Now().UnixMilli())
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
	checkError(write_err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODOS:
// create channel which accepts secretIds
// foreach over secret ids resolving/fetching secrets
// output channel writes secrets to file
