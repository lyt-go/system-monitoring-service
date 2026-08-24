package handler

import (
	"net/http"
	"strconv"

	"sysmonitor/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getStatsOverview)
	mux.HandleFunc("GET /api/stats/samples-by-metric", s.getSampleCountByMetric)
	mux.HandleFunc("GET /api/stats/samples-by-host", s.getSampleCountByHost)
	mux.HandleFunc("GET /api/stats/top-alert-metrics", s.getTopAlertMetrics)
}

func (s *Server) getStatsOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := s.svc.GetStatsOverview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, overview)
}

func (s *Server) getSampleCountByMetric(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.GetSampleCountByMetric()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) getSampleCountByHost(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.GetSampleCountByHost()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) getTopAlertMetrics(w http.ResponseWriter, r *http.Request) {
	n := 5
	q := r.URL.Query()
	if v := q.Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			httpx.BadRequest(w, "limit: 必须为整数")
			return
		}
		n = parsed
	} else if v := q.Get("n"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			httpx.BadRequest(w, "limit: 必须为整数")
			return
		}
		n = parsed
	}
	result, err := s.svc.GetTopAlertMetrics(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
