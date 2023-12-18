package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
)

// This file is an extension of https://pkg.go.dev/golang.org/x/net/webdav

// Dir is a custom webdav directory implementation that allows user configuration access for authentication.
// It extends the functionalities of the standard Dir by resolving paths based on user information and logging actions based on configuration settings.
type Dir struct {
	Config *Config
}

// resolveUser attempts to retrieve the username from the provided context.
// If the user is authenticated, their username is returned. Otherwise, an empty string is returned.
func (d Dir) resolveUser(ctx context.Context) string {
	authInfo := AuthFromContext(ctx)
	if authInfo != nil && authInfo.Authenticated {
		return authInfo.Username
	}

	return ""
}

// authorizationFromContext retrieves and formats the user's CRUD permissions based on the given context.
func (d Dir) authorizationFromContext(ctx context.Context) error {
	// Extract the authenticated user name from the provided context.
	user := d.resolveUser(ctx)
	// If no user is identified return an error
	if user == "" {
		return errors.New("no user identified")
	} else {
		// Format and validate the retrieved CRUD permissions for the identified user using the FormatCrud function.
		return FormatCrud(ctx, user, d.Config)
	}
}

// resolve builds the physical path for a given name based on user information and configuration settings.
// func (d Dir) resolve(ctx context.Context, name string) string {
// 	// Validate the name for any invalid characters or separators.
// 	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
// 		strings.Contains(name, "\x00") { // Null bytes are illegal in file names because they can be used to terminate strings prematurely and cause unexpected behavior.
// 		return ""
// 	}
// 	// Retrieve the base directory path from the configuration.
// 	dir := string(d.Config.Dir)
// 	// Use current directory if base directory is not set.
// 	if dir == "" {
// 		dir = "."
// 	}
// 	// Obtain authentication information from the context.
// 	authInfo := AuthFromContext(ctx)
// 	// Check if user is authenticated and has configured subdirectory.
// 	if authInfo != nil && authInfo.Authenticated {
// 		// Get user information from the configuration.
// 		userInfo := d.Config.Users[authInfo.Username]
// 		// If user has a configured subdirectory, append it to the path.
// 		if userInfo != nil && userInfo.Subdir != nil {
// 			return filepath.Join(dir, *userInfo.Subdir, filepath.FromSlash(path.Clean("/"+name)))
// 		}
// 	}
// 	// Build the final physical path by combining base directory and the provided name.
// 	return filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
// }

// Mkdir attempts to create a directory at the resolved physical path.
func (d Dir) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	// Resolve the physical path of the directory based on user information and configuration.
	if name = Resolve(ctx, name, d); name == "" {
		return os.ErrNotExist
	}
	// Get user authorization.
	err := d.authorizationFromContext(ctx)

	// Check for errors and return if any occur.
	if err != nil {
		return err
	}

	// resolve the user based on context.
	user := d.resolveUser(ctx)

	// Check for create permission.
	if !d.Config.Users[user].Crud.Create {
		if d.Config.Log.Create {
			log.WithField("user", user).Warn("unauthorized to create directory")
			return errors.New("unauthorized to create directory")
		} else {
			return nil
		}
	}

	// Create the directory using os.Mkdir.
	err = os.Mkdir(name, perm)
	// Check for errors and return if any occur.
	if err != nil {
		return err
	}
	// Log the directory creation action if logging is enabled in the configuration.
	if d.Config.Log.Create {
		log.WithFields(log.Fields{
			"path": name,
			"user": d.resolveUser(ctx),
		}).Info("Created directory")
	}

	return err
}

// OpenFile opens a file at the resolved physical path based on user permissions and returns a webdav.File object.
//
// This function takes a context (`ctx`), a file name (`name`), a flag (`flag`) indicating the access mode,
// and a permission mode (`perm`) for the file as input.
func (d Dir) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// Resolve the physical path of the file.
	if name = Resolve(ctx, name, d); name == "" {
		return nil, os.ErrNotExist
	}
	// Get user authorization.
	err := d.authorizationFromContext(ctx)

	// Check for errors and return if any occur.
	if err != nil {
		return nil, err
	}

	// resolve the user based on context.
	user := d.resolveUser(ctx)

	// Check for the file existence.
	_, err = os.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if !d.Config.Log.Create {
				log.WithFields(log.Fields{
					"path": name,
					"user": user,
				}).Warn("User does not have the permission to open a non-existant file they tried to create")
				return nil, errors.New("the file: " + name + " does not exist and user " + user + " has no write permission to create it")
			}
		}
	}

	// Check permissions based on access mode.
	if flag&os.O_RDONLY == 0 && !d.Config.Users[user].Crud.Read {
		return nil, errors.New("unauthorized to read file")
	}

	// Check if user has write permission, and also check if the operating system's file permissions allow writing.
	// This small section might seem a little counterintuitive, but if a user tries to create a file, they are also, automatically, trying
	// to open the that file. If they have read only permissions, they'll be able to open the any EXISTING file, but
	// if they have the permission of "read" ONLY and the file doesn't exist, they won't be able to create it, and
	// they shouldn't be able to open it, else an error will occur when the stats function inevitably runs on a non existsnt file.
	hasCreatePermission := d.Config.Users[user].Crud.Create
	if flag&(os.O_WRONLY|os.O_RDWR) != 0 && !hasCreatePermission {
		if !hasCreatePermission { // This user don't have the permission to create a file!
			if d.Config.Log.Create {
				log.WithField("user", user).Warn("unauthorized to create file")
			}
			return nil, nil
		} else { // This user has the permission to create a file, but the operating system's file permissions don't allow it.
			return nil, errors.New("unauthorized to write file based on the operating system's file permissions")
		}
	}

	// Open the file using os.OpenFile.
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	// Log the file opening action if configured.
	if d.Config.Log.Read {
		log.WithFields(log.Fields{
			"path": name,
			"user": user,
		}).Debug("Opened file")
	}
	// Return the opened file and nil error.
	return f, nil
}

// RemoveAll removes a file or directory at the resolved physical path based on user permissions.
func (d Dir) RemoveAll(ctx context.Context, name string) error {
	// Resolve the physical path of the file or directory.
	if name = Resolve(ctx, name, d); name == "" {
		return os.ErrNotExist
	}

	// Check if attempting to remove the virtual root directory.
	if name == filepath.Clean(string(d.Config.Dir)) {
		return errors.New("removing the virtual root directory is prohibited")
	}

	// Get user authorization.
	err := d.authorizationFromContext(ctx)

	// Check for errors and return if any occur.
	if err != nil {
		return err
	}

	// resolve the user based on context.
	user := d.resolveUser(ctx)

	// Check for delete permission.
	if !d.Config.Users[user].Crud.Delete {
		return errors.New("unauthorized to delete file or directory")
	}

	// Attempt to remove the file or directory using os.RemoveAll.
	err = os.RemoveAll(name)
	if err != nil {
		return err
	}

	// Log the deletion action if configured.
	if d.Config.Log.Delete {
		log.WithFields(log.Fields{
			"path": name,
			"user": user,
		}).Info("Deleted file or directory")
	}

	return nil
}

// Rename resolves the physical file and delegates this to an os.Rename execution
func (d Dir) Rename(ctx context.Context, oldName, newName string) error {
	// Resolve the physical paths of the old and new names.
	if oldName = Resolve(ctx, oldName, d); oldName == "" {
		return os.ErrNotExist
	}
	if newName = Resolve(ctx, newName, d); newName == "" {
		return os.ErrNotExist
	}

	// Check if attempting to rename the virtual root directory.
	if root := filepath.Clean(string(d.Config.Dir)); root == oldName || root == newName {
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	}

	// Get user authorization.
	err := d.authorizationFromContext(ctx)

	if err != nil {
		return err
	}

	// resolve the user based on context.
	user := d.resolveUser(ctx)

	// Check for rename permission.
	if !d.Config.Users[user].Crud.Update {
		return errors.New("unauthorized to rename file or directory")
	}

	// Attempt to rename the file or directory using os.Rename.
	err = os.Rename(oldName, newName)
	if err != nil {
		return err
	}

	// Log the rename action if configured.
	if d.Config.Log.Update {
		log.WithFields(log.Fields{
			"oldPath": oldName,
			"newPath": newName,
			"user":    user,
		}).Info("Renamed file or directory")
	}

	return nil
}

// Stat resolves the physical file and delegates this to an os.Stat execution
func (d Dir) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// 1. Resolve the provided path within the directory:
	name = Resolve(ctx, name, d)

	// 2. Check if the resolved path is empty (invalid).
	if name == "" {
		return nil, os.ErrNotExist
	}

	// 3. Determine the user accessing the file.
	user := d.resolveUser(ctx)

	// 4. Check if the user has read permission.
	if !d.Config.Users[user].Crud.Read {
		return nil, errors.New("unauthorized to read file")
	}

	// 5. Attempt to stat the resolved path.
	fileInfo, err := os.Stat(name)
	// 5.1 Handle different error cases:
	if err != nil {
		// File doesn't exist, and user is trying to create it when they don't have the permission to do so.
		if errors.Is(err, os.ErrNotExist) && d.Config.Users[user].Crud.Read && !d.Config.Users[user].Crud.Create {
			if d.Config.Log.Create { // Logging enabled for file creation
				log.WithFields(log.Fields{ // Log a slightly more detailed warning if file creation is not permitted.
					"path":  name,
					"user":  user,
					"crud":  d.Config.Users[user].Crud,
					"issue": "file does not exist and user does not have the write permission to create it",
				}).Warn("User does not have the write permission to create this file")
				return nil, nil
			}
		}
		// 5.2 Other errors (including the create permission issues) are passed along.
		return nil, err
	}

	// 6. If no errors, return the file information.
	return fileInfo, nil
}
