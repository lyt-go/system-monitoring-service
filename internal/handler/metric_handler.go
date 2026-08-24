package handler

import (
	"net/http"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/httpx"
)

func (s *Server) registerMetricRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/metrics", s.createMetric)
	mux.HandleFunc("GET /api/metrics", s.listMetrics)
	mux.HandleFunc("GET /api/metrics/{id}", s.getMetric)
	mux.HandleFunc("PUT /api/metrics/{id}", s.updateMetric)
	mux.HandleFunc("DELETE /api/metrics/{id}", s.deleteMetric)
	mux.HandleFunc("PATCH /api/metrics/{id}/status", s.transitionMetricStatus)
}

type createMetricRequest struct {
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (s *Server) createMetric(w http.ResponseWriter, r *http.Request) {
	var req createMetricRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMetric(model.Metric{Name: req.Name, Unit: req.Unit, Type: req.Type, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMetrics(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MetricFilter{
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMetrics(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMetric(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMetric(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type updateMetricRequest struct {
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) updateMetric(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMetricRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMetric(id, model.Metric{Name: req.Name, Unit: req.Unit, Type: req.Type, Description: req.Description, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMetric(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMetric(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionMetricStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.TransitionMetricStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}
