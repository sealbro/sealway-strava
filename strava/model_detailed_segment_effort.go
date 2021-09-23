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

type DetailedSegmentEffort struct {
	// The unique identifier of this effort
	Id int64 `bson:"id" json:"id,omitempty"`
	// The unique identifier of the activity related to this effort
	ActivityId int64 `bson:"activity_id" json:"activity_id,omitempty"`
	// The effort's elapsed time
	ElapsedTime int32 `bson:"elapsed_time" json:"elapsed_time,omitempty"`
	// The time at which the effort was started.
	StartDate time.Time `bson:"start_date" json:"start_date,omitempty"`
	// The time at which the effort was started in the local timezone.
	StartDateLocal time.Time `bson:"start_date_local" json:"start_date_local,omitempty"`
	// The effort's distance in meters
	Distance float32 `bson:"distance" json:"distance,omitempty"`
	// Whether this effort is the current best on the leaderboard
	IsKom bool `bson:"is_kom" json:"is_kom,omitempty"`
	// The name of the segment on which this effort was performed
	Name     string        `bson:"name" json:"name,omitempty"`
	Activity *MetaActivity `bson:"activity" json:"activity,omitempty"`
	Athlete  *MetaAthlete  `bson:"athlete" json:"athlete,omitempty"`
	// The effort's moving time
	MovingTime int32 `bson:"moving_time" json:"moving_time,omitempty"`
	// The start index of this effort in its activity's stream
	StartIndex int32 `bson:"start_index" json:"start_index,omitempty"`
	// The end index of this effort in its activity's stream
	EndIndex int32 `bson:"end_index" json:"end_index,omitempty"`
	// The effort's average cadence
	AverageCadence float32 `bson:"average_cadence" json:"average_cadence,omitempty"`
	// The average wattage of this effort
	AverageWatts float32 `bson:"average_watts" json:"average_watts,omitempty"`
	// For riding efforts, whether the wattage was reported by a dedicated recording device
	DeviceWatts bool `bson:"device_watts" json:"device_watts,omitempty"`
	// The heart heart rate of the athlete during this effort
	AverageHeartrate float32 `bson:"average_heartrate" json:"average_heartrate,omitempty"`
	// The maximum heart rate of the athlete during this effort
	MaxHeartrate float32         `bson:"max_heartrate" json:"max_heartrate,omitempty"`
	Segment      *SummarySegment `bson:"segment" json:"segment,omitempty"`
	// The rank of the effort on the global leaderboard if it belongs in the top 10 at the time of upload
	KomRank int32 `bson:"kom_rank" json:"kom_rank,omitempty"`
	// The rank of the effort on the athlete's leaderboard if it belongs in the top 3 at the time of upload
	PrRank int32 `bson:"pr_rank" json:"pr_rank,omitempty"`
	// Whether this effort should be hidden when viewed within an activity
	Hidden bool `bson:"hidden" json:"hidden,omitempty"`
}
