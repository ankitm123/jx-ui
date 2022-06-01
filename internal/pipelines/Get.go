package pipelines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/naming"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pipeline struct {
	Repository string
	Branch     string
	Build      string
	Status     string
	StartTime  *metav1.Time
	EndTime    *metav1.Time
}

func (ps *PipelineServer) PipelineHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner := vars["owner"]
		repo := vars["repo"]
		branch := vars["branch"]
		build := vars["build"]
		name := naming.ToValidName(owner + "-" + repo + "-" + branch + "-" + build)

		pipeline, err := ps.get(name)
		if err != nil {
			panic(err)
		}

		fmt.Println(pipeline)
		// ps.Server.Writer.JSON(w, http.StatusOK, pipeline) //nolint:errcheck
	}
}

func (ps PipelineServer) get(name string) (Pipeline, error) {
	pipeline := Pipeline{}
	pa, err := ps.Server.JXClient.
		Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return pipeline, err
	}
	pipeline.Repository = pa.Spec.GitRepository

	return pipeline, nil
}
