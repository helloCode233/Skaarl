package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/service"
	"{{ .ProjectName }}/pkg/helper/result"
)
// @wire:Handler
type {{ .StructName }}Handler struct {
	*Handler
	{{ .StructNameLowerFirst }}Service service.{{ .StructName }}Service
}

type create{{ .StructName }}Request struct {
	// TODO: 添加请求字段
}

type update{{ .StructName }}Request struct {
	// TODO: 添加请求字段
}

type list{{ .StructName }}sRequest struct {
	// TODO: 添加分页和过滤字段
}

func New{{ .StructName }}Handler(
    handler *Handler,
    {{ .StructNameLowerFirst }}Service service.{{ .StructName }}Service,
) *{{ .StructName }}Handler {
	return &{{ .StructName }}Handler{
		Handler:      handler,
		{{ .StructNameLowerFirst }}Service: {{ .StructNameLowerFirst }}Service,
	}
}

func (h *{{ .StructName }}Handler) Create{{ .StructName }}(ctx *gin.Context) {
	var req create{{ .StructName }}Request
	if err := ctx.ShouldBind(&req); err != nil {
		result.FailByErr(ctx, err)
		return
	}
	// TODO: 实现创建逻辑
}

func (h *{{ .StructName }}Handler) Get{{ .StructName }}(ctx *gin.Context) {
	id, err := h.GetIDFromParam(ctx)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	res, err := h.{{ .StructNameLowerFirst }}Service.Get{{ .StructName }}(ctx.Request.Context(), id)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, res)
}

func (h *{{ .StructName }}Handler) Update{{ .StructName }}(ctx *gin.Context) {
	id, err := h.GetIDFromParam(ctx)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	var req update{{ .StructName }}Request
	if err := ctx.ShouldBind(&req); err != nil {
		result.FailByErr(ctx, err)
		return
	}
	// TODO: 实现更新逻辑
}

func (h *{{ .StructName }}Handler) Delete{{ .StructName }}(ctx *gin.Context) {
	id, err := h.GetIDFromParam(ctx)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	if err := h.{{ .StructNameLowerFirst }}Service.Delete{{ .StructName }}(ctx.Request.Context(), id); err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.Success(ctx, res)
}

func (h *{{ .StructName }}Handler) List{{ .StructName }}s(ctx *gin.Context) {
	var req list{{ .StructName }}sRequest
	if err := ctx.ShouldBind(&req); err != nil {
		result.FailByErr(ctx, err)
		return
	}
	// TODO: 实现列表查询逻辑
}

func (h *{{ .StructName }}Handler) BatchGet{{ .StructName }}s(ctx *gin.Context) {
	// TODO: 实现批量获取逻辑
}

func (h *{{ .StructName }}Handler) BatchCreate{{ .StructName }}s(ctx *gin.Context) {
	// TODO: 实现批量创建逻辑
}

func (h *{{ .StructName }}Handler) BatchUpdate{{ .StructName }}s(ctx *gin.Context) {
	// TODO: 实现批量更新逻辑
}

func (h *{{ .StructName }}Handler) BatchDelete{{ .StructName }}s(ctx *gin.Context) {
	// TODO: 实现批量删除逻辑
}
