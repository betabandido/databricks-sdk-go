package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/betabandido/databricks-sdk-go/api/groups"
	"github.com/betabandido/databricks-sdk-go/client"
	"github.com/betabandido/databricks-sdk-go/models"
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

	endpoint := groups.Endpoint{
		Client: cl,
	}

	groups := getGroups(endpoint)
	for _, g := range groups {
		fmt.Printf("group: %s\n", g)
		printMembers(endpoint, g)
	}
}

func getGroups(endpoint groups.Endpoint) []string {
	resp, err := endpoint.List()
	if err != nil {
		panic(err)
	}

	return resp.GroupNames
}

func printMembers(endpoint groups.Endpoint, group string) {
	resp, err := endpoint.ListMembers(&models.GroupsListMembersRequest{
		GroupName: group,
	})
	if err != nil {
		panic(err)
	}

	for _, m := range resp.Members {
		fmt.Printf("  %s\n", m)
	}
}

type secrets struct {
	Domain string `json:"domain"`
	Token  string `json:"token"`
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
