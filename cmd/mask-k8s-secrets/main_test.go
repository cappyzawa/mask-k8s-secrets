package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	secretManifest = `apiVersion: v1
kind: Secret
metadata:
    name: secret-basic-auth
type: kubernetes.io/basic-auth
data:
    username: YWRtaW4=
    password: MWYyZDFlMmU2N2Rm
`
	maskedSecretManifest = `apiVersion: v1
data:
    password: '***'
    username: '***'
kind: Secret
metadata:
    name: secret-basic-auth
type: kubernetes.io/basic-auth
`
	deploymentManifest = `apiVersion: apps/v1
kind: Deployment
metadata:
    labels:
        app: nginx
    name: nginx-deployment
spec:
    replicas: 3
    selector:
        matchLabels:
            app: nginx
    template:
        metadata:
            labels:
                app: nginx
        spec:
            containers:
                - image: nginx:1.14.2
                  name: nginx
                  ports:
                    - containerPort: 80
`
)

var (
	deploymentAndSecretManifests       = deploymentManifest + "---\n" + secretManifest
	deploymentAndMaskedSecretManifests = deploymentManifest + "---\n" + maskedSecretManifest
)

func TestRun(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input    string
		output   string
		exitCode int
	}{
		// "only secret manifest": {
		// 	input:    secretManifest,
		// 	output:   maskedSecretManifest,
		// 	exitCode: 0,
		// },
		// "only deploymentManifest": {
		// 	input:    deploymentManifest,
		// 	output:   deploymentManifest,
		// 	exitCode: 0,
		// },
		"secret and deployment manifest": {
			input:    deploymentAndSecretManifests,
			output:   deploymentAndMaskedSecretManifests,
			exitCode: 0,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := new(bytes.Buffer)
			c := &cli{inStream: bytes.NewBuffer([]byte(tc.input)), outStream: actual, errStream: io.Discard}
			exitCode := c.run([]string{})
			if exitCode != tc.exitCode {
				t.Errorf("exitCode want %d, got %d", tc.exitCode, exitCode)
			}
			output := actual.String()
			if diff := cmp.Diff(tc.output, output); diff != "" {
				t.Errorf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
