package app

import (
	"context"
	"errors"
	"strings"
)

type CrudType struct {
	// Create, Read, Update, Delete
	Crud   string
	Create bool
	Read   bool
	Update bool
	Delete bool
}

// FormatCrud formats and validates a CRUD type string based on the provided context, user name, configuration.
// This function takes a context (`ctx`), user name (`name`), configuration (`c`), and a `CrudType` object (`crud`) as input.
func FormatCrud(ctx context.Context, name string, cfg *Config) error {
	// Check if user exists in config file and if crud exists in config file.
	if cfg.Users[name] != nil && cfg.Users[name].Crud != nil {
		crud := cfg.Users[name].Crud

		// Validate CRUD string length.
		if len(crud.Crud) < 1 || len(crud.Crud) > 4 {
			cfg.Users[name].Crud.Create = false
			cfg.Users[name].Crud.Read = false
			cfg.Users[name].Crud.Update = false
			cfg.Users[name].Crud.Delete = false
			return errors.New("invalid CRUD type string: length must be between 1 and 4")
		}

		// Convert CRUD string to lowercase and update the config.users.crud.crud string to be lowercase.
		cfg.Users[name].Crud.Crud = strings.ToLower(crud.Crud)

		// Initialize individual operation flags.
		var create, read, update, delete bool

		// Analyze each character and set corresponding flag.
		for _, ch := range cfg.Users[name].Crud.Crud {
			switch ch {
			case 'c':
				create = true
			case 'r':
				read = true
			case 'u':
				update = true
			case 'd':
				delete = true
			default:
				// Ignore invalid characters.
			}
		}
		// Update the fileds of the config.users.crud object.
		cfg.Users[name].Crud.Create = create
		cfg.Users[name].Crud.Read = read
		cfg.Users[name].Crud.Update = update
		cfg.Users[name].Crud.Delete = delete

		// Return formatted CrudType with updated flags.
		return nil
	} else {
		return errors.New("either user was not found in config file, or crud was not found in config file")
	}
}
