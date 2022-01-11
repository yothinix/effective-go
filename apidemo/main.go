package main

import (
	"apidemo/router"
	"apidemo/store"
	"apidemo/todo"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"golang.org/x/time/rate"
)

var (
	buildcommit = "dev"
	buildtime 	= time.Now().String()
)

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Println("please consider environment variables: %s", err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongoadmin:secret@localhost:27017"))
	if err != nil {
		panic("failed to connect database")
	}
	collection := client.Database("myapp").Collection("todos")

	r := router.NewMyRouter()

	//gormStore := store.NewGormStore(db)
	mongoStore := store.NewMongoDBStore(collection)

	handler := todo.NewTodoHandler(mongoStore)
	r.POST("/todos", handler.NewTask)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

    s := &http.Server{
		Addr:	":" + os.Getenv("PORT"),
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, precc Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}

}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
