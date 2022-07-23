package services

import (
	"github.com/modhanami/boinger/models"
)

type TimelineService interface {
	List() (timelineResponse, error)
}

type timelineService struct {
	userService  UserService
	boingService BoingService
}

type UserUidToUser map[string]models.User

type timelineResponse struct {
	Boings *[]models.Boing `json:"boings"`
	Users  *UserUidToUser  `json:"users"`
}

func newTimelineResponse(boings *[]models.Boing, users *[]models.User) timelineResponse {
	userUidToUserMap := makeUserUidToUserMap(users)

	return timelineResponse{
		Boings: boings,
		Users:  &userUidToUserMap,
	}
}

func NewTimelineService(userService UserService, boingService BoingService) TimelineService {
	return &timelineService{
		userService:  userService,
		boingService: boingService,
	}
}

func (s *timelineService) List() (timelineResponse, error) {
	boings, err := s.boingService.List()
	if err != nil {
		return timelineResponse{}, err
	}

	userUids := extractUserUidsFromBoings(&boings)
	users, err := s.userService.GetByUids(userUids)
	if err != nil {
		return timelineResponse{}, err
	}

	return newTimelineResponse(&boings, &users), nil
}

func extractUserUidsFromBoings(boings *[]models.Boing) []string {
	uidsSet := make(map[string]struct{})
	for _, boing := range *boings {
		uidsSet[boing.UserUid] = struct{}{}
	}

	uids := make([]string, len(uidsSet))
	i := 0
	for userUid := range uidsSet {
		uids[i] = userUid
		i++
	}

	return uids
}

func makeUserUidToUserMap(users *[]models.User) UserUidToUser {
	userUidToUser := make(UserUidToUser)
	for _, user := range *users {
		userUidToUser[user.Uid] = user
	}
	return userUidToUser
}
