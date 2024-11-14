package file

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	mlog "mall/log"
	"mall/middleware/auth"
)

const (
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "3vItiny21TIhkRF0wHbK"
	secretAccessKey = "F4tJFaTCVMtyKeafdhaehiVu9PUiPssS2GPk8HnS"
	useSSL          = false
	bucketName      = "mall"
)

var log *mlog.Log
var MinioClient *minio.Client

func Init(engine *gin.Engine) {
	log = mlog.NewLog("FileAPI")
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		panic(err)
	}
	MinioClient = client
	group := engine.Group("/File", auth.ParseToken)
	{
		group.POST("/Upload", Upload)
		group.GET("/Download/:filename", Download)
	}
}
