package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdministrator(t *testing.T) {
	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		user := TrapUser{
			TrapID: generateRandomUserName(),
		}
		err := adminRepo.AddAdministrator(user.TrapID)
		asr.NoError(err)

		admins, err := adminRepo.GetAdministratorList()
		asr.NoError(err)
		asr.Len(admins, 1)
		asr.Equal(user.TrapID, admins[0])

		flag, err := adminRepo.IsAdmin(user.TrapID)
		asr.NoError(err)
		asr.True(flag)

		testID := user.TrapID + "0"
		flag, err = adminRepo.IsAdmin(testID)
		asr.NoError(err)
		asr.False(flag)

		user.GiveIsUserAdmin(admins)
		asr.True(user.IsAdmin)

		err = adminRepo.RemoveAdministrator(user.TrapID)
		asr.NoError(err)

		admins, err = adminRepo.GetAdministratorList()
		asr.NoError(err)

		user.GiveIsUserAdmin(admins)
		asr.False(user.IsAdmin)
	})
}
