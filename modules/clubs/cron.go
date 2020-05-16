package clubs

import (
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
)

func (bot *module) removeNotVerified() {
	clubs, err := storage.Clubs.GetExpired()
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}

	tx, err := storage.Tx()
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}

	for _, club := range clubs {
		if err := club.DeleteMembers(tx); err != nil {
			go out.Err(true, err)
			go log.Error(err)
			tx.Rollback()
			return
		}

		if err := club.Delete(tx); err != nil {
			go out.Err(true, err)
			go log.Error(err)
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		go out.Err(true, err)
		go log.Error(err)
		return
	}
}
