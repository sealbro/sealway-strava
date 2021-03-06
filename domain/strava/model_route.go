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

type Route struct {
	Athlete *SummaryAthlete `bson:"athlete" json:"athlete,omitempty"`
	// The description of the route
	Description string `bson:"description" json:"description,omitempty"`
	// The route's distance, in meters
	Distance float32 `bson:"distance" json:"distance,omitempty"`
	// The route's elevation gain.
	ElevationGain float32 `bson:"elevation_gain" json:"elevation_gain,omitempty"`
	// The unique identifier of this route
	Id int64 `bson:"id" json:"id,omitempty"`
	// The unique identifier of the route in string format
	IdStr string       `bson:"id_str" json:"id_str,omitempty"`
	Map_  *PolylineMap `bson:"map" json:"map,omitempty"`
	// The name of this route
	Name string `bson:"name" json:"name,omitempty"`
	// Whether this route is private
	Private bool `bson:"private" json:"private,omitempty"`
	// Whether this route is starred by the logged-in athlete
	Starred bool `bson:"starred" json:"starred,omitempty"`
	// An epoch timestamp of when the route was created
	Timestamp int32 `bson:"timestamp" json:"timestamp,omitempty"`
	// This route's type (1 for ride, 2 for runs)
	Type_ int32 `bson:"type" json:"type,omitempty"`
	// This route's sub-type (1 for road, 2 for mountain bike, 3 for cross, 4 for trail, 5 for mixed)
	SubType int32 `bson:"sub_type" json:"sub_type,omitempty"`
	// The time at which the route was created
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	// The time at which the route was last updated
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
	// Estimated time in seconds for the authenticated athlete to complete route
	EstimatedMovingTime int32 `bson:"estimated_moving_time" json:"estimated_moving_time,omitempty"`
	// The segments traversed by this route
	Segments []SummarySegment `bson:"segments" json:"segments,omitempty"`
}
