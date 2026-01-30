package model

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	CommonFieldId = "id"
)

var (
	_ core.RecordProxy = (*User)(nil)
	_ core.RecordProxy = (*Config)(nil)
	_ core.RecordProxy = (*Activity)(nil)
	_ core.RecordProxy = (*Article)(nil)
	_ core.RecordProxy = (*Shield)(nil)
	_ core.RecordProxy = (*Vote)(nil)
	_ core.RecordProxy = (*VoteLog)(nil)
	_ core.RecordProxy = (*YearlyHistory)(nil)
	_ core.RecordProxy = (*RewardGroup)(nil)
	_ core.RecordProxy = (*Reward)(nil)
	_ core.RecordProxy = (*RewardDistribution)(nil)
	_ core.RecordProxy = (*RelArticle)(nil)
	_ core.RecordProxy = (*UserToken)(nil)
)
