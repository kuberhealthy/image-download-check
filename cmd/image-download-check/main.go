package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	nodecheck "github.com/kuberhealthy/kuberhealthy/v3/pkg/nodecheck"
	log "github.com/sirupsen/logrus"
)

// main loads configuration and executes the image download check.
func main() {
	// Enable nodecheck debug output for parity with v2 behavior.
	nodecheck.EnableDebugOutput()

	// Parse configuration from environment variables.
	cfg, err := parseConfig()
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Track how long the image pull takes.
	startTime := time.Now()

	// Download the image to validate the pull path.
	_, err = downloadImage(cfg)
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Calculate the pull duration.
	duration := time.Since(startTime)
	log.Infoln("image took this many seconds to download:", duration.Seconds())

	// Compare the pull duration against the configured limit.
	log.Infoln("checking to see if", duration, "<", cfg.TimeoutLimit)
	if duration < cfg.TimeoutLimit {
		log.Infoln("check passes, download duration is less than timeout limit.")
		reportSuccessAndExit()
		return
	}

	// Report failure when the pull exceeds the limit.
	log.Infoln("check fails, download duration is greater than timeout limit.")
	reportFailureAndExit(fmt.Errorf("check has failed, download duration is greater than timeout limit"))
}

// downloadImage pulls and inspects the configured image reference.
func downloadImage(cfg *CheckConfig) (v1.Image, error) {
	// Select the pull options based on the auth requirement.
	var options []crane.Option
	if cfg.LoginRequired {
		auth := &authn.Basic{
			Username: cfg.RegistryUsername,
			Password: cfg.RegistryPassword,
		}
		options = append(options, crane.WithAuth(auth))
	}

	// Pull the image.
	image, err := crane.Pull(cfg.FullImageURL, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to pull image %s: %w", cfg.FullImageURL, err)
	}
	log.Infoln("successfully downloaded image:", cfg.FullImageURL)

	// Save the image to /dev/null to mirror v2 behavior.
	err = crane.Save(image, "emptytag", "/dev/null")
	if err != nil {
		return nil, fmt.Errorf("failed to save image tarball: %w", err)
	}

	// Log the layer count for visibility.
	layers, err := image.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to read image layers: %w", err)
	}
	log.Infoln("layer count", len(layers))

	// Log the image size for visibility.
	size, err := image.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to read image size: %w", err)
	}
	log.Infoln("image size", size)

	return image, nil
}

// reportSuccessAndExit reports success to Kuberhealthy and exits.
func reportSuccessAndExit() {
	// Report the success to Kuberhealthy.
	err := checkclient.ReportSuccess()
	if err != nil {
		log.Errorln("Error reporting success to Kuberhealthy servers:", err)
		os.Exit(1)
	}
	log.Infoln("Successfully reported success to Kuberhealthy servers")

	// Exit after reporting success.
	os.Exit(0)
}

// reportFailureAndExit reports a failure to Kuberhealthy and exits.
func reportFailureAndExit(err error) {
	// Log the error locally.
	log.Errorln(err)

	// Report the failure to Kuberhealthy.
	reportErr := checkclient.ReportFailure([]string{err.Error()})
	if reportErr != nil {
		log.Errorln("Error reporting failure to Kuberhealthy servers:", reportErr)
		os.Exit(1)
	}
	log.Infoln("Successfully reported failure to Kuberhealthy servers")

	// Exit after reporting failure.
	os.Exit(0)
}
