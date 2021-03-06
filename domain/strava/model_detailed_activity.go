/*
 * Strava API v3
 *
 * The [Swagger Playground](https://developers.strava.com/playground) is the easiest way to familiarize yourself with the Strava API by submitting HTTP requests and observing the responses before you write any client code. It will show what a response will look like with different endpoints depending on the authorization scope you receive from your athletes. To use the Playground, go to https://www.strava.com/settings/api and change your “Authorization Callback Domain” to developers.strava.com. Please note, we only support Swagger 2.0. There is a known issue where you can only select one scope at a time. For more information, please check the section “client code” at https://developers.strava.com/docs.
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package strava

import (
	"time"
)

type DetailedActivity struct {
	// The unique identifier of the activity
	ID int64 `bson:"_id" json:"id,omitempty"`
	// The identifier provided at upload time
	ExternalID string `bson:"external_id" json:"external_id,omitempty"`
	// The identifier of the upload that resulted in this activity
	UploadID int64        `bson:"upload_id" json:"upload_id,omitempty"`
	Athlete  *MetaAthlete `bson:"athlete" json:"athlete,omitempty"`
	// The name of the activity
	Name string `bson:"name" json:"name,omitempty"`
	// The activity's distance, in meters
	Distance float64 `bson:"distance" json:"distance,omitempty"`
	// The activity's moving time, in seconds
	MovingTime int64 `bson:"moving_time" json:"moving_time,omitempty"`
	// The activity's elapsed time, in seconds
	ElapsedTime int64 `bson:"elapsed_time" json:"elapsed_time,omitempty"`
	// The activity's total elevation gain.
	TotalElevationGain float64 `bson:"total_elevation_gain" json:"total_elevation_gain,omitempty"`
	// The activity's highest elevation, in meters
	ElevHigh float64 `bson:"elev_high" json:"elev_high,omitempty"`
	// The activity's lowest elevation, in meters
	ElevLow float64 `bson:"elev_low" json:"elev_low,omitempty"`
	// model_activity_type.go List of ActivityType
	Type *string `bson:"type" json:"type,omitempty"`
	// The time at which the activity was started.
	StartDate time.Time `bson:"start_date" json:"start_date,omitempty"`
	// The time at which the activity was started in the local timezone.
	StartDateLocal time.Time `bson:"start_date_local" json:"start_date_local,omitempty"`
	// The timezone of the activity
	Timezone    string    `bson:"timezone" json:"timezone,omitempty"`
	StartLatlng []float64 `bson:"start_latlng" json:"start_latlng,omitempty"`
	EndLatlng   []float64 `bson:"end_latlng" json:"end_latlng,omitempty"`
	// The number of achievements gained during this activity
	AchievementCount int64 `bson:"achievement_count" json:"achievement_count,omitempty"`
	// The number of kudos given for this activity
	KudosCount int64 `bson:"kudos_count" json:"kudos_count,omitempty"`
	// The number of comments for this activity
	CommentCount int64 `bson:"comment_count" json:"comment_count,omitempty"`
	// The number of athletes for taking part in a group activity
	AthleteCount int64 `bson:"athlete_count" json:"athlete_count,omitempty"`
	// The number of Instagram photos for this activity
	PhotoCount int64 `bson:"photo_count" json:"photo_count,omitempty"`
	// The number of Instagram and Strava photos for this activity
	TotalPhotoCount int64        `bson:"total_photo_count" json:"total_photo_count,omitempty"`
	Map             *PolylineMap `bson:"map" json:"map,omitempty"`
	// Whether this activity was recorded on a training machine
	Trainer bool `bson:"trainer" json:"trainer,omitempty"`
	// Whether this activity is a commute
	Commute bool `bson:"commute" json:"commute,omitempty"`
	// Whether this activity was created manually
	Manual bool `bson:"manual" json:"manual,omitempty"`
	// Whether this activity is private
	Private bool `bson:"private" json:"private,omitempty"`
	// Whether this activity is flagged
	Flagged bool `bson:"flagged" json:"flagged,omitempty"`
	// The activity's workout type
	WorkoutType int64 `bson:"workout_type" json:"workout_type,omitempty"`
	// The unique identifier of the upload in string format
	UploadIDStr string `bson:"upload_id_str" json:"upload_id_str,omitempty"`
	// The activity's average speed, in meters per second
	AverageSpeed float64 `bson:"average_speed" json:"average_speed,omitempty"`
	// The activity's max speed, in meters per second
	MaxSpeed float64 `bson:"max_speed" json:"max_speed,omitempty"`
	// Whether the logged-in athlete has kudoed this activity
	HasKudoed bool `bson:"has_kudoed" json:"has_kudoed,omitempty"`
	// The id of the gear for the activity
	GearID string `bson:"gear_id" json:"gear_id,omitempty"`
	// The total work done in kilojoules during this activity. Rides only
	Kilojoules float64 `bson:"kilojoules" json:"kilojoules,omitempty"`
	// Average power output in watts during this activity. Rides only
	AverageWatts float64 `bson:"average_watts" json:"average_watts,omitempty"`
	// Whether the watts are from a power meter, false if estimated
	DeviceWatts bool `bson:"device_watts" json:"device_watts,omitempty"`
	// Rides with power meter data only
	MaxWatts int64 `bson:"max_watts" json:"max_watts,omitempty"`
	// Similar to Normalized Power. Rides with power meter data only
	WeightedAverageWatts int64 `bson:"weighted_average_watts" json:"weighted_average_watts,omitempty"`
	// The description of the activity
	Description string         `bson:"description" json:"description,omitempty"`
	Photos      *PhotosSummary `bson:"photos" json:"photos,omitempty"`
	Gear        *SummaryGear   `bson:"gear" json:"gear,omitempty"`
	// The number of kilocalories consumed during this activity
	Calories       float64                  `bson:"calories" json:"calories,omitempty"`
	SegmentEfforts []*DetailedSegmentEffort `bson:"segment_efforts" json:"segment_efforts,omitempty"`
	// The name of the device used to record the activity
	DeviceName string `bson:"device_name" json:"device_name,omitempty"`
	// The token used to embed a Strava activity
	EmbedToken string `bson:"embed_token" json:"embed_token,omitempty"`
	// The splits of this activity in metric units (for runs)
	SplitsMetric []*Split `bson:"splits_metric" json:"splits_metric,omitempty"`
	// The splits of this activity in imperial units (for runs)
	SplitsStandard []*Split                 `bson:"splits_standard" json:"splits_standard,omitempty"`
	Laps           []*Lap                   `bson:"laps" json:"laps,omitempty"`
	BestEfforts    []*DetailedSegmentEffort `bson:"best_efforts" json:"best_efforts,omitempty"`
}
