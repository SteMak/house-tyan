package clubs

import (
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
)

func (bot *module) removeNotVerified() {
	_, err := storage.Clubs.RemoveExpired()
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}

	// TODO: уведомить пользователей, что их клуб распущен
}
