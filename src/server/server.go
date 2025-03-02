package main

import (
	"log"
	"os"
	"time"

	"ThaiLy/graph/controller"
	"ThaiLy/graph/generated"
	resolver "ThaiLy/graph/resolver"
	"ThaiLy/server/client"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8081"

func main() {

	// Set gin to release mode
	// gin.SetMode(gin.ReleaseMode)
	// Load .env file
	// godotenv.Load()

	// fmt.Println("MySQL Name: ", os.Getenv("MYSQL_NAME"))

	// Create Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Chỉ cho phép origin cụ thể
		AllowMethods:     []string{"POST, GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	auth, err := client.NewGRPCClient("localhost:55555")
	if err != nil {
		log.Fatal("client auth error %v", err)
	}
	ctrl := controller.NewController(auth)
	// Create GraphQL handler
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{Ctrl: ctrl}}))

	// Add GraphQL transports (Options, GET, POST)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Set query cache
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Gin routes foyground and query handler
	r.GET("/", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))
	r.POST("/query", gin.WrapH(srv))
	// middlewares.RequireAuth, func(c *gin.Context) {
	// 	account, exists := c.Get("account")
	// 	if exists {
	// 		// Nếu account có, thêm account vào context
	// 		// Sử dụng custom key mà không cần ép kiểu (type assertion)
	// 		ctx := context.WithValue(c.Request.Context(), middlewares.AccountKey, account)
	// 		c.Request = c.Request.WithContext(ctx)
	// 	}
	// },

	// Run the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
