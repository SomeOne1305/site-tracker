package services

import (
	"context"
	"fmt"
	"visit-tracker/models"
	"visit-tracker/repository"

	"github.com/google/uuid"
)

type TrackingService struct {
	repo *repository.TrackingRepository
}

func NewTrackingService(repo *repository.TrackingRepository) *TrackingService {
	return &TrackingService{repo: repo}
}

// TrackPageView - Track visitor (PUBLIC)
func (s *TrackingService) TrackPageView(ctx context.Context, req models.TrackingRequest, ipAddr string) error {
	if req.ProjectID == "" || req.Path == "" || req.SessionID == "" {
		return fmt.Errorf("missing required fields")
	}

	return s.repo.TrackVisitor(ctx, req.ProjectID, req.Path, req.SessionID, ipAddr, req.UserAgent)
}

// GetLast7DaysChart - Last 7 days with growth (PROTECTED)
func (s *TrackingService) GetLast7DaysChart(ctx context.Context, projectID uuid.UUID) (*models.ChartData, error) {
	if projectID == uuid.Nil {
		return nil, fmt.Errorf("invalid project ID")
	}

	dailyData, err := s.repo.GetLast7DaysChart(ctx, projectID)
	if err != nil {
		return nil, err
	}

	chart := &models.ChartData{
		Type:       "daily",
		Period:     "last7days",
		DataPoints: []models.ChartDataPoint{},
		Summary:    models.ChartSummary{},
	}

	var totalViews, totalUnique int64
	var highestDay string
	var highestViews int64

	for _, d := range dailyData {
		chart.DataPoints = append(chart.DataPoints, models.ChartDataPoint{
			Label:  d.DayName,
			Value:  d.Pageviews,
			Date:   d.DateFormatted,
			Growth: d.Growth,
		})

		totalViews += d.Pageviews
		totalUnique += d.UniqueVisitors

		if d.Pageviews > highestViews {
			highestViews = d.Pageviews
			highestDay = d.DateFormatted
		}
	}

	if len(dailyData) > 0 {
		chart.Summary = models.ChartSummary{
			TotalViews:      totalViews,
			TotalUnique:     totalUnique,
			AveragePerDay:   totalViews / int64(len(dailyData)),
			HighestDay:      highestDay,
			HighestDayViews: highestViews,
		}

		if len(dailyData) > 1 {
			firstDay := dailyData[0].Pageviews
			lastDay := dailyData[len(dailyData)-1].Pageviews
			if firstDay > 0 {
				chart.Summary.GrowthPercent = float64(lastDay-firstDay) / float64(firstDay) * 100
			}
		}
	}

	return chart, nil
}

// GetLast12MonthsChart - Last 12 months with growth (PROTECTED)
func (s *TrackingService) GetLast12MonthsChart(ctx context.Context, projectID uuid.UUID) (*models.ChartData, error) {
	if projectID == uuid.Nil {
		return nil, fmt.Errorf("invalid project ID")
	}

	monthlyData, err := s.repo.GetLast12MonthsChart(ctx, projectID)
	if err != nil {
		return nil, err
	}

	chart := &models.ChartData{
		Type:       "monthly",
		Period:     "last12months",
		DataPoints: []models.ChartDataPoint{},
		Summary:    models.ChartSummary{},
	}

	var totalViews, totalUnique int64
	var highestMonth string
	var highestViews int64

	for _, m := range monthlyData {
		chart.DataPoints = append(chart.DataPoints, models.ChartDataPoint{
			Label:  m.MonthName,
			Value:  m.Pageviews,
			Date:   m.Month,
			Growth: m.Growth,
		})

		totalViews += m.Pageviews
		totalUnique += m.UniqueVisitors

		if m.Pageviews > highestViews {
			highestViews = m.Pageviews
			highestMonth = m.MonthName
		}
	}

	if len(monthlyData) > 0 {
		chart.Summary = models.ChartSummary{
			TotalViews:      totalViews,
			TotalUnique:     totalUnique,
			AveragePerDay:   totalViews / (int64(len(monthlyData)) * 30),
			HighestDay:      highestMonth,
			HighestDayViews: highestViews,
			GrowthPercent:   monthlyData[0].Growth,
		}
	}

	return chart, nil
}

// GetDashboardOverview - Complete dashboard (PROTECTED)
func (s *TrackingService) GetDashboardOverview(ctx context.Context, projectID uuid.UUID, projectName string) (*models.DashboardOverview, error) {
	if projectID == uuid.Nil {
		return nil, fmt.Errorf("invalid project ID")
	}

	totalViews, totalUnique, err := s.repo.GetTotalAllTime(ctx, projectID)
	if err != nil {
		return nil, err
	}

	currentMonth, err := s.repo.GetCurrentMonthStats(ctx, projectID)
	if err != nil {
		return nil, err
	}

	previousMonth, err := s.repo.GetPreviousMonthStats(ctx, projectID)
	if err != nil {
		return nil, err
	}

	currentWeek, err := s.repo.GetCurrentWeekStats(ctx, projectID)
	if err != nil {
		return nil, err
	}

	last7Days, err := s.GetLast7DaysChart(ctx, projectID)
	if err != nil {
		return nil, err
	}

	last12Months, err := s.GetLast12MonthsChart(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topPages, err := s.repo.GetTopPages(ctx, projectID, 5)
	if err != nil {
		return nil, err
	}

	comparison, err := s.repo.GetComparisonStats(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &models.DashboardOverview{
		ProjectID:              projectID,
		ProjectName:            projectName,
		TotalAllTime:           totalViews,
		UniqueAllTime:          totalUnique,
		CurrentMonthStats:      *currentMonth,
		PreviousMonthStats:     *previousMonth,
		CurrentWeekStats:       *currentWeek,
		Last7DaysChart:         *last7Days,
		Last12MonthsChart:      *last12Months,
		Top5Pages:              topPages,
		ComparisonCurrentMonth: comparison,
	}, nil
}

// GetTrackingRepo - Expose tracking repo for admin operations
func (s *TrackingService) GetTrackingRepo() *repository.TrackingRepository {
	return s.repo
}
