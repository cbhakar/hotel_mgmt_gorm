package store

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/lib/pq"
	"log"
)

const (
	host     = "raja.Db.elephantsql.com"
	port     = 5432
	user     = "stzrurfj"
	password = "mHLqxoPKfj2P0R5XD2AImPSr8Ozu7rWr"
	dbname   = "stzrurfj"
)

var Db *gorm.DB

func InititateDBConn() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	gormDB, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("error : ", err.Error())
		panic(err)
	}
	log.Println("database connection created")

	Db = gormDB
	Db.DropTableIfExists(&Room{})
	Db.DropTableIfExists(&Plan{})
	Db.DropTableIfExists(&Hotel{})

	Db.AutoMigrate(&Hotel{}, &Room{}, &Plan{})
	Db.Model(&Room{}).AddForeignKey("hotel_id", "hotels(hotel_id)", "CASCADE", "CASCADE")
	Db.Model(&Plan{}).AddForeignKey("hotel_id", "hotels(hotel_id)", "CASCADE", "CASCADE")

}

func init() {
	InititateDBConn()
}

type Hotel struct {
	HotelID     string         `gorm:"size:255;primary_key" json:"hotel_id"`
	Name        string         `gorm:"size:255;not null;unique" json:"name"`
	Country     string         `gorm:"size:50;not null" json:"country"`
	Address     string         `gorm:"size:500;not null" json:"address"`
	Latitude    float64        `gorm:"not null" json:"latitude"`
	Longitude   float64        `gorm:"not null" json:"longitude"`
	Telephone   string         `gorm:"size:50;not null" json:"telephone"`
	Amenities   pq.StringArray `gorm:"type:text[];size:255;not null" json:"amenities"`
	Description string         `gorm:"size:500;not null" json:"description"`
	RoomCount   int            `gorm:"not null" json:"room_count"`
	Currency    string         `gorm:"size:50;not null" json:"currency"`
}

type Cap struct {
	MaxAdults     int `gorm:"not null" json:"max_adults"`
	ExtraChildren int `gorm:"not null" json:"extra_children"`
}

type Capacity Cap

type Room struct {
	HotelID     string   `gorm:"foreignkey:Hotel_id" json:"hotel_id"`
	RoomID      string   `gorm:"size:255;primary_key" json:"room_id"`
	Description string   `gorm:"size:255;not null" json:"description"`
	Name        string   `gorm:"size:255;not null" json:"name"`
	HtlCapacity Capacity `gorm:"type:jsonb" json:"capacity"`
}

type CP struct {
	Type              string `gorm:"size:255;not null" json:"type"`
	ExpiresDaysBefore int    `gorm:"not null" json:"expires_days_before"`
}
type CancelationlPolicy []CP

type Plan struct {
	HotelID         string             `gorm:"foreignkey:Hotel_id" json:"hotel_id"`
	RatePlanID      string             `gorm:"size:255;primary_key" json:"rate_plan_id"`
	CclPolicy       CancelationlPolicy `gorm:"type:jsonb" json:"cancellation_policy"`
	Name            string             `gorm:"size:255;not null" json:"name"`
	OtherConditions pq.StringArray     `gorm:"type:text[];size:255;not null" json:"other_conditions"`
	MealPlan        string             `gorm:"size:255;not null" json:"meal_plan"`
}

type Data struct {
	Offers []struct {
		CmOfferID    string `json:"cm_offer_id"`
		Htl          Hotel  `json:"hotel"`
		Rm           Room   `json:"room"`
		RP           Plan   `json:"rate_plan"`
		OriginalData struct {
			GuaranteePolicy struct {
				Required bool `json:"Required"`
			} `json:"GuaranteePolicy"`
		} `json:"original_data"`
		Capacity struct {
			MaxAdults     int `json:"max_adults"`
			ExtraChildren int `json:"extra_children"`
		} `json:"capacity"`
		Number   int    `json:"number"`
		Price    int    `json:"price"`
		Currency string `json:"currency"`
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
		Fees     []struct {
			Type        string  `json:"type"`
			Description string  `json:"description"`
			Included    bool    `json:"included"`
			Percent     float64 `json:"percent"`
		} `json:"fees"`
	} `json:"offers"`
}

//using only basic validation due to time constraint
func (gr *Data) Validate() (err error) {

	if len(gr.Offers) < 1 {
		return errors.New("the 'offers' field is required")
	}
	return
}

func AddNewOffer(u Data) (err error) {
	if Db == nil {
		err = errors.New("unable to connect to database")
		log.Println("error : ", err.Error())
		return
	}

	tx := Db.Begin()

	for _, offer := range u.Offers {

		err = tx.Debug().Model(&Hotel{}).Create(&offer.Htl).Error
		if err != nil {
			tx.Rollback()
			log.Println("error inserting in Hotel table : ", err.Error())
			return
		}
		err = tx.Debug().Model(&Plan{}).Create(&offer.RP).Error
		if err != nil {
			tx.Rollback()
			log.Println("error inserting in Rate Plan table : ", err.Error())
			return
		}
		err = tx.Debug().Model(&Room{}).Create(&offer.Rm).Error
		if err != nil {
			tx.Rollback()
			log.Println("error inserting in Room table : ", err.Error())
			return
		}
	}
	defer tx.Commit()
	return
}

// value scanner for custom data types
func (p CancelationlPolicy) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	if err != nil {
		fmt.Println("err -> ", err)
	}
	return string(j), err
}

func (p *CancelationlPolicy) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i CancelationlPolicy
	err := json.Unmarshal(source, &i)
	if err != nil {
		fmt.Println("scan err -> ", err)
		return err
	}

	*p = i
	return nil
}

func (p Capacity) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	if err != nil {
		fmt.Println("err -> ", err)
	}
	return string(j), err
}

func (p *Capacity) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i Capacity
	err := json.Unmarshal(source, &i)
	if err != nil {
		fmt.Println("scan err -> ", err)
		return err
	}

	*p = i
	return nil
}

func CloseDbConn() (err error) {
	err = Db.Close()
	return
}
