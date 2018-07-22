package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/betabandido/databricks-sdk-go/api/clusters"
	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
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

	createCluster(endpoint, secrets.ClusterName)

	listClusters(endpoint)
}

func createCluster(endpoint clusters.Endpoint, clusterName string) {
	resp, err := endpoint.Create(&models.ClustersCreateRequest{
		ClusterName:  clusterName,
		SparkVersion: "4.2.x-scala2.11",
		NodeTypeId:   "Standard_D3_v2",
		NumWorkers:   1,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cluster %s created\n", resp.ClusterId)
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
	Domain      string `json:"domain"`
	Token       string `json:"token"`
	ClusterName string `json:"cluster_name"`
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
