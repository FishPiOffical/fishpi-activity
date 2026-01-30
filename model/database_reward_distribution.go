package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameRewardDistributions       = "rewardDistributions"
	RewardDistributionsFieldVoteId  = "voteId"
	RewardDistributionsFieldUserId  = "userId"
	RewardDistributionsFieldRank    = "rank"
	RewardDistributionsFieldPoint   = "point"
	RewardDistributionsFieldStatus  = "status"
	RewardDistributionsFieldMemo    = "memo"
	RewardDistributionsFieldCreated = "created"
	RewardDistributionsFieldUpdated = "updated"
)

type RewardDistribution struct {
	core.BaseRecordProxy
}

func NewRewardDistribution(record *core.Record) *RewardDistribution {
	rewardDistribution := new(RewardDistribution)
	rewardDistribution.SetProxyRecord(record)
	return rewardDistribution
}

func NewRewardDistributionFromCollection(collection *core.Collection) *RewardDistribution {
	record := core.NewRecord(collection)
	return NewRewardDistribution(record)
}

func (rd *RewardDistribution) VoteId() string {
	return rd.GetString(RewardDistributionsFieldVoteId)
}

func (rd *RewardDistribution) SetVoteId(value string) {
	rd.Set(RewardDistributionsFieldVoteId, value)
}

func (rd *RewardDistribution) UserId() string {
	return rd.GetString(RewardDistributionsFieldUserId)
}

func (rd *RewardDistribution) SetUserId(value string) {
	rd.Set(RewardDistributionsFieldUserId, value)
}

func (rd *RewardDistribution) Rank() int {
	return rd.GetInt(RewardDistributionsFieldRank)
}

func (rd *RewardDistribution) SetRank(value int) {
	rd.Set(RewardDistributionsFieldRank, value)
}

func (rd *RewardDistribution) Point() int {
	return rd.GetInt(RewardDistributionsFieldPoint)
}

func (rd *RewardDistribution) SetPoint(value int) {
	rd.Set(RewardDistributionsFieldPoint, value)
}

func (rd *RewardDistribution) Status() DistributionStatus {
	return MustParseDistributionStatus(rd.GetString(RewardDistributionsFieldStatus))
}

func (rd *RewardDistribution) SetStatus(value DistributionStatus) {
	rd.Set(RewardDistributionsFieldStatus, value)
}

func (rd *RewardDistribution) Memo() string {
	return rd.GetString(RewardDistributionsFieldMemo)
}

func (rd *RewardDistribution) SetMemo(value string) {
	rd.Set(RewardDistributionsFieldMemo, value)
}

func (rd *RewardDistribution) Created() types.DateTime {
	return rd.GetDateTime(RewardDistributionsFieldCreated)
}

func (rd *RewardDistribution) Updated() types.DateTime {
	return rd.GetDateTime(RewardDistributionsFieldUpdated)
}
