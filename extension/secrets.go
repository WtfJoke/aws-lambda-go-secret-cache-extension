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
	SecertId, SecretString string
}

func LoadSecrets() {
	// log current timestamp
	fmt.Println("Loading secrets...")
	fmt.Println("Current timestamp:", time.Now().UnixMilli())

	secretIds := readSecretIdsFromEnvironmentWhenStartsWithSecret()
	secrets :=  getSecretValuesFromListOfSecretIds(secretIds)
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

func getSecretValuesFromListOfSecretIds(secretIds []string) []Secret {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	checkError(err)
	client := secretsmanager.NewFromConfig(cfg)
	secretsChannel := make(chan Secret, len(secretIds))
	wg := sync.WaitGroup{}
	startTime := time.Now()
	for _, secretId := range secretIds {
		wg.Add(1)
		 go func(secretId string) { 
			fetchStartTime := time.Now()
			log.Print("Fetch secret: "+secretId+" started on ", time.Now().UnixMilli(), time.Since(startTime))
			secretsChannel <- getSecretValue(client, secretId) 
			log.Print("Fetch secret: "+secretId+" done on ", time.Now().UnixMilli(), time.Since(startTime), time.Since(fetchStartTime))
			wg.Done()
			}(secretId)
		
	}
	wg.Wait()
	log.Print("All secrets fetched on ", time.Now().UnixMilli(), time.Since(startTime))
	close(secretsChannel)
	loadedSecrets := collectSecrets(secretsChannel)
	return loadedSecrets
}

func getSecretValue(client *secretsmanager.Client, secretId string) Secret {
	fetchStartTime := time.Now()

	log.Print("Fetch secret: ", secretId, " on ", time.Now().UnixMilli())
	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	})
	log.Print("After fetching secret ", time.Now().UnixMilli(), time.Since(fetchStartTime))
	checkError(err)

	return Secret{secretId, aws.ToString(output.SecretString)}
}

func collectSecrets(secretsChannel chan Secret) []Secret {
	var secrets []Secret
	for secret := range secretsChannel {
		secrets = append(secrets, secret)
	}
	log.Print(len(secrets), " secrets loaded")
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

// TODOS:
// create channel which accepts secretIds
// foreach over secret ids resolving/fetching secrets
// output channel writes secrets to file
