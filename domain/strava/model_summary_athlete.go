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

type SummaryAthlete struct {
	// The unique identifier of the athlete
	Id int64 `bson:"id" json:"id,omitempty"`
	// Resource state, indicates level of detail. Possible values: 1 -> \"meta\", 2 -> \"summary\", 3 -> \"detail\"
	ResourceState int32 `bson:"resource_state" json:"resource_state,omitempty"`
	// The athlete's first name.
	Firstname string `bson:"firstname" json:"firstname,omitempty"`
	// The athlete's last name.
	Lastname string `bson:"lastname" json:"lastname,omitempty"`
	// URL to a 62x62 pixel profile picture.
	ProfileMedium string `bson:"profile_medium" json:"profile_medium,omitempty"`
	// URL to a 124x124 pixel profile picture.
	Profile string `bson:"profile" json:"profile,omitempty"`
	// The athlete's city.
	City string `bson:"city" json:"city,omitempty"`
	// The athlete's state or geographical region.
	State string `bson:"state" json:"state,omitempty"`
	// The athlete's country.
	Country string `bson:"country" json:"country,omitempty"`
	// The athlete's sex.
	Sex string `bson:"sex" json:"sex,omitempty"`
	// Deprecated.  Use summit field instead. Whether the athlete has any Summit subscription.
	Premium bool `bson:"premium" json:"premium,omitempty"`
	// Whether the athlete has any Summit subscription.
	Summit bool `bson:"summit" json:"summit,omitempty"`
	// The time at which the athlete was created.
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	// The time at which the athlete was last updated.
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}