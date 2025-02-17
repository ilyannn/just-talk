package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/charmbracelet/glamour"
	"os"
	"os/exec"
	"strings"
)

func runLoop(client *azopenai.Client, deploymentId string, recipes []Recipe) {
	var tools []azopenai.ChatCompletionsToolDefinitionClassification
	var recipesByTool = make(map[string]*Recipe)

	for _, recipe := range recipes {
		requiredArgs := make([]string, 0)
		properties := make(map[string]any)

		for _, arg := range recipe.Arguments {
			if !arg.Optional {
				requiredArgs = append(requiredArgs, arg.Name)
			}

			var description = ""
			if arg.Default != "" {
				description = fmt.Sprintf("The %s, defaulting to %s.", arg.Name, arg.Default)
			} else {
				description = fmt.Sprintf("The %s.", arg.Name)
			}

			if arg.Variadic {
				description = description + " Variadic; contains the rest of the arguments."
			}

			properties[arg.Name] = map[string]any{
				"type":        "string",
				"description": description,
			}
		}

		jsonBytes, err := json.Marshal(map[string]any{
			"required":   requiredArgs,
			"type":       "object",
			"properties": properties,
		})

		if err != nil {
			panic(err)
		}

		name := toValidName(recipe.Name)
		recipesByTool[name] = &recipe

		funcDef := &azopenai.ChatCompletionsFunctionToolDefinitionFunction{
			Name:        to.Ptr(name),
			Description: to.Ptr(recipe.Description),
			Parameters:  jsonBytes,
		}

		tools = append(tools, &azopenai.ChatCompletionsFunctionToolDefinition{
			Function: funcDef,
		})
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(PromptStyle.Render("*** What do you want to do (empty to quit)? "))

		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput == "" {
			break
		}

		log.Debugf("Starting LLM request...")
		resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
			DeploymentName: &deploymentId,
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestUserMessage{
					Content: azopenai.NewChatRequestUserMessageContent(userInput),
				},
			},
			Tools: tools,
		}, nil)

		if err != nil {
			log.With("err", err).Errorf("Error when accessing the LLM")
			continue
		} else {
			log.With("model", *resp.Model, "input_tokens", *resp.Usage.PromptTokens, "output_tokens", *resp.Usage.CompletionTokens).Debugf("LLM request done")
		}

		content := resp.Choices[0].Message.Content
		if content != nil && *content != "" {
			out, err := glamour.Render(*content, "auto")
			if err != nil {
				log.With("err", err).Errorf("Failed to render markdown")
			}
			fmt.Print(out)
		}

		toolCalls := resp.Choices[0].Message.ToolCalls
		for _, toolCall := range toolCalls {
			function := toolCall.(*azopenai.ChatCompletionsFunctionToolCall).Function

			var parameters map[string]string
			err = json.Unmarshal([]byte(*function.Arguments), &parameters)

			if err != nil {
				log.With("err", err).Errorf("Error when unpacking tool arguments")
				continue
			}

			recipe := recipesByTool[*function.Name]
			if recipe == nil {
				log.With("tool", *function.Name).Errorf("Unknown tool name")
				continue
			}

			var argLine = []string{recipe.Name}

			for _, arg := range recipe.Arguments {
				if arg.Variadic {
					split := strings.Split(parameters[arg.Name], " ")
					argLine = append(argLine, split...)
				} else {
					argLine = append(argLine, parameters[arg.Name])
				}
			}

			formattedCommand := strings.Join([]string{"just", strings.Join(argLine, " ")}, " ")

			log.With("command", formattedCommand).Infof("Running")
			cmd := exec.Command("just", argLine...)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				log.With("err", err).Errorf("Failed to run")
			}
		}
	}

	fmt.Println("Goodbye!")
}
