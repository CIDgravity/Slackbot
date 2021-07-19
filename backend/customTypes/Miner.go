package customTypes

import "database/sql"

type Miner struct {
	Id             int            `json:"id"`
	Address        sql.NullString `json:"address"`
	OwnerKey       sql.NullString `json:"owner_key"`
	IsActive       sql.NullBool   `json:"is_active"`
	SignProcessKey sql.NullString `json:"sign_process_key"`
	CreatedOn      sql.NullString `json:"created_on"`
}
