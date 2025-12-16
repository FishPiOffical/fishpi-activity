package application

import (
	"bless-activity/model"

	"github.com/pocketbase/pocketbase/core"
)

func (application *Application) registerHooks() {
	// Register hook for vote log creation to auto-populate valid field
	application.app.OnRecordCreate(model.DbNameVoteLogs).BindFunc(func(e *core.RecordEvent) error {
		voteLog := model.NewVoteLog(e.Record)

		// Get the user who is creating the vote
		userRecord, err := e.App.FindRecordById(model.DbNameUsers, voteLog.FromUserId())
		if err != nil {
			e.App.Logger().Warn("找不到投票用户", "userId", voteLog.FromUserId())
			// If user not found, mark as invalid
			voteLog.SetValid(model.VoteValidInvalid)
			return e.Next()
		}

		user := model.NewUser(userRecord)
		registeredAt := user.RegisteredAt()
		voteCreatedAt := voteLog.Created()

		// Check if vote is being created at least 3 months after registration
		threeMonthsAfterRegistration := registeredAt.Time().AddDate(0, 3, 0)

		if voteCreatedAt.Time().After(threeMonthsAfterRegistration) || voteCreatedAt.Time().Equal(threeMonthsAfterRegistration) {
			voteLog.SetValid(model.VoteValidValid)
		} else {
			voteLog.SetValid(model.VoteValidInvalid)
		}

		return e.Next()
	})
}
