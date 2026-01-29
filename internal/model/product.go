// Пакет model содержит структуры данных, используемые в приложении.
package model

// Product представляет товар одежды, хранится в MongoDB.
type Product struct {
	ID           string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string   `bson:"name" json:"name" validate:"required"`
	URL          string   `bson:"url" json:"url" validate:"required,url"`
	Price        int      `bson:"price" json:"price" validate:"required,min=0"`
	Brand        string   `bson:"brand" json:"brand" validate:"required"`
	PhotosURLs   []string `bson:"photos_urls" json:"photos_urls" validate:"required,dive,url"`
	Availability bool     `bson:"availability" json:"availability"`
	Description  string   `bson:"description" json:"description"`
}


