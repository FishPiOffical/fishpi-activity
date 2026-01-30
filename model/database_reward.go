package model

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	DbNameRewards             = "rewards"
	RewardsFieldRewardGroupId = "rewardGroupId"
	RewardsFieldName          = "name"
	RewardsFieldMin           = "min"
	RewardsFieldMax           = "max"
	RewardsFieldPoint         = "point"
	RewardsFieldShieldIds     = "shieldIds"
	RewardsFieldMore          = "more"
)

type Reward struct {
	core.BaseRecordProxy
}

func NewReward(record *core.Record) *Reward {
	reward := new(Reward)
	reward.SetProxyRecord(record)
	return reward
}

func NewRewardFromCollection(collection *core.Collection) *Reward {
	record := core.NewRecord(collection)
	return NewReward(record)
}

func (reward *Reward) RewardGroupId() string {
	return reward.GetString(RewardsFieldRewardGroupId)
}

func (reward *Reward) SetRewardGroupId(value string) {
	reward.Set(RewardsFieldRewardGroupId, value)
}

func (reward *Reward) Name() string {
	return reward.GetString(RewardsFieldName)
}

func (reward *Reward) SetName(value string) {
	reward.Set(RewardsFieldName, value)
}

func (reward *Reward) Min() int {
	return reward.GetInt(RewardsFieldMin)
}

func (reward *Reward) SetMin(value int) {
	reward.Set(RewardsFieldMin, value)
}

func (reward *Reward) Max() int {
	return reward.GetInt(RewardsFieldMax)
}

func (reward *Reward) SetMax(value int) {
	reward.Set(RewardsFieldMax, value)
}

func (reward *Reward) Point() int {
	return reward.GetInt(RewardsFieldPoint)
}

func (reward *Reward) SetPoint(value int) {
	reward.Set(RewardsFieldPoint, value)
}

func (reward *Reward) ShieldIds() []string {
	return reward.GetStringSlice(RewardsFieldShieldIds)
}

func (reward *Reward) SetShieldIds(value []string) {
	reward.Set(RewardsFieldShieldIds, value)
}

func (reward *Reward) More() string {
	return reward.GetString(RewardsFieldMore)
}

func (reward *Reward) SetMore(value string) {
	reward.Set(RewardsFieldMore, value)
}
