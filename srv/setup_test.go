package srv

import (
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type configTestCase struct {
	clientConfig      clientcmd.ClientConfig
	errStringContains string
}

var KUBECONFIG = `apiVersion: v1
kind: Config
preferences: {}
users:
- name: a
  user:
    client-certificate-data: ""
    client-key-data: ""
clusters:
- name: a
  cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:8080
contexts:
- name: a
  context:
    cluster: a
    user: a
current-context: a
`

func TestSetup_GetRestConfig(t *testing.T) {
	basic, err := clientcmd.NewClientConfigFromBytes([]byte(KUBECONFIG))
	assert.NilError(t, err)
	for i, ctc := range []configTestCase{
		{
			clientcmd.NewDefaultClientConfig(clientcmdapi.Config{}, &clientcmd.ConfigOverrides{}),
			"no configuration has been provided",
		},
		{
			basic,
			"",
		},
	} {
		s := &Setup{
			ClientConfig: ctc.clientConfig,
		}
		_, err := s.GetRestConfig()

		switch len(ctc.errStringContains) {
		case 0:
			if err != nil {
				t.Errorf("%d: unexpected error: %s", i, err.Error())
			}
		default:
			if err == nil {
				t.Errorf("%d: wrong error detected: %s (expected) != %s (actual)", i, ctc.errStringContains, err)
			}
			if !strings.Contains(err.Error(), ctc.errStringContains) {
				t.Errorf("%d: wrong error detected: %s (expected) != %s (actual)", i, ctc.errStringContains, err.Error())
			}
		}
	}
	sEmpty := &Setup{}
	sEmpty.ClientConfig, err = clientcmd.NewClientConfigFromBytes([]byte(""))
	assert.NilError(t, err)
	_, err = sEmpty.GetRestConfig()
	assert.ErrorContains(t, err, "no configuration has been provided")

	sEmpty = &Setup{}
	sEmpty.KubeCfgPath = filepath.Join("bad", "path", "file")
	_, err = sEmpty.GetRestConfig()
	assert.ErrorContains(t, err, "config file not found")
}
