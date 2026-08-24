package handler

import (
	"net/http"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/httpx"
)

func (s *Server) registerCollectorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/collectors", s.createCollector)
	mux.HandleFunc("GET /api/collectors", s.listCollectors)
	mux.HandleFunc("GET /api/collectors/{id}", s.getCollector)
	mux.HandleFunc("PUT /api/collectors/{id}", s.updateCollector)
	mux.HandleFunc("DELETE /api/collectors/{id}", s.deleteCollector)
	mux.HandleFunc("PATCH /api/collectors/{id}/status", s.transitionCollectorStatus)
}

type createCollectorRequest struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	IntervalSec int    `json:"interval_sec"`
}

func (s *Server) createCollector(w http.ResponseWriter, r *http.Request) {
	var req createCollectorRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCollector(model.Collector{Name: req.Name, Host: req.Host, IntervalSec: req.IntervalSec})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCollectors(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CollectorFilter{
		Status:  r.URL.Query().Get("status"),
		Host:    r.URL.Query().Get("host"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCollectors(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCollector(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateCollectorRequest struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	IntervalSec int    `json:"interval_sec"`
	Status      string `json:"status"`
}

func (s *Server) updateCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCollectorRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCollector(id, model.Collector{Name: req.Name, Host: req.Host, IntervalSec: req.IntervalSec, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCollector(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCollector(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) transitionCollectorStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req transitionStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.TransitionCollectorStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
