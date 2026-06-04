package repository

import (
	"context"
	"fmt"
	"visit-tracker/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackingRepository struct {
	pool *pgxpool.Pool
}

func NewTrackingRepository(pool *pgxpool.Pool) *TrackingRepository {
	return &TrackingRepository{pool: pool}
}

// TrackVisitor - Insert or increment visit count for a (path, session) pair.
// One row per unique (path_id, session_id). visit_count increments on every hit.
func (r *TrackingRepository) TrackVisitor(ctx context.Context, projectID, pathStr, sessionID, ipAddr, userAgent string) error {
	projID, err := uuid.Parse(projectID)
	if err != nil {
		return fmt.Errorf("invalid project id: %w", err)
	}

	query := `
		WITH project AS (
			SELECT id FROM projects WHERE id = $1
		),
		path_lookup AS (
			SELECT id FROM paths
			WHERE project_id = (SELECT id FROM project)
			AND path = $2
		),
		ensure_path AS (
			INSERT INTO paths (id, path, project_id, visit_count)
			SELECT gen_random_uuid(), $2, (SELECT id FROM project), 0
			WHERE NOT EXISTS (SELECT 1 FROM path_lookup)
			RETURNING id
		),
		resolved_path AS (
			SELECT id FROM path_lookup
			UNION ALL
			SELECT id FROM ensure_path
			LIMIT 1
		)
		INSERT INTO visitors (
			path_id,
			session_id,
			ip_address,
			user_agent,
			visit_time,
			first_visit,
			last_visit,
			visit_count
		)
		SELECT
			(SELECT id FROM resolved_path),
			$3,
			$4,
			$5,
			NOW(),
			NOW(),
			NOW(),
			1
		ON CONFLICT (path_id, session_id)
		DO UPDATE SET
			last_visit   = NOW(),
			visit_count  = visitors.visit_count + 1
	`
	_, err = r.pool.Exec(ctx, query, projID, pathStr, sessionID, ipAddr, userAgent)
	if err != nil {
		return fmt.Errorf("track visitor failed: %w", err)
	}
	return nil
}

// GetLast7DaysChart - Live daily stats for the last 7 days directly from visitors.
// Uses generate_series so days with zero visits still appear.
func (r *TrackingRepository) GetLast7DaysChart(ctx context.Context, projectID uuid.UUID) ([]models.DailyChartData, error) {
	query := `
		WITH days AS (
			SELECT generate_series(
				CURRENT_DATE - INTERVAL '6 days',
				CURRENT_DATE,
				'1 day'::interval
			)::date AS day
		),
		daily AS (
			SELECT
				DATE(v.first_visit)             AS day,
				SUM(v.visit_count)              AS pageviews,
				COUNT(DISTINCT v.session_id)    AS unique_visitors
			FROM visitors v
			JOIN paths p ON v.path_id = p.id
			WHERE p.project_id = $1
			  AND DATE(v.first_visit) >= CURRENT_DATE - INTERVAL '6 days'
			GROUP BY DATE(v.first_visit)
		),
		filled AS (
			SELECT
				d.day,
				COALESCE(dv.pageviews, 0)        AS pageviews,
				COALESCE(dv.unique_visitors, 0)  AS unique_visitors,
				LAG(COALESCE(dv.pageviews, 0)) OVER (ORDER BY d.day) AS prev_pageviews
			FROM days d
			LEFT JOIN daily dv ON d.day = dv.day
		)
		SELECT
			day,
			TO_CHAR(day, 'YYYY-MM-DD')   AS date_formatted,
			pageviews,
			unique_visitors,
			TO_CHAR(day, 'Dy')           AS day_name,
			CASE
				WHEN prev_pageviews IS NULL OR prev_pageviews = 0 THEN 0::numeric
				ELSE ROUND(((pageviews - prev_pageviews)::numeric / prev_pageviews::numeric * 100), 2)
			END AS growth
		FROM filled
		ORDER BY day ASC
	`
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("get last 7 days chart failed: %w", err)
	}
	defer rows.Close()

	var results []models.DailyChartData
	for rows.Next() {
		var d models.DailyChartData
		if err := rows.Scan(&d.Date, &d.DateFormatted, &d.Pageviews, &d.UniqueVisitors, &d.DayName, &d.Growth); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
}

// GetLast12MonthsChart - Live monthly stats for the last 12 months directly from visitors.
// Uses generate_series so months with zero visits still appear.
func (r *TrackingRepository) GetLast12MonthsChart(ctx context.Context, projectID uuid.UUID) ([]models.MonthlyChartData, error) {
	query := `
		WITH months AS (
			SELECT generate_series(
				DATE_TRUNC('month', CURRENT_DATE - INTERVAL '11 months'),
				DATE_TRUNC('month', CURRENT_DATE),
				'1 month'::interval
			)::date AS month_start
		),
		monthly AS (
			SELECT
				DATE_TRUNC('month', v.first_visit)::date  AS month_start,
				SUM(v.visit_count)                        AS pageviews,
				COUNT(DISTINCT v.session_id)              AS unique_visitors,
				COUNT(DISTINCT DATE(v.first_visit))       AS active_days
			FROM visitors v
			JOIN paths p ON v.path_id = p.id
			WHERE p.project_id = $1
			  AND v.first_visit >= DATE_TRUNC('month', CURRENT_DATE - INTERVAL '11 months')
			GROUP BY DATE_TRUNC('month', v.first_visit)
		),
		filled AS (
			SELECT
				m.month_start,
				COALESCE(mv.pageviews, 0)       AS pageviews,
				COALESCE(mv.unique_visitors, 0) AS unique_visitors,
				COALESCE(mv.active_days, 0)     AS active_days,
				LAG(COALESCE(mv.pageviews, 0)) OVER (ORDER BY m.month_start) AS prev_pageviews
			FROM months m
			LEFT JOIN monthly mv ON m.month_start = mv.month_start
		)
		SELECT
			TO_CHAR(month_start, 'YYYY-MM')                AS month,
			TO_CHAR(month_start, 'Mon YYYY')               AS month_name,
			pageviews,
			unique_visitors,
			CASE
				WHEN active_days > 0
				THEN ROUND(pageviews::numeric / active_days, 0)::bigint
				ELSE 0
			END                                            AS avg_daily,
			CASE
				WHEN prev_pageviews IS NULL OR prev_pageviews = 0 THEN 0::numeric
				ELSE ROUND(((pageviews - prev_pageviews)::numeric / prev_pageviews::numeric * 100), 2)
			END                                            AS growth
		FROM filled
		ORDER BY month_start ASC
	`
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("get last 12 months chart failed: %w", err)
	}
	defer rows.Close()

	var results []models.MonthlyChartData
	for rows.Next() {
		var m models.MonthlyChartData
		if err := rows.Scan(&m.Month, &m.MonthName, &m.Pageviews, &m.UniqueVisitors, &m.AvgDailyPageviews, &m.Growth); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

// GetCurrentMonthStats - Live stats for the current calendar month.
func (r *TrackingRepository) GetCurrentMonthStats(ctx context.Context, projectID uuid.UUID) (*models.MonthlyChartData, error) {
	query := `
		SELECT
			TO_CHAR(CURRENT_DATE, 'YYYY-MM')                        AS month,
			TO_CHAR(CURRENT_DATE, 'Mon YYYY')                       AS month_name,
			COALESCE(SUM(v.visit_count), 0)                         AS pageviews,
			COUNT(DISTINCT v.session_id)                            AS unique_visitors,
			CASE
				WHEN COUNT(DISTINCT DATE(v.first_visit)) > 0
				THEN ROUND(COALESCE(SUM(v.visit_count), 0)::numeric
				     / COUNT(DISTINCT DATE(v.first_visit)), 0)::bigint
				ELSE 0
			END                                                     AS avg_daily,
			0::numeric                                              AS growth
		FROM visitors v
		JOIN paths p ON v.path_id = p.id
		WHERE p.project_id = $1
		  AND v.first_visit >= DATE_TRUNC('month', CURRENT_DATE)
	`
	var m models.MonthlyChartData
	err := r.pool.QueryRow(ctx, query, projectID).Scan(
		&m.Month, &m.MonthName, &m.Pageviews, &m.UniqueVisitors, &m.AvgDailyPageviews, &m.Growth,
	)
	if err != nil {
		return nil, fmt.Errorf("get current month stats failed: %w", err)
	}
	return &m, nil
}

// GetPreviousMonthStats - Live stats for the previous calendar month.
func (r *TrackingRepository) GetPreviousMonthStats(ctx context.Context, projectID uuid.UUID) (*models.MonthlyChartData, error) {
	query := `
		SELECT
			TO_CHAR(DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month'), 'YYYY-MM')     AS month,
			TO_CHAR(DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month'), 'Mon YYYY')    AS month_name,
			COALESCE(SUM(v.visit_count), 0)                                                AS pageviews,
			COUNT(DISTINCT v.session_id)                                                   AS unique_visitors,
			CASE
				WHEN COUNT(DISTINCT DATE(v.first_visit)) > 0
				THEN ROUND(COALESCE(SUM(v.visit_count), 0)::numeric
				     / COUNT(DISTINCT DATE(v.first_visit)), 0)::bigint
				ELSE 0
			END                                                                            AS avg_daily,
			0::numeric                                                                     AS growth
		FROM visitors v
		JOIN paths p ON v.path_id = p.id
		WHERE p.project_id = $1
		  AND v.first_visit >= DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')
		  AND v.first_visit  < DATE_TRUNC('month', CURRENT_DATE)
	`
	var m models.MonthlyChartData
	err := r.pool.QueryRow(ctx, query, projectID).Scan(
		&m.Month, &m.MonthName, &m.Pageviews, &m.UniqueVisitors, &m.AvgDailyPageviews, &m.Growth,
	)
	if err != nil {
		return nil, fmt.Errorf("get previous month stats failed: %w", err)
	}
	return &m, nil
}

// GetCurrentWeekStats - Live stats for the current ISO week (Mon–Sun).
func (r *TrackingRepository) GetCurrentWeekStats(ctx context.Context, projectID uuid.UUID) (*models.WeeklyChartData, error) {
	query := `
		SELECT
			TO_CHAR(DATE_TRUNC('week', CURRENT_DATE), 'YYYY-"W"WW')                          AS week,
			TO_CHAR(DATE_TRUNC('week', CURRENT_DATE), 'YYYY-MM-DD')                           AS start_date,
			TO_CHAR(DATE_TRUNC('week', CURRENT_DATE) + INTERVAL '6 days', 'YYYY-MM-DD')       AS end_date,
			COALESCE(SUM(v.visit_count), 0)                                                   AS pageviews,
			COUNT(DISTINCT v.session_id)                                                      AS unique_visitors,
			CASE
				WHEN COUNT(DISTINCT DATE(v.first_visit)) > 0
				THEN ROUND(COALESCE(SUM(v.visit_count), 0)::numeric
				     / COUNT(DISTINCT DATE(v.first_visit)), 0)::bigint
				ELSE 0
			END                                                                               AS avg_daily,
			0::numeric                                                                        AS growth
		FROM visitors v
		JOIN paths p ON v.path_id = p.id
		WHERE p.project_id = $1
		  AND v.first_visit >= DATE_TRUNC('week', CURRENT_DATE)
	`
	var w models.WeeklyChartData
	err := r.pool.QueryRow(ctx, query, projectID).Scan(
		&w.Week, &w.StartDate, &w.EndDate, &w.Pageviews, &w.UniqueVisitors, &w.AvgDailyPageviews, &w.Growth,
	)
	if err != nil {
		return nil, fmt.Errorf("get current week stats failed: %w", err)
	}
	return &w, nil
}

// GetTotalAllTime - Lifetime totals across all dates.
func (r *TrackingRepository) GetTotalAllTime(ctx context.Context, projectID uuid.UUID) (int64, int64, error) {
	query := `
		SELECT
			COALESCE(SUM(v.visit_count), 0) AS total_views,
			COUNT(DISTINCT v.session_id)    AS unique_visitors
		FROM visitors v
		JOIN paths p ON v.path_id = p.id
		WHERE p.project_id = $1
	`
	var totalViews, totalUnique int64
	err := r.pool.QueryRow(ctx, query, projectID).Scan(&totalViews, &totalUnique)
	if err != nil {
		return 0, 0, fmt.Errorf("get total all time failed: %w", err)
	}
	return totalViews, totalUnique, nil
}

// GetTopPages - Top N pages by pageviews in the last 30 days.
func (r *TrackingRepository) GetTopPages(ctx context.Context, projectID uuid.UUID, limit int) ([]models.PageStats, error) {
	query := `
		SELECT
			p.path,
			COALESCE(SUM(v.visit_count), 0) AS pageviews,
			COUNT(DISTINCT v.session_id)    AS unique_visitors
		FROM visitors v
		JOIN paths p ON v.path_id = p.id
		WHERE p.project_id = $1
		  AND v.first_visit >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY p.path
		ORDER BY pageviews DESC
		LIMIT $2
	`
	rows, err := r.pool.Query(ctx, query, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("get top pages failed: %w", err)
	}
	defer rows.Close()

	var results []models.PageStats
	for rows.Next() {
		var pg models.PageStats
		if err := rows.Scan(&pg.Path, &pg.Pageviews, &pg.UniqueVisitors); err != nil {
			return nil, err
		}
		results = append(results, pg)
	}
	return results, nil
}

// GetComparisonStats - Current month vs previous month pageviews.
func (r *TrackingRepository) GetComparisonStats(ctx context.Context, projectID uuid.UUID) (*models.ComparisonStats, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE
				WHEN v.first_visit >= DATE_TRUNC('month', CURRENT_DATE)
				THEN v.visit_count ELSE 0
			END), 0) AS current_month,
			COALESCE(SUM(CASE
				WHEN v.first_visit >= DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')
				 AND v.first_visit  < DATE_TRUNC('month', CURRENT_DATE)
				THEN v.visit_count ELSE 0
			END), 0) AS previous_month
		FROM visitors v
		JOIN paths p ON v.path_id = p.id
		WHERE p.project_id = $1
		  AND v.first_visit >= DATE_TRUNC('month', CURRENT_DATE - INTERVAL '1 month')
	`
	var current, previous int64
	err := r.pool.QueryRow(ctx, query, projectID).Scan(&current, &previous)
	if err != nil {
		return nil, fmt.Errorf("get comparison stats failed: %w", err)
	}

	comparison := &models.ComparisonStats{
		Current:  current,
		Previous: previous,
		Change:   current - previous,
	}
	if previous > 0 {
		comparison.Percent = float64(comparison.Change) / float64(previous) * 100
	}
	return comparison, nil
}

// AggregateDaily - No longer needed. All stats are computed live from the visitors table.
// Kept so the handler/admin route compiles without changes. Safe to call — it's a no-op.
func (r *TrackingRepository) AggregateDaily(ctx context.Context) error {
	return nil
}
