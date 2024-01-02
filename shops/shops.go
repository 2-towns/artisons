package shops

import (
	"context"
	"errors"
	"fmt"
	"gifthub/conf"
	"gifthub/db"
	"gifthub/users"
	"gifthub/validators"
	"log/slog"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Shop struct {
	Logo string `validate:"required"`

	// The lastname is not used.
	// The firstname is the shop name.
	Address users.Address
}

func Get(c context.Context) (Shop, error) {
	slog.LogAttrs(c, slog.LevelInfo, "get the shop info")

	data, err := db.Redis.HGetAll(context.Background(), "shop").Result()
	if err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot get the shop info", slog.String("error", err.Error()))
		return Shop{}, errors.New("error_http_general")
	}

	return Shop{
		Logo: data["logo"],
		Address: users.Address{
			Lastname:      data["address_lastname"],
			Firstname:     data["address_firstname"],
			City:          data["address_city"],
			Street:        data["address_street"],
			Complementary: data["address_complementary"],
			Zipcode:       data["address_zipcode"],
			Phone:         data["address_phone"],
		},
	}, nil
}

func (s Shop) Save(c context.Context) error {
	l := slog.With(slog.String("name", s.Address.Firstname))
	l.LogAttrs(c, slog.LevelInfo, "trying to save the shop")

	s.Address.Lastname = "None"

	if err := validators.V.Struct(s); err != nil {
		slog.LogAttrs(c, slog.LevelError, "cannot validate the admin", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input_%s_required", low)
	}

	now := time.Now()
	_, err := db.Redis.HSet(context.Background(), "shop",
		"logo", fmt.Sprintf("%s/%s", conf.ImgProxy.Path, s.Logo),
		"address_lastname", s.Address.Lastname,
		"address_firstname", s.Address.Firstname,
		"address_street", s.Address.Street,
		"address_city", s.Address.City,
		"address_complementary", s.Address.Complementary,
		"address_zipcode", s.Address.Zipcode,
		"address_phone", s.Address.Phone,
		"updated_at", now.Unix(),
		"created_at", now.Unix(),
	).Result()

	if err != nil {
		l.LogAttrs(c, slog.LevelError, "cannot save the shop", slog.String("error", err.Error()))
		return errors.New("error_http_general")
	}

	return nil
}
