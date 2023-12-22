package config

import "fmt"

type Validate interface {
	validate() error
}

func (c *Config) fullValidation() error {
	toValidate := []Validate{c, &c.CallerConfig, &c.SheetConfig, &c.NotificatorConfig}
	for _, item := range toValidate {
		if err := item.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CallerConfig) validate() error {
	switch {
	case c.Query == "":
		return emptyFieldError("query")
	case c.CountResult == 0:
		return emptyFieldError("count_result")
	case len(c.YouTubeApiKey) == 0:
		return emptyFieldError("api_keys")
	case c.TimeZone == "":
		return emptyFieldError("time_zone")
	default:
		return nil
	}
}

func (c *SheetConfig) validate() error {
	switch {
	case c.Name == "":
		return emptyFieldError("name")
	case c.SpreadsheetID == "":
		return emptyFieldError("id")
	case c.SheetID < 0:
		return fmt.Errorf("invalid sheet ID")
	default:
		return nil
	}
}

func (c *NotificationConfig) validate() error {
	switch {
	case c.TelegramBotApiKey == "":
		return emptyFieldError("telegram_api_key")
	case c.ReceiverID == 0:
		return emptyFieldError("receiver_id")
	default:
		return nil
	}
}

func (c *Config) validate() error {
	switch {
	case c.RequestInterval == 0:
		return emptyFieldError("interval_in_seconds")
	default:
		return nil
	}
}

func emptyFieldError(field string) error {
	return fmt.Errorf("%s is empty", field)
}
