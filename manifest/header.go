package manifest

import "time"

type manifestHeader struct {
	ID          string    `json:"ID"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}
