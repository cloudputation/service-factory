package service


type RequestBody struct {
	ServiceSpecs
}

type ServiceSpecs struct {
	Service Service `hcl:"service,block"`
}

type Service struct {
  Name	string		`hcl:"name"`
  Group	string		`hcl:"group"`
  Port	string		`hcl:"port"`
  Tags	[]string	`hcl:"tags"`
  Template  			`hcl:"template,block"`
  Repo      			`hcl:"repo,block"`
  Network      		`hcl:"network,block"`
}

type Template struct {
	TemplateURL  string `hcl:"template_url"`
	TemplateName string `hcl:"template"`
}

type Repo struct {
	ProviderURL     string `hcl:"provider_url"`
	NamespaceID     string `hcl:"namespace_id"`
	Token           string `hcl:"token"`
	RunnerID        string `hcl:"runner_id"`
	RegistryToken   string `hcl:"registry_token"`
	RepositoryOwner string `hcl:"repository_owner"`
}

type Network struct {
	AuthoritativeServer	string `hcl:"authoritative_server"`
	ClientHostname			string `hcl:"client_hostname"`
}
