# image-download-check

The `image-download-check` pulls a configured container image and verifies the download time stays under a configured threshold. If the pull duration exceeds the limit, the check reports failure to Kuberhealthy.

## Configuration

Set these environment variables in the `HealthCheck` spec:

- `FULL_IMAGE_URL` (required): image reference to pull (for example, `nginx:1.21`).
- `TIMEOUT_LIMIT` (required): maximum pull duration as a Go duration (for example, `180s`).
- `LOGIN_REQUIRED` (optional): set to `true` when the registry requires authentication.
- `REGISTRY_USERNAME` (optional): registry username when `LOGIN_REQUIRED` is `true`.
- `REGISTRY_PASSWORD` (optional): registry password when `LOGIN_REQUIRED` is `true`.

## Build

- `just build` builds the container image locally.
- `just test` runs unit tests.
- `just binary` builds the binary in `bin/`.

## Example HealthCheck

Apply the example below or the provided `healthcheck.yaml`:

```yaml
apiVersion: kuberhealthy.github.io/v2
kind: HealthCheck
metadata:
  name: image-download-check
  namespace: kuberhealthy
spec:
  runInterval: 10m
  timeout: 25m
  podSpec:
    spec:
      containers:
        - name: image-download-check
          image: kuberhealthy/image-download-check:sha-<short-sha>
          imagePullPolicy: IfNotPresent
          env:
            - name: FULL_IMAGE_URL
              value: "nginx:1.21"
            - name: TIMEOUT_LIMIT
              value: "180s"
            - name: LOGIN_REQUIRED
              value: "true"
            - name: REGISTRY_USERNAME
              value: "username"
            - name: REGISTRY_PASSWORD
              value: "password"
          resources:
            requests:
              cpu: 15m
              memory: 15Mi
            limits:
              cpu: 25m
      restartPolicy: Always
      terminationGracePeriodSeconds: 5
```
