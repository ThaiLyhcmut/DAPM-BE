package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"ThaiLy/graph/controller"
	"ThaiLy/graph/generated"
	"ThaiLy/graph/helper"
	"ThaiLy/graph/resolver"

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
		AllowMethods:     []string{"POST", "GET"},
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
	auth, err := client.NewGRPCAuthClient("localhost:55555")
	if err != nil {
		log.Fatalf("client auth error %v", err)
	}
	equipment, err := client.NewGRPCEquipmentClient("localhost:55556")
	if err != nil {
		log.Fatalf("client equipment error %v", err)
	}
	kafka, err := client.NewGRPCKafkaClient("localhost:55557")
	if err != nil {
		log.Fatalf("client equipment error %v", err)
	}
	ctrl := controller.NewController(auth, equipment, kafka)
	// Create GraphQL handler
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{Ctrl: ctrl}}))

	// Add GraphQL transports (Options, GET, POST)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{})
	// Set query cache
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Gin routes foyground and query handler
	r.GET("/", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))
	r.Any("/query", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
				Claims, err := helper.ParseJWT(token)
				if Claims == nil || err != nil {
					c.JSON(401, gin.H{"message": "Invalid Authorization header"})
					c.Abort()
					return
				}
				ctx := context.WithValue(c.Request.Context(), helper.Auth, Claims)
				c.Request = c.Request.WithContext(ctx)
			}
		}

		// Kiểm tra nếu là WebSocket, thì dùng srv.ServeHTTP trực tiếp
		if strings.Contains(c.GetHeader("Upgrade"), "websocket") {
			srv.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Nếu là HTTP bình thường thì vẫn dùng WrapH
		gin.WrapH(srv)(c)
	})

	// Run the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
