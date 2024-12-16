package rankings_service

import (
	configs "gym-badges-api/config/gym-badges-server"
	userDAO "gym-badges-api/internal/repository/user"
	"gym-badges-api/models"
	"gym-badges-api/tools/utils"

	log "github.com/sirupsen/logrus"
)

func NewRankingsService(userDAO userDAO.IUserDAO) IRankingsService {
	return &rankingsService{
		UserDAO: userDAO,
	}
}

type rankingsService struct {
	UserDAO userDAO.IUserDAO
}

func (r *rankingsService) GetGlobalRanking(userID string, page int32, ctxLog *log.Entry) (*models.GetRankingResponse, error) {

	ctxLog.Debugf("RANKINGS_SERVICE: Processing GetGlobalRanking for user: %s", userID)

	offset := int64(page-1) * int64(configs.Basic.RankingsPageSize)
	size := configs.Basic.RankingsPageSize
	firstRank := offset + 1
	lastRank := firstRank + int64(size) - 1

	users, err := r.UserDAO.GetUsersOrderedByExp(offset, size, ctxLog)
	if err != nil {
		return nil, err
	}

	user, selfRank, err := r.UserDAO.GetUserWithGlobalRank(userID, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.GetRankingResponse{
		Ranking: mapRanking(users, firstRank),
	}

	// User not in rank
	if selfRank < firstRank || selfRank > lastRank {
		response.Yourself = &models.RakingUser{
			UserID: userID,
			Name:   user.Name,
			// Image:  user.Image,
			Level:  int64(utils.CalcLevel(user.Experience)),
			Rank:   selfRank,
			Streak: user.Streak,
		}
	}

	return &response, nil
}

func mapRanking(users []*userDAO.User, firstRank int64) []*models.RakingUser {

	ranking := make([]*models.RakingUser, len(users))

	for i, u := range users {
		ranking[i] = &models.RakingUser{
			UserID: u.ID,
			Name:   u.Name,
			// Image:  u.Image,
			Level:  int64(utils.CalcLevel(u.Experience)),
			Rank:   firstRank + int64(i),
			Streak: u.Streak,
		}
	}

	return ranking
}
