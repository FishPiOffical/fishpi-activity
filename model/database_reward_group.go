package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameRewardGroups       = "rewardGroups"
	RewardGroupsFieldName    = "name"
	RewardGroupsFieldCreated = "created"
	RewardGroupsFieldUpdated = "updated"
)

type RewardGroup struct {
	core.BaseRecordProxy
}

func NewRewardGroup(record *core.Record) *RewardGroup {
	rewardGroup := new(RewardGroup)
	rewardGroup.SetProxyRecord(record)
	return rewardGroup
}

func NewRewardGroupFromCollection(collection *core.Collection) *RewardGroup {
	record := core.NewRecord(collection)
	return NewRewardGroup(record)
}

func (rewardGroup *RewardGroup) Name() string {
	return rewardGroup.GetString(RewardGroupsFieldName)
}

func (rewardGroup *RewardGroup) SetName(value string) {
	rewardGroup.Set(RewardGroupsFieldName, value)
}

func (rewardGroup *RewardGroup) Created() types.DateTime {
	return rewardGroup.GetDateTime(RewardGroupsFieldCreated)
}

func (rewardGroup *RewardGroup) Updated() types.DateTime {
	return rewardGroup.GetDateTime(RewardGroupsFieldUpdated)
}
