job "sealway-strava" {
  datacenters = ["home"]
  type = "service"

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
      port "app" {
        to = 8080
      }

      dns {
        servers = ["172.17.0.1", "192.168.1.1"]
      }
    }

    service {
      name = "integration-strava"
      tags = ["wss", "http", "sealway", "api", "private", "internal"]
      port = "app"

      check {
        type     = "http"
        port     = "app"
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

        ports = ["app"]

        labels {
          from_nomad = "yes"
        }
      }

      template {
        data = <<EOH
SEALWAY_Services__Strava__Client={{with secret "applications/prod/default/Services/Strava"}}{{.Data.data.client_id}}{{end}}
SEALWAY_Services__Strava__Secret={{with secret "applications/prod/default/Services/Strava"}}{{.Data.data.client_secret}}{{end}}
SEALWAY_ConnectionStrings__Mongo__Connection={{ key "applications/prod/default/ConnectionStrings/Mongo" }}
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
