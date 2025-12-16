//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

// ConfigKey
/*
ENUM(
fishpi // 摸鱼派
)
*/
type ConfigKey string

// DistributionStatus
/*
ENUM(
pending      // 待发放
distributing // 发放中
failed       // 发放失败
success      // 发放成功
)
*/
type DistributionStatus string

// VoteValid
/*
ENUM(
valid   // 有效
invalid // 无效
)
*/
type VoteValid string
