package pkg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Item struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Quantity    int        `json:"quantity" db:"quantity"`
	Description string     `json:"description" db:"description"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

func (i *Item) Bind(r *http.Request) error {
	err1 := validate.Validate(
		&validators.StringIsPresent{Name: "name", Field: i.Name, Message: fmt.Sprintf("The %s value is required", "name")},
		&validators.StringIsPresent{Name: "description", Field: i.Description, Message: fmt.Sprintf("The %s value is required", "description")},
		&validators.IntIsPresent{Name: "quantity", Field: i.Quantity, Message: fmt.Sprintf("The %s value is required", "quantity")},
		&validators.StringLengthInRange{Name: "description", Field: i.Description, Min: 2, Max: 500, Message: fmt.Sprintf("The %s value must be between %d and %d characters", "description", 2, 500)},
		&validators.StringLengthInRange{Name: "name", Field: i.Name, Min: 2, Max: 50, Message: fmt.Sprintf("The %s value must be between %d and %d characters", "name", 2, 50)},
	)

	if err1.HasAny() {
		return err1
	}
	return nil
}

type DeletedComment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Comment   string    `json:"comment" db:"comment"`
	Item      string    `json:"item" db:"item"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (d *DeletedComment) Bind(r *http.Request) error {
	err1 := validate.Validate(
		&validators.StringIsPresent{Name: "name", Field: d.Comment, Message: fmt.Sprintf("The %s value is required", "name")},
		&validators.StringIsPresent{Name: "item", Field: d.Item, Message: fmt.Sprintf("The %s value is required", "description")},
		&validators.StringLengthInRange{Name: "name", Field: d.Comment, Min: 2, Max: 500, Message: fmt.Sprintf("The %s value must be between %d and %d characters", "name", 2, 500)},
		&validators.FuncValidator{
			Name: "item",
			Fn: func() bool {
				if _, err := uuid.FromString(d.Item); err != nil {
					return false
				}
				return true
			},
			Field:   d.Item,
			Message: fmt.Sprintf("The %s value is required", "item"),
		},
	)

	if err1.HasAny() {
		return err1
	}

	return nil
}
