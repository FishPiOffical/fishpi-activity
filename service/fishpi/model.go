package fishpi

import "fmt"

type Config struct {
	BaseUrl        string `json:"base_url"`
	ApiKey         string `json:"api_key"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	MfaCode        string `json:"mfa_code"`
	Totp           string `json:"totp"`
	GoldFingerKey  string `json:"gold_finger_key"`
	MetalFingerKey string `json:"metal_finger_key"`
}

type UserInfoResult struct {
	Msg  string    `json:"msg"`
	Code int       `json:"code"`
	Data *UserInfo `json:"data"`
}

type UserInfo struct {
	UserAvatarURL string `json:"userAvatarURL"`
	UserNickname  string `json:"userNickname"`
	UserName      string `json:"userName"`
}

func (userInfo *UserInfo) Name() string {
	return fmt.Sprintf("(%s)%s", userInfo.UserName, userInfo.UserNickname)
}

type EditPointReq struct {
	UserName string // 用户名 8888
	Point    int    // 积分操作 -1
	Memo     string // 请求来源 游戏《我要当学霸》 原因 (8888)开摆购买道具 交易内容 学时x1
}

type EditPointReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type GetUserReply struct {
	UserCity           string `json:"userCity"`
	UserOnlineFlag     bool   `json:"userOnlineFlag"`
	UserPoint          int    `json:"userPoint"`
	UserAppRole        string `json:"userAppRole"`
	UserIntro          string `json:"userIntro"`
	UserNo             string `json:"userNo"`
	OnlineMinute       int    `json:"onlineMinute"`
	UserAvatarURL      string `json:"userAvatarURL"`
	UserNickname       string `json:"userNickname"`
	OId                string `json:"oId"`
	UserName           string `json:"userName"`
	CardBg             string `json:"cardBg"`
	AllMetalOwned      string `json:"allMetalOwned"`
	FollowingUserCount int    `json:"followingUserCount"`
	UserAvatarURL20    string `json:"userAvatarURL20"`
	SysMetal           string `json:"sysMetal"`
	Mbti               string `json:"mbti"`
	CanFollow          string `json:"canFollow"`
	UserRole           string `json:"userRole"`
	UserAvatarURL210   string `json:"userAvatarURL210"`
	FollowerCount      int    `json:"followerCount"`
	UserURL            string `json:"userURL"`
	UserAvatarURL48    string `json:"userAvatarURL48"`
}

type GetChatroomNodeGetResponse struct {
	Msg       string                         `json:"msg"`
	Code      int                            `json:"code"`
	Data      string                         `json:"data"`
	ApiKey    string                         `json:"apiKey"`
	Avaliable []*GetChatroomNodeGetAvailable `json:"avaliable"`
}

type GetChatroomNodeGetAvailable struct {
	Node   string `json:"node"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Online int    `json:"online"`
}

type PostChatroomSendRequest struct {
	ApiKey  string `json:"apiKey"`
	Client  string `json:"client"`
	Content string `json:"content"`
}

type PostChatroomSendResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type GetApiArticlesTagResponse struct {
	Msg  string                        `json:"msg"`
	Code int                           `json:"code"`
	Data GetApiArticlesTagResponseData `json:"data"`
}

type GetApiArticlesTagResponseData struct {
	Articles   []*GetApiArticlesTagResponseArticle `json:"articles"`
	Pagination struct {
		PaginationPageCount int   `json:"paginationPageCount"`
		PaginationPageNums  []int `json:"paginationPageNums"`
	} `json:"pagination"`
	Tag struct {
		TagShowSideAd      int           `json:"tagShowSideAd"`
		TagIconPath        string        `json:"tagIconPath"`
		TagStatus          int           `json:"tagStatus"`
		TagBadCnt          int           `json:"tagBadCnt"`
		TagRandomDouble    float64       `json:"tagRandomDouble"`
		TagTitle           string        `json:"tagTitle"`
		IsReserved         bool          `json:"isReserved"`
		OId                string        `json:"oId"`
		TagURI             string        `json:"tagURI"`
		TagAd              string        `json:"tagAd"`
		TagGoodCnt         int           `json:"tagGoodCnt"`
		TagCSS             string        `json:"tagCSS"`
		TagCommentCount    int           `json:"tagCommentCount"`
		TagDescriptionText string        `json:"tagDescriptionText"`
		TagFollowerCount   int           `json:"tagFollowerCount"`
		TagRelatedTags     []interface{} `json:"tagRelatedTags"`
		TagDomains         []interface{} `json:"tagDomains"`
		TagSeoTitle        string        `json:"tagSeoTitle"`
		TagLinkCount       int           `json:"tagLinkCount"`
		TagSeoDesc         string        `json:"tagSeoDesc"`
		TagReferenceCount  int           `json:"tagReferenceCount"`
		TagSeoKeywords     string        `json:"tagSeoKeywords"`
		TagDescription     string        `json:"tagDescription"`
	} `json:"tag"`
}

type GetApiArticlesTagResponseArticle struct {
	ArticleShowInList   int    `json:"articleShowInList"`
	ArticleCreateTime   string `json:"articleCreateTime"`
	ArticleAuthorId     string `json:"articleAuthorId"`
	ArticleBadCnt       int    `json:"articleBadCnt"`
	ArticleParticipants []struct {
		ArticleParticipantURL          string `json:"articleParticipantURL"`
		CommentId                      string `json:"commentId"`
		OId                            string `json:"oId"`
		ArticleParticipantName         string `json:"articleParticipantName"`
		ArticleParticipantThumbnailURL string `json:"articleParticipantThumbnailURL"`
	} `json:"articleParticipants"`
	ArticleLatestCmtTime        string  `json:"articleLatestCmtTime"`
	ArticleGoodCnt              int     `json:"articleGoodCnt"`
	ArticleQnAOfferPoint        int     `json:"articleQnAOfferPoint"`
	ArticleThumbnailURL         string  `json:"articleThumbnailURL"`
	ArticleStickRemains         int     `json:"articleStickRemains"`
	TimeAgo                     string  `json:"timeAgo"`
	ArticleUpdateTimeStr        string  `json:"articleUpdateTimeStr"`
	ArticleAuthorName           string  `json:"articleAuthorName"`
	ArticleType                 int     `json:"articleType"`
	Offered                     bool    `json:"offered"`
	ArticleCreateTimeStr        string  `json:"articleCreateTimeStr"`
	ArticleViewCount            int     `json:"articleViewCount"`
	ArticleAuthorThumbnailURL20 string  `json:"articleAuthorThumbnailURL20"`
	ArticleWatchCnt             int     `json:"articleWatchCnt"`
	ArticlePreviewContent       string  `json:"articlePreviewContent"`
	ArticleTitleEmoj            string  `json:"articleTitleEmoj"`
	ArticleTitleEmojUnicode     string  `json:"articleTitleEmojUnicode"`
	ArticleAuthorThumbnailURL48 string  `json:"articleAuthorThumbnailURL48"`
	ArticleCommentCount         int     `json:"articleCommentCount"`
	ArticleCollectCnt           int     `json:"articleCollectCnt"`
	ArticleTitle                string  `json:"articleTitle"`
	ArticleLatestCmterName      string  `json:"articleLatestCmterName"`
	ArticleTags                 string  `json:"articleTags"`
	OId                         string  `json:"oId"`
	CmtTimeAgo                  string  `json:"cmtTimeAgo"`
	ArticleStick                float64 `json:"articleStick"`
	ArticleTagObjs              []struct {
		TagShowSideAd     int     `json:"tagShowSideAd"`
		TagIconPath       string  `json:"tagIconPath"`
		TagStatus         int     `json:"tagStatus"`
		TagBadCnt         int     `json:"tagBadCnt"`
		TagRandomDouble   float64 `json:"tagRandomDouble"`
		TagTitle          string  `json:"tagTitle"`
		OId               string  `json:"oId"`
		TagURI            string  `json:"tagURI"`
		TagAd             string  `json:"tagAd"`
		TagGoodCnt        int     `json:"tagGoodCnt"`
		TagCSS            string  `json:"tagCSS"`
		TagCommentCount   int     `json:"tagCommentCount"`
		TagFollowerCount  int     `json:"tagFollowerCount"`
		TagSeoTitle       string  `json:"tagSeoTitle"`
		TagLinkCount      int     `json:"tagLinkCount"`
		TagSeoDesc        string  `json:"tagSeoDesc"`
		TagReferenceCount int     `json:"tagReferenceCount"`
		TagSeoKeywords    string  `json:"tagSeoKeywords"`
		TagDescription    string  `json:"tagDescription"`
	} `json:"articleTagObjs"`
	ArticleLatestCmtTimeStr      string                                  `json:"articleLatestCmtTimeStr"`
	ArticleAnonymous             int                                     `json:"articleAnonymous"`
	ArticleThankCnt              int                                     `json:"articleThankCnt"`
	ArticleUpdateTime            string                                  `json:"articleUpdateTime"`
	ArticleStatus                int                                     `json:"articleStatus"`
	ArticleHeat                  int                                     `json:"articleHeat"`
	ArticlePerfect               int                                     `json:"articlePerfect"`
	ArticleAuthorThumbnailURL210 string                                  `json:"articleAuthorThumbnailURL210"`
	ArticlePermalink             string                                  `json:"articlePermalink"`
	ArticleAuthor                *GetApiArticlesTagResponseArticleAuthor `json:"articleAuthor"`
}

type GetApiArticlesTagResponseArticleAuthor struct {
	UserOnlineFlag              bool   `json:"userOnlineFlag"`
	OnlineMinute                int    `json:"onlineMinute"`
	UserPointStatus             int    `json:"userPointStatus"`
	UserFollowerStatus          int    `json:"userFollowerStatus"`
	UserGuideStep               int    `json:"userGuideStep"`
	UserOnlineStatus            int    `json:"userOnlineStatus"`
	ChatRoomPictureStatus       int    `json:"chatRoomPictureStatus"`
	UserTags                    string `json:"userTags"`
	UserCommentStatus           int    `json:"userCommentStatus"`
	UserTimezone                string `json:"userTimezone"`
	UserURL                     string `json:"userURL"`
	UserForwardPageStatus       int    `json:"userForwardPageStatus"`
	UserUAStatus                int    `json:"userUAStatus"`
	UserIndexRedirectURL        string `json:"userIndexRedirectURL"`
	UserLatestArticleTime       int64  `json:"userLatestArticleTime"`
	UserTagCount                int    `json:"userTagCount"`
	UserNickname                string `json:"userNickname"`
	UserListViewMode            int    `json:"userListViewMode"`
	UserAvatarType              int    `json:"userAvatarType"`
	UserSubMailStatus           int    `json:"userSubMailStatus"`
	UserJoinPointRank           int    `json:"userJoinPointRank"`
	UserAppRole                 int    `json:"userAppRole"`
	UserAvatarViewMode          int    `json:"userAvatarViewMode"`
	UserStatus                  int    `json:"userStatus"`
	UserWatchingArticleStatus   int    `json:"userWatchingArticleStatus"`
	UserProvince                string `json:"userProvince"`
	UserNo                      int    `json:"userNo"`
	UserAvatarURL               string `json:"userAvatarURL"`
	UserFollowingTagStatus      int    `json:"userFollowingTagStatus"`
	UserLanguage                string `json:"userLanguage"`
	UserJoinUsedPointRank       int    `json:"userJoinUsedPointRank"`
	UserFollowingArticleStatus  int    `json:"userFollowingArticleStatus"`
	UserKeyboardShortcutsStatus int    `json:"userKeyboardShortcutsStatus"`
	UserReplyWatchArticleStatus int    `json:"userReplyWatchArticleStatus"`
	UserCommentViewMode         int    `json:"userCommentViewMode"`
	UserBreezemoonStatus        int    `json:"userBreezemoonStatus"`
	UserUsedPoint               int    `json:"userUsedPoint"`
	UserArticleStatus           int    `json:"userArticleStatus"`
	UserPoint                   int    `json:"userPoint"`
	UserCommentCount            int    `json:"userCommentCount"`
	UserIntro                   string `json:"userIntro"`
	UserMobileSkin              string `json:"userMobileSkin"`
	UserListPageSize            int    `json:"userListPageSize"`
	OId                         string `json:"oId"`
	UserName                    string `json:"userName"`
	UserGeoStatus               int    `json:"userGeoStatus"`
	UserSkin                    string `json:"userSkin"`
	UserNotifyStatus            int    `json:"userNotifyStatus"`
	UserFollowingUserStatus     int    `json:"userFollowingUserStatus"`
	UserArticleCount            int    `json:"userArticleCount"`
	Mbti                        string `json:"mbti"`
	UserRole                    string `json:"userRole"`
}
