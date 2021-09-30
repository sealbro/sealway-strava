job "sealway-strava" {
  datacenters = ["home"]
  type = "service"

  update {
    max_parallel = 1

    min_healthy_time = "2m"

    healthy_deadline = "5m"

    progress_deadline = "10m"

    auto_revert = false

    canary = 0
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
        servers = ["172.17.0.1", "192.168.1.55", "192.168.1.1"]
      }
    }

    service {
      name = "integration-strava"
      tags = ["wss", "ui", "http", "sealway", "api", "internal"]
      port = "app"

      check {
        type     = "http"
        port     = "app"
        interval = "1m"
        timeout  = "30s"
        path     = "/health"
      }
    }

    restart {
      attempts = 20
      interval = "24h"

      delay = "7m"

      mode = "fail"
    }

    ephemeral_disk {
      size = 300
    }

    task "container" {
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
SEALWAY_Services__Strava__Client={{ key "client_id" }}
SEALWAY_Services__Strava__Secret={{ key "client_secret" }}
SEALWAY_ConnectionStrings__Mongo__Connection={{ key "mongodb" }}
EOH

        destination = "secrets/file.env"
        env         = true
      }

      env {
        PORT = "8080"
        SLUG = "integration-strava"
      }

      resources {
        cpu    = 200
        memory = 50
      }
    }
  }
}
