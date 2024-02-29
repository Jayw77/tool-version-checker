# Version Checker Application

## Overview

The Version Checker is a Go application designed to monitor and compare the versions of various tools against their latest available versions.
It utilizes a configuration file (`config/config.yaml`) to specify which tools to check, along with optional settings for fetch intervals, current versions,
and comments for each tool. The application periodically fetches the latest version data for these tools and serves a webpage displaying this information,
making it easy to see which tools are up to date and which need attention.

## Features

- **Periodic Version Check**: Automatically fetches tool version data at specified intervals.
- **Custom Configuration**: Supports overriding default settings for fetch intervals, current versions, and comments through `config.yaml`.
- **Web Interface**: Provides a web interface to view the current and latest versions of tools.

## Prerequisites

Before running the application, ensure you have the following installed:

- Go (if running locally)
- Docker (if running via Docker)

## Default supported tools

By default some endpoints for common tools are setup.  If you want to use an endpoint that isn't already configured, you can either use custom
as displayed below or do a PR to add support for the tool of choice by amending [here](builtin_endpoints.go).

## Configuration - config.yaml

\*\*Note: if you don't pass a `config/config.yaml` then it will use the `default_config.yaml` which has test data only.\*\*

```yaml
fetchInterval: 10 # defined in minutes

endpoints:
  - name: Grafana
    type: grafana
    url: https://grafana.xxxxx.dev
  - name: Something custom
    type: custom
    url: https://mycustomshiz.something.com
    custom:
      myVersion:
        endpoint: /stuff/version
        jsonKey: data.version
      latestVersion:
        endpoint: https://api.github.com/repos/customs/randomshiz/releases/latest
        jsonKey: tag_name
  - name: custom working
    type: custom
    url: https://api.github.com
    custom:
      myVersion:
        endpoint: /repos/prometheus/prometheus/releases/latest
        jsonKey: tag_name
      latestVersion:
        endpoint: https://api.github.com/repos/prometheus/prometheus/releases/latest
        jsonKey: tag_name
```

## Output / Frontend

To be added later.

## Running Locally

To run the application locally, follow these steps:

- Clone the repository and navigate to the project directory.
- Ensure config.yaml is in the root directory of the project.
  Run the application:

```bash
go run .
```

The server will start on port 8080. Access the web interface by navigating to http://localhost:8080 in your web browser.

## Running with Docker

The application is also available as a Docker image on DockerHub (jayw77/version-checker). To run it using Docker, follow these steps:

### Pull the Docker image:

`docker pull jayw77/version-checker:latest`

### Run the Docker container, ensuring to mount the config.yaml to the container:

```bash
docker run -d -p 8080:8080 -v $(pwd)/config/config.yaml:/config/config.yaml jayw77/version-checker:latest
```

Replace $(pwd) with the path to the directory containing your config.yaml file. The server will start on port 8080, and you can access the web interface as described above.

## Dockerhub Releases

A github workflow is setup to automatically release to dockerhub on a successful run on main. It will update the latest tag and tag the SHA.
In addition, when creating a release using semver, it will automatically create a tag to match the release.
