package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	server "jx-ui/internal/server"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	pipelines "jx-ui/internal/pipelines"
)

type Backend struct {
	Writer    *render.Render
	pipelines interface {
		List() ([]pipelines.PipelineList, error)
		// get(name string) (pipelines.Pipeline, error)
	}
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

// ToDo: make it configurable
const defaultNamespace = "jx"

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 	// otherwise, use http.FileServer to serve the static dir
// 	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
// }

// PipelinesHandler function
// func (s *server) PipelinesHandler(w http.ResponseWriter, r *http.Request) {
// 	pa, err := s.jxClient.
// 		List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		// Todo: improve error handling!
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	sort.Slice(pa.Items, func(i, j int) bool {
// 		return pa.Items[j].Spec.StartedTimestamp.Before(pa.Items[i].Spec.StartedTimestamp)
// 	})

// 	s.render.JSON(w, http.StatusOK, pa.Items) //nolint:errcheck
// }

// // PipelineHandler function
// func (s *server) PipelineHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	owner := vars["owner"]
// 	repo := vars["repo"]
// 	branch := vars["branch"]
// 	build := vars["build"]
// 	name := naming.ToValidName(owner + "-" + repo + "-" + branch + "-" + build)
// 	method := r.Method
// 	if method == "GET" {
// 		pa, err := s.jxClient.
// 			Get(context.Background(), name, metav1.GetOptions{})
// 		if err != nil {
// 			// Todo: improve error handling!
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}

// 		s.render.JSON(w, http.StatusOK, pa) //nolint:errcheck
// 	} else {
// 		pa, err := s.jxClient.
// 			Get(context.Background(), name, metav1.GetOptions{})
// 		if err != nil {
// 			// Todo: improve error handling!
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}

// 		prName := pa.Labels["tekton.dev/pipeline"]

// 		pr, err := s.tknClient.TektonV1beta1().PipelineRuns("jx").Get(context.Background(), prName, metav1.GetOptions{})
// 		if err != nil {
// 			panic(err)
// 		}

// 		if pr.Status.CompletionTime == nil {
// 			pr.Spec.Status = pipelineapi.PipelineRunSpecStatusCancelled
// 		}
// 		_, err = s.tknClient.TektonV1beta1().PipelineRuns("jx").Update(context.Background(), pr, metav1.UpdateOptions{})
// 		if err != nil {
// 			panic(err)
// 		}
// 		s.render.JSON(w, http.StatusOK, "pipeline "+name+" stopped") //nolint:errcheck
// 	}
// }

// // PipelineLogHandler returns the logs for a given pipeline
// func (s *server) PipelineLogHandler(w http.ResponseWriter, r *http.Request) {
// 	setSSEReponseHeaders(w)

// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		fmt.Println("response writer does not implement http flusher")
// 	}

// 	ctx := context.Background()
// 	vars := mux.Vars(r)
// 	owner := vars["owner"]
// 	repo := vars["repo"]
// 	branch := vars["branch"]
// 	build := vars["build"]

// 	paName := fmt.Sprintf("%s-%s-%s-%s",
// 		naming.ToValidName(owner),
// 		naming.ToValidName(repo),
// 		naming.ToValidName(branch),
// 		build)

// 	baseName := fmt.Sprintf("%s/%s/%s #%s",
// 		naming.ToValidName(owner),
// 		naming.ToValidName(repo),
// 		naming.ToValidName(branch),
// 		strings.ToLower(build))

// 	logger := tektonlog.TektonLogger{
// 		JXClient:     s.jxIface,
// 		TektonClient: s.tknClient,
// 		KubeClient:   s.kubeClient,
// 		Namespace:    defaultNamespace,
// 	}

// 	pa, err := s.jxClient.
// 		Get(context.Background(), paName, metav1.GetOptions{})
// 	if err != nil {
// 		// Todo: improve error handling!
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	triggerContext := pa.Spec.Context
// 	name := fmt.Sprintf("%s %s", baseName, naming.ToValidName(triggerContext))

// 	filter := tektonlog.BuildPodInfoFilter{
// 		Owner:      owner,
// 		Repository: repo,
// 		Branch:     branch,
// 		Build:      build,
// 	}

// 	_, _, prMap, err := logger.GetTektonPipelinesWithActivePipelineActivity(ctx, &filter)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	prList := prMap[name]

// 	for k := range logger.GetRunningBuildLogs(ctx, pa, prList, name) {
// 		select {
// 		case <-r.Context().Done():
// 			return
// 		default:
// 			fmt.Fprintf(w, "data: %s\n\n", k.Line)
// 			flusher.Flush()
// 		}
// 	}
// }

// // PipelineArchivedLogHandler returns the logs from Long term storage
// func (s *server) PipelineArchivedLogHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	owner := vars["owner"]
// 	repo := vars["repo"]
// 	branch := vars["branch"]
// 	build := vars["build"]
// 	name := naming.ToValidName(owner + "-" + repo + "-" + branch + "-" + build)

// 	logger := tektonlog.TektonLogger{
// 		JXClient:     s.jxIface,
// 		TektonClient: s.tknClient,
// 		KubeClient:   s.kubeClient,
// 		Namespace:    defaultNamespace,
// 	}

// 	pa, err := s.jxClient.
// 		Get(context.Background(), name, metav1.GetOptions{})
// 	if err != nil {
// 		// Todo: improve error handling!
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	logs := []string{}

// 	// Read for archived logs if builds are not running
// 	for line := range logger.StreamPipelinePersistentLogs(pa.Spec.BuildLogsURL) {
// 		logs = append(logs, line.Line)
// 	}

// 	s.render.JSON(w, http.StatusOK, logs) //nolint:errcheck
// }

// // RepositoriesHandler function
// func (s *server) RepositoriesHandler(w http.ResponseWriter, r *http.Request) {
// 	repo, err := s.srClient.
// 		List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		// Todo: improve error handling!
// 		panic(err)
// 	}

// 	s.render.JSON(w, http.StatusOK, repo.Items) //nolint:errcheck
// }

func registerRoutes(router *mux.Router, backend *Backend) *mux.Router {
	router.HandleFunc("/api/v1/pipelines", backend.PipelinesHandler())
	// router.HandleFunc("/api/v1/pipelines/{owner}/{repo}/{branch}/{build}", server.PipelineHandler).Methods("GET", "POST")
	// router.HandleFunc("/api/v1/logs/{owner}/{repo}/{branch}/{build}", server.PipelineLogHandler)
	// router.HandleFunc("/api/v1/logs_archived/{owner}/{repo}/{branch}/{build}", server.PipelineArchivedLogHandler)
	// router.HandleFunc("/api/v1/repositories", server.RepositoriesHandler)
	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)
	return router
}

func main() {
	s, err := server.CreateServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	backend := &Backend{
		pipelines: pipelines.PipelineServer{
			Server: s,
		},
	}
	backend.Writer = render.New()
	router := mux.NewRouter()
	registerRoutes(router, backend)
	// This is only required for dev mode
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	s.Server = &http.Server{
		Handler: c.Handler(router),
		Addr:    "127.0.0.1:8080", // Todo: Make it configurable
	}

	log.Fatal(s.Server.ListenAndServe())
}

func (b *Backend) PipelinesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pipelines, err := b.pipelines.List()
		if err != nil {
			panic(err)
		}
		b.Writer.JSON(w, http.StatusOK, pipelines)
	}
}

// setSSEReponseHeaders sets the response headers to appropriate values for Server Sent Events (SSE)
func setSSEReponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
