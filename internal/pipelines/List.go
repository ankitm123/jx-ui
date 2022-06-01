package pipelines

import (
	"context"
	"sort"

	server "jx-ui/internal/server"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineServer struct {
	Server *server.Server
}

type PipelineList struct {
	Owner      string
	Repository string
	Branch     string
	Build      string
	Status     string
	StartTime  *metav1.Time
	EndTime    *metav1.Time
}

// func (ps *PipelineServer) PipelinesHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		pipelines, err := ps.List()
// 		if err != nil {
// 			panic(err)
// 		}

// 		ps.Server.Writer.JSON(w, http.StatusOK, pipelines) //nolint:errcheck
// 	}
// }

func (ps PipelineServer) List() ([]PipelineList, error) {
	pList := PipelineList{}
	pLists := []PipelineList{}

	pa, err := ps.Server.JXClient.
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return pLists, err
	}

	sort.Slice(pa.Items, func(i, j int) bool {
		return pa.Items[j].Spec.StartedTimestamp.Before(pa.Items[i].Spec.StartedTimestamp)
	})

	for k := range pa.Items {
		pList.Owner = pa.Items[k].Spec.GitOwner
		pList.Repository = pa.Items[k].Spec.GitRepository
		pList.Branch = pa.Items[k].Spec.GitBranch
		pList.Build = pa.Items[k].Spec.Build
		pList.Status = pa.Items[k].Spec.Status.String()
		pList.StartTime = pa.Items[k].Spec.StartedTimestamp
		pList.EndTime = pa.Items[k].Spec.CompletedTimestamp
		pLists = append(pLists, pList)
	}

	return pLists, nil
}
