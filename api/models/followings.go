package models

type UnFollowRequest struct {
	FollowerId  string `json:"follower_id,omitempty"`
	FollowingId string `json:"following_id,omitempty"`
}

type FollowRequest struct {
	FollowerId  string `json:"follower_id,omitempty"`
	FollowingId string `json:"following_id,omitempty"`
}

type GetListFollowingRequest struct {
	Id     string  `json:"id,omitempty"`
	Page   int32   `json:"page,omitempty"`
	Limit  int32   `json:"limit,omitempty"`
	Search *string `json:"search,omitempty"`
}

type GetListFollowerRequest struct {
	Id     string  `json:"id,omitempty"`
	Page   int32   `json:"page,omitempty"`
	Limit  int32   `json:"limit,omitempty"`
	Search *string `json:"search,omitempty"`
}
