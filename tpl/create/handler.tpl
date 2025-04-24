package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/model"
    "{{ .ProjectName }}/internal/service"
    "{{ .ProjectName }}/pkg/helper/result"
	"strings"
)

// @wire:Handler
type {{ .StructName }}Handler struct {
	service service.{{ .StructName }}Service
}

// New{{ .StructName }}Handler 创建控制器
func New{{ .StructName }}Handler(service service.{{ .StructName }}Service) *{{ .StructName }}Handler {
	return &{{ .StructName }}Handler{
		service: service,
	}
}

type IdsRequest struct {
	ids []string
}

// @Summary 创建
// @Description 创建新的记录
// @Tags {{ .StructNameLowerFirst }}
// @Accept  json
// @Produce  json
// @Param {{ .StructNameLowerFirst }} body model.{{ .StructName }} true "信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s [post]
func (c *{{ .StructName }}Handler) Create{{ .StructName }}(ctx *gin.Context) {
	var {{ .StructNameLowerFirst }} model.{{ .StructName }}
	if err := ctx.ShouldBindJSON(&{{ .StructNameLowerFirst }}); err != nil {
		result.FailByErr(ctx, err)
		return
	}

	err := c.service.Create{{ .StructName }}(&{{ .StructNameLowerFirst }})
	if err != nil {
		result.FailByErr(ctx, err)
	} else {
		result.Success(ctx, {{ .StructNameLowerFirst }})

	}

}

// @Summary 获取
// @Description 根据ID获取记录
// @Tags {{ .StructNameLowerFirst }}
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/{id} [get]
func (c *{{ .StructName }}Handler) Get{{ .StructName }}(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		result.ServerError(ctx, "IDs are required")
		return
	}

	{{ .StructNameLowerFirst }}, err := c.service.Get{{ .StructName }}(id)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	if {{ .StructNameLowerFirst }} == nil {
		result.ServerError(ctx, "{{ .StructName }} not found")

		return
	}

	result.Success(ctx, {{ .StructNameLowerFirst }})
}

// @Summary 更新
// @Description 更新记录
// @Tags {{ .StructNameLowerFirst }}
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param {{ .StructNameLowerFirst }} body model.{{ .StructName }} true "信息"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/{id} [post]
func (c *{{ .StructName }}Handler) Update{{ .StructName }}(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		result.ServerError(ctx, "IDs are required")
		return
	}

	var {{ .StructNameLowerFirst }} model.{{ .StructName }}
	if err := ctx.ShouldBindJSON(&{{ .StructNameLowerFirst }}); err != nil {
		result.FailByErr(ctx, err)
		return
	}

	err := c.service.Update{{ .StructName }}(&{{ .StructNameLowerFirst }})
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, {{ .StructNameLowerFirst }})
}

// @Summary 删除
// @Description 根据ID删除记录
// @Tags {{ .StructNameLowerFirst }}
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/{id} [get]
func (c *{{ .StructName }}Handler) Delete{{ .StructName }}(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		result.ServerError(ctx, "IDs are required")
		return
	}

	err := c.service.Delete{{ .StructName }}(id)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, "")
}

// @Summary 列表
// @Description 获取列表
// @Tags {{ .StructNameLowerFirst }}
// @Produce  json
// @Param filter query map[string]interface{} false "过滤条件"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s [post]
func (c *{{ .StructName }}Handler) List{{ .StructName }}s(ctx *gin.Context) {
	filter := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		result.FailByErr(ctx, err)
		return
	}
	{{ .StructNameLowerFirst }}s, err := c.service.List{{ .StructName }}s(filter)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, {{ .StructNameLowerFirst }}s)
}

// @Summary 批量获取
// @Description 批量获取记录
// @Tags {{ .StructNameLowerFirst }}
// @Produce  json
// @Param ids query []string true "ID列表"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/batch [get]
func (c *{{ .StructName }}Handler) BatchGet{{ .StructName }}s(ctx *gin.Context) {
	var idsReq IdsRequest
	if err := ctx.ShouldBindJSON(&idsReq); err != nil {
		result.FailByErr(ctx, err)
		return
	}

	{{ .StructNameLowerFirst }}s, err := c.service.BatchGet{{ .StructName }}s(idsReq.ids)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, {{ .StructNameLowerFirst }}s)
}

// @Summary 批量创建
// @Description 批量创建记录
// @Tags {{ .StructNameLowerFirst }}
// @Accept  json
// @Produce  json
// @Param {{ .StructNameLowerFirst }}s body []model.{{ .StructName }} true "列表"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/batch [post]
func (c *{{ .StructName }}Handler) BatchCreate{{ .StructName }}s(ctx *gin.Context) {
	var {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}
	if err := ctx.ShouldBindJSON(&{{ .StructNameLowerFirst }}s); err != nil {
		result.FailByErr(ctx, err)
		return
	}

	err := c.service.BatchCreate{{ .StructName }}s({{ .StructNameLowerFirst }}s)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, {{ .StructNameLowerFirst }}s)
}

// @Summary 批量更新
// @Description 批量更新记录
// @Tags {{ .StructNameLowerFirst }}
// @Accept  json
// @Produce  json
// @Param {{ .StructNameLowerFirst }}s body []model.{{ .StructName }} true "列表"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/batch [post]
func (c *{{ .StructName }}Handler) BatchUpdate{{ .StructName }}s(ctx *gin.Context) {
	var {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}
	if err := ctx.ShouldBindJSON(&{{ .StructNameLowerFirst }}s); err != nil {
		result.FailByErr(ctx, err)
		return
	}

	err := c.service.BatchUpdate{{ .StructName }}s({{ .StructNameLowerFirst }}s)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, {{ .StructNameLowerFirst }}s)
}

// @Summary 批量删除
// @Description 批量删除记录
// @Tags {{ .StructNameLowerFirst }}
// @Produce  json
// @Param ids query []string true "ID列表"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /{{ .StructNameLowerFirst }}s/batch [post]
func (c *{{ .StructName }}Handler) BatchDelete{{ .StructName }}s(ctx *gin.Context) {
    var idsReq IdsRequest
    if err := ctx.ShouldBindJSON(&idsReq); err != nil {
        result.FailByErr(ctx, err)
        return
    }
	err := c.service.BatchDelete{{ .StructName }}s(idsReq.ids)
	if err != nil {
		result.FailByErr(ctx, err)
	} else {
		result.Success(ctx, "")
	}
}
