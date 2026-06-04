package models

import (
	"time"

	"github.com/google/uuid"
)

// ===== REQUESTS =====
type TrackingRequest struct {
	ProjectID string `json:"project_id" form:"project_id" query:"project_id" validate:"required,uuid"`
	Path      string `json:"path" form:"path" query:"path" validate:"required"`
	SessionID string `json:"session_id" form:"session_id" query:"session_id" validate:"required"`
	UserAgent string `json:"user_agent" form:"user_agent" query:"user_agent"`
	Referrer  string `json:"referrer" form:"referrer" query:"referrer"`
}

// ===== CHART RESPONSES =====
type ChartDataPoint struct {
	Label  string  `json:"label"`
	Value  int64   `json:"value"`
	Date   string  `json:"date,omitempty"`
	Growth float64 `json:"growth,omitempty"`
}

type ChartSummary struct {
	TotalViews      int64   `json:"total_views"`
	TotalUnique     int64   `json:"total_unique"`
	AveragePerDay   int64   `json:"average_per_day"`
	HighestDay      string  `json:"highest_day"`
	HighestDayViews int64   `json:"highest_day_views"`
	GrowthPercent   float64 `json:"growth_percent"`
}

type ChartData struct {
	Type       string           `json:"type"`
	Period     string           `json:"period"`
	DataPoints []ChartDataPoint `json:"data_points"`
	Summary    ChartSummary     `json:"summary"`
}

// ===== DAILY DATA =====
type DailyChartData struct {
	Date           time.Time `json:"date"`
	DateFormatted  string    `json:"date_formatted"`
	Pageviews      int64     `json:"pageviews"`
	UniqueVisitors int64     `json:"unique_visitors"`
	DayName        string    `json:"day_name"`
	Growth         float64   `json:"growth_percent"`
}

// ===== MONTHLY DATA =====
type MonthlyChartData struct {
	Month             string  `json:"month"`
	MonthName         string  `json:"month_name"`
	Pageviews         int64   `json:"pageviews"`
	UniqueVisitors    int64   `json:"unique_visitors"`
	AvgDailyPageviews int64   `json:"avg_daily_pageviews"`
	Growth            float64 `json:"growth_percent"`
}

// ===== WEEKLY DATA =====
type WeeklyChartData struct {
	Week              string  `json:"week"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
	Pageviews         int64   `json:"pageviews"`
	UniqueVisitors    int64   `json:"unique_visitors"`
	AvgDailyPageviews int64   `json:"avg_daily_pageviews"`
	Growth            float64 `json:"growth_percent"`
}

// ===== COMPARISON =====
type ComparisonStats struct {
	Current  int64   `json:"current"`
	Previous int64   `json:"previous"`
	Change   int64   `json:"change"`
	Percent  float64 `json:"percent"`
}

// ===== DASHBOARD =====
type DashboardOverview struct {
	ProjectID              uuid.UUID        `json:"project_id"`
	ProjectName            string           `json:"project_name"`
	TotalAllTime           int64            `json:"total_all_time"`
	UniqueAllTime          int64            `json:"unique_all_time"`
	CurrentMonthStats      MonthlyChartData `json:"current_month_stats"`
	PreviousMonthStats     MonthlyChartData `json:"previous_month_stats"`
	CurrentWeekStats       WeeklyChartData  `json:"current_week_stats"`
	Last7DaysChart         ChartData        `json:"last_7_days_chart"`
	Last12MonthsChart      ChartData        `json:"last_12_months_chart"`
	Top5Pages              []PageStats      `json:"top_5_pages"`
	ComparisonCurrentMonth *ComparisonStats `json:"comparison_current_month"`
}
