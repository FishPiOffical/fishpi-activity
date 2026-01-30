package model

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	DbNameShields         = "shields"
	ShieldsFieldText      = "text"
	ShieldsFieldImg       = "img"
	ShieldsFieldUrl       = "url"
	ShieldsFieldBackcolor = "backcolor"
	ShieldsFieldFontcolor = "fontcolor"
	ShieldsFieldVer       = "ver"
	ShieldsFieldScale     = "scale"
	ShieldsFieldSize      = "size"
	ShieldsFieldBorder    = "border"
	ShieldsFieldBarLen    = "barlen"
	ShieldsFieldFontsize  = "fontsize"
	ShieldsFieldBarRadius = "barradius"
	ShieldsFieldShadow    = "shadow"
	ShieldsFieldAnime     = "anime"
	ShieldsFieldCreated   = "created"
	ShieldsFieldUpdated   = "updated"
)

// Shield wrapper type
type Shield struct {
	core.BaseRecordProxy
}

func NewShield(record *core.Record) *Shield {
	shield := new(Shield)
	shield.SetProxyRecord(record)
	return shield
}

func NewShieldFromCollection(collection *core.Collection) *Shield {
	record := core.NewRecord(collection)
	return NewShield(record)
}

func (shield *Shield) Text() string {
	return shield.GetString(ShieldsFieldText)
}

func (shield *Shield) SetText(value string) {
	shield.Set(ShieldsFieldText, value)
}

func (shield *Shield) Img() string {
	return shield.GetString(ShieldsFieldImg)
}

func (shield *Shield) SetImg(value any) {
	shield.Set(ShieldsFieldImg, value)
}

func (shield *Shield) Url() string {
	return shield.GetString(ShieldsFieldUrl)
}

func (shield *Shield) SetUrl(value string) {
	shield.Set(ShieldsFieldUrl, value)
}

func (shield *Shield) Backcolor() string {
	return shield.GetString(ShieldsFieldBackcolor)
}

func (shield *Shield) SetBackcolor(value string) {
	shield.Set(ShieldsFieldBackcolor, value)
}

func (shield *Shield) Fontcolor() string {
	return shield.GetString(ShieldsFieldFontcolor)
}

func (shield *Shield) SetFontcolor(value string) {
	shield.Set(ShieldsFieldFontcolor, value)
}

func (shield *Shield) Ver() string {
	return shield.GetString(ShieldsFieldVer)
}

func (shield *Shield) SetVer(value string) {
	shield.Set(ShieldsFieldVer, value)
}

func (shield *Shield) Scale() string {
	return shield.GetString(ShieldsFieldScale)
}

func (shield *Shield) SetScale(value string) {
	shield.Set(ShieldsFieldScale, value)
}

func (shield *Shield) Size() string {
	return shield.GetString(ShieldsFieldSize)
}

func (shield *Shield) SetSize(value string) {
	shield.Set(ShieldsFieldSize, value)
}

func (shield *Shield) Border() string {
	return shield.GetString(ShieldsFieldBorder)
}

func (shield *Shield) SetBorder(value string) {
	shield.Set(ShieldsFieldBorder, value)
}

func (shield *Shield) BarLen() string {
	return shield.GetString(ShieldsFieldBarLen)
}

func (shield *Shield) SetBarLen(value string) {
	shield.Set(ShieldsFieldBarLen, value)
}

func (shield *Shield) Fontsize() string {
	return shield.GetString(ShieldsFieldFontsize)
}

func (shield *Shield) SetFontsize(value string) {
	shield.Set(ShieldsFieldFontsize, value)
}

func (shield *Shield) BarRadius() string {
	return shield.GetString(ShieldsFieldBarRadius)
}

func (shield *Shield) SetBarRadius(value string) {
	shield.Set(ShieldsFieldBarRadius, value)
}

func (shield *Shield) Shadow() string {
	return shield.GetString(ShieldsFieldShadow)
}

func (shield *Shield) SetShadow(value string) {
	shield.Set(ShieldsFieldShadow, value)
}

func (shield *Shield) Anime() string {
	return shield.GetString(ShieldsFieldAnime)
}

func (shield *Shield) SetAnime(value string) {
	shield.Set(ShieldsFieldAnime, value)
}

func (shield *Shield) Created() types.DateTime {
	return shield.GetDateTime(ShieldsFieldCreated)
}

func (shield *Shield) Updated() types.DateTime {
	return shield.GetDateTime(ShieldsFieldUpdated)
}
