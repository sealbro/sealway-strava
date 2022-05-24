/*
 * Strava API v3
 *
 * The [Swagger Playground](https://developers.strava.com/playground) is the easiest way to familiarize yourself with the Strava API by submitting HTTP requests and observing the responses before you write any client code. It will show what a response will look like with different endpoints depending on the authorization scope you receive from your athletes. To use the Playground, go to https://www.strava.com/settings/api and change your “Authorization Callback Domain” to developers.strava.com. Please note, we only support Swagger 2.0. There is a known issue where you can only select one scope at a time. For more information, please check the section “client code” at https://developers.strava.com/docs.
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package strava

// A set of rolled-up statistics and totals for an athlete
type ActivityStats struct {
	// The longest distance ridden by the athlete.
	BiggestRideDistance float64 `bson:"biggest_ride_distance" json:"biggest_ride_distance,omitempty"`
	// The highest climb ridden by the athlete.
	BiggestClimbElevationGain float64 `bson:"biggest_climb_elevation_gain" json:"biggest_climb_elevation_gain,omitempty"`
	// The recent (last 4 weeks) ride stats for the athlete.
	RecentRideTotals *ActivityTotal `bson:"recent_ride_totals" json:"recent_ride_totals,omitempty"`
	// The recent (last 4 weeks) run stats for the athlete.
	RecentRunTotals *ActivityTotal `bson:"recent_run_totals" json:"recent_run_totals,omitempty"`
	// The recent (last 4 weeks) swim stats for the athlete.
	RecentSwimTotals *ActivityTotal `bson:"recent_swim_totals" json:"recent_swim_totals,omitempty"`
	// The year to date ride stats for the athlete.
	YtdRideTotals *ActivityTotal `bson:"ytd_ride_totals" json:"ytd_ride_totals,omitempty"`
	// The year to date run stats for the athlete.
	YtdRunTotals *ActivityTotal `bson:"ytd_run_totals" json:"ytd_run_totals,omitempty"`
	// The year to date swim stats for the athlete.
	YtdSwimTotals *ActivityTotal `bson:"ytd_swim_totals" json:"ytd_swim_totals,omitempty"`
	// The all time ride stats for the athlete.
	AllRideTotals *ActivityTotal `bson:"all_ride_totals" json:"all_ride_totals,omitempty"`
	// The all time run stats for the athlete.
	AllRunTotals *ActivityTotal `bson:"all_run_totals" json:"all_run_totals,omitempty"`
	// The all time swim stats for the athlete.
	AllSwimTotals *ActivityTotal `bson:"all_swim_totals" json:"all_swim_totals,omitempty"`
}