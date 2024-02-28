// backend/models/agents.go

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	// "github.com/go-sql-driver/mysql"
	// Prefered my MySQL driver
	_ "log"

	"gorm.io/gorm"
)

// Agents represents the schema of the agents table
type Agents struct {
	ID           uint                  `gorm:"primaryKey" json:"id"`
	FirstName    string                `gorm:"size:255;not null" json:"first_name" binding:"required"`
	LastName     string                `gorm:"size:255;not null" json:"last_name" binding:"required"`
	Email        string                `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	Credentials  AgentLoginCredentials `gorm:"embedded;foreignKey:AgentLoginCredentialsID" json:"agent_credentials,omitempty"`
	Phone        *string               `gorm:"size:20" json:"phone,omitempty" binding:"omitempty,e164"`
	PositionID   uint                  `gorm:"index;type:int unsigned" json:"position_id,omitempty"`
	DepartmentID uint                  `gorm:"index;type:int unsigned" json:"department_id,omitempty"`
	IsActive     bool                  `gorm:"default:true" json:"is_active"`
	ProfilePic   *string               `gorm:"size:255" json:"profile_pic,omitempty"`
	LastLoginAt  *time.Time            `json:"last_login_at,omitempty"`
	TeamID       *uint                 `gorm:"type:int unsigned" json:"team_id,omitempty"`
	SupervisorID *uint                 `gorm:"type:int unsigned" json:"supervisor_id,omitempty"`
	Roles        []Role                `gorm:"many2many:agent_roles;" json:"roles"`
	Biography    string                `json:"biography,omitempty"`
	UserID       uint                  `gorm:"index;type:int; unsigned" json:"user_id,omitempty"`
	AgentProfile AgentProfile          `gorm:"foreignKey:AgentID" json:"_"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	DeletedAt    gorm.DeletedAt        `gorm:"index" json:"deleted_at,omitempty"`
}

func (Agents) TableName() string {
	return "agents"
}

// Implement driver.Valuer for Agents
func (a Agents) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Implement driver.Scanner for Agents
func (a *Agents) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Agents scan")
	}

	return json.Unmarshal(data, &a)
}

type AgentProfile struct {
	AgentID         uint   `gorm:"primaryKey;autoIncrement:false" json:"agent_id"`
	Bio             string `gorm:"type:text" json:"bio,omitempty"`
	AvatarURL       string `gorm:"type:text" json:"avatar_url,omitempty"`
	Preferences     string `gorm:"type:text" json:"preferences,omitempty"`
	PrivacySettings string `gorm:"type:text" json:"privacy_settings,omitempty"`
}

func (AgentProfile) TableName() string {
	return "agent_profiles"
}

// Implement driver.Valuer for AgentProfile
func (ap AgentProfile) Value() (driver.Value, error) {
	return json.Marshal(ap)
}

// Implement driver.Scanner for AgentProfile
func (ap *AgentProfile) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for AgentProfile scan")
	}

	return json.Unmarshal(data, &ap)
}

// Unit represents the schema of the unit table
type Unit struct {
	gorm.Model
	UnitName string  `gorm:"size:255;not null" json:"unit_name"`
	Emoji    *string `gorm:"size:255" json:"emoji,omitempty"`
}

func (Unit) TableName() string {
	return "units"
}

// Implement driver.Valuer for Unit
func (u Unit) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Implement driver.Scanner for Unit
func (u *Unit) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Unit scan")
	}

	return json.Unmarshal(data, &u)
}

// Permission represents the schema of the permission table
type Permission struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

// Implement driver.Valuer for Permission
func (p Permission) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Implement driver.Scanner for Permission
func (p *Permission) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Permission scan")
	}

	return json.Unmarshal(data, &p)
}

// Team represents the schema of the teams table
type Teams struct {
	gorm.Model
	TeamName         string  `gorm:"size:255;not null" json:"team_name"`
	Emoji            *string `gorm:"size:255" json:"emoji,omitempty"`
	TeamPermissionID *uint   `gorm:"type:int unsigned" json:"team_permission_id,omitempty"`
}

func (Teams) TableName() string {
	return "teams"
}

// Implement driver.Valuer for Teams
func (t Teams) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Implement driver.Scanner for Teams
func (t *Teams) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Teams scan")
	}

	return json.Unmarshal(data, &t)
}

// TeamPermission links 'teams' with their 'permissions'.
type TeamPermission struct {
	gorm.Model
	TeamID      uint          `gorm:"not null;index:idx_team_id,unique" json:"team_id"`
	Permissions []*Permission `gorm:"many2many:team_permissions_permissions;" json:"permissions,omitempty"`
}

func (TeamPermission) TableName() string {
	return "team_permission"
}

// Implement driver.Valuer for TeamPermission
func (tp TeamPermission) Value() (driver.Value, error) {
	return json.Marshal(tp)
}

// Implement driver.Scanner for TeamPermission
func (tp *TeamPermission) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for TeamPermission scan")
	}

	return json.Unmarshal(data, &tp)
}

// Role represents the schema of the role table
type Role struct {
	gorm.Model
	RoleName    string  `gorm:"size:255;not null" json:"role_name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	Users       []Users `gorm:"many2many:user_roles;" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

// Implement driver.Valuer for Role
func (r Role) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Implement driver.Scanner for Role
func (r *Role) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Role scan")
	}

	return json.Unmarshal(data, &r)
}

// RoleBase represents a foundational role structure that may be used for additional role metadata
type RoleBase struct {
	gorm.Model
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (RoleBase) TableName() string {
	return "role_bases"
}

// Implement driver.Valuer for RoleBase
func (rb RoleBase) Value() (driver.Value, error) {
	return json.Marshal(rb)
}

// Implement driver.Scanner for RoleBase
func (rb *RoleBase) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for RoleBase scan")
	}

	return json.Unmarshal(data, &rb)
}

// RolePermission links roles with permissions in a many-to-many relationship
type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"not null;index:idx_role_permission,unique" json:"role_id"`
	PermissionID uint `gorm:"not null;index:idx_role_permission,unique" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// Implement driver.Valuer for RolePermission
func (rp RolePermission) Value() (driver.Value, error) {
	return json.Marshal(rp)
}

// Implement driver.Scanner for RolePermission
func (rp *RolePermission) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for RolePermission scan")
	}

	return json.Unmarshal(data, &rp)
}

// AgentRole links agents with roles in a many-to-many relationship
type AgentRole struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	AgentID   uint            `gorm:"not null;index:idx_agent_role,unique" json:"agent_id"`
	RoleID    uint            `gorm:"not null;index:idx_agent_role,unique" json:"role_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentRole) TableName() string {
	return "agent_roles"
}

// Implement driver.Valuer for AgentRole
func (ar AgentRole) Value() (driver.Value, error) {
	return json.Marshal(ar)
}

// Implement driver.Scanner for AgentRole
func (ar *AgentRole) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for AgentRole scan")
	}

	return json.Unmarshal(data, &ar)
}

// AgentTrainingSession represents the schema of the agent_training_sessions table
type AgentTrainingSession struct {
	gorm.Model
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location" gorm:"type:varchar(255)"`
	TrainerID   uint      `json:"trainer_id" gorm:"index;not null"`
	Attendees   []Agents  `gorm:"many2many:agent_training_attendees;" json:"attendees"`
}

func (AgentTrainingSession) TableName() string {
	return "agent_training_sessions"
}

// GetAttendeesCount returns the count of attendees for the training session
func (ats *AgentTrainingSession) GetAttendeesCount() int {
	return len(ats.Attendees)
}

// IsTrainerAvailable checks if the trainer for the session is available
func (ats *AgentTrainingSession) IsTrainerAvailable() bool {
	// Implement the logic to check trainer availability
	// You might want to query the schedule or availability
	return true
}

// AddAttendee adds an agent to the list of attendees for the training session
func (ats *AgentTrainingSession) AddAttendee(agent *Agents) {
	ats.Attendees = append(ats.Attendees, *agent)
}

// UserAgent represents the relationship between a user and an agent
type UserAgent struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null;index:idx_user_agent,user_id" json:"user_id"`
	AgentID   uint            `gorm:"not null;index:idx_user_agent,agent_id" json:"agent_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserAgent) TableName() string {
	return "user_agents"
}

// IsActive checks if the user-agent relationship is active
func (ua *UserAgent) IsActive() bool {
	// Implement the logic to check if the relationship is active
	return true
}

// TerminateRelationship terminates the user-agent relationship
func (ua *UserAgent) TerminateRelationship() {
	// Implement the logic to terminate the relationship
}

// TeamAgent represents the relationship between a team and an agent
type TeamAgent struct {
	gorm.Model
	TeamID  uint `gorm:"not null;index:idx_team_agent,team_id" json:"team_id"`
	AgentID uint `gorm:"not null;index:idx_team_agent,agent_id" json:"agent_id"`
}

func (TeamAgent) TableName() string {
	return "team_agents"
}

// IsLeader checks if the agent is the leader of the team
func (ta *TeamAgent) IsLeader() bool {
	// Implement the logic to check if the agent is the leader
	return true
}

// ChangeTeam changes the team for the agent
func (ta *TeamAgent) ChangeTeam(newTeamID uint) {
	// Implement the logic to change the team
}

// AgentPermission represents the relationship between an agent and their granted permissions
type AgentPermission struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	AgentID      uint            `json:"agent_id" gorm:"index;not null"`
	PermissionID uint            `json:"permission_id" gorm:"index;not null"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentPermission) TableName() string {
	return "agent_permissions"
}

// SearchCriteria represents the schema of the search_criteria table
type SearchCriteria struct {
	gorm.Model
	Name       string `json:"name,omitempty"`
	Role       string `json:"role,omitempty"`
	Department string `json:"department,omitempty"`
}

func (SearchCriteria) TableName() string {
	return "search_criteria"
}

// AgentSchedule represents the schema of the agent_schedules table
type AgentSchedule struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ShiftType string    `json:"shift_type" gorm:"type:varchar(100);not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
}

// AgentShift represents the schema of the agent_shifts table
type AgentShift struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	ShiftDate time.Time `json:"shift_date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	ShiftType string    `json:"shift_type" gorm:"type:varchar(100);not null"`
}

// AgentSkill represents the schema of the agent_skills table
type AgentSkill struct {
	gorm.Model
	AgentID   uint   `json:"agent_id" gorm:"index;not null"`
	SkillName string `json:"skill_name" gorm:"type:varchar(255);not null"`
	Level     int    `json:"level" gorm:"not null"`
}

// AgentFeedback represents the schema of the agent_feedbacks table
type AgentFeedback struct {
	gorm.Model
	AgentID      uint      `json:"agent_id" gorm:"index;not null"`
	FeedbackType string    `json:"feedback_type" gorm:"type:varchar(100);not null"`
	Score        int       `json:"score" gorm:"type:int;not null"`
	Comments     string    `json:"comments" gorm:"type:text"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

// AgentKPI represents the schema of the agent_kpis table
type AgentKPI struct {
	gorm.Model
	AgentID     uint    `json:"agent_id" gorm:"index;not null"`
	KPIName     string  `json:"kpi_name" gorm:"type:varchar(255);not null"`
	Value       float64 `json:"value" gorm:"type:decimal(10,2);not null"`
	TargetValue float64 `json:"target_value" gorm:"type:decimal(10,2)"`
	Period      string  `json:"period" gorm:"type:varchar(100);not null"`
}

// AgentOnboarding represents the schema of the agent_onboardings table
type AgentOnboarding struct {
	gorm.Model
	AgentID        uint       `json:"agent_id" gorm:"index;not null"`
	OnboardingStep string     `json:"onboarding_step" gorm:"type:varchar(255);not null"`
	Status         string     `json:"status" gorm:"type:varchar(100);not null"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// AgentTeam represents the schema of the agent_teams table
type AgentTeam struct {
	gorm.Model
	TeamName    string   `json:"team_name" gorm:"type:varchar(255);not null;unique"`
	Description string   `json:"description" gorm:"type:text"`
	LeaderID    uint     `json:"leader_id" gorm:"index"`
	Members     []Agents `gorm:"many2many:agent_teams_members;" json:"-"`
}

// AgentTrainingRecord represents the schema of the agent_training_records table
type AgentTrainingRecord struct {
	gorm.Model
	AgentID          uint      `json:"agent_id" gorm:"index;not null"`
	TrainingModuleID uint      `json:"training_module_id" gorm:"index;not null"`
	Score            int       `json:"score"`
	Feedback         string    `json:"feedback" gorm:"type:text"`
	CompletedAt      time.Time `json:"completed_at"`
}

// AgentAvailability represents the schema of the agent_availabilities table
type AgentAvailability struct {
	gorm.Model
	AgentID       uint       `json:"agent_id" gorm:"index;not null"`
	Availability  string     `json:"availability" gorm:"type:varchar(100);not null"`
	LastUpdated   time.Time  `json:"last_updated"`
	NextAvailable *time.Time `json:"next_available,omitempty"`
}

// AgentContactInfo represents the schema of the agent_contact_infos table
type AgentContactInfo struct {
	gorm.Model
	AgentID      uint   `json:"agent_id" gorm:"index;not null"`
	ContactType  string `json:"contact_type" gorm:"type:varchar(100);not null"`
	ContactValue string `json:"contact_value" gorm:"type:varchar(255);not null"`
}

// AgentLoginActivity represents the schema of the agent_login_activities table
type AgentLoginActivity struct {
	gorm.Model
	AgentID    uint       `json:"agent_id" gorm:"index;not null"`
	LoginTime  time.Time  `json:"login_time"`
	LogoutTime *time.Time `json:"logout_time,omitempty"`
	IP         string     `json:"ip" gorm:"type:varchar(45)"`
}

// AgentVacation represents the schema of the agent_vacations table
type AgentVacation struct {
	gorm.Model
	AgentID    uint      `json:"agent_id" gorm:"index;not null"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Reason     string    `json:"reason" gorm:"type:text"`
	ApprovedBy uint      `json:"approved_by"`
}

// CustomerInteraction represents the schema of the customer_interactions table
type CustomerInteraction struct {
	gorm.Model
	CustomerID      uint      `json:"customer_id" gorm:"index;not null"`
	AgentID         uint      `json:"agent_id" gorm:"index;not null"`
	Channel         string    `json:"channel" gorm:"type:varchar(100);not null"`
	Content         string    `json:"content" gorm:"type:text;not null"`
	InteractionTime time.Time `json:"interaction_time"`
}

// FeedbackReview represents the schema of the feedback_reviews table
type FeedbackReview struct {
	gorm.Model
	FeedbackID uint      `json:"feedback_id" gorm:"index;not null"`
	ReviewerID uint      `json:"reviewer_id" gorm:"index;not null"`
	Review     string    `json:"review" gorm:"type:text"`
	ReviewedAt time.Time `json:"reviewed_at"`
}

// AgentTicketAssignment represents the schema of the agent_ticket_assignments table
type AgentTicketAssignment struct {
	gorm.Model
	TicketID   uint      `json:"ticket_id" gorm:"index;not null"`
	AgentID    uint      `json:"agent_id" gorm:"index;not null"`
	AssignedAt time.Time `json:"assigned_at"`
}

// AgentTrainingModule represents the schema of the agent_training_modules table
type AgentTrainingModule struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	ModuleType  string `json:"module_type" gorm:"type:varchar(100);not null"`
	Duration    int    `json:"duration"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}

// AgentCertification represents the schema of the agent_certifications table
type AgentCertification struct {
	gorm.Model
	AgentID       uint       `json:"agent_id" gorm:"index;not null"`
	Certification string     `json:"certification" gorm:"type:varchar(255);not null"`
	IssuedBy      string     `json:"issued_by" gorm:"type:varchar(255)"`
	IssuedDate    time.Time  `json:"issued_date"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
}

// AgentPerformanceReview represents the schema of the agent_performance_reviews table
type AgentPerformanceReview struct {
	gorm.Model
	AgentID    uint      `json:"agent_id" gorm:"index;not null"`
	ReviewDate time.Time `json:"review_date"`
	Score      float64   `json:"score"`
	Feedback   string    `json:"feedback" gorm:"type:text"`
}

// AgentLeaveRequest represents the schema of the agent_leave_requests table
type AgentLeaveRequest struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	LeaveType string    `json:"leave_type" gorm:"type:varchar(100);not null"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status" gorm:"type:varchar(100);not null"`
}

// AgentScheduleOverride represents the schema of the agent_schedule_overrides table
type AgentScheduleOverride struct {
	gorm.Model
	AgentID   uint      `json:"agent_id" gorm:"index;not null"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Reason    string    `json:"reason" gorm:"type:text"`
}

// AgentSkillSet represents the schema of the agent_skill_sets table
type AgentSkillSet struct {
	gorm.Model
	AgentID uint   `json:"agent_id" gorm:"index;not null"`
	Skill   string `json:"skill" gorm:"type:varchar(255);not null"`
	Level   string `json:"level" gorm:"type:varchar(100);not null"`
}

// AgentEvent represents the schema of the agent_events table
type AgentEvent struct {
	gorm.Model
	Title       *string    `gorm:"size:255;not null" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	ActionType  *string    `json:"action_type"`
	StartTime   *time.Time `json:"start_time"`
	Details     *string    `gorm:"type:text" json:"details"`
	Timestamp   time.Time  `json:"time_stamp"`
	AllDay      bool       `json:"all_day"`
	Location    *string    `gorm:"size:255" json:"location"`
	AgentID     uint       `gorm:"not null;index" json:"user_id"`
	Agents      []Agents   `gorm:"foreignKey:UserID" json:"-"`
}

func (AgentEvent) TableName() string {
	return "agent_events"
}

// AgentActivityLog represents the schema of the agent_activity_logs table
type AgentActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `json:"agent_id"`
	Activity  string    `json:"activity"`
	Timestamp time.Time `json:"timestamp"`
}

func (AgentActivityLog) TableName() string {
	return "agent_activity_logs"
}

// AgentSchedule represents the schema of the agent_schedules table
func (AgentSchedule) TableName() string {
	return "agent_schedules"
}

// AgentShift represents the schema of the agent_shifts table
func (AgentShift) TableName() string {
	return "agent_shifts"
}

// AgentSkill represents the schema of the agent_skills table
func (AgentSkill) TableName() string {
	return "agent_skills"
}

// AgentFeedback represents the schema of the agent_feedbacks table
func (AgentFeedback) TableName() string {
	return "agent_feedbacks"
}

// AgentKPI represents the schema of the agent_kpis table
func (AgentKPI) TableName() string {
	return "agent_kpis"
}

// AgentOnboarding represents the schema of the agent_onboardings table
func (AgentOnboarding) TableName() string {
	return "agent_onboardings"
}

// AgentTeam represents the schema of the agent_teams table
func (AgentTeam) TableName() string {
	return "agent_teams"
}

// AgentTrainingRecord represents the schema of the agent_training_records table
func (AgentTrainingRecord) TableName() string {
	return "agent_training_records"
}

// AgentAvailability represents the schema of the agent_availabilities table
func (AgentAvailability) TableName() string {
	return "agent_availabilities"
}

// AgentContactInfo represents the schema of the agent_contact_infos table
func (AgentContactInfo) TableName() string {
	return "agent_contact_infos"
}

// AgentLoginActivity represents the schema of the agent_login_activities table
func (AgentLoginActivity) TableName() string {
	return "agent_login_activities"
}

// AgentVacation represents the schema of the agent_vacations table
func (AgentVacation) TableName() string {
	return "agent_vacations"
}

// CustomerInteraction represents the schema of the customer_interactions table
func (CustomerInteraction) TableName() string {
	return "customer_interactions"
}

// FeedbackReview represents the schema of the feedback_reviews table
func (FeedbackReview) TableName() string {
	return "feedback_reviews"
}

// AgentTicketAssignment represents the schema of the agent_ticket_assignments table
func (AgentTicketAssignment) TableName() string {
	return "agent_ticket_assignments"
}

// AgentTrainingModule represents the schema of the agent_training_modules table
func (AgentTrainingModule) TableName() string {
	return "agent_training_modules"
}

// AgentCertification represents the schema of the agent_certifications table
func (AgentCertification) TableName() string {
	return "agent_certifications"
}

// AgentPerformanceReview represents the schema of the agent_performance_reviews table
func (AgentPerformanceReview) TableName() string {
	return "agent_performance_reviews"
}

// AgentLeaveRequest represents the schema of the agent_leave_requests table
func (AgentLeaveRequest) TableName() string {
	return "agent_leave_requests"
}

// AgentScheduleOverride represents the schema of the agent_schedule_overrides table
func (AgentScheduleOverride) TableName() string {
	return "agent_schedule_overrides"
}

// AgentSkillSet represents the schema of the agent_skill_sets table
func (AgentSkillSet) TableName() string {
	return "agent_skill_sets"
}

// GetTrainingSessionDuration returns the duration of the training session
func (ats *AgentTrainingSession) GetTrainingSessionDuration() time.Duration {
	return ats.EndDate.Sub(ats.StartDate)
}

// ChangeTrainer changes the trainer for the training session
func (ats *AgentTrainingSession) ChangeTrainer(newTrainerID uint) {
	ats.TrainerID = newTrainerID
}

// IsValid checks if the search criteria is valid
func (sc *SearchCriteria) IsValid() bool {
	// Implement validation logic for search criteria
	return true
}

// ClearCriteria resets all search criteria fields
func (sc *SearchCriteria) ClearCriteria() {
	sc.Name = ""
	sc.Role = ""
	sc.Department = ""
}

// GetShiftDuration returns the duration of the agent's shift
func (as *AgentShift) GetShiftDuration() time.Duration {
	return as.EndTime.Sub(as.StartTime)
}

// ShiftIsOver checks if the agent's shift is over
func (as *AgentShift) ShiftIsOver() bool {
	return time.Now().After(as.EndTime)
}

// GetSkillInfo returns information about the agent's skill
func (as *AgentSkill) GetSkillInfo() string {
	return fmt.Sprintf("Agent ID: %d, Skill: %s, Level: %d", as.AgentID, as.SkillName, as.Level)
}

// IncreaseSkillLevel increases the skill level of the agent
func (as *AgentSkill) IncreaseSkillLevel() {
	as.Level++
}

// IsPositiveFeedback checks if the feedback is positive
func (af *AgentFeedback) IsPositiveFeedback() bool {
	return af.Score > 3
}

// RequestFeedbackRevision requests a revision for the agent's feedback
func (af *AgentFeedback) RequestFeedbackRevision() {
	// Implement logic to request revision
}

// IsKPIAchieved checks if the agent has achieved the KPI
func (ak *AgentKPI) IsKPIAchieved() bool {
	return ak.Value >= ak.TargetValue
}

// SetTargetValue sets a new target value for the KPI
func (ak *AgentKPI) SetTargetValue(newTargetValue float64) {
	ak.TargetValue = newTargetValue
}

// IsOnboardingCompleted checks if the agent has completed onboarding
func (ao *AgentOnboarding) IsOnboardingCompleted() bool {
	return ao.Status == "completed"
}

// UpdateOnboardingStatus updates the status of the agent's onboarding
func (ao *AgentOnboarding) UpdateOnboardingStatus(newStatus string) {
	ao.Status = newStatus
}

// GetTeamMembersCount returns the number of members in the team
func (at *AgentTeam) GetTeamMembersCount() int {
	return len(at.Members)
}

// IsTrainingPassed checks if the agent has passed the training
func (atr *AgentTrainingRecord) IsTrainingPassed() bool {
	return atr.Score >= 70
}

// ProvideTrainingFeedback provides feedback for the agent's training
func (atr *AgentTrainingRecord) ProvideTrainingFeedback(feedback string) {
	atr.Feedback = feedback
}

// IsAvailableNow checks if the agent is currently available
func (aa *AgentAvailability) IsAvailableNow() bool {
	return aa.Availability == "available"
}

// SetNextAvailable sets the next available time for the agent
func (aa *AgentAvailability) SetNextAvailable(nextAvailable time.Time) {
	aa.NextAvailable = &nextAvailable
}

// GetContactInfo returns the contact information of the agent
func (aci *AgentContactInfo) GetContactInfo() string {
	return fmt.Sprintf("Agent ID: %d, Type: %s, Value: %s", aci.AgentID, aci.ContactType, aci.ContactValue)
}

// UpdateContactValue updates the contact value for the agent
func (aci *AgentContactInfo) UpdateContactValue(newValue string) {
	aci.ContactValue = newValue
}

// IsLoggedIn checks if the agent is currently logged in
func (ala *AgentLoginActivity) IsLoggedIn() bool {
	return ala.LogoutTime == nil
}

// Logout logs out the agent and sets the logout time
func (ala *AgentLoginActivity) Logout() {
	logoutTime := time.Now()
	ala.LogoutTime = &logoutTime
}

// IsOnVacation checks if the agent is currently on vacation
func (av *AgentVacation) IsOnVacation() bool {
	return time.Now().After(av.StartDate) && time.Now().Before(av.EndDate)
}

// GetInteractionSummary returns a summary of the customer interaction
func (ci *CustomerInteraction) GetInteractionSummary() string {
	return fmt.Sprintf("Customer ID: %d, Agent ID: %d, Channel: %s", ci.CustomerID, ci.AgentID, ci.Channel)
}

// GetAssignmentDetails returns details of the ticket assignment
func (ata *AgentTicketAssignment) GetAssignmentDetails() string {
	return fmt.Sprintf("Ticket ID: %d, Agent ID: %d", ata.TicketID, ata.AgentID)
}

// ReassignTicket reassigns the ticket to another agent
func (ata *AgentTicketAssignment) ReassignTicket(newAgentID uint) {
	ata.AgentID = newAgentID
}

// IsModuleActive checks if the training module is currently active
func (atm *AgentTrainingModule) IsModuleActive() bool {
	return atm.IsActive
}

// DeactivateModule deactivates the training module
func (atm *AgentTrainingModule) DeactivateModule() {
	atm.IsActive = false
}

// IsCertificationExpired checks if the agent's certification has expired
func (ac *AgentCertification) IsCertificationExpired() bool {
	if ac.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*ac.ExpiryDate)
}

// RenewCertification renews the agent's certification
func (ac *AgentCertification) RenewCertification(newExpiryDate time.Time) {
	ac.ExpiryDate = &newExpiryDate
}

// IsLeaveRequestApproved checks if the leave request is approved
func (alr *AgentLeaveRequest) IsLeaveRequestApproved() bool {
	return alr.Status == "approved"
}

// ApproveLeaveRequest approves the leave request
func (alr *AgentLeaveRequest) ApproveLeaveRequest() {
	alr.Status = "approved"
}

// GetOverrideDetails returns details of the schedule override
func (aso *AgentScheduleOverride) GetOverrideDetails() string {
	return fmt.Sprintf("Agent ID: %d, Start Date: %s, End Date: %s", aso.AgentID, aso.StartDate, aso.EndDate)
}

// GetSkillSetDetails returns details of the agent's skill set
func (ass *AgentSkillSet) GetSkillSetDetails() string {
	return fmt.Sprintf("Agent ID: %d, Skill: %s, Level: %s", ass.AgentID, ass.Skill, ass.Level)
}

// UpdateSkillLevel updates the skill level for the agent
func (ass *AgentSkillSet) UpdateSkillLevel(newLevel string) {
	ass.Level = newLevel
}

// GetEventDetails returns details of the agent's event
func (ae *AgentEvent) GetEventDetails() string {
	return fmt.Sprintf("Agent ID: %d, Title: %s, Start Time: %s", ae.AgentID, *ae.Title, ae.StartTime)
}

// UpdateEventDetails updates the details of the agent's event
func (ae *AgentEvent) UpdateEventDetails(newDetails string) {
	ae.Details = &newDetails
}

// LogActivity logs an activity for the agent
func (aal *AgentActivityLog) LogActivity(activity string) {
	aal.Activity = activity
	aal.Timestamp = time.Now()
}

// GetActivityDetails returns details of the agent's activity log
func (aal *AgentActivityLog) GetActivityDetails() string {
	return fmt.Sprintf("Agent ID: %d, Activity: %s, Timestamp: %s", aal.AgentID, aal.Activity, aal.Timestamp)
}

// AgentStorage interface consolidates all operations related to agents.
type AgentStorage interface {
	CreateAgent(agent *Agents) error
	UpdateAgent(agent *Agents) error
	DeleteAgent(agentID uint) error
	GetAllAgents() ([]*Agents, error)
	GetAgentByID(agentID uint) (*Agents, error)
	AssignRolesToAgent(agentID uint, roleNames []string) error
	RevokeRoleFromAgent(agentID uint, roleName string) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	LogAgentActivity(activity *AgentActivityLog) error
	AddAgentSchedule(schedule *AgentSchedule) error
	UpdateAgentSchedule(scheduleID uint, newSchedule *AgentSchedule) error
	DeleteAgentSchedule(scheduleID uint) error
	ListAgentSchedules(agentID uint) ([]AgentSchedule, error)
	AddAgentSkill(skill *AgentSkill) error
	UpdateAgentSkill(skillID uint, newSkill *AgentSkill) error
	DeleteAgentSkill(skillID uint) error
	ListAgentSkills(agentID uint) ([]AgentSkill, error)
	SubmitAgentFeedback(feedback *AgentFeedback) error
	ListAgentFeedback(agentID uint) ([]AgentFeedback, error)
	UpdateAgentFeedback(feedbackID uint, newFeedback *AgentFeedback) error
	DeleteAgentFeedback(feedbackID uint) error
	AddAgentCertification(certification *AgentCertification) error
	UpdateAgentCertification(certificationID uint, newCertification *AgentCertification) error
	DeleteAgentCertification(certificationID uint) error
	ListAgentCertifications(agentID uint) ([]AgentCertification, error)
	RequestAgentLeave(leaveRequest *AgentLeaveRequest) error
	UpdateAgentLeaveRequest(leaveRequestID uint, newLeaveRequest *AgentLeaveRequest) error
	CancelAgentLeaveRequest(leaveRequestID uint) error
	ListAgentLeaveRequests(agentID uint) ([]AgentLeaveRequest, error)
	UpdateAgentAvailability(availability *AgentAvailability) error
	ApproveAgentLeaveRequest(leaveRequestID uint) error
	DenyAgentLeaveRequest(leaveRequestID uint) error
	ListAgentTrainingRecords(agentID uint) ([]AgentTrainingRecord, error)
	DeleteAgentTrainingRecord(recordID uint) error
	SubmitAgentLeaveRequest(leaveRequest *AgentLeaveRequest) error
	UpdateAgentLeaveRequestStatus(requestID uint, newStatus string) error
	UpdateAgentCertificationDetails(certificationID uint, newDetails *AgentCertification) error
	ListAgentActivities(agentID uint) ([]AgentActivityLog, error)
	CreateAgentTrainingSession(trainingSession *AgentTrainingSession) error
	UpdateAgentTrainingSession(trainingSessionID uint, newTrainingSession *AgentTrainingSession) error
	DeleteAgentTrainingSession(trainingSessionID uint) error
	ListAgentTrainingSessions(agentID uint) ([]AgentTrainingSession, error)
	EnrollAgentInTrainingSession(agentID, sessionID uint) error
	RecordAgentPerformanceReview(review *AgentPerformanceReview) error
	UpdateAgentPerformanceReview(reviewID uint, newReview *AgentPerformanceReview) error
	ListAgentPerformanceReviews(agentID uint) ([]AgentPerformanceReview, error)
	DeleteAgentPerformanceReview(reviewID uint) error
	RecordCustomerInteraction(interaction *CustomerInteraction) error
	UpdateCustomerInteraction(interactionID uint, newInteraction *CustomerInteraction) error
	ListCustomerInteractionsByAgent(agentID uint) ([]CustomerInteraction, error)
	DeleteCustomerInteraction(interactionID uint) error
	AssignAgentToTrainingModule(agentID, moduleID uint) error
	UpdateAgentTrainingRecord(recordID uint, newRecord *AgentTrainingRecord) error
	RecordAgentLoginActivity(activity *AgentLoginActivity) error
	UpdateAgentProfilePic(agentID uint, profilePicURL string) error
	AddOrUpdateAgentSkillSet(agentID uint, skills []AgentSkillSet) error
	LogAction(agentID uint, actionType, details string) error
	AssignPermissionsToAgent2(agentID uint, permissionNames []string) error
	AddAgentToTeam(agentID, teamID uint) error
	AssignMultipleAgentsToTeam(agentIDs []uint, teamID uint) error
	UnassignAgentsFromTeam(agentIDs []uint, teamID uint) error
	UpdateTeamPermissions(teamID uint, permissionIDs []uint) error
	AssignPermissionsToTeamByPermissionName(teamID uint, permissionName string) error
	RevokePermissionFromTeamByPermissionName(teamID uint, permissionName string) error
	AddPermissionToRoleByPermissionName(roleName string, permissionName string) error
	RemovePermissionFromRoleByPermissionName(roleName string, permissionName string) error
}

// UnitStorage interface outlines operations for managing units.
type UnitStorage interface {
	CreateUnit(unit *Unit) error
	DeleteUnit(unitID uint) error
	UpdateUnit(unit *Unit) error
	GetUnits() ([]*Unit, error)
	GetUnitByID(unitID uint) (*Unit, error)
	GetUnitByNumber(unitNumber int) (*Unit, error)
}

// TeamStorage interface defines operations for team management.
type TeamStorage interface {
	CreateTeam(team *Teams) error
	UpdateTeam(team *Teams) error
	DeleteTeam(id uint) error
	GetTeamByID(id uint) (*Teams, error)
	GetTeams() ([]*Teams, error)
	GetTeamByNumber(teamNumber int) (*Teams, error)
}

// RoleStorage interface specifies methods for role management.
type RoleStorage interface {
	CreateRole(role *Role) error
	DeleteRole(roleID uint) error
	UpdateRole(role *Role) error
	GetRoles() ([]*Role, error)
	GetRoleByID(roleID uint) (*Role, error)
	GetRoleByNumber(roleNumber int) (*Role, error)
	AssignRolesToAgent(agentID uint, roleNames []string) error
	GetAgentRoles(agentID uint) ([]*Role, error)
	RevokeRoleFromAgent(agentID uint, roleName string) error
	GetRoleByName(name string) (*Role, error)
	AssignPermissionsToRole(roleName string, permissionNames []string) error
	RevokePermissionFromRole(roleName string, permissionName string) error
	GetRolePermissions(roleID uint) ([]*Permission, error)
	GetPermissionsByRole(roleName string) ([]*Permission, error)
}

// PermissionStorage interface encapsulates operations related to permissions.
type PermissionStorage interface {
	GetAllPermissions() ([]*Permission, error)
	CreatePermission(permission *Permission) error
	UpdatePermission(permission *Permission) error
	DeletePermission(id uint) error
	GetPermissionByID(id uint) (*Permission, error)
	GetPermissionByName(name string) (*Permission, error)
	GetPermissions() ([]*Permission, error)
}

// AgentDBModel handles database operations for Agent
type AgentDBModel struct {
	DB             *gorm.DB
	log            Logger
	EventPublisher EventPublisherImpl
}

// NewAgentDBModel creates a new instance of AgentDBModel
func NewAgentDBModel(db *gorm.DB, log Logger, eventPublisher EventPublisherImpl) *AgentDBModel {
	return &AgentDBModel{
		DB:             db,
		log:            log,
		EventPublisher: eventPublisher,
	}
}

func (repo *AgentDBModel) CreateAgent2(agent *Agents) error {
	return repo.DB.Create(agent).Error
}

// CreateAgent adds a new user to the database with transactional integrity.
func (db *AgentDBModel) CreateAgent(agent *Agents) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(agent).Error; err != nil {
			return err
		}
		return nil
	})
}

func (repo *AgentDBModel) UpdateAgent(agent *Agents) error {
	return repo.DB.Save(agent).Error
}

func (repo *AgentDBModel) DeleteAgent(id uint) error {
	return repo.DB.Delete(&Agents{}, id).Error
}

// GetAgentByID retrieves an agent by their ID with related roles preloaded.
func (repo *AgentDBModel) GetAgentByID(id uint) (*Agents, error) {
	var agent Agents
	err := repo.DB.Preload("Roles").Where("id = ?", id).First(&agent).Error
	if err != nil {
		// Assuming a logging function is available
		repo.log.Error("Error retrieving agent by ID: %v", err)
		return nil, err
	}
	return &agent, nil
}

// GetAllAgents retrieves all agents with enhanced error handling and logging.
func (repo *AgentDBModel) GetAllAgents() ([]*Agents, error) {
	var agents []*Agents
	if err := repo.DB.Preload("Roles").Find(&agents).Error; err != nil {
		repo.log.Error("Failed to retrieve all agents: %v", err)
		return nil, err
	}
	return agents, nil
}

// GetAgentRoles fetches roles associated with an agent with enhanced error handling.
func (repo *AgentDBModel) GetAgentRoles(agentID uint) ([]*Role, error) {
	var roles []*Role
	if err := repo.DB.Table("agent_roles").Select("roles.*").
		Joins("join roles on roles.id = agent_roles.role_id").
		Where("agent_roles.agent_id = ?", agentID).Scan(&roles).Error; err != nil {
		repo.log.Error("Failed to retrieve agent roles: %v", err)
		return nil, err
	}
	return roles, nil
}

// LogAgentActivity encapsulates the logging of an agent's activity within a transaction.
func (db *AgentDBModel) LogAgentActivity(activity *AgentActivityLog) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(activity).Error; err != nil {
			db.log.Error("Failed to log agent activity: %v", err)
			return err
		}
		return nil
	})
}

// GetAgentByNumber retrieves an Agent by their agent number with proper error handling and logging.
func (as *AgentDBModel) GetAgentByNumber(agentNumber int) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("agent_number = ?", agentNumber).First(&agent).Error
	if err != nil {
		as.log.Error("Error retrieving agent by number: %v", err)
		return nil, err
	}
	return &agent, nil
}

// AssignRoleToAgent assigns a single role to an agent, ensuring transactional integrity and logging.
func (db *AgentDBModel) AssignRoleToAgent(agentID uint, roleIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var agent Agents
		if err := tx.First(&agent, agentID).Error; err != nil {
			db.log.Error("Error finding agent by ID during role assignment: %v", err)
			return err
		}
		if err := tx.Model(&agent).Association("Roles").Replace(roleIDs); err != nil {
			db.log.Error("Error assigning role to agent: %v", err)
			return err
		}
		return nil
	})
}

// AssignRolesToAgent assigns roles to an agent by role names.
func (repo *AgentDBModel) AssignRolesToAgent(agentID uint, roleNames []string) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		var roles []Role
		if err := tx.Where("name IN ?", roleNames).Find(&roles).Error; err != nil {
			repo.log.Error("Error finding roles by names: %v", err)
			return err
		}
		if err := tx.Model(&Agents{ID: agentID}).Association("Roles").Replace(roles); err != nil {
			repo.log.Error("Error assigning roles to agent: %v", err)
			return err
		}
		return nil
	})
}

// RevokeRoleFromAgent removes a role from an agent with improved transactional handling and error logging.
func (repo *AgentDBModel) RevokeRoleFromAgent(agentID uint, roleName string) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			repo.log.Error("Failed to find role by name: %v", err)
			return err
		}
		if err := tx.Model(&Agents{ID: agentID}).Association("Roles").Delete(&role); err != nil {
			repo.log.Error("Failed to revoke role from agent: %v", err)
			return err
		}
		return nil
	})
}

// AssignRolesToAgentByIDOrName assigns roles to an agent by either ID or name, with full transaction support and logging.
func (db *AgentDBModel) AssignRolesToAgentByIDOrName(agentID uint, roleIdentifiers []interface{}) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var agent Agents
		if err := tx.First(&agent, agentID).Error; err != nil {
			db.log.Error("Failed to find agent during role assignment: %v", err)
			return err
		}
		var roles []Role
		for _, identifier := range roleIdentifiers {
			switch id := identifier.(type) {
			case uint:
				var role Role
				if err := tx.Where("id = ?", id).First(&role).Error; err != nil {
					db.log.Error("Failed to find role by ID: %v", err)
					return err
				}
				roles = append(roles, role)
			case string:
				var role Role
				if err := tx.Where("name = ?", id).First(&role).Error; err != nil {
					db.log.Error("Failed to find role by name: %v", err)
					return err
				}
				roles = append(roles, role)
			default:
				db.log.Error("Unsupported role identifier type")
				return fmt.Errorf("unsupported role identifier type")
			}
		}
		if err := tx.Model(&agent).Association("Roles").Replace(roles); err != nil {
			db.log.Error("Failed to assign roles to agent: %v", err)
			return err
		}
		return nil
	})
}

// GetRolesByAgent retrieves all roles associated with an Agent.
func (as *AgentDBModel) GetRolesByAgent(agentID uint) ([]Role, error) {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}
	return agent.Roles, nil
}

// RevokePermissionFromRole demonstrates enhanced transactional integrity and error handling when revoking a permission from a role.
func (as *AgentDBModel) RevokePermissionFromRole(roleName string, permissionName string) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		var permission Permission
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			as.log.Error("Error finding role by name: %v", err)
			return err
		}
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			as.log.Error("Error finding permission by name: %v", err)
			return err
		}
		if err := tx.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).Delete(&RolePermission{}).Error; err != nil {
			as.log.Error("Error revoking permission from role: %v", err)
			return err
		}
		return nil
	})
}

// CreateUnit handles the creation of a new unit with transactional integrity and logging.
func (as *AgentDBModel) CreateUnit(unit *Unit) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(unit).Error; err != nil {
			as.log.Error("Failed to create unit: %v", err)
			return err
		}
		return nil
	})
}

// UpdateUnit manages the update of an existing unit within a transaction, including error handling and logging.
func (as *AgentDBModel) UpdateUnit(unit *Unit) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(unit).Error; err != nil {
			as.log.Error("Failed to update unit: %v", err)
			return err
		}
		return nil
	})
}

// DeleteUnit removes a unit from the database with robust transactional support and comprehensive logging.
func (as *AgentDBModel) DeleteUnit(id uint) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Unit{}, id).Error; err != nil {
			as.log.Error("Failed to delete unit: %v", err)
			return err
		}
		return nil
	})
}

// GetUnits retrieves all units from the database, with improved error handling and logging.
func (as *AgentDBModel) GetUnits() ([]*Unit, error) {
	var units []*Unit
	if err := as.DB.Find(&units).Error; err != nil {
		as.log.Error("Failed to retrieve units: %v", err)
		return nil, err
	}
	return units, nil
}

// GetUnitByID fetches a unit by its ID, incorporating error handling and logging for improved robustness.
func (as *AgentDBModel) GetUnitByID(id uint) (*Unit, error) {
	var unit Unit
	if err := as.DB.Where("id = ?", id).First(&unit).Error; err != nil {
		as.log.Error("Failed to retrieve unit by ID: %v", err)
		return nil, err
	}
	return &unit, nil
}

// CreateTeam adds a new team to the database, ensuring transactional integrity and logging any errors encountered.
func (db *AgentDBModel) CreateTeam(team *Teams) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(team).Error; err != nil {
			db.log.Error("Failed to create team: %v", err)
			return err
		}
		return nil
	})
}

// UpdateTeam updates an existing team's details, wrapping the operation in a transaction for integrity and logging errors.
func (db *AgentDBModel) UpdateTeam(team *Teams) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(team).Error; err != nil {
			db.log.Error("Failed to update team: %v", err)
			return err
		}
		return nil
	})
}

// AssignAgentToTeam manages the assignment of an agent to a team, ensuring the operation is transactional and errors are logged.
func (db *AgentDBModel) AssignAgentToTeam(agentID, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var teamAgent TeamAgent
		err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&teamAgent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			teamAgent = TeamAgent{AgentID: agentID, TeamID: teamID}
			if err := tx.Create(&teamAgent).Error; err != nil {
				db.log.Error("Failed to assign agent to team: %v", err)
				return err
			}
		} else if err != nil {
			db.log.Error("Unexpected error during agent-team assignment: %v", err)
			return err
		}
		return nil
	})
}

// RemoveAgentFromTeam orchestrates the removal of an agent from a team within a transaction, complete with error handling and logging.
func (db *AgentDBModel) RemoveAgentFromTeam(agentID, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error; err != nil {
			db.log.Error("Failed to remove agent from team: %v", err)
			return err
		}
		return nil
	})
}

// GetTeamByID retrieves a team by its ID, incorporating error checks and logging for improved maintainability.
func (as *AgentDBModel) GetTeamByID(id uint) (*Teams, error) {
	var team Teams
	if err := as.DB.Where("id = ?", id).First(&team).Error; err != nil {
		as.log.Error("Failed to retrieve team by ID: %v", err)
		return nil, err
	}
	return &team, nil
}

// DeleteTeam deletes a team from the database.
func (as *AgentDBModel) DeleteTeam(id uint) error {
	return as.DB.Delete(&Teams{}, id).Error
}

// GetTeams retrieves all teams from the database.
func (as *AgentDBModel) GetTeams() ([]*Teams, error) {
	var teams []*Teams
	err := as.DB.Find(&teams).Error
	return teams, err
}

// GetTeamByID retrieves a team by its ID.
func (as *AgentDBModel) GetTeamByID2(id uint) (*Teams, error) {
	var team Teams
	err := as.DB.Where("id = ?", id).First(&team).Error
	return &team, err
}

// AssignAgentToTeam2 enhances the assignment of an agent to a team with additional checks and logging.
func (db *AgentDBModel) AssignAgentToTeam2(agentID, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var teamAgent TeamAgent
		// Check if the agent is already assigned to the team to prevent duplicate entries
		err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&teamAgent).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// The assignment does not exist; proceed to create a new one
			teamAgent = TeamAgent{AgentID: agentID, TeamID: teamID}
			if err := tx.Create(&teamAgent).Error; err != nil {
				db.log.Error("Failed to assign agent to team: %v", err)
				return err // Error encountered while creating the assignment; rollback the transaction
			}
		} else if err != nil {
			db.log.Error("Unexpected error during agent-team assignment check: %v", err)
			return err // An unexpected error occurred; rollback the transaction
		}
		// The assignment already exists or has been successfully created
		return nil
	})
}

// RemoveAgentFromTeam2 refines the process of removing an agent from a specific team, with enhanced error handling and logging.
func (db *AgentDBModel) RemoveAgentFromTeam2(agentID, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error; err != nil {
			db.log.Error("Failed to remove agent from team: %v", err)
			return err // Error encountered while removing the assignment; rollback the transaction
		}
		// Successfully removed the assignment
		return nil
	})
}

// RemoveAgentFromTeams facilitates the safe removal of an agent from multiple teams with proper transactional control and logging.
func (db *AgentDBModel) RemoveAgentFromTeams(agentID uint, teamIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, teamID := range teamIDs {
			if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).Delete(&TeamAgent{}).Error; err != nil {
				db.log.Error("Failed to remove agent from team %d: %v", teamID, err)
				return err // Rollback transaction on error
			}
		}
		return nil // Commit transaction if all deletions succeed
	})
}

// GetTeamsByAgent enhances the retrieval of all teams associated with an agent, including error handling and logging.
func (model *AgentDBModel) GetTeamsByAgent(agentID uint) ([]Teams, error) {
	var teams []Teams
	err := model.DB.Joins("JOIN team_agents ON teams.id = team_agents.team_id").
		Where("team_agents.agent_id = ?", agentID).Find(&teams).Error
	if err != nil {
		model.log.Error("Failed to retrieve teams by agent: %v", err)
		return nil, err
	}
	return teams, nil
}

// UpdateAgentPermissions enhances the update of an agent's permissions, ensuring transactional integrity, error handling, and logging.
func (db *AgentDBModel) UpdateAgentPermissions(agentID uint, newPermissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// First, remove any existing permissions that are not in the newPermissionIDs list
		if err := tx.Where("agent_id = ? AND permission_id NOT IN ?", agentID, newPermissionIDs).Delete(&AgentPermission{}).Error; err != nil {
			db.log.Error("Failed to update agent permissions: %v", err)
			return err // Rollback on error
		}

		// Next, add new permissions from newPermissionIDs that the agent does not already have
		for _, permissionID := range newPermissionIDs {
			var existing AgentPermission
			result := tx.Where("agent_id = ? AND permission_id = ?", agentID, permissionID).First(&existing)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Permission does not exist for this agent; proceed to add it
				if err := tx.Create(&AgentPermission{AgentID: agentID, PermissionID: permissionID}).Error; err != nil {
					db.log.Error("Failed to add new permission to agent: %v", err)
					return err // Rollback on error
				}
			} // Ignore any found records, as we do not need to add existing permissions
		}

		return nil // Commit the transaction
	})
}

// AssignRolesToAgent2 refines the process of assigning roles to an agent, with comprehensive transaction management, error handling, and logging.
func (as *AgentDBModel) AssignRolesToAgent2(agentID uint, roleNames []string) error {
	if len(roleNames) == 0 {
		return fmt.Errorf("no role names provided")
	}

	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		as.log.Error("Failed to retrieve agent by ID during role assignment: %v", err)
		return err
	}

	roles := make([]Role, len(roleNames))
	for i, roleName := range roleNames {
		role, err := as.GetRoleByName(roleName)
		if err != nil {
			as.log.Error("Failed to retrieve role by name '%s': %v", roleName, err)
			return err
		}
		roles[i] = *role
	}

	agent.Roles = roles

	return as.UpdateAgent(agent)
}

// RevokeRoleFromAgent2 enhances the process of revoking a role from an agent, with error handling and logging.
func (as *AgentDBModel) RevokeRoleFromAgent2(agentID uint, roleName string) error {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		as.log.Error("Failed to retrieve agent by ID during role revocation: %v", err)
		return err
	}

	var updatedRoles []Role
	for _, role := range agent.Roles {
		if role.RoleName != roleName {
			updatedRoles = append(updatedRoles, role)
		}
	}
	agent.Roles = updatedRoles

	return as.UpdateAgent(agent)
}

// GetAgentPermissions enhances the retrieval of all permissions associated with an agent's roles, including error handling and logging.
func (as *AgentDBModel) GetAgentPermissions(agentID uint) ([]*Permission, error) {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		as.log.Error("Failed to retrieve agent by ID during permission retrieval: %v", err)
		return nil, err
	}

	var permissions []*Permission
	for _, role := range agent.Roles {
		rolePermissions, err := as.GetPermissionsByRole(role.RoleName)
		if err != nil {
			as.log.Error("Failed to retrieve permissions by role during agent permission retrieval: %v", err)
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions, nil
}

// GetPermissionsByAgent enhances the retrieval of all permissions associated with an agent, ensuring efficient query handling and error logging.
func (model *AgentDBModel) GetPermissionsByAgent(agentID uint) ([]Permission, error) {
	var permissions []Permission
	err := model.DB.Table("permissions").
		Joins("JOIN agent_permissions ON permissions.id = agent_permissions.permission_id").
		Where("agent_permissions.agent_id = ?", agentID).Scan(&permissions).Error
	if err != nil {
		model.log.Error("Failed to retrieve permissions for agent %d: %v", agentID, err)
		return nil, err
	}
	return permissions, nil
}

// CreateRolePermission creates a new role-permission association with proper error handling and transactional integrity.
func (as *AgentDBModel) CreateRolePermission(roleID uint, permissionID uint) error {
	rolePermission := RolePermission{RoleID: roleID, PermissionID: permissionID}
	return as.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&rolePermission).Error; err != nil {
			as.log.Error("Failed to create role-permission association: %v", err)
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// GrantPermissionToAgent grants a specific permission to an agent, ensuring the process is managed within a transaction for data integrity.
func (model *AgentDBModel) GrantPermissionToAgent(agentID, permissionID uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		agentPermission := AgentPermission{AgentID: agentID, PermissionID: permissionID}
		if err := tx.FirstOrCreate(&agentPermission, AgentPermission{AgentID: agentID, PermissionID: permissionID}).Error; err != nil {
			model.log.Error("Failed to grant permission to agent: %v", err)
			return err // Rollback on error
		}
		return nil // Commit if successful
	})
}

// RevokePermissionFromAgent refines the process of revoking a specific permission from an agent, incorporating transactional control and enhanced logging.
func (model *AgentDBModel) RevokePermissionFromAgent(agentID, permissionID uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("agent_id = ? AND permission_id = ?", agentID, permissionID).Delete(&AgentPermission{}).Error; err != nil {
			model.log.Error("Failed to revoke permission from agent: %v", err)
			return err // Rollback on error
		}
		return nil // Commit if successful
	})
}

// UpdateAgentRolesAndPermissions updates both roles and permissions for an agent, ensuring that all operations are performed within a single transaction for consistency.
func (db *AgentDBModel) UpdateAgentRolesAndPermissions(agentID uint, newRoleIDs, newPermissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Update roles
		if err := db.updateAgentRolesWithinTransaction(tx, agentID, newRoleIDs); err != nil {
			db.log.Error("Failed to update agent roles: %v", err)
			return err
		}
		// Update permissions
		if err := db.updateAgentPermissionsWithinTransaction(tx, agentID, newPermissionIDs); err != nil {
			db.log.Error("Failed to update agent permissions: %v", err)
			return err
		}
		return nil // Commit transaction if both updates succeed
	})
}

func (db *AgentDBModel) updateAgentRolesWithinTransaction(tx *gorm.DB, agentID uint, newRoleIDs []uint) error {
	// Logic to update agent's roles based on newRoleIDs
	// This function assumes that roles are directly associated with the agent and updates accordingly.
	if err := tx.Model(&Agents{ID: agentID}).Association("Roles").Replace(newRoleIDs); err != nil {
		db.log.Error("Failed to update roles within transaction: %v", err)
		return err
	}
	return nil
}

func (db *AgentDBModel) updateAgentPermissionsWithinTransaction(tx *gorm.DB, agentID uint, newPermissionIDs []uint) error {
	// Similar logic to UpdateAgentPermissions example provided earlier
	// Remove old permissions not in the new list and add new permissions
	if err := tx.Where("agent_id = ? AND permission_id NOT IN ?", agentID, newPermissionIDs).Delete(&AgentPermission{}).Error; err != nil {
		db.log.Error("Failed to remove old permissions within transaction: %v", err)
		return err
	}
	for _, permissionID := range newPermissionIDs {
		if err := tx.FirstOrCreate(&AgentPermission{AgentID: agentID, PermissionID: permissionID}).Error; err != nil {
			db.log.Error("Failed to add new permission within transaction: %v", err)
			return err
		}
	}
	return nil
}

// AssignPermissionsToAgent refines the process of assigning permissions to an agent by ensuring the operation is performed within a transaction.
func (repo *AgentDBModel) AssignPermissionsToAgent(agentID uint, permissionNames []string, publish func(event interface{}) error) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		var permissions []Permission
		if err := tx.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
			repo.log.Error("Failed to find permissions during assignment: %v", err)
			return err
		}

		for _, perm := range permissions {
			ap := AgentPermission{
				AgentID:      agentID,
				PermissionID: perm.ID,
			}
			// Avoid duplicating permissions
			if err := tx.FirstOrCreate(&ap, AgentPermission{AgentID: agentID, PermissionID: perm.ID}).Error; err != nil {
				repo.log.Error("Failed to assign permission to agent within transaction: %v", err)
				return err
			}

			// Asynchronously publish an event for each permission assignment
			go func(perm Permission) {
				event := struct {
					AgentID      uint
					PermissionID uint
					Message      string
				}{
					AgentID:      agentID,
					PermissionID: perm.ID,
					Message:      "Permission assigned to agent",
				}

				if err := publish(event); err != nil {
					// Log error without failing the transaction
					repo.log.Error("Failed to publish permission assignment event: %v", err)
				}
			}(perm)
		}
		return nil
	})
}

// AssignPermissionsToAgent assigns permissions to an agent's roles.
func (as *AgentDBModel) AssignPermissionsToAgent3(agentID uint, permissionNames []string) error {
	if len(permissionNames) == 0 {
		return fmt.Errorf("no permission names provided")
	}

	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {
		permission, err := as.GetPermissionByName(permissionName)
		if err != nil {
			// If permission doesn't exist, create it
			newPermission := &Permission{
				Name:        permissionName,
				Description: nil,
			}
			err := as.CreatePermission(newPermission)
			if err != nil {
				return err
			}
			permission, err = as.GetPermissionByName(permissionName)
			if err != nil {
				return err
			}
		}

		// Check if the permission is already assigned
		assigned := false
		for _, role := range agent.Roles {
			rolePermissions, err := as.GetRolePermissionPairs(role.ID)
			if err != nil {
				return err
			}
			for _, rp := range rolePermissions {
				if rp.ID == permission.ID {
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}

		// If not assigned, assign it to all roles of the agent
		if !assigned {
			for i := range agent.Roles {
				rolePermissions := &RolePermission{
					RoleID:       agent.Roles[i].ID,
					PermissionID: permission.ID,
				}
				err := as.CreateRolePermission(rolePermissions.RoleID, rolePermissions.PermissionID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RevokePermissionFromAgent revokes a permission from an agent.
func (as *AgentDBModel) RevokePermissionFromAgent2(agentID uint, permissionName string) error {
	agent, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	for _, role := range agent.Roles {
		err := as.RevokePermissionFromRole(role.RoleName, permission.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateAgentWithRolesAndPermissions updates agent details, roles, and permissions atomically.
func (model *AgentDBModel) UpdateAgentWithRolesAndPermissions(agent *Agents, roleIDs, permissionIDs []uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// Update agent details
		if err := tx.Save(agent).Error; err != nil {
			return err
		}

		// Update agent roles
		if err := tx.Model(&agent).Association("Roles").Replace(roleIDs); err != nil {
			return err
		}

		// Update agent permissions
		for _, permissionID := range permissionIDs {
			if err := tx.Model(&AgentPermission{}).Where("agent_id = ?", agent.ID).Update("permission_id", permissionID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// AssignAgentToMultipleTeams assigns an agent to multiple teams, ensuring no duplicates.
func (model *AgentDBModel) AssignAgentToMultipleTeams(agentID uint, teamIDs []uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		for _, teamID := range teamIDs {
			// Check if the assignment already exists
			exists := tx.Model(&TeamAgent{}).Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&TeamAgent{}).Error
			if errors.Is(exists, gorm.ErrRecordNotFound) {
				// Create new assignment if it does not exist
				if err := tx.Create(&TeamAgent{AgentID: agentID, TeamID: teamID}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// UpdateAgentPermissions updates the list of permissions for an agent.
func (db *AgentDBModel) UpdateAgentPermissions2(agentID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Remove permissions not in the new list
		if err := tx.Where("agent_id = ? AND permission_id NOT IN (?)", agentID, permissionIDs).Delete(&AgentPermission{}).Error; err != nil {
			return err
		}

		// Find existing permissions to avoid duplicates
		var existingPermissions []AgentPermission
		tx.Where("agent_id = ?", agentID).Find(&existingPermissions)

		// Map for quick lookup
		existingMap := make(map[uint]bool)
		for _, p := range existingPermissions {
			existingMap[p.PermissionID] = true
		}

		// Add new permissions
		for _, pid := range permissionIDs {
			if !existingMap[pid] {
				if err := tx.Create(&AgentPermission{AgentID: agentID, PermissionID: pid}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// contains checks if a slice contains a specific uint value.
func Contains(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// containsPermission checks if a slice of AgentPermission contains a specific PermissionID.
func ContainsPermission(s []AgentPermission, e uint) bool {
	for _, a := range s {
		if a.PermissionID == e {
			return true
		}
	}
	return false
}

// ... Rest of your code ...

func (as *AgentDBModel) Create(agent *Agents) error {
	return as.DB.Create(agent).Error
}

func (as *AgentDBModel) Update(agent *Agents) error {
	return as.DB.Save(agent).Error
}

func (as *AgentDBModel) Delete(agentID uint) error {
	return as.DB.Delete(&Agents{}, agentID).Error
}

func (as *AgentDBModel) GetAll() ([]*Agents, error) {
	var agents []*Agents
	err := as.DB.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func (as *AgentDBModel) GetByID(agentID uint) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("id = ?", agentID).First(&agent).Error
	return &agent, err
}

func (as *AgentDBModel) GetByNumber(agentNumber int) (*Agents, error) {
	var agent Agents
	err := as.DB.Where("agent_id = ?", agentNumber).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// AssignRoles assigns roles to an agent.
func (as *AgentDBModel) AssignRoles(agentID uint, roleNames []string) error {
	if len(roleNames) == 0 {
		return fmt.Errorf("no role names provided")
	}

	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	var roles []Role
	for _, roleName := range roleNames {
		role, err := as.GetRoleByName(roleName)
		if err != nil {
			return err
		}
		roles = append(roles, *role)
	}

	agent.Roles = roles

	return as.Update(agent)
}

// RevokeRole revokes a role from an agent.
func (as *AgentDBModel) RevokeRole(agentID uint, roleName string) error {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	var updatedRoles []Role
	for _, role := range agent.Roles {
		if role.RoleName != roleName {
			updatedRoles = append(updatedRoles, role)
		}
	}
	agent.Roles = updatedRoles

	return as.Update(agent)
}

// GetPermissions retrieves all permissions associated with an agent's roles.
func (as *AgentDBModel) GetPermissions(agentID uint) ([]*Permission, error) {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return nil, err
	}

	var permissions []*Permission
	for _, role := range agent.Roles {
		rolePermissions, err := as.GetRolePermissions(role.ID)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions, nil
}

// GetRolePermissions retrieves all permissions associated with a role.
func (as *AgentDBModel) GetRolePermissions(roleID uint) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Where("rolePermissions.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetPermissionsByRole retrieves all permissions associated with a role.
func (as *AgentDBModel) GetPermissionsByRole(roleName string) ([]*Permission, error) {
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return nil, err
	}

	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err = as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Where("rolePermissions.role_id = ?", role.ID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// AssignPermissions assigns permissions to an agent's roles.
func (as *AgentDBModel) AssignPermissions(agentID uint, permissionNames []string) error {
	if len(permissionNames) == 0 {
		return fmt.Errorf("no permission names provided")
	}

	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {
		permission, err := as.GetPermissionByName(permissionName)
		if err != nil {
			// If permission doesn't exist, create it
			newPermission := &Permission{
				Name:        permissionName,
				Description: nil,
			}
			err := as.CreatePermission(newPermission)
			if err != nil {
				return err
			}
			permission, err = as.GetPermissionByName(permissionName)
			if err != nil {
				return err
			}
		}

		// Check if the permission is already assigned
		assigned := false
		for _, role := range agent.Roles {
			rolePermissions, err := as.GetRolePermissionPairs(role.ID)
			if err != nil {
				return err
			}
			for _, rp := range rolePermissions {
				if rp.ID == permission.ID {
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}

		// If not assigned, assign it to all roles of the agent
		if !assigned {
			for i := range agent.Roles {
				rolePermissions := &RolePermission{
					RoleID:       agent.Roles[i].ID,
					PermissionID: permission.ID,
				}
				err := as.CreateRolePermission(rolePermissions.PermissionID, rolePermissions.PermissionID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RevokePermission revokes a permission from an agent.
func (as *AgentDBModel) RevokePermission(agentID uint, permissionName string) error {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	for _, role := range agent.Roles {
		err := as.RevokePermissionFromRole(role.RoleName, permission.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAssignedRoles retrieves all roles assigned to an agent.
func (as *AgentDBModel) GetAssignedRoles(agentID uint) ([]*Role, error) {
	agent, err := as.GetByID(agentID)
	if err != nil {
		return nil, err
	}

	var roles []*Role
	for _, role := range agent.Roles {
		roles = append(roles, &role)
	}

	return roles, nil
}

// GetAssignedPermissions retrieves all permissions assigned to an agent.
func (as *AgentDBModel) GetAssignedPermissions(agentID uint) ([]*Permission, error) {
	permissions, err := as.GetPermissions(agentID)
	if err != nil {
		return nil, err
	}

	var assignedPermissions []*Permission
	for i, permission := range permissions {
		assignedPermissions = append(assignedPermissions, permission)
		fmt.Printf("%v permission added", assignedPermissions[i].Name)
	}

	return assignedPermissions, nil
}

// CreateRole adds a new role to the database with transactional integrity.
func (db *AgentDBModel) CreateRole(role *Role) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// AssignPermissionsToRole assigns permissions to a role, handling it transactionally to ensure data integrity.
func (db *AgentDBModel) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err // Error fetching role, rollback
		}

		// Assuming `permissionIDs` is a slice of permission IDs to be associated with the role
		if err := tx.Model(&role).Association("Permissions").Replace(permissionIDs); err != nil {
			return err // Error assigning permissions, rollback
		}
		return nil // Success, commit the transaction
	})
}

// UpdateRole updates an existing role's details, including its permissions, within a transaction.
func (db *AgentDBModel) UpdateRoleWithPermissions(role *Role, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Update role details
		if err := tx.Save(role).Error; err != nil {
			return err // Rollback on error
		}

		// Update permissions associated with the role
		if err := tx.Model(role).Association("Permissions").Replace(permissionIDs); err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// GetRoleByID retrieves a role by its ID.
func (as *AgentDBModel) GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := as.DB.Where("id = ?", id).First(&role).Error
	return &role, err
}

// GetRoleByName retrieves a role by its name.
func (as *AgentDBModel) GetRoleByName(name string) (*Role, error) {
	var role Role
	err := as.DB.Where("name = ?", name).First(&role).Error
	return &role, err
}

// DeleteRole deletes a role from the database.
func (as *AgentDBModel) DeleteRole(id uint) error {
	return as.DB.Delete(&Role{}, id).Error
}

// CreatePermission creates a new permission.
func (as *AgentDBModel) CreatePermission(permission *Permission) error {
	return as.DB.Create(permission).Error
}

// GetPermissionByID retrieves a permission by its ID.
func (as *AgentDBModel) GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("id = ?", id).First(&permission).Error
	return &permission, err
}

// GetPermissionByName retrieves a permission by its name.
func (as *AgentDBModel) GetPermissionByName(name string) (*Permission, error) {
	var permission Permission
	err := as.DB.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

// UpdatePermission updates the details of an existing permission.
func (as *AgentDBModel) UpdatePermission(permission *Permission) error {
	return as.DB.Save(permission).Error
}

// DeletePermission deletes a permission from the database.
func (as *AgentDBModel) DeletePermission(id uint) error {
	return as.DB.Delete(&Permission{}, id).Error
}

// GetRolePermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetRolePermissionPairs(roleID uint) ([]*RolePermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var rolePermissions []*RolePermission
	err := as.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

// CreateUserRole creates a new user-role association.
func (as *AgentDBModel) CreateUserRoleBase(agentID uint, roleID uint) error {
	userRole := AgentRole{AgentID: agentID, RoleID: roleID}
	return as.DB.Create(&userRole).Error
}

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateRoleBasePermission(roleID uint, permissionID uint) error {
	rolePermission := RolePermission{RoleID: roleID, PermissionID: permissionID}
	return as.DB.Create(&rolePermission).Error
}

// GetUserRoles returns a list of roles assigned to a user.
func (as *AgentDBModel) GetAgentRoles2(agentID uint) ([]*Role, error) {
	var roles []*Role
	err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
		Where("userRoles.user_id = ?", agentID).
		Find(&roles).Error
	return roles, err
}

// GetRolesByUser retrieves all roles associated with a user.
func (as *AgentDBModel) GetRoleByAgent(agentID uint) ([]*Role, error) {
	// Implement logic to retrieve roles associated with the user.
	// You might need to join userRoles and roles tables.

	var roles []*Role
	err := as.DB.Joins("JOIN userRoles ON roles.id = userRoles.role_id").
		Where("userRoles.user_id = ?", agentID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRoles retrieves all roles from the database.
func (as *AgentDBModel) GetRoles() ([]*Role, error) {
	var roles []*Role
	err := as.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByNumber retrieves a role by its role number.
func (as *AgentDBModel) GetRoleByNumber(roleNumber int) (*Role, error) {
	var role Role
	err := as.DB.Where("role_id = ?", roleNumber).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAllRoles retrieves all roles.
func (as *AgentDBModel) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	err := as.DB.Find(&roles).Error
	return roles, err
}

// GetUnitByNumber retrieves a unit by its unit number.
func (as *AgentDBModel) GetTeamByNumber(teamNumber int) (*Teams, error) {
	var team Teams
	err := as.DB.Where("team_id = ?", teamNumber).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

// AssignPermissionsToRole assigns multiple permissions to a role.
func (as *AgentDBModel) AssignPermissionsToRoleBase(roleName string, permissionNames []string) error {
	var permission Permission
	var errors []error

	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	for _, permissionName := range permissionNames {

		p, er := as.GetPermissionByName(permissionName)
		if er != nil {
			permis := &Permission{
				Name:        permissionName,
				Description: nil,
			}
			err := as.CreatePermission(permis)
			if err != nil {
				return err
			}
			p, er := as.GetPermissionByName(permissionName)
			if er != nil {
				return er
			}
			permission = *p
		} else {
			permission = *p
		}

		erro := as.AssociatePermissionWithRole(role.ID, permission.Name, permission.ID)
		if err != nil {
			errors = append(errors, erro)
			continue
		}
	}
	fmt.Println(errors)
	return nil
}

// AssociatePermissionWithRole associates a permission with a role.
func (as *AgentDBModel) AssociatePermissionWithRole(roleID uint, permissionName string, permissionID uint) error {
	role, err := as.GetRoleByID(roleID)
	var permission *Permission
	if err != nil {
		return err
	}

	if permissionID <= 0 {
		p, err := as.GetPermissionByName(permissionName)
		if err != nil {
			return err
		}
		permission = p

	} else if permissionName == "" {
		p, err := as.GetPermissionByID(permissionID)
		if err != nil {
			return err
		}
		permission = p

	} else if permissionName == "" || permissionID <= 0 {
		return fmt.Errorf("the permission you selected does not exist")

	} else if permissionName != "" || permissionID > 0 {
		p, err := as.GetPermissionByID(permissionID)
		if err != nil {
			return err
		}
		permission = p
	}

	permittedRole := as.CreateRoleBasePermission(role.ID, permission.ID)

	// Implement logic to associate the permission with the role.
	// You might have a separate table or method to manage role-permission relationships.
	// Example: rolePermissions table with columns (roleID, permissionID).

	if permittedRole != nil {
		return permittedRole
	}

	return nil
}

// AssociatePermissionWithRole associates a permission with a team.
func (as *AgentDBModel) GrantPermissionToTeam(permission *Permission, teamID uint) error {
	tp, err := as.GetTeamPermission(teamID)
	if err != nil {
		er := as.CreateTeamPermission(teamID, permission)
		if er != nil {
			return fmt.Errorf("unable to assign permission to team")
		}
	}
	tp.Permissions = append(tp.Permissions, permission)

	return nil
}

// GetTeamPermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetTeamPermission(teamID uint) (*TeamPermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var teamPermissions *TeamPermission
	err := as.DB.Where("team_id = ?", teamID).Find(&teamPermissions).Error
	if err != nil {
		return nil, err
	}

	return teamPermissions, nil
}

// GetTeamPermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetTeamPermissions(teamID uint) ([]*TeamPermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var teamPermissions []*TeamPermission
	err := as.DB.Where("team_id = ?", teamID).Find(&teamPermissions).Error
	if err != nil {
		return nil, err
	}

	return teamPermissions, nil
}

// CreateRolePermission creates a new role-permission association.
func (as *AgentDBModel) CreateTeamPermission(teamID uint, permission *Permission) error {
	var p []*Permission
	p = append(p, permission)

	teamPermission := TeamPermission{TeamID: teamID, Permissions: p}
	return as.DB.Create(&teamPermission).Error
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) RevokePermissionFromTeam(teamID uint, permissionID uint) error {
	role, err := as.GetTeamByID(teamID)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err = as.DB.Where("team_id = ? AND permission = ?", role.ID, permission).
		Delete(&TeamPermission{}).Error
	if err != nil {
		return err
	}

	return nil
}

// GetRolePermissionPairs retrieves all role-permission pairs associated with a role.
func (as *AgentDBModel) GetRolePermissionPairs2(roleID uint) ([]*RolePermission, error) {
	// Implement logic to retrieve role-permission pairs associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var rolePermissions []*RolePermission
	err := as.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

// GetRolePermissions retrieves all permissions associated with a role.
func (as *AgentDBModel) GetRolePermissions2(roleID uint) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Where("rolePermissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetAllPermissions retrieves all permissions.
func (as *AgentDBModel) GetAllPermissions() ([]*Permission, error) {
	var permissions []*Permission
	err := as.DB.Find(&permissions).Error
	return permissions, err
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) RevokePermissionFromRole2(roleName string, permissionName string) error {
	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	permission, err := as.GetPermissionByName(permissionName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err = as.DB.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
		Delete(&RolePermission{}).Error
	if err != nil {
		return err
	}

	return nil
}

// GetPermissionsByRole retrieves all permissions associated with a role.
func (as *AgentDBModel) GetPermissionsByRole2(roleName string) ([]*Permission, error) {
	// Implement logic to retrieve permissions associated with the role.
	// You might need to join rolePermissions and permissions tables.

	var permissions []*Permission
	err := as.DB.Joins("JOIN rolePermissions ON permissions.id = rolePermissions.permission_id").
		Joins("JOIN roles ON roles.id = rolePermissions.role_id").
		Where("roles.name = ?", roleName).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// RevokeRoleFromAgent revokes a role from an agent.
func (as *AgentDBModel) RevokeRoleBaseFromAgent(agentID uint, roleName string) error {
	_, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the role from the agent.
	// Example: Delete the corresponding agent-role relationship record.
	err = as.DB.Where("agent_id = ? AND role_id = ?", agentID, role.ID).Delete(&AgentRole{}).Error
	if err != nil {
		return err
	}

	return nil
}

// RevokeRoleFromUser revokes a role from a user.
func (as *AgentDBModel) RevokeRoleBaseFromAgent2(agentID uint, roleName string) error {
	_, err := as.GetAgentByID(agentID)
	if err != nil {
		return err
	}

	role, err := as.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	// Implement logic to revoke the role from the user.
	// Example: Delete the corresponding user-role relationship record.
	err = as.DB.Where("user_id = ? AND role_id = ?", agentID, role.ID).Delete(&AgentRole{}).Error
	if err != nil {
		return err
	}

	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////
// UserAgentRepository implementations
func (as *AgentDBModel) CreateUserAgent(agentID uint, userID uint) error {
	// Implement logic to create a user-agent relationship in the database.
	// Example: Create a record in the userAgents table.
	userAgent := UserAgent{AgentID: agentID, UserID: userID}
	return as.DB.Create(&userAgent).Error
}

func (as *AgentDBModel) GetAgentsByUser(userID uint) ([]*Agents, error) {
	// Implement logic to retrieve agents associated with a user.
	// Example: Join userAgents and agents tables to get the agents by user.
	var agents []*Agents
	err := as.DB.Joins("JOIN userAgents ON agents.id = userAgents.agent_id").
		Where("userAgents.user_id = ?", userID).
		Find(&agents).Error
	return agents, err
}

func (as *AgentDBModel) GetUsersByAgent(agentID uint) ([]*Users, error) {
	// Implement logic to retrieve users associated with an agent.
	// Example: Join userAgents and users tables to get the users by agent.
	var users []*Users
	err := as.DB.Joins("JOIN userAgents ON users.id = userAgents.user_id").
		Where("userAgents.agent_id = ?", agentID).
		Find(&users).Error
	return users, err
}

func (as *AgentDBModel) RemoveUserAgentRelationship(agentID uint, userID uint) error {
	// Implement logic to remove a user-agent relationship from the database.
	// Example: Delete the record from the userAgents table.
	return as.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&UserAgent{}).Error
}

// UserPermissionRepository implementations
func (as *AgentDBModel) AssignPermissionToAgent(agentID uint, permissionID uint) error {
	// Implement logic to assign a permission to a user.
	// Example: Create a record in the userPermissions table.
	agentPermission := AgentPermission{AgentID: agentID, PermissionID: permissionID}
	return as.DB.Create(&agentPermission).Error
}

func (as *AgentDBModel) GetAgentsByTeam(teamID uint) ([]*Agents, error) {
	// Implement logic to retrieve agents associated with a team.
	// Example: Join teamAgents and agents tables to get agents by team.
	var agents []*Agents
	err := as.DB.Joins("JOIN teamAgents ON agents.id = teamAgents.agent_id").
		Where("teamAgents.team_id = ?", teamID).
		Find(&agents).Error
	return agents, err
}

func (as *AgentDBModel) GetTeamByName(name string) (*Teams, error) {
	// Implement logic to retrieve a team by its name.
	var team Teams
	err := as.DB.Where("name = ?", name).First(&team).Error
	return &team, err
}

// AgentTeamRepository implementations
func (as *AgentDBModel) GetAgentTeams(agentID uint) ([]*Teams, error) {
	// Implement logic to retrieve teams associated with an agent.
	// Example: Join teamAgents and teams tables to get teams by agent.
	var teams []*Teams
	err := as.DB.Joins("JOIN teamAgents ON teams.id = teamAgents.team_id").
		Where("teamAgents.agent_id = ?", agentID).
		Find(&teams).Error
	return teams, err
}

// Define more repository interfaces and their implementations as needed.

// Example of a custom query method to retrieve agents with specific criteria.
func (as *AgentDBModel) GetAgentsWithCriteria(criteria string) ([]*Agents, error) {
	// Implement a custom query to retrieve agents based on specific criteria.
	var agents []*Agents
	err := as.DB.Where(criteria).Find(&agents).Error
	return agents, err
}

// UserRepository implementations
func (as *AgentDBModel) AddAgentToUser(userID uint, agentID uint) error {
	// Implement logic to associate an agent with a user.
	userAgent := UserAgent{UserID: userID, AgentID: agentID}
	return as.DB.Create(&userAgent).Error
}

// RevokePermissionFromRole revokes a permission from a role.
func (as *AgentDBModel) DeleteTeamAgent(teamID uint, agentID uint) error {

	// Implement logic to revoke the permission from the role.
	// Example: Delete the corresponding role-permission relationship record.
	err := as.DB.Where("team_id = ? AND agent_id = ?", teamID, agentID).Delete(&TeamAgent{}).Error
	if err != nil {
		return err
	}

	return nil
}

// AssignMultipleAgentsToTeam facilitates the assignment of multiple agents to a specific team.
func (db *AgentDBModel) AssignMultipleAgentsToTeam(agentIDs []uint, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, agentID := range agentIDs {
			var teamAgent TeamAgent
			if err := tx.Where("agent_id = ? AND team_id = ?", agentID, teamID).First(&teamAgent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Agent is not yet assigned to the team, proceed with the assignment
					teamAgent = TeamAgent{AgentID: agentID, TeamID: teamID}
					if err := tx.Create(&teamAgent).Error; err != nil {
						return err // Rollback on error
					}
				} else {
					return err // Rollback on unexpected error
				}
			}
			// If the agent is already assigned, skip to the next agent without doing anything
		}
		return nil // Commit the transaction
	})
}

// UnassignAgentsFromTeam handles the removal of multiple agents from a team.
func (db *AgentDBModel) UnassignAgentsFromTeam(agentIDs []uint, teamID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Attempt to remove each specified agent from the team
		if err := tx.Where("agent_id IN ? AND team_id = ?", agentIDs, teamID).Delete(&TeamAgent{}).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit the transaction
	})
}

// UpdateTeamPermissions updates a team's permissions.
func (db *AgentDBModel) UpdateTeamPermissions(teamID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Find existing TeamPermission relationship
		var teamPermission TeamPermission
		err := tx.Where("team_id = ?", teamID).First(&teamPermission).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err // Rollback on unexpected error
		}

		// Update the permissions associated with the team
		if err := tx.Model(&teamPermission).Association("Permissions").Replace(permissionIDs); err != nil {
			return err // Rollback on error
		}

		return nil // Commit the transaction
	})
}

// PerformComplexOperationWithTransaction demonstrates a complex operation involving multiple steps.
func (db *AgentDBModel) PerformComplexOperationWithTransaction(agentID uint, teamID uint, permissionIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Example step 1: Assign an agent to a team
		if err := db.AssignAgentToTeam(agentID, teamID); err != nil {
			return err // Rollback on error
		}

		// Example step 2: Update agent's permissions
		if err := db.UpdateAgentPermissions(agentID, permissionIDs); err != nil {
			return err // Rollback on error
		}

		// Additional steps can be added here, each protected by the transaction
		// ...

		return nil // Commit the transaction if all steps are successful
	})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////

// UpdateRole updates an existing role's details.
func (db *AgentDBModel) UpdateRole(role *Role) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			db.log.Error("Failed to update role: %v", err)
			return err // Error updating role, rollback
		}
		return nil // Success, commit the transaction
	})
}

// GetUnitByNumber retrieves a unit by its unit number.
func (db *AgentDBModel) GetUnitByNumber(unitNumber int) (*Unit, error) {
	var unit Unit
	err := db.DB.Where("unit_number = ?", unitNumber).First(&unit).Error
	if err != nil {
		db.log.Error("Failed to retrieve unit by number: %v", err)
		return nil, err
	}
	return &unit, nil
}

// AssignRolesToAgentByRoleID assigns roles to an agent by role IDs.
func (db *AgentDBModel) AssignRolesToAgentByRoleID(agentID uint, roleIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var agent Agents
		if err := tx.First(&agent, agentID).Error; err != nil {
			db.log.Error("Failed to find agent during role assignment by ID: %v", err)
			return err
		}
		if err := tx.Model(&agent).Association("Roles").Replace(roleIDs); err != nil {
			db.log.Error("Failed to assign roles to agent by ID: %v", err)
			return err
		}
		return nil
	})
}

// RevokeRolesFromAgentByRoleIDs revokes roles from an agent by role IDs.
func (db *AgentDBModel) RevokeRolesFromAgentByRoleIDs(agentID uint, roleIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var agent Agents
		if err := tx.First(&agent, agentID).Error; err != nil {
			db.log.Error("Failed to find agent during role revocation by ID: %v", err)
			return err
		}
		if err := tx.Model(&agent).Association("Roles").Delete(roleIDs); err != nil {
			db.log.Error("Failed to revoke roles from agent by ID: %v", err)
			return err
		}
		return nil
	})
}

// RevokeRoleFromAgentByRoleName revokes a role from an agent by role name.
func (db *AgentDBModel) RevokeRoleFromAgentByRoleName(agentID uint, roleName string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			db.log.Error("Failed to find role by name during revocation: %v", err)
			return err
		}
		if err := tx.Model(&Agents{ID: agentID}).Association("Roles").Delete(&role); err != nil {
			db.log.Error("Failed to revoke role from agent by name: %v", err)
			return err
		}
		return nil
	})
}

// AssignPermissionsToAgentByPermissionNames assigns permissions to an agent by permission names.
func (db *AgentDBModel) AssignPermissionsToAgentByPermissionNames(agentID uint, permissionNames []string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var permissions []Permission
		if err := tx.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
			db.log.Error("Failed to find permissions by name during assignment to agent: %v", err)
			return err
		}
		for _, perm := range permissions {
			ap := AgentPermission{
				AgentID:      agentID,
				PermissionID: perm.ID,
			}
			if err := tx.FirstOrCreate(&ap, AgentPermission{AgentID: agentID, PermissionID: perm.ID}).Error; err != nil {
				db.log.Error("Failed to assign permission to agent by name within transaction: %v", err)
				return err
			}
		}
		return nil
	})
}

// RevokePermissionFromAgentByPermissionName revokes a permission from an agent by permission name.
func (db *AgentDBModel) RevokePermissionFromAgentByPermissionName(agentID uint, permissionName string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var permission Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			db.log.Error("Failed to find permission by name during revocation from agent: %v", err)
			return err
		}
		if err := tx.Where("agent_id = ? AND permission_id = ?", agentID, permission.ID).Delete(&AgentPermission{}).Error; err != nil {
			db.log.Error("Failed to revoke permission from agent by name: %v", err)
			return err
		}
		return nil
	})
}

// AssignPermissionsToTeamByPermissionNames assigns permissions to a team by permission names.
func (db *AgentDBModel) AssignPermissionsToTeamByPermissionNames(teamID uint, permissionNames []string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var permissions []*Permission
		if err := tx.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
			db.log.Error("Failed to find permissions by name during assignment to team: %v", err)
			return err
		}
		for _, perm := range permissions {
			perm := append(permissions, perm)
			tp := TeamPermission{
				TeamID:      teamID,
				Permissions: perm,
			}
			if err := tx.FirstOrCreate(&tp, TeamPermission{TeamID: teamID, Permissions: perm}).Error; err != nil {
				db.log.Error("Failed to assign permission to team by name within transaction: %v", err)
				return err
			}
		}
		return nil
	})
}

// AssignPermissionsToTeamByPermissionName assigns a single permission to a team by the permission's name, ensuring transactional integrity and comprehensive logging.
func (db *AgentDBModel) AssignPermissionsToTeamByPermissionName(teamID uint, permissionName string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Find the permission by name
		var permission Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			db.log.Error(fmt.Sprintf("Failed to find permission by name '%s' during assignment to team: %v", permissionName, err))
			return err // Error finding permission, rollback
		}

		// Find the team and preload its permissions
		var team TeamPermission
		if err := tx.Preload("Permissions").Where("team_id = ?", teamID).First(&team).Error; err != nil {
			db.log.Error(fmt.Sprintf("Failed to find team by ID '%d' during permission assignment: %v", teamID, err))
			return err // Error finding team, rollback
		}

		// Check if the permission is already assigned to avoid duplicates
		alreadyAssigned := false
		for _, p := range team.Permissions {
			if p.ID == permission.ID {
				alreadyAssigned = true
				break
			}
		}

		if !alreadyAssigned {
			// Append the permission to the team's permissions and save
			team.Permissions = append(team.Permissions, &permission)
			if err := tx.Save(&team).Error; err != nil {
				db.log.Error(fmt.Sprintf("Failed to assign permission '%s' to team within transaction: %v", permissionName, err))
				return err // Error assigning permission to team, rollback
			}
		} else {
			db.log.Info(fmt.Sprintf("Permission '%s' already assigned to team; skipping reassignment.", permissionName))
		}

		return nil // Success, commit the transaction
	})
}

// RevokePermissionFromTeamByName revokes a permission from a team by permission name.
func (db *AgentDBModel) RevokePermissionFromTeamByPermissionName(teamID uint, permissionName string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var permission Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			db.log.Error("Failed to find permission by name during revocation from team: %v", err)
			return err
		}
		if err := tx.Where("team_id = ? AND permission_id = ?", teamID, permission.ID).Delete(&TeamPermission{}).Error; err != nil {
			db.log.Error("Failed to revoke permission from team by name: %v", err)
			return err
		}
		return nil
	})
}

// AddPermissionToRoleByPermissionName adds a permission to a role by permission name.
func (db *AgentDBModel) AddPermissionToRoleByPermissionName(roleName string, permissionName string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			db.log.Error("Failed to find role by name during permission addition: %v", err)
			return err
		}
		var permission Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			db.log.Error("Failed to find permission by name during addition to role: %v", err)
			return err
		}
		rp := RolePermission{
			RoleID:       role.ID,
			PermissionID: permission.ID,
		}
		if err := tx.Create(&rp).Error; err != nil {
			db.log.Error("Failed to add permission to role by name: %v", err)
			return err
		}
		return nil
	})
}

// RemovePermissionFromRoleByPermissionName removes a permission from a role by permission name.
func (db *AgentDBModel) RemovePermissionFromRoleByPermissionName(roleName string, permissionName string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			db.log.Error("Failed to find role by name during permission removal: %v", err)
			return err
		}
		var permission Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			db.log.Error("Failed to find permission by name during removal from role: %v", err)
			return err
		}
		if err := tx.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).Delete(&RolePermission{}).Error; err != nil {
			db.log.Error("Failed to remove permission from role by name: %v", err)
			return err
		}
		return nil
	})
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////////////

// LogAction records an action performed by an agent for auditing purposes.
func (db *AgentDBModel) LogAgentAction(agentID uint, actionType string, details string) error {
	actionLog := AgentEvent{
		AgentID:    agentID,
		ActionType: &actionType,
		Details:    &details,
		Timestamp:  time.Now(),
	}

	return db.DB.Create(&actionLog).Error
}

// GetAgentActions retrieves the recent actions performed by a specific agent.
func (db *AgentDBModel) GetAgentActions(agentID uint, limit int) ([]AgentEvent, error) {
	var agentActions []AgentEvent
	if err := db.DB.Where("agent_id = ?", agentID).Order("timestamp desc").Limit(limit).Find(&agentActions).Error; err != nil {
		return nil, err
	}
	return agentActions, nil
}

// GetAllActions retrieves recent actions performed by all agents.
func (db *AgentDBModel) GetAllActions(limit int) ([]AgentEvent, error) {
	var allActions []AgentEvent
	if err := db.DB.Order("timestamp desc").Limit(limit).Find(&allActions).Error; err != nil {
		return nil, err
	}
	return allActions, nil
}

// SearchAgents searches for agents based on the provided criteria with pagination and sorting options.
func (db *AgentDBModel) SearchAgents(criteria SearchCriteria, page, pageSize int, sortBy string, sortOrder string) ([]Agents, error) {
	var agents []Agents
	query := db.DB

	// Apply search criteria
	if criteria.Name != "" {
		query = query.Where("name LIKE ?", "%"+criteria.Name+"%")
	}
	if criteria.Role != "" {
		query = query.Where("role = ?", criteria.Role)
	}
	if criteria.Department != "" {
		query = query.Where("department = ?", criteria.Department)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Apply sorting
	if sortBy != "" {
		orderBy := sortBy + " " + sortOrder
		query = query.Order(orderBy)
	}

	// Execute the query
	if err := query.Find(&agents).Error; err != nil {
		return nil, err
	}

	return agents, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////

// AddAgentSchedule adds a new schedule for an agent, ensuring data integrity with transactional support.
func (db *AgentDBModel) AddAgentSchedule(schedule *AgentSchedule) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(schedule).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentSchedule updates an existing schedule for an agent, handling changes within a transaction.
func (db *AgentDBModel) UpdateAgentSchedule(scheduleID uint, newSchedule *AgentSchedule) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var schedule AgentSchedule
		if err := tx.Where("id = ?", scheduleID).First(&schedule).Error; err != nil {
			return err // Schedule not found, rollback
		}
		schedule.StartDate = newSchedule.StartDate
		schedule.EndDate = newSchedule.EndDate
		schedule.ShiftType = newSchedule.ShiftType
		schedule.IsActive = newSchedule.IsActive
		if err := tx.Save(&schedule).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// DeleteAgentSchedule removes a specific schedule for an agent, ensuring the operation is transactionally safe.
func (db *AgentDBModel) DeleteAgentSchedule(scheduleID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentSchedule{}, scheduleID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentSchedules retrieves all schedules for a specific agent, ensuring accurate and up-to-date information.
func (db *AgentDBModel) ListAgentSchedules(agentID uint) ([]AgentSchedule, error) {
	var schedules []AgentSchedule
	if err := db.DB.Where("agent_id = ?", agentID).Find(&schedules).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return schedules, nil
}

// AddAgentSkill assigns a new skill to an agent, wrapped in a transaction for data consistency.
func (db *AgentDBModel) AddAgentSkill(skill *AgentSkill) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(skill).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentSkill updates details of an existing skill for an agent within a transaction for safety.
func (db *AgentDBModel) UpdateAgentSkill(skillID uint, newSkill *AgentSkill) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var skill AgentSkill
		if err := tx.Where("id = ?", skillID).First(&skill).Error; err != nil {
			return err // Skill not found, rollback
		}
		skill.SkillName = newSkill.SkillName
		skill.Level = newSkill.Level
		if err := tx.Save(&skill).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// DeleteAgentSkill removes a skill from an agent's profile, ensuring the operation's atomicity.
func (db *AgentDBModel) DeleteAgentSkill(skillID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentSkill{}, skillID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentSkills fetches all skills associated with a specific agent, providing a comprehensive view.
func (db *AgentDBModel) ListAgentSkills(agentID uint) ([]AgentSkill, error) {
	var skills []AgentSkill
	if err := db.DB.Where("agent_id = ?", agentID).Find(&skills).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return skills, nil
}

// SubmitAgentFeedback records new feedback for an agent, wrapped in a transaction for consistency and reliability.
func (db *AgentDBModel) SubmitAgentFeedback(feedback *AgentFeedback) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(feedback).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// ListAgentFeedback retrieves all feedback entries for a specific agent, providing insights into performance and areas for improvement.
func (db *AgentDBModel) ListAgentFeedback(agentID uint) ([]AgentFeedback, error) {
	var feedbacks []AgentFeedback
	if err := db.DB.Where("agent_id = ?", agentID).Find(&feedbacks).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return feedbacks, nil
}

// UpdateAgentFeedback allows modifications to an existing feedback record, ensuring data integrity through a transaction.
func (db *AgentDBModel) UpdateAgentFeedback(feedbackID uint, newFeedback *AgentFeedback) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var feedback AgentFeedback
		if err := tx.Where("id = ?", feedbackID).First(&feedback).Error; err != nil {
			return err // Feedback not found, rollback
		}
		feedback.Score = newFeedback.Score
		feedback.Comments = newFeedback.Comments
		if err := tx.Save(&feedback).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// DeleteAgentFeedback removes a feedback record from an agent's profile, executed within a transaction for safety.
func (db *AgentDBModel) DeleteAgentFeedback(feedbackID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentFeedback{}, feedbackID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// AddAgentCertification associates a new certification with an agent.
func (db *AgentDBModel) AddAgentCertification(certification *AgentCertification) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(certification).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentCertification updates details of an existing certification for an agent, ensuring atomicity with a transaction.
func (db *AgentDBModel) UpdateAgentCertification(certificationID uint, newCertification *AgentCertification) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var certification AgentCertification
		if err := tx.Where("id = ?", certificationID).First(&certification).Error; err != nil {
			return err // Certification not found, rollback
		}
		certification.Certification = newCertification.Certification
		certification.IssuedBy = newCertification.IssuedBy
		certification.IssuedDate = newCertification.IssuedDate
		certification.ExpiryDate = newCertification.ExpiryDate
		if err := tx.Save(&certification).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// DeleteAgentCertification removes a certification record from an agent's profile.
func (db *AgentDBModel) DeleteAgentCertification(certificationID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentCertification{}, certificationID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentCertifications retrieves all certifications associated with a specific agent.
func (db *AgentDBModel) ListAgentCertifications(agentID uint) ([]AgentCertification, error) {
	var certifications []AgentCertification
	if err := db.DB.Where("agent_id = ?", agentID).Find(&certifications).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return certifications, nil
}

// RequestAgentLeave submits a new leave request for an agent, ensuring transactional integrity for the operation.
func (db *AgentDBModel) RequestAgentLeave(leaveRequest *AgentLeaveRequest) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(leaveRequest).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentLeaveRequest allows modifications to an existing leave request, maintaining data consistency with a transaction.
func (db *AgentDBModel) UpdateAgentLeaveRequest(leaveRequestID uint, newLeaveRequest *AgentLeaveRequest) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var leaveRequest AgentLeaveRequest
		if err := tx.Where("id = ?", leaveRequestID).First(&leaveRequest).Error; err != nil {
			return err // Leave request not found, rollback
		}
		leaveRequest.Status = newLeaveRequest.Status
		if err := tx.Save(&leaveRequest).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// CancelAgentLeaveRequest handles the cancellation of an agent's leave request, executed within a transaction for safety.
func (db *AgentDBModel) CancelAgentLeaveRequest(leaveRequestID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentLeaveRequest{}, leaveRequestID).Error; err != nil {
			return err // Error during cancellation, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentLeaveRequests retrieves all leave requests submitted by a specific agent.
func (db *AgentDBModel) ListAgentLeaveRequests(agentID uint) ([]AgentLeaveRequest, error) {
	var requests []AgentLeaveRequest
	if err := db.DB.Where("agent_id = ?", agentID).Find(&requests).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return requests, nil
}

// UpdateAgentAvailability updates an agent's availability status with a transaction.
func (db *AgentDBModel) UpdateAgentAvailability(availability *AgentAvailability) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(availability).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// ApproveAgentLeaveRequest updates the status of an agent's leave request to "Approved".
func (db *AgentDBModel) ApproveAgentLeaveRequest(leaveRequestID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var leaveRequest AgentLeaveRequest
		if err := tx.Where("id = ?", leaveRequestID).First(&leaveRequest).Error; err != nil {
			return err // If leave request not found, rollback
		}
		leaveRequest.Status = "Approved"
		if err := tx.Save(&leaveRequest).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// DenyAgentLeaveRequest updates the status of an agent's leave request to "Denied".
func (db *AgentDBModel) DenyAgentLeaveRequest(leaveRequestID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var leaveRequest AgentLeaveRequest
		if err := tx.Where("id = ?", leaveRequestID).First(&leaveRequest).Error; err != nil {
			return err // If leave request not found, rollback
		}
		leaveRequest.Status = "Denied"
		if err := tx.Save(&leaveRequest).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentTrainingRecords retrieves all training records for a specific agent.
func (db *AgentDBModel) ListAgentTrainingRecords(agentID uint) ([]AgentTrainingRecord, error) {
	var records []AgentTrainingRecord
	if err := db.DB.Where("agent_id = ?", agentID).Find(&records).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return records, nil
}

// DeleteAgentTrainingRecord removes a training record from the database.
func (db *AgentDBModel) DeleteAgentTrainingRecord(recordID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentTrainingRecord{}, recordID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// SubmitAgentLeaveRequest allows an agent to submit a request for leave.
func (db *AgentDBModel) SubmitAgentLeaveRequest(leaveRequest *AgentLeaveRequest) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(leaveRequest).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentLeaveRequestStatus updates the status of an existing leave request.
func (db *AgentDBModel) UpdateAgentLeaveRequestStatus(requestID uint, newStatus string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var request AgentLeaveRequest
		if err := tx.Where("id = ?", requestID).First(&request).Error; err != nil {
			return err // If request not found, rollback
		}
		request.Status = newStatus
		if err := tx.Save(&request).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// UpdateAgentCertificationDetails updates details of an existing certification for an agent.
func (db *AgentDBModel) UpdateAgentCertificationDetails(certificationID uint, newDetails *AgentCertification) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var certification AgentCertification
		if err := tx.Where("id = ?", certificationID).First(&certification).Error; err != nil {
			return err // If certification not found, rollback
		}
		certification.Certification = newDetails.Certification
		certification.IssuedBy = newDetails.IssuedBy
		certification.IssuedDate = newDetails.IssuedDate
		certification.ExpiryDate = newDetails.ExpiryDate
		if err := tx.Save(&certification).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentActivities retrieves a list of activities for a specific agent.
func (db *AgentDBModel) ListAgentActivities(agentID uint) ([]AgentActivityLog, error) {
	var activities []AgentActivityLog
	if err := db.DB.Where("agent_id = ?", agentID).Find(&activities).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return activities, nil
}

// CreateAgentTrainingSession adds a new training session for agents, ensuring data consistency with a transaction.
func (db *AgentDBModel) CreateAgentTrainingSession(trainingSession *AgentTrainingSession) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trainingSession).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentTrainingSession updates details of an existing training session.
func (db *AgentDBModel) UpdateAgentTrainingSession(trainingSessionID uint, newTrainingSession *AgentTrainingSession) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var trainingSession AgentTrainingSession
		if err := tx.Where("id = ?", trainingSessionID).First(&trainingSession).Error; err != nil {
			return err // If training session not found, rollback
		}
		trainingSession.Title = newTrainingSession.Title
		trainingSession.Description = newTrainingSession.Description
		trainingSession.StartDate = newTrainingSession.StartDate
		trainingSession.EndDate = newTrainingSession.EndDate
		trainingSession.Location = newTrainingSession.Location
		trainingSession.TrainerID = newTrainingSession.TrainerID
		if err := tx.Save(&trainingSession).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// DeleteAgentTrainingSession removes a training session from the database.
func (db *AgentDBModel) DeleteAgentTrainingSession(trainingSessionID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentTrainingSession{}, trainingSessionID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentTrainingSessions retrieves all training sessions an agent is enrolled in.
func (db *AgentDBModel) ListAgentTrainingSessions(agentID uint) ([]AgentTrainingSession, error) {
	var trainingSessions []AgentTrainingSession
	if err := db.DB.Where("attendees @> ?", []uint{agentID}).Find(&trainingSessions).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return trainingSessions, nil
}

// EnrollAgentInTrainingSession adds an agent to a training session's attendee list.
func (db *AgentDBModel) EnrollAgentInTrainingSession(agentID, sessionID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Retrieve training session
		var session AgentTrainingSession
		if err := tx.First(&session, sessionID).Error; err != nil {
			return err // If session not found, rollback
		}
		// Append agent to session's attendees list if not already included
		if !contains(session.Attendees, agentID) {
			session.Attendees = append(session.Attendees, Agents{ID: agentID})
			if err := tx.Save(&session).Error; err != nil {
				return err // Error adding attendee, rollback
			}
		}
		return nil // Commit if no error
	})
}

// RecordAgentPerformanceReview records a new performance review for an agent.
func (db *AgentDBModel) RecordAgentPerformanceReview(review *AgentPerformanceReview) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentPerformanceReview updates an existing performance review for an agent.
func (db *AgentDBModel) UpdateAgentPerformanceReview(reviewID uint, newReview *AgentPerformanceReview) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var review AgentPerformanceReview
		if err := tx.Where("id = ?", reviewID).First(&review).Error; err != nil {
			return err // If review not found, rollback
		}
		review.Score = newReview.Score
		review.Feedback = newReview.Feedback
		review.ReviewDate = newReview.ReviewDate
		if err := tx.Save(&review).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// ListAgentPerformanceReviews retrieves all performance reviews for a specific agent.
func (db *AgentDBModel) ListAgentPerformanceReviews(agentID uint) ([]AgentPerformanceReview, error) {
	var reviews []AgentPerformanceReview
	if err := db.DB.Where("agent_id = ?", agentID).Find(&reviews).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return reviews, nil
}

// DeleteAgentPerformanceReview removes a performance review record.
func (db *AgentDBModel) DeleteAgentPerformanceReview(reviewID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AgentPerformanceReview{}, reviewID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// RecordCustomerInteraction logs an interaction between an agent and a customer.
func (db *AgentDBModel) RecordCustomerInteraction(interaction *CustomerInteraction) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(interaction).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateCustomerInteraction updates details of an existing customer interaction.
func (db *AgentDBModel) UpdateCustomerInteraction(interactionID uint, newInteraction *CustomerInteraction) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var interaction CustomerInteraction
		if err := tx.Where("id = ?", interactionID).First(&interaction).Error; err != nil {
			return err // If interaction not found, rollback
		}
		interaction.Channel = newInteraction.Channel
		interaction.Content = newInteraction.Content
		interaction.InteractionTime = newInteraction.InteractionTime
		if err := tx.Save(&interaction).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// ListCustomerInteractionsByAgent retrieves all customer interactions associated with a specific agent.
func (db *AgentDBModel) ListCustomerInteractionsByAgent(agentID uint) ([]CustomerInteraction, error) {
	var interactions []CustomerInteraction
	if err := db.DB.Where("agent_id = ?", agentID).Order("interaction_time desc").Find(&interactions).Error; err != nil {
		return nil, err // Handle potential retrieval errors
	}
	return interactions, nil
}

// DeleteCustomerInteraction removes a customer interaction record from the database.
func (db *AgentDBModel) DeleteCustomerInteraction(interactionID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&CustomerInteraction{}, interactionID).Error; err != nil {
			return err // Error during deletion, rollback
		}
		return nil // Commit if no error
	})
}

// AssignAgentToTrainingModule assigns an agent to a specific training module and records the assignment.
func (db *AgentDBModel) AssignAgentToTrainingModule(agentID, moduleID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var trainingRecord AgentTrainingRecord
		trainingRecord.AgentID = agentID
		trainingRecord.TrainingModuleID = moduleID
		trainingRecord.CompletedAt = time.Now() // Assuming immediate completion for simplicity
		if err := tx.Create(&trainingRecord).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentTrainingRecord updates a training record for an agent, for example, to record a score or feedback.
func (db *AgentDBModel) UpdateAgentTrainingRecord(recordID uint, newRecord *AgentTrainingRecord) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var record AgentTrainingRecord
		if err := tx.Where("id = ?", recordID).First(&record).Error; err != nil {
			return err // If record not found, rollback
		}
		record.Score = newRecord.Score
		record.Feedback = newRecord.Feedback
		record.CompletedAt = newRecord.CompletedAt
		if err := tx.Save(&record).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// contains checks if the slice of Agents contains an agent with the given ID.
func contains(agents []Agents, agentID uint) bool {
	for _, agent := range agents {
		if agent.ID == agentID {
			return true
		}
	}
	return false
}

// RecordAgentLoginActivity logs an agent's login activity to the system.
func (db *AgentDBModel) RecordAgentLoginActivity(activity *AgentLoginActivity) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(activity).Error; err != nil {
			return err // Rollback on error
		}
		return nil // Commit if no error
	})
}

// UpdateAgentProfilePic updates the profile picture URL of an agent.
func (db *AgentDBModel) UpdateAgentProfilePic(agentID uint, profilePicURL string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var agent Agents
		if err := tx.Where("id = ?", agentID).First(&agent).Error; err != nil {
			return err // If agent not found, rollback
		}
		agent.ProfilePic = &profilePicURL
		if err := tx.Save(&agent).Error; err != nil {
			return err // Error updating, rollback
		}
		return nil // Commit if no error
	})
}

// AddAgentSkillSet assigns new skills or updates existing skills for an agent.
func (db *AgentDBModel) AddOrUpdateAgentSkillSet(agentID uint, skills []AgentSkillSet) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, skill := range skills {
			var existingSkill AgentSkillSet
			// Check if the skill already exists for the agent
			result := tx.Where("agent_id = ? AND skill = ?", agentID, skill.Skill).First(&existingSkill)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Skill does not exist, create new
				skill.AgentID = agentID
				if err := tx.Create(&skill).Error; err != nil {
					return err // Rollback on error
				}
			} else {
				// Skill exists, update level
				existingSkill.Level = skill.Level
				if err := tx.Save(&existingSkill).Error; err != nil {
					return err // Rollback on error
				}
			}
		}
		return nil // Commit if no error
	})
}

// //////////////////////////////////////////////////////////////////////////////////
func (service *AgentDBModel) LogCustomerInteraction(interaction CustomerInteraction) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&interaction).Error; err != nil {
			// Transaction will automatically roll back on return with error
			return err
		}
		// If no error, transaction will be committed
		return nil
	})
}

func (repo *AgentDBModel) AssignPermissionsToAgent2(agentID uint, permissionNames []string) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		var permissions []Permission
		if err := tx.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
			return err
		}

		for _, perm := range permissions {
			ap := AgentPermission{
				AgentID:      agentID,
				PermissionID: perm.ID,
			}
			// Avoid duplicating permissions
			if err := tx.FirstOrCreate(&ap, AgentPermission{AgentID: agentID, PermissionID: perm.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *AgentDBModel) AddAgentToTeam(agentID, teamID uint) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		teamAgent := TeamAgent{
			TeamID:  teamID,
			AgentID: agentID,
		}
		return tx.FirstOrCreate(&teamAgent, TeamAgent{TeamID: teamID, AgentID: agentID}).Error
	})
}

// LogAction records an action performed by an agent for auditing purposes with improved logging.
func (db *AgentDBModel) LogAction(agentID uint, actionType, details string) error {
	actionLog := AgentActivityLog{
		AgentID:   agentID,
		Activity:  actionType,
		Timestamp: time.Now(),
	}
	if err := db.DB.Create(&actionLog).Error; err != nil {
		db.log.Error("Error logging agent action: %v", err)
		return err
	}
	return nil
}
