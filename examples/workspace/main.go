package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/betabandido/databricks-sdk-go/api/workspace"
	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
	"io/ioutil"
)

func main() {
	secrets := loadSecrets()

	cl, err := client.NewClient(secrets.Domain, secrets.Token)
	if err != nil {
		panic(err)
	}

	endpoint := workspace.Endpoint{
		Client: cl,
	}

	importNotebook(endpoint, secrets.NotebookPath)

	printNotebookStatus(endpoint, secrets.NotebookPath)

	deleteNotebook(endpoint, secrets.NotebookPath)
}

func importNotebook(endpoint workspace.Endpoint, path string) {
	language := models.PYTHON
	content := base64.StdEncoding.EncodeToString([]byte("print('hello world')"))

	err := endpoint.Import(&models.WorkspaceImportRequest{
		Path:     path,
		Language: &language,
		Content:  content,
	})
	if err != nil {
		panic(err)
	}
}

func printNotebookStatus(endpoint workspace.Endpoint, path string) {
	resp, err := endpoint.GetStatus(&models.WorkspaceGetStatusRequest{
		Path: path,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Language %s\n", *resp.Language)
}

func deleteNotebook(endpoint workspace.Endpoint, path string) {
	err := endpoint.Delete(&models.WorkspaceDeleteRequest{
		Path: path,
	})
	if err != nil {
		panic(err)
	}
}

type secrets struct {
	Domain       string `json:"domain"`
	Token        string `json:"token"`
	NotebookPath string `json:"notebook_path"`
}

func loadSecrets() *secrets {
	content, err := ioutil.ReadFile("../secrets.json")
	if err != nil {
		panic(err)
	}

	var sc secrets
	err = json.Unmarshal(content, &sc)
	if err != nil {
		panic(err)
	}

	return &sc
}
