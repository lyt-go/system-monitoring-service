package handler

import (
	"net/http"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/httpx"
)

func (s *Server) registerAlertRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alerts", s.createAlert)
	mux.HandleFunc("GET /api/alerts", s.listAlerts)
	mux.HandleFunc("GET /api/alerts/{id}", s.getAlert)
	mux.HandleFunc("PUT /api/alerts/{id}", s.updateAlert)
	mux.HandleFunc("DELETE /api/alerts/{id}", s.deleteAlert)
	mux.HandleFunc("POST /api/alerts/batch-status", s.batchUpdateAlertStatus)
}

type createAlertRequest struct {
	MetricID    string `json:"metric_id"`
	ThresholdID string `json:"threshold_id"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}

func (s *Server) createAlert(w http.ResponseWriter, r *http.Request) {
	var req createAlertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAlert(model.Alert{MetricID: req.MetricID, ThresholdID: req.ThresholdID, Level: req.Level, Message: req.Message})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertFilter{
		MetricID: r.URL.Query().Get("metric_id"),
		Level:    r.URL.Query().Get("level"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAlerts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAlert(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAlertRequest struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (s *Server) updateAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAlert(id, model.Alert{Level: req.Level, Message: req.Message, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlert(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchUpdateAlertStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

func (s *Server) batchUpdateAlertStatus(w http.ResponseWriter, r *http.Request) {
	var req batchUpdateAlertStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchUpdateAlertStatus(req.IDs, req.Status); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
