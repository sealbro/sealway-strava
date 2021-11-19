// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AthleteToken struct {
	AthleteID int64  `json:"athlete_id"`
	Refresh   string `json:"refresh"`
}

type DetailedActivity struct {
	ID                   int64                    `json:"id"`
	ExternalID           string                   `json:"external_id"`
	UploadID             int64                    `json:"upload_id"`
	Athlete              *MetaAthlete             `json:"athlete"`
	Name                 string                   `json:"name"`
	Distance             float64                  `json:"distance"`
	MovingTime           int64                    `json:"moving_time"`
	ElapsedTime          int64                    `json:"elapsed_time"`
	TotalElevationGain   float64                  `json:"total_elevation_gain"`
	ElevHigh             float64                  `json:"elev_high"`
	ElevLow              float64                  `json:"elev_low"`
	Type                 *string                  `json:"type"`
	StartDate            time.Time                `json:"start_date"`
	StartDateLocal       time.Time                `json:"start_date_local"`
	Timezone             string                   `json:"timezone"`
	StartLatlng          []float64                `json:"start_latlng"`
	EndLatlng            []float64                `json:"end_latlng"`
	AchievementCount     int64                    `json:"achievement_count"`
	KudosCount           int64                    `json:"kudos_count"`
	CommentCount         int64                    `json:"comment_count"`
	AthleteCount         int64                    `json:"athlete_count"`
	PhotoCount           int64                    `json:"photo_count"`
	TotalPhotoCount      int64                    `json:"total_photo_count"`
	Map                  *PolylineMap             `json:"map"`
	Trainer              bool                     `json:"trainer"`
	Commute              bool                     `json:"commute"`
	Manual               bool                     `json:"manual"`
	Private              bool                     `json:"private"`
	Flagged              bool                     `json:"flagged"`
	WorkoutType          int64                    `json:"workout_type"`
	UploadIDStr          string                   `json:"upload_id_str"`
	AverageSpeed         float64                  `json:"average_speed"`
	MaxSpeed             float64                  `json:"max_speed"`
	HasKudoed            bool                     `json:"has_kudoed"`
	GearID               string                   `json:"gear_id"`
	Kilojoules           float64                  `json:"kilojoules"`
	AverageWatts         float64                  `json:"average_watts"`
	DeviceWatts          bool                     `json:"device_watts"`
	MaxWatts             int64                    `json:"max_watts"`
	WeightedAverageWatts int64                    `json:"weighted_average_watts"`
	Description          string                   `json:"description"`
	Photos               *PhotosSummary           `json:"photos"`
	Gear                 *SummaryGear             `json:"gear"`
	Calories             float64                  `json:"calories"`
	SegmentEfforts       []*DetailedSegmentEffort `json:"segment_efforts"`
	DeviceName           string                   `json:"device_name"`
	EmbedToken           string                   `json:"embed_token"`
	SplitsMetric         []*Split                 `json:"splits_metric"`
	SplitsStandard       []*Split                 `json:"splits_standard"`
	Laps                 []*Lap                   `json:"laps"`
	BestEfforts          []*DetailedSegmentEffort `json:"best_efforts"`
}

type DetailedSegmentEffort struct {
	ID               int64           `json:"id"`
	ActivityID       int64           `json:"activity_id"`
	ElapsedTime      int64           `json:"elapsed_time"`
	StartDate        time.Time       `json:"start_date"`
	StartDateLocal   time.Time       `json:"start_date_local"`
	Distance         float64         `json:"distance"`
	IsKom            bool            `json:"is_kom"`
	Name             string          `json:"name"`
	Activity         *MetaActivity   `json:"activity"`
	Athlete          *MetaAthlete    `json:"athlete"`
	MovingTime       int64           `json:"moving_time"`
	StartIndex       int64           `json:"start_index"`
	EndIndex         int64           `json:"end_index"`
	AverageCadence   float64         `json:"average_cadence"`
	AverageWatts     float64         `json:"average_watts"`
	DeviceWatts      bool            `json:"device_watts"`
	AverageHeartrate float64         `json:"average_heartrate"`
	MaxHeartrate     float64         `json:"max_heartrate"`
	Segment          *SummarySegment `json:"segment"`
	KomRank          int64           `json:"kom_rank"`
	PrRank           int64           `json:"pr_rank"`
	Hidden           bool            `json:"hidden"`
}

type Lap struct {
	ID                 int64         `json:"id"`
	Activity           *MetaActivity `json:"activity"`
	Athlete            *MetaAthlete  `json:"athlete"`
	AverageCadence     float64       `json:"average_cadence"`
	AverageSpeed       float64       `json:"average_speed"`
	Distance           float64       `json:"distance"`
	ElapsedTime        int64         `json:"elapsed_time"`
	StartIndex         int64         `json:"start_index"`
	EndIndex           int64         `json:"end_index"`
	LapIndex           int64         `json:"lap_index"`
	MaxSpeed           float64       `json:"max_speed"`
	MovingTime         int64         `json:"moving_time"`
	Name               string        `json:"name"`
	PaceZone           int64         `json:"pace_zone"`
	Split              int64         `json:"split"`
	StartDate          time.Time     `json:"start_date"`
	StartDateLocal     time.Time     `json:"start_date_local"`
	TotalElevationGain float64       `json:"total_elevation_gain"`
}

type MetaActivity struct {
	ID int64 `json:"id"`
}

type MetaAthlete struct {
	ID int64 `json:"id"`
}

type NewAthleteToken struct {
	AthleteID int64  `json:"athlete_id"`
	Refresh   string `json:"refresh"`
}

type PhotosSummary struct {
	Count   int64                 `json:"count"`
	Primary *PhotosSummaryPrimary `json:"primary"`
}

type PhotosSummaryPrimary struct {
	ID       int64    `json:"id"`
	Source   int64    `json:"source"`
	UniqueID string   `json:"unique_id"`
	Urls     []string `json:"urls"`
}

type PolylineMap struct {
	ID              string `json:"id"`
	Polyline        string `json:"polyline"`
	SummaryPolyline string `json:"summary_polyline"`
}

type Split struct {
	AverageSpeed        float64 `json:"average_speed"`
	Distance            float64 `json:"distance"`
	ElapsedTime         int64   `json:"elapsed_time"`
	ElevationDifference float64 `json:"elevation_difference"`
	PaceZone            int64   `json:"pace_zone"`
	MovingTime          int64   `json:"moving_time"`
	Split               int64   `json:"split"`
}

type SummaryGear struct {
	ID            string  `json:"id"`
	ResourceState int64   `json:"resource_state"`
	Primary       bool    `json:"primary"`
	Name          string  `json:"name"`
	Distance      float64 `json:"distance"`
}

type SummaryPrSegmentEffort struct {
	PrActivityID  int64     `json:"pr_activity_id"`
	PrElapsedTime int64     `json:"pr_elapsed_time"`
	PrDate        time.Time `json:"pr_date"`
	EffortCount   int64     `json:"effort_count"`
}

type SummarySegment struct {
	ID                  int64                   `json:"id"`
	Name                string                  `json:"name"`
	ActivityType        string                  `json:"activity_type"`
	Distance            float64                 `json:"distance"`
	AverageGrade        float64                 `json:"average_grade"`
	MaximumGrade        float64                 `json:"maximum_grade"`
	ElevationHigh       float64                 `json:"elevation_high"`
	ElevationLow        float64                 `json:"elevation_low"`
	StartLatlng         []float64               `json:"start_latlng"`
	EndLatlng           []float64               `json:"end_latlng"`
	ClimbCategory       int64                   `json:"climb_category"`
	City                string                  `json:"city"`
	State               string                  `json:"state"`
	Country             string                  `json:"country"`
	Private             bool                    `json:"private"`
	AthletePrEffort     *SummarySegmentEffort   `json:"athlete_pr_effort"`
	AthleteSegmentStats *SummaryPrSegmentEffort `json:"athlete_segment_stats"`
}

type SummarySegmentEffort struct {
	ID             int64     `json:"id"`
	ActivityID     int64     `json:"activity_id"`
	ElapsedTime    int64     `json:"elapsed_time"`
	StartDate      time.Time `json:"start_date"`
	StartDateLocal time.Time `json:"start_date_local"`
	Distance       float64   `json:"distance"`
	IsKom          bool      `json:"is_kom"`
}
