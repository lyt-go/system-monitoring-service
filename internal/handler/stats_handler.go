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
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	if n <= 0 {
		n = 5
	}
	result, err := s.svc.GetTopAlertMetrics(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
