package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdministrator(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		user := User{
			TrapId: generateRandomUserName(),
		}
		err := adminRepo.AddAdministrator(user.TrapId)
		asr.NoError(err)

		admins, err := adminRepo.GetAdministratorList()
		asr.NoError(err)

		user.GiveIsUserAdmin(admins)
		asr.True(user.IsAdmin)

		err = adminRepo.RemoveAdministrator(user.TrapId)
		asr.NoError(err)

		admins, err = adminRepo.GetAdministratorList()
		asr.NoError(err)

		user.GiveIsUserAdmin(admins)
		asr.False(user.IsAdmin)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		user := User{
			TrapId: generateRandomUserName(),
		}
		err := adminRepo.AddAdministrator(user.TrapId)
		asr.NoError(err)

		err = adminRepo.AddAdministrator(user.TrapId)
		asr.Error(err)
	})
}
