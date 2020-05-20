package main

import (
	"log"
	"os"
	"path"
	"sort"

	"github.com/aovlllo/lazylab/internal/gitlab"
	"github.com/aovlllo/lazylab/internal/ui"
)

func die(msg string, e error) {
	if e != nil {
		log.Fatalf("%s: %v", msg, e)
	}
}

func loadConfig() gitlab.ClientConfig {
	userHome, err := os.UserHomeDir()
	die("could not obtain user home directory", err)

	configPath := path.Join(userHome, ".lazylab")

	configFile, err := os.Open(configPath)
	die("could not read config file", err)

	config, err := gitlab.LoadConfig(configFile)
	die("could not parse config", err)

	return config
}

type MRAdder interface {
	AddMR(m *gitlab.MergeRequest)
}

func getRefreshFn(client *gitlab.Client, opts *gitlab.ListMergeRequestsOptions, m MRAdder) func() {
	return func() {
		mrs, _, _ := client.MergeRequests.ListMergeRequests(opts)
		sort.Slice(mrs, func(i, j int) bool {
			return mrs[i].UpdatedAt.After(*mrs[j].UpdatedAt)
		})
		for _, mr := range mrs {
			m.AddMR(mr)
		}

	}
}

func draw(app *ui.Application) {
	config := loadConfig()
	client, err := gitlab.NewClient(config)
	die("could not create client", err)
	state, scope := "opened", "all"

	mlAuthor := ui.NewMRList(app)
	refreshAuthor := getRefreshFn(client, &gitlab.ListMergeRequestsOptions{AuthorID: &config.UserID, State: &state, Scope: &scope}, mlAuthor)
	mlAuthor.SetRefresh(refreshAuthor)
	app.AddSection("Author", mlAuthor)
	//refreshAuthor()

	mlApprover := ui.NewMRList(app)
	refreshApprover := getRefreshFn(client, &gitlab.ListMergeRequestsOptions{ApproverIDs: []int{config.UserID}, State: &state, Scope: &scope}, mlApprover)
	mlApprover.SetRefresh(refreshApprover)
	app.AddSection("Approver", mlApprover)
	//refreshApprover()

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	draw(ui.NewApplication())
}
