job "sealway-strava" {
  datacenters = ["lan"]

  type = "service"

  // affinity {
  //   attribute = "${meta.os_architecture}"
  //   value     = "amd64"
  //   weight    = 100
  // }

  constraint {
    attribute = "${attr.cpu.arch}"
    value     = "amd64"
  }

  update {
    max_parallel = 1
    health_check = "checks"
    min_healthy_time = "30s"
    healthy_deadline = "5m"
  }

  migrate {
    max_parallel = 1
    health_check = "checks"
    min_healthy_time = "2m"
    healthy_deadline = "5m"
  }

  group "strava-app" {
    count = 1

    network {
      port "app-http" {
        to = 8080
        host_network = "private"
      }
    }

    service {
      name = "integration-strava"
      tags = ["wss", "http", "sealway", "api", "private", "internal"]
      port = "app-http"

      check {
        type     = "http"
        port     = "app-http"
        interval = "30s"
        timeout  = "5s"
        path     = "/healthz"

        check_restart {
          limit = 3
          grace = "90s"
          ignore_warnings = true
        }
      }
    }

    restart {
      attempts = 20
      interval = "30m"
      delay = "1m"
      mode = "fail"
    }

    task "sealway-strava" {
      driver = "docker"

      config {
        image = "sealway/strava"
        force_pull = true

        ports = ["app-http"]

        labels {
          from_nomad = "yes"
        }

        logging {
          type = "loki"
          config {
            loki-pipeline-stages = <<EOH
- static_labels:
    app: sealway-strava
- json:
    expressions:
      time: ts_orig
- timestamp:
    source: time
    format: RFC3339
EOH
          }
        }
      }

      template {
        data = <<EOH
STRAVA_CLIENT={{with secret "applications/prod/default/Services/Strava"}}{{.Data.data.client_id}}{{end}}
STRAVA_SECRET={{with secret "applications/prod/default/Services/Strava"}}{{.Data.data.client_secret}}{{end}}
MONGO_CONNECTION=mongodb://mongo.service.consul
EOH

        destination = "secrets/file.env"
        env         = true
      }

      vault {
        policies = ["nomad-server"]
        env = false
      }

      env {
        PORT = "8080"
        SLUG = "integration-strava"
      }

      resources {
        cpu    = 100
        memory = 32
      }
    }
  }
}
