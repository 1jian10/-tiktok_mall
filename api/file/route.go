package file

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go"
	"io"
	"mall/model"
	"mall/util"
	"net/http"
)

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	if file == nil {
		log.Info("no file upload")
		util.Response(c, model.BADREQUEST, "no file upload")
		return
	}
	temp, err := file.Open()
	if err != nil {
		log.Error(err.Error())
		util.Response(c, model.BADREQUEST, "open file failed")
		return
	}
	log.Info("get file:" + file.Filename)
	defer temp.Close()

	contentType, err := mimetype.DetectReader(temp)
	if err != nil {
		log.Error("detect type failed:" + err.Error())
		util.Response(c, model.BADREQUEST, "can not detect type:"+err.Error())
		return
	}
	FileName := uuid.New().String()
	_, err = MinioClient.PutObject(bucketName, FileName, temp, file.Size, minio.PutObjectOptions{ContentType: contentType.String()})
	if err != nil {
		log.Error("upload file to minio:" + err.Error())
		util.Response(c, model.BADREQUEST, "upload file to minio failed")
		return
	}
	path := "http://1jian10.cn:23333/File/Downlaod/" + FileName
	util.Response(c, model.OK, "", gin.H{"FilePath": path})
}

func Download(c *gin.Context) {
	var n Name
	if err := c.ShouldBindUri(&n); err != nil {
		log.Info("bind uri fail:" + err.Error())
		util.Response(c, model.BADREQUEST, "uri can not bind")
		return
	}
	obj, err := MinioClient.GetObject(bucketName, n.FileName, minio.GetObjectOptions{})
	if err != nil {
		log.Error("get object fail:" + err.Error())
		util.Response(c, model.BADREQUEST, err.Error())
		return
	}
	defer obj.Close()
	info, err := MinioClient.StatObject(bucketName, n.FileName, minio.StatObjectOptions{})
	if err != nil {
		log.Error("get object info fail:" + err.Error())
		util.Response(c, model.BADREQUEST, err.Error())
		return
	}
	buf := make([]byte, info.Size)
	_, err = obj.Read(buf)
	if err != nil && err != io.EOF {
		log.Error("read file fail:" + err.Error())
		util.Response(c, model.ERROR, err.Error())
		return
	}
	c.Data(http.StatusOK, info.ContentType, buf)
}
