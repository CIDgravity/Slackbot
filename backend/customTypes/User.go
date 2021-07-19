package customTypes

import "database/sql"

type User struct {
	Id   			int    				`json:"id"`
	Email 			sql.NullString 		`json:"email"`
	GithubUser 		sql.NullString 		`json:"github_user"`
	SlackUserId		sql.NullString 		`json:"slack_user_id"`
	SlackTeamId		sql.NullString 		`json:"slack_team_id"`
	SlackUsername	sql.NullString 		`json:"slack_username"`
}
