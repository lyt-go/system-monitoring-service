package handler

import (
	"net/http"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/httpx"
)

func (s *Server) registerThresholdRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/thresholds", s.createThreshold)
	mux.HandleFunc("GET /api/thresholds", s.listThresholds)
	mux.HandleFunc("GET /api/thresholds/{id}", s.getThreshold)
	mux.HandleFunc("PUT /api/thresholds/{id}", s.updateThreshold)
	mux.HandleFunc("DELETE /api/thresholds/{id}", s.deleteThreshold)
	mux.HandleFunc("PATCH /api/thresholds/{id}/status", s.transitionThresholdStatus)
}

type createThresholdRequest struct {
	MetricID string  `json:"metric_id"`
	Operator string  `json:"operator"`
	Value    float64 `json:"value"`
}

func (s *Server) createThreshold(w http.ResponseWriter, r *http.Request) {
	var req createThresholdRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateThreshold(model.Threshold{MetricID: req.MetricID, Operator: req.Operator, Value: req.Value})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listThresholds(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ThresholdFilter{
		MetricID: r.URL.Query().Get("metric_id"),
		Operator: r.URL.Query().Get("operator"),
		Status:   r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListThresholds(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getThreshold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetThreshold(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateThresholdRequest struct {
	MetricID string  `json:"metric_id"`
	Operator string  `json:"operator"`
	Value    float64 `json:"value"`
	Status   string  `json:"status"`
}

func (s *Server) updateThreshold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateThresholdRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateThreshold(id, model.Threshold{MetricID: req.MetricID, Operator: req.Operator, Value: req.Value, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteThreshold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteThreshold(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) transitionThresholdStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.TransitionThresholdStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}
