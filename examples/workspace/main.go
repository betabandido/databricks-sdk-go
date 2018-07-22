package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/betabandido/databricks-sdk-go/api/workspace"
	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
	"io/ioutil"
	"path"
)

func main() {
	flag.Parse() // required to suppress warnings from glog

	// TODO: use a better logging library
	//flag.Lookup("logtostderr").Value.Set("true")
	//flag.Lookup("stderrthreshold").Value.Set("INFO")

	secrets := loadSecrets()

	cl, err := client.NewClient(client.Options{
		Domain: &secrets.Domain, Token: &secrets.Token,
	})
	if err != nil {
		panic(err)
	}

	endpoint := workspace.Endpoint{
		Client: cl,
	}

	mkdirs(endpoint, secrets.DirPath)

	notebookPath := path.Join(secrets.DirPath, secrets.NotebookName)

	importNotebook(endpoint, notebookPath)
	printNotebookStatus(endpoint, notebookPath)

	listPath(endpoint, secrets.DirPath)

	deletePath(endpoint, notebookPath)
	deletePath(endpoint, secrets.DirPath)
}

func mkdirs(endpoint workspace.Endpoint, path string) {
	err := endpoint.Mkdirs(&models.WorkspaceMkdirsRequest{
		Path: path,
	})
	if err != nil {
		panic(err)
	}
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

func deletePath(endpoint workspace.Endpoint, path string) {
	err := endpoint.Delete(&models.WorkspaceDeleteRequest{
		Path: path,
	})
	if err != nil {
		panic(err)
	}
}

func listPath(endpoint workspace.Endpoint, path string) {
	resp, err := endpoint.List(&models.WorkspaceListRequest{
		Path: path,
	})
	if err != nil {
		panic(err)
	}

	for _, obj := range resp.Objects {
		fmt.Printf("path: %v, lang: %v, type: %v\n", obj.Path, *obj.Language, *obj.ObjectType)
	}
}

type secrets struct {
	Domain       string `json:"domain"`
	Token        string `json:"token"`
	NotebookName string `json:"notebook_name"`
	DirPath      string `json:"dir_path"`
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
