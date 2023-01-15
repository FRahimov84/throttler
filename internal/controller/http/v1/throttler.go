package v1

import (
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
		h.GET("/:UUID", r.getRequestByID)
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
// @Router      /throttler [post]
func (r *throttlerroutes) newRequest(c *gin.Context) {
	// TODO: parse request structure from body
	uuid, err := r.t.NewRequest(c.Request.Context(), entity.Request{Status: "new"})
	if err != nil {
		if err == entity.RequestStatusErr {
			r.l.Info("http - v1 - newRequest - uc.NewRequest", zap.Error(err))
			errorResponse(c, http.StatusBadRequest, "bad request status")
			return
		}

		r.l.Error("http - v1 - newRequest - uc.NewRequest", zap.Error(err))
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
func (r *throttlerroutes) getRequestByID(c *gin.Context) {
	reqUUID, err := entity.ParseUUID(c.Param("UUID"))
	if err != nil {
		r.l.Info("http - v1 - getRequestByID - ParseUUID", zap.Any("uuid", c.Param("UUID")), zap.Error(err))
		errorResponse(c, http.StatusBadRequest, "bad request uuid")

		return
	}

	request, err := r.t.GetRequestByID(c.Request.Context(), reqUUID)
	if err != nil {
		if err == entity.RepoNotFoundErr {
			r.l.Info("http - v1 - getRequestByID - uc.GetRequestByID", zap.Error(err))
			errorResponse(c, http.StatusNotFound, "request not found")
			return
		}
		r.l.Error("http - v1 - getRequestByID - uc.GetRequestByID", zap.Error(err))
		errorResponse(c, http.StatusInternalServerError, "db error")
		return
	}

	c.JSON(http.StatusOK, getReqResp{request})
}
