package models

type Schedule struct {
	ID        int    `json:"id"`
	CourseID  int    `json:"course_id"`
	TeacherID *int   `json:"teacher_id,omitempty"`
	ProxyPool []int32  `json:"proxy_pool,omitempty"`
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Room      string `json:"room"`
}