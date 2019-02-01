package cloud

import (
	"github.com/pkg/errors"
	"github.com/tcnksm/go-input"
)

func Ask(question, def string, required bool) string {
	ask := input.DefaultUI()
	ans, err := ask.Ask(question, &input.Options{
		Default:  def,
		Loop:     required,
		Required: required,
		ValidateFunc: func(s string) error {
			if required {
				if s == "" {
					return errors.New("empty input detected- input is required")
				}
			}
			if len(s) > 50 {
				return errors.New("input must be 50 characters or less")
			}
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	return ans
}

func Select(question string, opts []string, def string, required bool) string {
	ask := input.DefaultUI()
	ans, err := ask.Select(question, opts, &input.Options{
		Default:  def,
		Loop:     required,
		Required: required,
		ValidateFunc: func(s string) error {
			if required {
				if s == "" {
					return errors.New("empty input detected- input is required")
				}
			}
			if len(s) > 50 {
				return errors.New("input must be 50 characters or less")
			}
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	return ans
}

func SelectEnv() Env {
	ask := input.DefaultUI()
	ans, err := ask.Select("Please supply an environment to run in", []string{Local.String(), Gcloud.String(), Aws.String()}, &input.Options{
		Default:  Local.String(),
		Loop:     true,
		Required: true,
		ValidateFunc: func(s string) error {
			if s == "" {
				return errors.New("empty input detected- input is required")
			}
			if len(s) > 10 {
				return errors.New("input must be 10 characters or less")
			}
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	var env Env
	switch ans {
	case "local", "locl", "Local", "LOCAL":
		env = Local
	case "gcp", "gcloud", "google", "GCP", "Gcloud":
		env = Gcloud
	case "aws", "AWS", "amazon":
		env = Aws
	}
	return env
}
