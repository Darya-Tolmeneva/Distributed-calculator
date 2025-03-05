package application

import (
	"Distributed_calculator/internal/services/agent"
	"Distributed_calculator/internal/transport"
	"log"
	"net/http"
	"os"
	"sync"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = ":8080"
	} else {
		config.Addr = ":" + config.Addr
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) RunServer() error {
	ports := []string{"50050", "50051", "50052", "50053", "50054"}

	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			log.Printf("Starting agent on port %s", port)
			agent.Run(port)
		}(port)
	}

	log.Printf("Starting server on %s", a.config.Addr)
	http.HandleFunc("/api/v1/calculate", transport.PutExpressionHandler)
	http.HandleFunc("/api/v1/expressions", transport.GetExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/{id}", transport.GetExpressionHandler)
	http.HandleFunc("/internal/task", transport.HandleTaskResult)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Not found: %s", r.URL.Path)
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
	})

	return http.ListenAndServe(a.config.Addr, nil)
}
