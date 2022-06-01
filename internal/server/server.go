package internal

import (
	"fmt"
	"net/http"

	internal "jx-ui/internal/kube"

	// "github.com/gorilla/mux"
	"github.com/jenkins-x/jx-api/v4/pkg/client/clientset/versioned"
	jenkinsxv1 "github.com/jenkins-x/jx-api/v4/pkg/client/clientset/versioned/typed/jenkins.io/v1"

	// "github.com/rs/cors"

	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxclient"
	tknclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type Server struct {
	// addr       string
	Server     *http.Server
	JXInterface    versioned.Interface
	JXClient   jenkinsxv1.PipelineActivityInterface
	SrClient   jenkinsxv1.SourceRepositoryInterface
	TknClient  tknclient.Interface
	KubeClient kubernetes.Interface
}

// ToDo: make it configurable
const defaultNamespace = "jx"

func CreateServer() (*Server, error) {
	s := &Server{}
	config, err := internal.GetKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve kubeconfig")
	}

	jxClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create jx client with the supplied config")
	}
	s.JXClient = jxClient.JenkinsV1().PipelineActivities(defaultNamespace)

	s.SrClient = jxClient.JenkinsV1().SourceRepositories(defaultNamespace)

	tknClient, err := tknclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create tekton client")
	}
	s.TknClient = tknClient

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create kube client")
	}
	s.KubeClient = kubeClient

	s.JXInterface, err = jxclient.LazyCreateJXClient(s.JXInterface)
	if err != nil {
		return nil, fmt.Errorf("cannot lazy create jx client")
	}
	return s, nil
}
