//go:generate go-enum --marshal --names --values --ptr --mustparse
package model

// ConfigKey
/*
ENUM(
fishpi // 摸鱼派
)
*/
type ConfigKey string

// PointStatus
/*
ENUM(
pending // 待发放
success // 发放成功
failed  // 发放失败
)
*/
type PointStatus string
