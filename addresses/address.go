package addresses

import (
	"artisons/db"
	"artisons/validators"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Address struct {
	Lastname      string `validate:"required"`
	Firstname     string `validate:"required"`
	City          string `validate:"required"`
	Street        string `validate:"required"`
	Complementary string
	Zipcode       string `validate:"required"`
	Phone         string `validate:"required"`
}

func (a Address) Validate(ctx context.Context) error {
	if err := validators.V.Struct(a); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the user", slog.String("error", err.Error()))
		field := err.(validator.ValidationErrors)[0]
		low := strings.ToLower(field.Field())
		return fmt.Errorf("input:%s", low)
	}

	return nil
}

func (a Address) Save(ctx context.Context, key string) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "saving the address")

	if key == "" {
		slog.LogAttrs(ctx, slog.LevelError, "cannot validate the key while it is empty")
		return errors.New("something went wrong")
	}

	if _, err := db.Redis.HSet(ctx, key,
		"firstname", a.Firstname,
		"lastname", a.Lastname,
		"complementary", a.Complementary,
		"city", a.City,
		"phone", a.Phone,
		"zipcode", a.Zipcode,
		"street", a.Street,
	).Result(); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "cannot store the user", slog.String("error", err.Error()))
		return errors.New("something went wrong")
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "the address is saved")

	return nil
}
