package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/FRahimov84/throttler/internal/entity"
	"github.com/FRahimov84/throttler/internal/usecase"
)

type throttlerroutes struct {
	t usecase.Throttler
	l *zap.Logger
}

func newThrottlerRoutes(handler *gin.RouterGroup, t usecase.Throttler, l *zap.Logger) {
	r := &throttlerroutes{t, l}

	h := handler.Group("/throttler")
	{
		h.POST("", r.newRequest)
		h.GET("/:REQ-UUID", r.getRequestByUUID)
	}
}

type newReqResp struct {
	UUID string `json:"uuid"`
}

// newRequest
// @Summary     New request
// @Description Add new request for external svc
// @ID          NewReq
// @Tags  	    throttler
// @Accept      json
// @Produce     json
// @Success     200 {object} newReqResp
// @Failure     500 {object} response
// @Router      /translation [post]
func (r *throttlerroutes) newRequest(c *gin.Context) {
	uuid, err := r.t.NewRequest(c.Request.Context())
	if err != nil {
		r.l.Error("http - v1 - newRequest", zap.Error(err))
		errorResponse(c, http.StatusInternalServerError, "db error")
		return
	}

	c.JSON(http.StatusOK, newReqResp{uuid.String()})
}

type getReqResp struct {
	Req entity.Request `json:"request"`
}

// getRequestByUUID
// @Summary     Request By UUID
// @Description Return request by uuid
// @ID          GetRequest
// @Tags  	    throttler
// @Param		uuid	path	string	true	"request uuid"
// @Produce     json
// @Success     200 {object} getReqResp
// @Failure     500 {object} response
// @Router      /throttler/{uuid} [get]
func (r *throttlerroutes) getRequestByUUID(c *gin.Context) {
	reqUUID := entity.ParseUUID(c.Param("REQ-UUID"))
	if reqUUID == entity.EmptyUUID {
		r.l.Error("http - v1 - getRequestByUUID", zap.Error(fmt.Errorf("bad uuid in request")))
		errorResponse(c, http.StatusBadRequest, "invalid request uuid")

		return
	}

	request, err := r.t.RequestByUUID(c.Request.Context(), reqUUID)
	if err != nil {
		r.l.Error("http - v1 - getRequestByUUID", zap.Error(err))
		errorResponse(c, http.StatusInternalServerError, "db error")
		return
	}

	c.JSON(http.StatusOK, getReqResp{request})
}
