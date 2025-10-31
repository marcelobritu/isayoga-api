package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EnrollmentID        primitive.ObjectID `json:"enrollment_id" bson:"enrollment_id"`
	MercadoPagoID       string             `json:"mercado_pago_id" bson:"mercado_pago_id"`
	Status              string             `json:"status" bson:"status"`
	AmountInCents       int64              `json:"amount_in_cents" bson:"amount_in_cents"`
	PaymentMethod       string             `json:"payment_method" bson:"payment_method"`
	PreferenceID        string             `json:"preference_id" bson:"preference_id"`
	InitPointURL        string             `json:"init_point_url" bson:"init_point_url"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewPayment(enrollmentID primitive.ObjectID, amountInCents int64) *Payment {
	now := time.Now()
	return &Payment{
		ID:            primitive.NewObjectID(),
		EnrollmentID:  enrollmentID,
		Status:        "pending",
		AmountInCents: amountInCents,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (p *Payment) UpdateFromMercadoPago(mpID, status, paymentMethod string) {
	p.MercadoPagoID = mpID
	p.Status = status
	p.PaymentMethod = paymentMethod
	p.UpdatedAt = time.Now()
}

func (p *Payment) SetPreference(preferenceID, initPointURL string) {
	p.PreferenceID = preferenceID
	p.InitPointURL = initPointURL
	p.UpdatedAt = time.Now()
}

