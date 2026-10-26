package service

import (
	"errors"
	"school-management-system/internal/models"
	"school-management-system/internal/repository"
)

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepository
}

func NewScheduleService(scheduleRepo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{scheduleRepo: scheduleRepo}
}

// SetMasterSchedule allows an admin to map a course, primary teacher, and proxy pool to a time slot
func (s *ScheduleService) SetMasterSchedule(courseID, teacherID int, proxyPool []int, dayOfWeek, startTime, endTime string) error {
	if courseID == 0 || teacherID == 0 {
		return errors.New("course ID and primary teacher ID are required")
	}

	var int32Pool []int32
	for _, id := range proxyPool {
		int32Pool = append(int32Pool, int32(id))
	}

	schedule := &models.Schedule{
		CourseID:  courseID,
		TeacherID: &teacherID,
		ProxyPool: int32Pool,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return s.scheduleRepo.CreateSchedule(schedule)
}