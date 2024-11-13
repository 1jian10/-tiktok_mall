package Product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/api"
	"mall/service/product/proto/product"
	"net/http"
)

func List(c *gin.Context) {
	req := product.ListProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, ListResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}

	resp, err := ProductClient.ListProducts(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, ListResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, ListResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
}

func Get(c *gin.Context) {
	req := product.GetProductReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, GetResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	resp, err := ProductClient.GetProduct(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, GetResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, GetResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
}

func Search(c *gin.Context) {
	req := product.SearchProductsReq{}
	if err := c.ShouldBindUri(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, SearchResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	resp, err := ProductClient.SearchProducts(c, &req)
	if err != nil {
		c.JSON(http.StatusOK, SearchResp{
			Status: api.Status{
				Code:     api.ERROR,
				ErrorMsg: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, SearchResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
	return
}

func Create(c *gin.Context) {
	req := product.CreateProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error("can not bind req:" + fmt.Sprint(&req))
		c.JSON(http.StatusOK, CreateResp{
			Status: api.Status{
				Code:     api.BADREQUEST,
				ErrorMsg: "json can not bind",
			},
		})
		return
	}
	resp, _ := ProductClient.CreateProducts(c, &req)
	c.JSON(http.StatusOK, CreateResp{
		Status: api.Status{
			Code: api.OK,
		},
		Data: resp,
	})
}

// Update service not finish
func Update(c *gin.Context) {
	req := product.UpdateProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	resp, _ := ProductClient.UpdateProducts(c, &req)
	c.JSON(http.StatusOK, resp)
}

// Delete service not finish
func Delete(c *gin.Context) {
	req := product.DeleteProductsReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	resp, _ := ProductClient.DeleteProducts(c, &req)
	c.JSON(http.StatusOK, resp)
}
