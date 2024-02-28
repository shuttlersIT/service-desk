// backend/models/incident_db_model.go

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ContentArticle struct {
	gorm.Model
	AuthorID    uint       `gorm:"index" json:"author_id"`                  // ID of the author/user
	Title       string     `gorm:"type:varchar(255);not null" json:"title"` // Article title
	Body        string     `gorm:"type:text;not null" json:"body"`          // Article body/content
	CategoryID  uint       `gorm:"index" json:"category_id"`                // Category of the article
	PublishedAt *time.Time `json:"published_at,omitempty"`                  // Optional publication date
}

func (ContentArticle) TableName() string {
	return "content_articles"
}

// Implement driver.Valuer for ContentArticle
func (ca ContentArticle) Value() (driver.Value, error) {
	return json.Marshal(ca)
}

// Implement driver.Scanner for ContentArticle
func (ca *ContentArticle) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ContentArticle scan")
	}

	return json.Unmarshal(data, &ca)
}

type ContentCategory struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);unique;not null" json:"name"` // Category name
	Description string `gorm:"type:text" json:"description,omitempty"`        // Category description
}

func (ContentCategory) TableName() string {
	return "content_categories"
}

// Implement driver.Valuer for ContentCategory
func (cc ContentCategory) Value() (driver.Value, error) {
	return json.Marshal(cc)
}

// Implement driver.Scanner for ContentCategory
func (cc *ContentCategory) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ContentCategory scan")
	}

	return json.Unmarshal(data, &cc)
}

type UserInteraction struct {
	gorm.Model
	UserID          uint   `gorm:"index;not null" json:"user_id"`                      // ID of the user
	InteractionType string `gorm:"type:varchar(100);not null" json:"interaction_type"` // Type of interaction (e.g., view, like)
	EntityID        uint   `gorm:"index;not null" json:"entity_id"`                    // ID of the entity interacted with
	EntityType      string `gorm:"type:varchar(100);not null" json:"entity_type"`      // Type of entity (e.g., article, product)
}

func (UserInteraction) TableName() string {
	return "user_interactions"
}

// Implement driver.Valuer for UserInteraction
func (ui UserInteraction) Value() (driver.Value, error) {
	return json.Marshal(ui)
}

// Implement driver.Scanner for UserInteraction
func (ui *UserInteraction) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for UserInteraction scan")
	}

	return json.Unmarshal(data, &ui)
}

type ContentMedia struct {
	gorm.Model
	ArticleID   uint   `gorm:"index" json:"article_id"`                      // ID of the associated article
	MediaType   string `gorm:"type:varchar(100);not null" json:"media_type"` // Type of media (e.g., image, video)
	URL         string `gorm:"type:text;not null" json:"url"`                // URL to access the media
	Description string `gorm:"type:text" json:"description,omitempty"`       // Media description
}

func (ContentMedia) TableName() string {
	return "content_media"
}

// Implement driver.Valuer for ContentMedia
func (cm ContentMedia) Value() (driver.Value, error) {
	return json.Marshal(cm)
}

// Implement driver.Scanner for ContentMedia
func (cm *ContentMedia) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ContentMedia scan")
	}

	return json.Unmarshal(data, &cm)
}

type ArticleTag struct {
	gorm.Model
	ArticleID uint `gorm:"index;not null" json:"article_id"` // ID of the article
	TagID     uint `gorm:"index;not null" json:"tag_id"`     // ID of the tag
}

func (ArticleTag) TableName() string {
	return "article_tags"
}

// Implement driver.Valuer for ArticleTag
func (at ArticleTag) Value() (driver.Value, error) {
	return json.Marshal(at)
}

// Implement driver.Scanner for ArticleTag
func (at *ArticleTag) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ArticleTag scan")
	}

	return json.Unmarshal(data, &at)
}

type UserJourneyEvent struct {
	ID        uint   `gorm:"primaryKey"`
	SessionID string `gorm:"type:varchar(255);not null" json:"session_id"` // Session identifier for correlating events
	UserID    uint   `gorm:"index;not null" json:"user_id"`                // User associated with the event, if authenticated
	EventName string `gorm:"type:varchar(255);not null" json:"event_name"` // Name of the event
	Details   string `gorm:"type:text" json:"details,omitempty"`           // JSON encoded details about the event
	CreatedAt time.Time
}

func (UserJourneyEvent) TableName() string {
	return "user_journey_events"
}

// Implement driver.Valuer for UserJourneyEvent
func (uje UserJourneyEvent) Value() (driver.Value, error) {
	return json.Marshal(uje)
}

// Implement driver.Scanner for UserJourneyEvent
func (uje *UserJourneyEvent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for UserJourneyEvent scan")
	}

	return json.Unmarshal(data, &uje)
}

type LocalizedContent struct {
	ID             uint   `gorm:"primaryKey"`
	ContentID      uint   `gorm:"index;not null" json:"content_id"`          // ID of the original content
	Language       string `gorm:"type:varchar(10);not null" json:"language"` // ISO language code for the translation
	TranslatedText string `gorm:"type:text;not null" json:"translated_text"` // Translated or localized text
	CreatedAt      time.Time
}

func (LocalizedContent) TableName() string {
	return "localized_contents"
}

// Implement driver.Valuer for LocalizedContent
func (lc LocalizedContent) Value() (driver.Value, error) {
	return json.Marshal(lc)
}

// Implement driver.Scanner for LocalizedContent
func (lc *LocalizedContent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for LocalizedContent scan")
	}

	return json.Unmarshal(data, &lc)
}

type Tenant struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);not null" json:"name"`           // Name of the tenant
	APIKey    string `gorm:"type:varchar(255);unique;not null" json:"api_key"` // Unique API key for the tenant
	IsActive  bool   `gorm:"default:true" json:"is_active"`                    // Flag to activate or deactivate the tenant
	CreatedAt time.Time
}

func (Tenant) TableName() string {
	return "tenants"
}

// Implement driver.Valuer for Tenant
func (t Tenant) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Implement driver.Scanner for Tenant
func (t *Tenant) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Tenant scan")
	}

	return json.Unmarshal(data, &t)
}

type TenantScopedModel struct {
	TenantID uint `gorm:"index;not null" json:"tenant_id"` // Foreign key linking to the Tenant
	// Include fields common to tenant-scoped models here
}

func (TenantScopedModel) TableName() string {
	return "tenant_scoped_model"
}

// Implement driver.Valuer for TenantScopedModel
func (tsm TenantScopedModel) Value() (driver.Value, error) {
	return json.Marshal(tsm)
}

// Implement driver.Scanner for TenantScopedModel
func (tsm *TenantScopedModel) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for TenantScopedModel scan")
	}

	return json.Unmarshal(data, &tsm)
}

type ContentCollection struct {
	ID          uint          `gorm:"primaryKey"`
	Title       string        `gorm:"type:varchar(255);not null" json:"title"`        // Title of the collection
	Description string        `gorm:"type:text" json:"description,omitempty"`         // Description of what the collection represents
	CuratorID   uint          `gorm:"index;not null" json:"curator_id"`               // ID of the user or system curating the collection
	Contents    []ContentItem `gorm:"many2many:collection_contents;" json:"contents"` // Items part of the collection
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ContentCollection) TableName() string {
	return "content_collections"
}

// Implement driver.Valuer for ContentCollection
func (cc ContentCollection) Value() (driver.Value, error) {
	return json.Marshal(cc)
}

// Implement driver.Scanner for ContentCollection
func (cc *ContentCollection) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ContentCollection scan")
	}

	return json.Unmarshal(data, &cc)
}

type ContentItem struct {
	ID          uint   `gorm:"primaryKey"`
	Type        string `gorm:"type:varchar(100);not null" json:"type"` // Type of content (e.g., article, video, image)
	ReferenceID uint   `gorm:"not null" json:"reference_id"`           // ID of the actual content, interpretation depends on Type
	Metadata    string `gorm:"type:text" json:"metadata,omitempty"`    // JSON for additional metadata about the item
}

func (ContentItem) TableName() string {
	return "content_items"
}

// Implement driver.Valuer for ContentItem
func (ci ContentItem) Value() (driver.Value, error) {
	return json.Marshal(ci)
}

// Implement driver.Scanner for ContentItem
func (ci *ContentItem) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ContentItem scan")
	}

	return json.Unmarshal(data, &ci)
}

type BehaviorEvent struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`                // ID of the user performing the event
	EventType  string    `gorm:"type:varchar(255);not null" json:"event_type"` // Type of event (e.g., click, view)
	EventData  string    `gorm:"type:text" json:"event_data,omitempty"`        // JSON encoded data providing context about the event
	OccurredAt time.Time `json:"occurred_at"`                                  // Timestamp when the event occurred
}

func (BehaviorEvent) TableName() string {
	return "behavior_events"
}

// Implement driver.Valuer for BehaviorEvent
func (be BehaviorEvent) Value() (driver.Value, error) {
	return json.Marshal(be)
}

// Implement driver.Scanner for BehaviorEvent
func (be *BehaviorEvent) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for BehaviorEvent scan")
	}

	return json.Unmarshal(data, &be)
}

type SystemMetric struct {
	ID                uint      `gorm:"primaryKey"`
	MetricType        string    `gorm:"type:varchar(100);not null" json:"metric_type"` // Type of metric (e.g., CPU, memory usage)
	SystemMetricValue float64   `json:"value"`                                         // Metric value
	RecordedAt        time.Time `json:"recorded_at"`                                   // Timestamp when the metric was recorded
}

func (SystemMetric) TableName() string {
	return "system_metrics"
}

// Implement driver.Valuer for SystemMetric
func (sm SystemMetric) Value() (driver.Value, error) {
	return json.Marshal(sm)
}

// Implement driver.Scanner for SystemMetric
func (sm *SystemMetric) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for SystemMetric scan")
	}

	return json.Unmarshal(data, &sm)
}

type UserPreferenceVector struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"`   // ID of the user the preferences belong to
	Preferences string    `gorm:"type:text;not null" json:"preferences"` // JSON-encoded preferences
	UpdatedAt   time.Time `json:"updated_at"`                            // Timestamp when the preferences were last updated
}

func (UserPreferenceVector) TableName() string {
	return "user_preference_vectors"
}

// Implement driver.Valuer for UserPreferenceVector
func (upv UserPreferenceVector) Value() (driver.Value, error) {
	return json.Marshal(upv)
}

// Implement driver.Scanner for UserPreferenceVector
func (upv *UserPreferenceVector) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for UserPreferenceVector scan")
	}

	return json.Unmarshal(data, &upv)
}

type ContentRecommendation struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`    // ID of the user to whom the recommendation is made
	ContentID     uint      `gorm:"index;not null" json:"content_id"` // ID of the recommended content
	Score         float64   `json:"score"`                            // Relevance score of the recommendation
	RecommendedAt time.Time `json:"recommended_at"`                   // Timestamp when the recommendation was made
}

func (ContentRecommendation) TableName() string {
	return "content_recommendations"
}

// Implement driver.Valuer for ContentRecommendation
func (cr ContentRecommendation) Value() (driver.Value, error) {
	return json.Marshal(cr)
}

// Implement driver.Scanner for ContentRecommendation
func (cr *ContentRecommendation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ContentRecommendation scan")
	}

	return json.Unmarshal(data, &cr)
}

type Community struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`      // Name of the community
	Description string  `gorm:"type:text" json:"description,omitempty"`      // Description of the community's purpose
	OwnerID     uint    `gorm:"index;not null" json:"owner_id"`              // ID of the user who owns or manages the community
	Members     []Users `gorm:"many2many:community_members;" json:"members"` // Users who are members of the community
	CreatedAt   time.Time
}

func (Community) TableName() string {
	return "communities"
}

// Implement driver.Valuer for Community
func (c Community) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Implement driver.Scanner for Community
func (c *Community) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for Community scan")
	}

	return json.Unmarshal(data, &c)
}

type CommunityPost struct {
	ID          uint      `gorm:"primaryKey"`
	CommunityID uint      `gorm:"index;not null" json:"community_id"`      // ID of the community where the post is made
	AuthorID    uint      `gorm:"index;not null" json:"author_id"`         // ID of the user who authored the post
	Title       string    `gorm:"type:varchar(255);not null" json:"title"` // Title of the post
	Body        string    `gorm:"type:text;not null" json:"body"`          // Body content of the post
	PostedAt    time.Time `json:"posted_at"`                               // Timestamp when the post was made
}

func (CommunityPost) TableName() string {
	return "community_posts"
}

// Implement driver.Valuer for CommunityPost
func (cp CommunityPost) Value() (driver.Value, error) {
	return json.Marshal(cp)
}

// Implement driver.Scanner for CommunityPost
func (cp *CommunityPost) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for CommunityPost scan")
	}

	return json.Unmarshal(data, &cp)
}

type FeatureFlag struct {
	ID             uint   `gorm:"primaryKey"`
	FeatureName    string `gorm:"type:varchar(255);unique;not null" json:"feature_name"` // Name of the feature
	IsEnabled      bool   `gorm:"not null;default:false" json:"is_enabled"`              // Whether the feature is enabled
	RolloutPercent uint   `json:"rollout_percent"`                                       // Percentage of users the feature is enabled for
}

func (FeatureFlag) TableName() string {
	return "feature_flags"
}

// Implement driver.Valuer for FeatureFlag
func (ff FeatureFlag) Value() (driver.Value, error) {
	return json.Marshal(ff)
}

// Implement driver.Scanner for FeatureFlag
func (ff *FeatureFlag) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for FeatureFlag scan")
	}

	return json.Unmarshal(data, &ff)
}

type KnowledgeBaseArticle struct {
	gorm.Model
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	AuthorID    uint      `json:"author_id" gorm:"index;not null"`
	PublishedAt time.Time `json:"published_at"`
	Category    string    `json:"category" gorm:"type:varchar(100);not null"`
	Tags        []Tag     `gorm:"many2many:knowledge_base_article_tags;" json:"tags"`
}

func (KnowledgeBaseArticle) TableName() string {
	return "knowledge_base_article"
}

// Implement driver.Valuer for KnowledgeBaseArticle
func (kba KnowledgeBaseArticle) Value() (driver.Value, error) {
	return json.Marshal(kba)
}

// Implement driver.Scanner for KnowledgeBaseArticle
func (kba *KnowledgeBaseArticle) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for KnowledgeBaseArticle scan")
	}

	return json.Unmarshal(data, &kba)
}

type ServiceDeskAnnouncement struct {
	gorm.Model
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Message     string     `json:"message" gorm:"type:text;not null"`
	PublishedAt time.Time  `json:"published_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Target      string     `json:"target" gorm:"type:varchar(100);not null"` // E.g., "Agents", "Customers", "All"
}

func (ServiceDeskAnnouncement) TableName() string {
	return "service_desk_announcement"
}

// Implement driver.Valuer for ServiceDeskAnnouncement
func (sda ServiceDeskAnnouncement) Value() (driver.Value, error) {
	return json.Marshal(sda)
}

// Implement driver.Scanner for ServiceDeskAnnouncement
func (sda *ServiceDeskAnnouncement) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for ServiceDeskAnnouncement scan")
	}

	return json.Unmarshal(data, &sda)
}

type FeedbackRequest struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	RequestText string    `json:"request_text" gorm:"type:text;not null"`
	SubmittedAt time.Time `json:"submitted_at"`
	Status      string    `json:"status" gorm:"type:varchar(100);not null"` // E.g., "Requested", "Completed"
}

func (FeedbackRequest) TableName() string {
	return "feedback_request"
}

// Implement driver.Valuer for FeedbackRequest
func (fr FeedbackRequest) Value() (driver.Value, error) {
	return json.Marshal(fr)
}

// Implement driver.Scanner for FeedbackRequest
func (fr *FeedbackRequest) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for FeedbackRequest scan")
	}

	return json.Unmarshal(data, &fr)
}
