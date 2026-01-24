//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

// ConfigKey
/*
ENUM(
fishpi // 摸鱼派
)
*/
type ConfigKey string

// VoteValid
/*
ENUM(
valid   // 有效
invalid // 无效
)
*/
type VoteValid string
