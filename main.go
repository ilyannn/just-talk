package main

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	_log "github.com/charmbracelet/log"
	"os"
	"os/exec"
	"strings"
)

import _ "github.com/joho/godotenv/autoload"

func prepareRecipes() []Recipe {
	cmd := exec.Command("just", "--list", "--unsorted")
	out, err := cmd.Output()

	if err != nil {
		var logWithContext *_log.Logger
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.Stderr != nil {
			// convert the stderr bytes to a string
			var stderr string
			stderr = strings.TrimSpace(string(exitErr.Stderr))
			logWithContext = log.With("stderr", stderr)
		} else {
			logWithContext = log.With("err", err)
		}
		logWithContext.Fatalf("Failed to get %s recipes", JustBold)
	}

	recipes := parseJustOutput(string(out))

	var recipesNames []string
	for _, recipe := range recipes {
		recipesNames = append(recipesNames, recipe.Name)
	}

	log.With("names", recipesNames).Infof("Found %d %s recipes", len(recipesNames),
		JustBold)

	return recipes
}

func prepareClient() *azopenai.Client {
	openAIAPIKey := os.Getenv("JUST_TALK_OPENAI_API_KEY")
	openAIDeploymentEndpoint := os.Getenv("JUST_TALK_AZURE_ENDPOINT")

	if openAIAPIKey == "" {
		println("Please set the environment variables JUST_TALK_OPENAI_API_KEY (and possibly JUST_TALK_AZURE_ENDPOINT and/or JUST_TALK_OPENAI_MODEL_ID); you can use .env file.")
		log.Fatalf("JUST_TALK_OPENAI_API_KEY is not set")
	}

	keyCredential := azcore.NewKeyCredential(openAIAPIKey)

	var client *azopenai.Client
	var err error

	if openAIDeploymentEndpoint != "" {
		// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
		client, err = azopenai.NewClientWithKeyCredential(openAIDeploymentEndpoint, keyCredential, nil)
		log.With("endpoint", openAIDeploymentEndpoint).Debugf("Connecting to Azure")
	} else {
		// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
		endpoint := "https://api.openai.com/v1"
		client, err = azopenai.NewClientForOpenAI(endpoint, keyCredential, nil)
		log.With("endpoint", endpoint).Debugf("Connecting to OpenAI")
	}

	if err != nil {
		log.With("err", err).Fatalf("Failed to create a new client")
	}

	return client
}

func main() {
	recipes := prepareRecipes()
	client := prepareClient()

	deploymentId, present := os.LookupEnv("JUST_TALK_AZURE_DEPLOYMENT_ID")

	if !present {
		log.Warn("JUST_TALK_AZURE_DEPLOYMENT_ID is not set.")
	}

	runLoop(client, deploymentId, recipes)
}
