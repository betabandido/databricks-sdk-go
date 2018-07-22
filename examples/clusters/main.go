package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/betabandido/databricks-sdk-go/api/clusters"
	"github.com/betabandido/databricks-sdk-go/client"
	"io/ioutil"
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

	endpoint := clusters.Endpoint{
		Client: cl,
	}

	listClusters(endpoint)
}

func listClusters(endpoint clusters.Endpoint) {
	resp, err := endpoint.List()
	if err != nil {
		panic(err)
	}

	for _, c := range resp.Clusters {
		fmt.Printf("id: %s, name: %s, creator: %s, spark-version: %s\n",
			c.ClusterId, c.ClusterName, c.CreatorUserName, c.SparkVersion)
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
