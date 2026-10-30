# Image Download Check

Kuberhealthy's image download check

## What it is
This repository builds the container image used by Kuberhealthy to run the image-download-check check.

## Image
- `docker.io/kuberhealthy/image-download-check`
- `ghcr.io/kuberhealthy/image-download-check`
- Tags: `latest` and short git SHA for `main` pushes; `vX.Y.Z` for releases.

## Quick start
- Apply the example manifest: `kubectl apply -f healthcheck.yaml`
- Edit the manifest to set any required inputs for your environment.

## Build locally
- `docker build -f ./Containerfile -t kuberhealthy/image-download-check:dev .`

## Contributing
Issues and PRs are welcome. Please keep changes focused and add a short README update when behavior changes.

## License
See `LICENSE`.
