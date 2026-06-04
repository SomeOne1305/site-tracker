package handlers

import (
	"fmt"
	"sync"
	"visit-tracker/models"
	"visit-tracker/services"
	"visit-tracker/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type TrackingHandler struct {
	service     *services.TrackingService
	scriptCache map[string]string
	scriptMutex sync.RWMutex
}

func NewTrackingHandler(service *services.TrackingService) *TrackingHandler {
	return &TrackingHandler{
		service:     service,
		scriptCache: make(map[string]string),
	}
}

// getMinifiedCode - Generate minified tracking code
func (h *TrackingHandler) getMinifiedCode(projectID, host string) string {
	return fmt.Sprintf(`!function(){let a=localStorage.getItem('s%s');if(!a){const b=[navigator.userAgent,navigator.language,screen.width+'x'+screen.height].join('|');a=btoa(b).substring(0,32);localStorage.setItem('s%s',a)}function c(){const d=new URLSearchParams({project_id:'%s',path:window.location.pathname,session_id:a,user_agent:navigator.userAgent,referrer:document.referrer});fetch('http://%s/api/track',{method:'POST',headers:{'Content-Type':'application/x-www-form-urlencoded'},body:d,keepalive:true}).catch(()=>{})}document.readyState==='loading'?document.addEventListener('DOMContentLoaded',c):c();window.addEventListener('popstate',c);const e=window.history.pushState;window.history.pushState=function(){e.apply(window.history,arguments);c()}}();`,
		projectID, projectID, projectID, host)
}

// GetEmbedScript - GET /embed/:projectID - Returns <script>code</script>
func (h *TrackingHandler) GetEmbedScript(c fiber.Ctx) error {
	projectID := c.Params("projectID")

	if _, err := uuid.Parse(projectID); err != nil {
		return c.Status(fiber.StatusNotFound).SendString("")
	}

	h.scriptMutex.RLock()
	if cached, exists := h.scriptCache[projectID]; exists {
		h.scriptMutex.RUnlock()
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=86400")
		return c.SendString(cached)
	}
	h.scriptMutex.RUnlock()

	code := h.getMinifiedCode(projectID, c.Hostname())
	scriptTag := fmt.Sprintf(`<script>%s</script>`, code)

	h.scriptMutex.Lock()
	h.scriptCache[projectID] = scriptTag
	h.scriptMutex.Unlock()

	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Set("Cache-Control", "public, max-age=86400")
	return c.SendString(scriptTag)
}

// GetScriptFile - GET /t/:projectID.js - Returns JavaScript file
func (h *TrackingHandler) GetScriptFile(c fiber.Ctx) error {
	projectID := c.Params("projectID")

	if _, err := uuid.Parse(projectID); err != nil {
		return c.Status(fiber.StatusNotFound).SendString("")
	}

	code := h.getMinifiedCode(projectID, c.Hostname())

	c.Set("Content-Type", "application/javascript; charset=utf-8")
	c.Set("Cache-Control", "public, max-age=86400")
	return c.SendString(code)
}

// TrackPageView - POST /api/track - Track visitor (PUBLIC - OPEN)
func (h *TrackingHandler) TrackPageView(c fiber.Ctx) error {
	var req models.TrackingRequest
	fmt.Print("Request come")
	if c.Get("Content-Type") == "application/json" {
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}
	} else {
		if err := c.Bind().Body(&req); err != nil {
			if err := c.Bind().Query(&req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
			}
		}
	}

	if req.ProjectID == "" || req.Path == "" || req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing fields"})
	}

	ipAddr := c.IP()
	err := h.service.TrackPageView(c.Context(), req, ipAddr)
	if err != nil {
		fmt.Println(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetDashboard - GET /api/v1/projects/:projectID/dashboard (PROTECTED)
func (h *TrackingHandler) GetDashboard(c fiber.Ctx) error {
	projectIDStr := c.Params("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project ID"})
	}

	projectName := c.Query("name", "Project")

	overview, err := h.service.GetDashboardOverview(c.Context(), projectID, projectName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(overview)
}

// GetChartDaily - GET /api/v1/projects/:projectID/chart/daily (PROTECTED)
func (h *TrackingHandler) GetChartDaily(c fiber.Ctx) error {
	projectIDStr := c.Params("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project ID"})
	}

	chart, err := h.service.GetLast7DaysChart(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(chart)
}

// GetChartMonthly - GET /api/v1/projects/:projectID/chart/monthly (PROTECTED)
func (h *TrackingHandler) GetChartMonthly(c fiber.Ctx) error {
	projectIDStr := c.Params("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project ID"})
	}

	chart, err := h.service.GetLast12MonthsChart(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(chart)
}

// TriggerAggregation - POST /api/v1/admin/aggregate (PROTECTED - ADMIN)
func (h *TrackingHandler) TriggerAggregation(c fiber.Ctx) error {
	err := h.service.GetTrackingRepo().AggregateDaily(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "aggregation triggered successfully"})
}

func (h *TrackingHandler) GetAPIEndpoints(c fiber.Ctx) error {
	projectID := c.Params("projectID")
	if _, err := uuid.Parse(projectID); err != nil {
		return c.Status(fiber.StatusNotFound).SendString("")
	}
	apiKey, _ := utils.EncryptID(projectID)
	endpoints := fiber.Map{
		"daily_chart":   fmt.Sprintf("/api/v1/projects/%s/chart/daily", projectID),
		"monthly_chart": fmt.Sprintf("/api/v1/projects/%s/chart/monthly", projectID),
		"dashboard":     fmt.Sprintf("/api/v1/projects/%s/dashboard", projectID),
		"api_key":       apiKey,
	}
	return c.JSON(endpoints)
}
