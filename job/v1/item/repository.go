package item

import (
	"database/sql"
	"time"
)

type itemRepository struct {
	tx *sql.Tx
}

func NewItemRepository(transaction *sql.Tx) ItemRepository {
	return &itemRepository{
		tx: transaction,
	}
}

type ItemRepository interface {
	GetZohoUpdated(string) (*time.Time, error)
}

func (r *itemRepository) GetZohoUpdated(zohoID string) (*time.Time, error) {
	var res time.Time
	row := r.tx.QueryRow(`SELECT zoho_updated FROM items WHERE zoho_id = ? LIMIT 1`, zohoID)
	err := row.Scan(&res)
	return &res, err
}

// func (r *itemRepository) UpdateToken(id int64, info Token) error {
// 	_, err := r.tx.Exec(`
// 		Update tokens SET
// 		access_token = ?,
// 		api_domain = ?,
// 		token_type = ?,
// 		expires_time = ?
// 		WHERE id = ?
// 	`, info.AccessToken, info.ApiDomain, info.TokenType, info.ExpiresTime, id)
// 	if err != nil {
// 		msg := "更新失败:" + err.Error()
// 		return errors.New(msg)
// 	}
// 	return nil
// }
