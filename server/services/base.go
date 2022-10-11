package services

import (
	"final/server/views"

	"github.com/gin-gonic/gin"
)

func WriteJsonResponse(ctx *gin.Context, payload *views.Response) {
	ctx.JSON(payload.Status, payload)
}

func WriteJsonResponse_Succes(ctx *gin.Context, data *views.Resp_Register_Success) {
	ctx.JSON(data.Status, data)
}

func WriteJsonResponse_Login(ctx *gin.Context, data *views.Resp_Login) {
	ctx.JSON(data.Status, data)
}

func WriteJsonResponse_Put(ctx *gin.Context, data *views.Resp_Put) {
	ctx.JSON(data.Status, data)
}

func WriteJsonResponse_Delete(ctx *gin.Context, data *views.Resp_Delete) {
	ctx.JSON(data.Status, data)
}

// //// Foto JSON RESPONSE //////
func WriteJsonResponse_PostPhoto(ctx *gin.Context, data *views.Resp_Post_Photo) {
	ctx.JSON(data.Status, data)
}
