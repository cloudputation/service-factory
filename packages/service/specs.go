package service


type RequestBody struct {
	ServiceSpecs
	TerraformDir string `json:"terraform_dir"`
}

type ServiceSpecs struct {
    Repo      Repo      `hcl:"repo,block"`
    Scheduler Scheduler `hcl:"scheduler,block"`
}

type Repo struct {
		Provider				string `hcl:"provider_url"`
    TemplateName    string `hcl:"template"`
    NamespaceID     string `hcl:"namespace_id"`
    Token           string `hcl:"token"`
    RunnerID        string `hcl:"runner_id"`
    RegistryToken   string `hcl:"registry_token"`
    RepositoryOwner string `hcl:"repository_owner"`
}

type Scheduler struct {
    Nomad Nomad `hcl:"nomad,block"`
}

type Nomad struct {
    ServiceName         string   	`hcl:"name"`
    ServiceGroup        string   	`hcl:"group"`
    ServicePort         string 		`hcl:"port"`
    ServiceTags         []string	`hcl:"tags"`
    ServiceType         string		`hcl:"type"`
    ServiceTargetClient string		`hcl:"target_client"`
		ServerAddress				string		`hcl:"server_address"`
}
