package route

import (
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/database"
	"github.com/Artexus/api-widyabhuvana/src/middleware"
	v1 "github.com/Artexus/api-widyabhuvana/src/route/v1"
	"github.com/gin-gonic/gin"
)

var client *firestore.Client

func InitRouter(engine *gin.Engine) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		var err error
		client, err = database.InitDB()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
	wg.Wait()

	middleware := middleware.New(client)
	v1 := v1.New(middleware)

	v1.InitRouter(engine, client)
}
