package Product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/model"
	"mall/service/product/proto/product"
	"mall/util"
	"strconv"
)

func List(c *gin.Context) {

	c.DefaultQuery("page_size", "10")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Info("get page:" + err.Error())
		util.Response(c, model.BADREQUEST, err.Error())
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		log.Info("get page size:" + err.Error())
		util.Response(c, model.BADREQUEST, err.Error())
		return
	}
	req := product.ListProductsReq{
		Page:     uint32(page),
		PageSize: uint32(size),
	}

	resp, err := ProductClient.ListProducts(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}

func Get(c *gin.Context) {
	str := c.Query("id")
	if str == "" {
		log.Info("no id in get request")
		util.Response(c, model.BADREQUEST, "there is no id")
		return
	}
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Info("get product id:" + err.Error())
		util.Response(c, model.BADREQUEST, err.Error())
		return
	}
	req := product.GetProductReq{
		Id: uint32(id),
	}
	resp, err := ProductClient.GetProduct(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}

func Search(c *gin.Context) {
	m := c.DefaultQuery("method", "name")
	query := c.Query("query")
	if query == "" {
		log.Info("no query in search")
		util.Response(c, model.BADREQUEST, "need query")
		return
	}
	req := product.SearchProductsReq{
		Method: m,
		Query:  query,
	}
	resp, err := ProductClient.SearchProducts(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
	return
}

func Create(c *gin.Context) {
	req := product.CreateProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Info("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	resp, _ := ProductClient.CreateProducts(c, &req)
	util.Response(c, model.OK, "", resp)
}

func Update(c *gin.Context) {
	req := product.UpdateProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	resp, err := ProductClient.UpdateProducts(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}

func Delete(c *gin.Context) {
	req := product.DeleteProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		util.Response(c, model.BADREQUEST, "json can not bind")
		return
	}
	resp, err := ProductClient.DeleteProducts(c, &req)
	if err != nil {
		util.Response(c, model.ERROR, err.Error())
		return
	}
	util.Response(c, model.OK, "", resp)
}
