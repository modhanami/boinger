//go:build exclude

package services

import (
	"github.com/modhanami/boinger/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimelineService_List(t *testing.T) {
	service, userMock, boingMock := initTimelineServiceWithMocks()
	users := makeUsersFixture()
	boings := makeBoingsFixture()

	userMock.GetByUidsFunc = func(_ []string) ([]models.User, error) {
		return users, nil
	}
	boingMock.ListFunc = func() ([]models.Boing, error) {
		return boings, nil
	}
	expectedResult := newTimelineResponse(&boings, &users)

	result, err := service.List()

	assert.NoError(t, err)
	assert.Equal(t, result, expectedResult)
}

func makeUsersFixture() []models.User {
	return []models.User{
		{
			Username: "user1",
			Password: "pass1",
		},
		{
			Username: "user2",
			Password: "pass2",
		},
	}
}

func makeBoingsFixture() []models.Boing {
	return []models.Boing{
		{
			Id:        1,
			Uid:       "b1",
			UserId:    1,
			UserUid:   "a1",
			Text:      "boing1",
			CreatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			Id:        2,
			Uid:       "b2",
			UserId:    1,
			UserUid:   "a1",
			Text:      "boing2",
			CreatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			Id:        3,
			Uid:       "b3",
			UserId:    2,
			UserUid:   "a2",
			Text:      "boing3",
			CreatedAt: time.Now().AddDate(0, 0, -1),
		},
	}
}

func initTimelineServiceWithMocks() (TimelineService, *userServiceMock, *boingServiceMock) {
	userServiceMock := &userServiceMock{}
	boingServiceMock := &boingServiceMock{}

	return NewTimelineService(userServiceMock, boingServiceMock), userServiceMock, boingServiceMock
}

type userServiceMock struct {
	UserService
	GetByUidsFunc func([]string) ([]models.User, error)
}

func (s *userServiceMock) GetByUids(uids []string) ([]models.User, error) {
	return s.GetByUidsFunc(uids)
}

type boingServiceMock struct {
	BoingService
	ListFunc func() ([]models.Boing, error)
}

func (s *boingServiceMock) List() ([]models.Boing, error) {
	return s.ListFunc()
}
