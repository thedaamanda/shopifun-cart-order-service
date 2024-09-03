package routes

import (
	"cart-order-service/config"
	"cart-order-service/handlers/cart"
	"cart-order-service/util/middleware"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Routes struct {
	Router *http.ServeMux
	Cart   *cart.Handler
}

func URLRewriter(baseURLPath string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, baseURLPath)

		next.ServeHTTP(w, r)
	}
}

func (r *Routes) SetupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL != "" && baseURL != "/" {
		r.Router.HandleFunc(baseURL+"/", URLRewriter(baseURL, r.Router))
	}
}

func (r *Routes) cartRoutes() {
	r.Router.HandleFunc("GET /cart/{user_id}", middleware.ApplyMiddleware(r.Cart.GetCartByUserID, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) SetupRouter() {
	r.Router = http.NewServeMux()
	r.SetupBaseURL()
	r.cartRoutes()
}

func (r *Routes) Run(port string) {
	r.SetupRouter()

	log.Printf("[Running-Success] clients on localhost on port :%s", port)
	srv := &http.Server{
		Handler:      r.Router,
		Addr:         "localhost:" + port,
		WriteTimeout: config.WriteTimeout() * time.Second,
		ReadTimeout:  config.ReadTimeout() * time.Second,
	}

	log.Panic(srv.ListenAndServe())
}
