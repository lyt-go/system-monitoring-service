package handler

import (
	"net/http"
	"time"

	"sysmonitor/internal/model"
	"sysmonitor/pkg/httpx"
)

func (s *Server) registerSampleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/samples", s.createSample)
	mux.HandleFunc("GET /api/samples", s.listSamples)
	mux.HandleFunc("GET /api/samples/{id}", s.getSample)
	mux.HandleFunc("DELETE /api/samples/{id}", s.deleteSample)
	mux.HandleFunc("POST /api/samples/batch-delete", s.batchDeleteSamples)
}

type createSampleRequest struct {
	MetricID  string    `json:"metric_id"`
	Value     float64   `json:"value"`
	Host      string    `json:"host"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Server) createSample(w http.ResponseWriter, r *http.Request) {
	var req createSampleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sa, err := s.svc.CreateSample(model.Sample{MetricID: req.MetricID, Value: req.Value, Host: req.Host, Timestamp: req.Timestamp})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sa)
}

func (s *Server) listSamples(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SampleFilter{
		MetricID: r.URL.Query().Get("metric_id"),
		Host:     r.URL.Query().Get("host"),
	}
	items, total, err := s.svc.ListSamples(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sa, err := s.svc.GetSample(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sa)
}

func (s *Server) deleteSample(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSample(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchDeleteSamplesRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteSamples(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteSamplesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchDeleteSamples(req.IDs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
