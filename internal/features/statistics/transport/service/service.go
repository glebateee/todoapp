package statistics_service

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface{}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) StatisticsService {
	return StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
