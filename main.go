package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

import _ "github.com/joho/godotenv/autoload"

func prepareRecipes() []Recipe {
	cmd := exec.Command("just", "--list", "--unsorted")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run just: %v", err)
	}

	recipes := parseJustOutput(string(out))

	var recipesNames []string
	for _, recipe := range recipes {
		recipesNames = append(recipesNames, recipe.Name)
	}

	log.Printf("Known just recipes are: %s.", strings.Join(recipesNames, ", "))

	return recipes
}

func prepareClient() *azopenai.Client {
	openAIAPIKey := os.Getenv("JUST_TALK_OPENAI_API_KEY")
	openAIDeploymentEndpoint := os.Getenv("JUST_TALK_AZURE_ENDPOINT")

	if openAIAPIKey == "" {
		log.Fatal("Please set the environment variable JUST_TALK_OPENAI_API_KEY (and possibly JUST_TALK_AZURE_ENDPOINT and/or JUST_TALK_OPENAI_MODEL_ID); you can use .env file.")
	}

	keyCredential := azcore.NewKeyCredential(openAIAPIKey)

	var client *azopenai.Client
	var err error

	if openAIDeploymentEndpoint != "" {
		// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
		client, err = azopenai.NewClientWithKeyCredential(openAIDeploymentEndpoint, keyCredential, nil)
	} else {
		// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
		client, err = azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)
	}

	if err != nil {
		log.Fatalf("Failed to create a new client: %v", err)
	}

	return client
}

func main() {
	recipes := prepareRecipes()
	client := prepareClient()

	deploymentId, present := os.LookupEnv("JUST_TALK_AZURE_DEPLOYMENT_ID")

	if !present {
	}

	runLoop(client, deploymentId, recipes)
}
