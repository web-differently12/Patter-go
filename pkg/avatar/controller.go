package avatar

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lynxflow/patter-go/pkg/core"
)

type AvatarController struct {
	avatarService AvatarEngineService
}

func NewAvatarController(svc AvatarEngineService) *AvatarController {
	return &AvatarController{avatarService: svc}
}

// StreamRealtimeAvatar godoc
// @Summary      Générer un flux vidéo Avatar interactif en temps réel (WaveSpeed API)
// @Tags         Gateway - Avatar Engine
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body StreamAvatarRequest true "Paramètres de streaming Avatar"
// @Success      200 {object} core.APIResponse{data=AvatarResponse}
// @Router       /api/v1/gateway/avatar/stream [post]
func (ctrl *AvatarController) StreamRealtimeAvatar(c *gin.Context) {
	var req StreamAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.avatarService.StreamRealtimeAvatar(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}

// RenderMassAvatarVideo godoc
// @Summary      Rendu vidéo Avatar asynchrone en masse à bas coût (Local GPU Pool)
// @Tags         Gateway - Avatar Engine
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID header string true "ID du Tenant White-Label"
// @Param        request body RenderAvatarVideoRequest true "Paramètres de rendu masse"
// @Success      200 {object} core.APIResponse{data=AvatarResponse}
// @Router       /api/v1/gateway/avatar/render [post]
func (ctrl *AvatarController) RenderMassAvatarVideo(c *gin.Context) {
	var req RenderAvatarVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tenantID := core.GetTenantID(c)
	resp, err := ctrl.avatarService.RenderMassAvatarVideo(c.Request.Context(), tenantID, req)
	if err != nil {
		core.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	core.Success(c, resp)
}
