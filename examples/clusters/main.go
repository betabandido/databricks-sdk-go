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

const (
	SparkVersion = "4.2.x-scala2.11"
	NodeTypeId   = "Standard_D3_v2"
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

	clusterId := createCluster(endpoint, secrets.ClusterName)
	printClusterInfo(endpoint, clusterId)

	renameCluster(endpoint, clusterId, secrets.ClusterName+"-renamed")
	printClusterInfo(endpoint, clusterId)

	listClusters(endpoint)

	permanentlyDeleteCluster(endpoint, clusterId)
}

func createCluster(endpoint clusters.Endpoint, clusterName string) string {
	fmt.Println("Creating cluster")

	resp, err := endpoint.CreateSync(&models.ClustersCreateRequest{
		ClusterName:  clusterName,
		SparkVersion: SparkVersion,
		NodeTypeId:   NodeTypeId,
		NumWorkers:   1,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cluster %s created\n", resp.ClusterId)

	return resp.ClusterId
}

func renameCluster(endpoint clusters.Endpoint, clusterId string, name string) {
	fmt.Println("Renaming cluster")

	err := endpoint.EditSync(&models.ClustersEditRequest{
		ClusterId:    clusterId,
		ClusterName:  name,
		SparkVersion: SparkVersion,
		NodeTypeId:   NodeTypeId,
		NumWorkers:   1,
	})
	if err != nil {
		panic(err)
	}
}

func printClusterInfo(endpoint clusters.Endpoint, clusterId string) {
	resp, err := endpoint.Get(&models.ClustersGetRequest{
		ClusterId: clusterId,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("id: %s, name: %s, creator: %s, spark-version: %s, state: %s\n",
		resp.ClusterId, resp.ClusterName, resp.CreatorUserName, resp.SparkVersion, *resp.State)
}

func listClusters(endpoint clusters.Endpoint) {
	resp, err := endpoint.List()
	if err != nil {
		panic(err)
	}

	for _, c := range resp.Clusters {
		fmt.Printf("id: %s, name: %s, creator: %s, spark-version: %s, state: %s\n",
			c.ClusterId, c.ClusterName, c.CreatorUserName, c.SparkVersion, *c.State)
	}
}

func permanentlyDeleteCluster(endpoint clusters.Endpoint, clusterId string) {
	err := endpoint.PermanentDelete(&models.ClustersPermanentDeleteRequest{
		ClusterId: clusterId,
	})
	if err != nil {
		panic(err)
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
