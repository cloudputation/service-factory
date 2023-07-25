

curl -X POST -H "Content-Type: application/json" \
     -d '{
         	"repo_url": "https://gitlab.com/franksrobins/cookie-cutter-api.git",
          "namespace_id": "69638879",
         	"terraform_dir": "terraform/",
         	"git_token": "glpat-HuEekH9zXTi8DbWkyLzo",
         	"runner_id": "20665767",
         	"nomad_host": "10.100.200.241",
         	"service_name": "my-first-service",
          "service_group": "tests-apis",
          "service_port": "9999",
          "service_tag": "SF-Managed",
          "service_type": "service",
          "target_nomad_client": "tower2",
          "service_repository_owner": "franksrobins",
          "service_registry_token": "glpat-kuuxn3oXxwBsCMkxssDY"
         }' \
     http://10.100.200.248:48840/run
