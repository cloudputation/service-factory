service {
  name    = "test-service11"
  group   = "test-apis"
  port    = "9900"
  tags    = ["SF-Managed"]

  template {
    template_url  = "gitlab.com"
    template      = "franksrobins/cookie-cutter-api"
  }

  repository {
    provider          = "gitlab"
    namespace_id      = "69638879"
    runner_id         = "20665767"
    registry_token    = "glpat-kuuxn3oXxwBsCMkxssDY"
    repository_owner  = "franksrobins"
  }

  secrets {
    vault {
      vault_host    = "10.100.200.241"
    }
  }

  network {
    consul {
      consul_host       = "10.100.200.241"
      repo_service_name = "service-repository"
    }
  }

  scheduler {
    nomad {
      nomad_host    = "10.100.200.241"
      target_host   = "comm1"
    }
  }

  pipeline {
    atlantis {
      atlantis_host = "10.100.200.243"
      webhook {
        secret = "EC7F8C30-3384-4310-B65B-DCE6FA6DF268"
        gitlab {
          triggers {
            push_events                 = true
            issues_events               = false
            tag_push_events             = false
            note_events                 = false
            job_events                  = false
            pipeline_events             = false
            wiki_page_events            = false
            merge_requests_events       = true
            confidential_issues_events  = false
          }
        }
      }
    }
  }

  code_suggestion {
    open_ai {
      token       = "EC7F8C30-3384-4310-B65B-DCE6FA6DF268"
      prompt_dir  = "/etc/service-factory/prompt.d"
      webhook {
        secret = "EC7F8C30-3384-4310-B65B-DCE6FA6DF268"
        triggers {
          push_events   = true
          issues_events = false
        }
      }
    }
  }

}
