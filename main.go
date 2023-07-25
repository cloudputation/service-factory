package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

type EnvVars struct {
	CookieCutterGitUrl      string `json:"repo_url"`
	GitToken                string `json:"git_token"`
	RunnerID                string `json:"runner_id"`

  NomadHost               string `json:"nomad_host"`
  ServiceName             string `json:"service_name"`
  ServiceGroup            string `json:"service_group"`
  ServicePort             string `json:"service_port"`
  ServiceTag              string `json:"service_tag"`
  ServiceType             string `json:"service_type"`
  TargetNomadClient       string `json:"target_nomad_client"`
  ServiceRepositoryOwner  string `json:"service_repository_owner"`
  ServiceRegistryToken    string `json:"service_registry_token"`
}

type RequestBody struct {
	EnvVars
	TerraformDir string `json:"terraform_dir"`
}

type ResponseBody struct {
	Message  string `json:"message"`
}

func main() {

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var body RequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusBadRequest)
			return
		}
    var service_dir = body.ServiceName

  	cookie_cutter_files := map[string]string{
  		"cookie-cutter-api/service.py.tpl": service_dir+"/service.py",
  		"cookie-cutter-api/Dockerfile.tpl": service_dir+"/Dockerfile",
  		"cookie-cutter-api/nomad/nomad.tf.tpl": service_dir+"/nomad/nomad.tf",
  		"cookie-cutter-api/nomad/jobspec/service.nomad.tpl": service_dir+"/nomad/jobspec/service.nomad",
      "cookie-cutter-api/.gitlab-ci.yml.tpl": service_dir+"/.gitlab-ci.yml",
      "cookie-cutter-api/API_VERSION": service_dir+"/API_VERSION",
      "cookie-cutter-api/requirements.txt": service_dir+"/requirements.txt",
  	}

		err = createServiceDir(body.EnvVars)
		if err != nil {
			http.Error(w, "Failed to create service directory:"+err.Error(),
				http.StatusInternalServerError)
			return
		}

		err = gitClone(body.CookieCutterGitUrl, body.EnvVars)
		if err != nil {
			http.Error(w, "Failed to clone the repository: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

    for tplFile, outputFile := range cookie_cutter_files {
			err = renderTemplate(tplFile, outputFile, body.EnvVars)
			if err != nil {
				http.Error(w, "Failed to render the template: "+err.Error(),
					http.StatusInternalServerError)
				return
			}
		}

		err = runTerraform(body.TerraformDir, body.EnvVars)
		if err != nil {
			http.Error(w, "Failed to run terraform: "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		err = cleanServiceDir(body.CookieCutterGitUrl, body.EnvVars)
		if err != nil {
			http.Error(w, "Failed to clean service directory:"+err.Error(),
				http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ResponseBody{Message: "Operation completed successfully"})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":48840", nil)
}

func createServiceDir(envVars EnvVars) error {
	cmd := exec.Command("mkdir", "-p", envVars.ServiceName+"/nomad/jobspec")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}


func cleanServiceDir(CookieCutterGitUrl string,envVars EnvVars) error {
	cmd := exec.Command("rm", "-fr", envVars.ServiceName, CookieCutterGitUrl)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}


func gitClone(CookieCutterGitUrl string, envVars EnvVars) error {
	cmd := exec.Command("git", "clone", CookieCutterGitUrl)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func renderTemplate(templateFilename, outputFilename string, envVars EnvVars) error {
	tmpl, err := template.ParseFiles(templateFilename)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = tmpl.Execute(outputFile, envVars)
	if err != nil {
		return err
	}
	return nil
}

func runTerraform(terraformDir string, envVars EnvVars) error {
  // var filesJson []byte
  // filesJson, err := json.Marshal(cookieCutterFiles)
  // if err != nil {
  //     return err
  // }
	cmd := exec.Command("terraform", "-chdir="+terraformDir, "apply", "-auto-approve", "-var", "repo_name="+envVars.ServiceName, "-var", "gitlab_token="+envVars.GitToken, "-var", "runner_id="+envVars.RunnerID)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
