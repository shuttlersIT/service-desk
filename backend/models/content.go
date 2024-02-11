// backend/models/incident_db_model.go

package models

import (
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

type ContentCategory struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);unique;not null" json:"name"` // Category name
	Description string `gorm:"type:text" json:"description,omitempty"`        // Category description
}

func (ContentCategory) TableName() string {
	return "content_categories"
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

type ArticleTag struct {
	gorm.Model
	ArticleID uint `gorm:"index;not null" json:"article_id"` // ID of the article
	TagID     uint `gorm:"index;not null" json:"tag_id"`     // ID of the tag
}

func (ArticleTag) TableName() string {
	return "article_tags"
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

type TenantScopedModel struct {
	TenantID uint `gorm:"index;not null" json:"tenant_id"` // Foreign key linking to the Tenant
	// Include fields common to tenant-scoped models here
}

func (TenantScopedModel) TableName() string {
	return "tenant_scoped_model"
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

type ContentItem struct {
	ID          uint   `gorm:"primaryKey"`
	Type        string `gorm:"type:varchar(100);not null" json:"type"` // Type of content (e.g., article, video, image)
	ReferenceID uint   `gorm:"not null" json:"reference_id"`           // ID of the actual content, interpretation depends on Type
	Metadata    string `gorm:"type:text" json:"metadata,omitempty"`    // JSON for additional metadata about the item
}

func (ContentItem) TableName() string {
	return "content_items"
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

type SystemMetric struct {
	ID         uint      `gorm:"primaryKey"`
	MetricType string    `gorm:"type:varchar(100);not null" json:"metric_type"` // Type of metric (e.g., CPU, memory usage)
	Value      float64   `json:"value"`                                         // Metric value
	RecordedAt time.Time `json:"recorded_at"`                                   // Timestamp when the metric was recorded
}

func (SystemMetric) TableName() string {
	return "system_metrics"
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

type FeatureFlag struct {
	ID             uint   `gorm:"primaryKey"`
	FeatureName    string `gorm:"type:varchar(255);unique;not null" json:"feature_name"` // Name of the feature
	IsEnabled      bool   `gorm:"not null;default:false" json:"is_enabled"`              // Whether the feature is enabled
	RolloutPercent uint   `json:"rollout_percent"`                                       // Percentage of users the feature is enabled for
}

func (FeatureFlag) TableName() string {
	return "feature_flags"
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
