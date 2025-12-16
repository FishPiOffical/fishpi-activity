package application

import (
	"bless-activity/model"

	"github.com/pocketbase/pocketbase/core"
)

type fixBugHandler func(e *core.BootstrapEvent) error

func (application *Application) fixBug(e *core.BootstrapEvent) error {
	list := []fixBugHandler{
		application.fixExample,
		//application.fixVoteLogValid,
	}

	for _, handler := range list {
		if err := handler(e); err != nil {
			return err
		}
	}

	return nil
}

func (application *Application) fixExample(*core.BootstrapEvent) error {
	return nil
}

func (application *Application) fixVoteLogValid(e *core.BootstrapEvent) error {
	app := e.App

	// Find all vote logs without valid field
	voteLogs, err := app.FindRecordsByFilter(
		model.DbNameVoteLogs,
		"valid = '' || valid = null",
		"",
		0,
		0,
	)
	if err != nil {
		return err
	}

	app.Logger().Info("修复投票日志 valid 字段", "count", len(voteLogs))

	for _, voteLogRecord := range voteLogs {
		voteLog := model.NewVoteLog(voteLogRecord)

		// Get the user who created the vote
		userRecord, err := app.FindRecordById(model.DbNameUsers, voteLog.FromUserId())
		if err != nil {
			app.Logger().Warn("找不到投票用户", "userId", voteLog.FromUserId(), "voteLogId", voteLogRecord.Id)
			continue
		}

		user := model.NewUser(userRecord)
		registeredAt := user.RegisteredAt()
		voteCreatedAt := voteLog.Created()

		// Check if vote was created at least 3 months after registration
		threeMonthsAfterRegistration := registeredAt.Time().AddDate(0, 3, 0)

		if voteCreatedAt.Time().After(threeMonthsAfterRegistration) || voteCreatedAt.Time().Equal(threeMonthsAfterRegistration) {
			voteLog.SetValid(model.VoteValidValid)
		} else {
			voteLog.SetValid(model.VoteValidInvalid)
		}

		if err := app.Save(voteLogRecord); err != nil {
			app.Logger().Error("保存投票日志失败", "voteLogId", voteLogRecord.Id, "err", err)
			return err
		}
	}

	app.Logger().Info("修复投票日志 valid 字段完成")
	return nil
}
