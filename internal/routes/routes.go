package routes

// import (
// 	"github.com/gorilla/mux"
// 	server "jx-ui/internal/server"
// )
// type spaHandler struct {
// 	staticPath string
// 	indexPath  string
// }

// func RegisterRoutes(router *mux.Router, server *server.Server) *mux.Router {
// 	router.HandleFunc("/api/v1/pipelines", server.PipelinesHandler)
// 	router.HandleFunc("/api/v1/pipelines/{owner}/{repo}/{branch}/{build}", server.PipelineHandler).Methods("GET", "POST")
// 	router.HandleFunc("/api/v1/logs/{owner}/{repo}/{branch}/{build}", server.PipelineLogHandler)
// 	router.HandleFunc("/api/v1/logs_archived/{owner}/{repo}/{branch}/{build}", server.PipelineArchivedLogHandler)
// 	router.HandleFunc("/api/v1/repositories", server.RepositoriesHandler)
// 	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
// 	router.PathPrefix("/").Handler(spa)
// 	return router
// }
