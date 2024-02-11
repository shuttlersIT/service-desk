// backend/models/users.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	gorm.Model
	UserID                 uint              `gorm:"index;not null" json:"user_id"` // Associated user
	FirstName              string            `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName               string            `gorm:"type:varchar(255);not null" json:"last_name"`
	Email                  string            `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone                  *string           `gorm:"type:varchar(20);" json:"phone,omitempty"` // Phone number in E.164 format
	Assets                 []AssetAssignment `json:"assets" gorm:"foreignKey:UserID"`
	IsActive               bool              `gorm:"default:true" json:"is_active"`
	Projects               []Project         `gorm:"many2many:project_members;" json:"projects"`
	ProfilePic             string            `gorm:"size:255" json:"profile_pic"`
	PasswordHashed         string            `json:"-"`                                   // Excluded from JSON responses
	ResetPasswordRequestID *uint             `json:"reset_password_request_id,omitempty"` // Assuming this is optional
	Processed              bool              `json:"processed,omitempty"`
	TransactionLogID       *uint             `gorm:"index" json:"transaction_log_id,omitempty"` // Optional reference to a transaction log
	Amount                 float64           `gorm:"type:decimal(10,2);not null" json:"amount"` // Transaction amount
	Type                   string            `gorm:"type:varchar(100);not null" json:"type"`    // Type of transaction, e.g., "purchase", "refund"
}

func (Transaction) TableName() string {
	return "transaction"
}

type PaymentDetails struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"index;not null" json:"user_id"`             // User associated with the payment
	Amount      float64     `gorm:"type:decimal(10,2);not null" json:"amount"` // Amount of the payment
	Status      string      `gorm:"type:varchar(100);not null" json:"status"`  // Status of the payment, e.g., "completed", "pending"
	PaymentInfo PaymentInfo `gorm:"embedded;embeddedPrefix:payment_"`          // Embedded struct for payment info
}

func (PaymentDetails) TableName() string {
	return "payment_details"
}

type PaymentInfo struct {
	Method        string `gorm:"type:varchar(100);not null" json:"method"`         // Payment method, e.g., "credit_card", "paypal"
	TransactionID string `gorm:"type:varchar(255);not null" json:"transaction_id"` // Identifier for the payment transaction
}

type TransactionLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`                    // User associated with the transaction
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`        // Amount of the transaction
	TransactionID string    `gorm:"type:varchar(255);not null" json:"transaction_id"` // External or internal transaction ID
	Status        string    `gorm:"type:varchar(100);not null" json:"status"`         // Status of the transaction, e.g., "completed", "failed"
	CreatedAt     time.Time `json:"created_at"`
}

func (TransactionLog) TableName() string {
	return "transaction_logs"
}

type OrderDetails struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	UserID       uint         `gorm:"index;not null" json:"user_id"`                   // User who placed the order
	TotalAmount  float64      `gorm:"type:decimal(10,2);not null" json:"total_amount"` // Total amount of the order
	Status       string       `gorm:"type:varchar(100);not null" json:"status"`        // Current status of the order
	PaymentInfo  PaymentInfo  `gorm:"embedded;embeddedPrefix:payment_"`                // Embedded struct for payment info
	Items        []OrderItem  `gorm:"foreignKey:OrderID" json:"items"`                 // Associated order items
	ShipmentInfo ShipmentInfo `gorm:"embedded;embeddedPrefix:shipment_"`               // Embedded struct for shipment info
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (OrderDetails) TableName() string {
	return "order_details"
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	OrderID   uint `gorm:"index;not null" json:"order_id"` // Associated order
	ProductID uint `gorm:"not null" json:"product_id"`     // Product being ordered
	Quantity  int  `gorm:"not null" json:"quantity"`       // Quantity of the product ordered
}

func (OrderItem) TableName() string {
	return "order_items"
}

type ShipmentInfo struct {
	Address    string `gorm:"type:text;not null" json:"address"`             // Shipment address
	Carrier    string `gorm:"type:varchar(100);not null" json:"carrier"`     // Carrier service used for shipment
	TrackingID string `gorm:"type:varchar(255);not null" json:"tracking_id"` // Tracking ID for shipment
}

func (ShipmentInfo) TableName() string {
	return "shipment_info"
}

type Expenses struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	UserID        uint    `gorm:"index;not null" json:"user_id"`             // User who incurred the expense
	AgentID       uint    `json:"agent_id"`                                  // Optional agent ID for expense association
	TransactionID uint    `gorm:"index;not null" json:"transaction_id"`      // Associated transaction for the expense
	Amount        float64 `gorm:"type:decimal(10,2);not null" json:"amount"` // Amount of the expense
	Status        string  `gorm:"type:varchar(100);not null" json:"status"`  // Status of the expense, e.g., "pending", "reimbursed"
}

func (Expenses) TableName() string {
	return "expenses"
}

type UserWallet struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	UserID  uint    `gorm:"index;not null" json:"user_id"`              // Foreign key for the user
	Version int     `gorm:"not null" json:"version"`                    // Version of the wallet for optimistic locking
	Balance float64 `gorm:"type:decimal(10,2);not null" json:"balance"` // Current balance of the wallet
}

// TableName sets the table name for the UserWallet model.
func (UserWallet) TableName() string {
	return "user_wallet"
}

type TransactionStorage interface {
	CreateUser(*Users) error
	DeleteUser(int) error
	UpdateUser(*Users) error
	GetUsers() ([]*Users, error)
	GetUserByID(int) (*Users, error)
	GetUserByNumber(int) (*Users, error)
}

// UserModel handles database operations for User
type TransactionDBModel struct {
	DB *gorm.DB
}

// NewUserModel creates a new instance of UserModel
func NewTransactionDBModel(db *gorm.DB) *TransactionDBModel {
	return &TransactionDBModel{
		DB: db,
	}
}

// ProcessPayment processes a payment and updates related financial records with transactional integrity.
func (db *UserDBModel) ProcessPayment(paymentDetails PaymentDetails) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Record the payment.
		if err := tx.Create(&paymentDetails).Error; err != nil {
			return err
		}

		// Step 2: Update the user's balance or subscription status.
		if err := tx.Model(&Users{}).Where("id = ?", paymentDetails.UserID).
			Update("subscription_status", "active").Error; err != nil {
			return err
		}

		// Step 3: Log the transaction for auditing purposes.
		transactionLog := TransactionLog{
			UserID:        paymentDetails.UserID,
			Amount:        paymentDetails.Amount,
			TransactionID: paymentDetails.ID,
			Status:        "completed",
		}
		if err := tx.Create(&transactionLog).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetDepartmentExpenseReport generates an expense report for departments, potentially using transactional reads for consistency.
func (db *UserDBModel) GetDepartmentExpenseReport(departmentID uint) (*Expenses, error) {
	var report Expenses

	// Transactional read might be used for ensuring a consistent snapshot in high concurrency scenarios.
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Table("expenses").
			Select("department_id, SUM(amount) as total_expense").
			Where("department_id = ?", departmentID).
			Group("department_id").
			Scan(&report).Error
	})

	return &report, err
}

// UpdateUserBalanceOptimistically performs an optimistic locking update to a user's balance.
func (db *UserDBModel) UpdateUserBalanceOptimistically(userID uint, amount float64) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var userWallet UserWallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&userWallet).Error; err != nil {
			return err
		}

		// Assume there's a version or timestamp field for optimistic concurrency control.
		originalVersion := userWallet.Version

		userWallet.Balance += amount
		userWallet.Version++

		// Attempt to save, ensuring the version hasn't changed in the meantime.
		result := tx.Model(&UserWallet{}).Where("id = ? AND version = ?", userID, originalVersion).Save(&userWallet)
		if result.RowsAffected == 0 {
			// No rows affected implies a concurrent modification.
			return fmt.Errorf("concurrent update error")
		}
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

// RefundPurchase handles refund operations, which might involve compensating transactions if external systems are involved.
func (db *UserDBModel) RefundPurchase(purchaseID uint) error {
	// Assume initial database update to mark refund as initiated.
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Transaction{}).Where("id = ?", purchaseID).Update("status", "refund_initiated").Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Call to external refund service.
	refundResult, err := externalRefundService.ProcessRefund(purchaseID)
	if err != nil {
		// Handle external service failure, potentially initiating a compensating transaction.
		// This might involve updating the status back, notifying administrators, or logging for manual intervention.
		return err
	}

	// Finalize refund process based on refundResult.
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Update based on refund result.
		return nil
	})
}

// ExecuteOrderSaga orchestrates the creation of an order, including inventory check, payment processing, and shipment.
func (db *UserDBModel) ExecuteOrderSaga(order OrderDetails) error {
	// Step 1: Reserve inventory.
	if err := Assets.ReserveInventory(order.Items); err != nil {
		return err
	}

	// Step 2: Process payment.
	if err := paymentService.ProcessPayment(order.PaymentInfo); err != nil {
		// Compensating transaction: release inventory.
		Assets.ReleaseInventory(order.Items)
		return err
	}

	// Step 3: Initiate shipment.
	if err := shippingService.ScheduleShipment(order.ShipmentInfo); err != nil {
		// Compensating transactions: release inventory and refund payment.
		Assets.ReleaseInventory(order.Items)
		paymentService.RefundPayment(order.PaymentInfo)
		return err
	}

	// Finalize order in the database.
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		return nil
	})
}

/*
// CompletePurchaseTransaction handles a purchase operation that involves external payment services.
func (db *UserDBModel) CompletePurchaseTransaction(purchase Transaction) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Update internal records to indicate a pending purchase.
		if err := tx.Model(&Transaction{}).Create(&purchase).Error; err != nil {
			return err
		}

		// Step 2: Call external payment service.
		paymentResult, err := externalPaymentService.ProcessPayment(Transaction.Type)
		if err != nil {
			// Log or handle the external service error.
			// Decide whether to rollback based on business logic.
			return err
		}

		// Step 3: Update the purchase record based on the payment result.
		if err := tx.Model(&Transaction{}).Where("id = ?", purchase.ID).Updates(map[string]interface{}{"status": paymentResult.Status}).Error; err != nil {
			return err
		}

		// Additional steps based on payment result.
		return nil
	})
}
*/
