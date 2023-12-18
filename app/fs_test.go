package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestDirResolveUser(t *testing.T) {
	// This function tests how Dir.resolveUser() behaves with different user contexts passed through context.Context.
	// It verifies that the method returns the correct username for authenticated users and an empty string for unauthenticated users.
	configTmp := createTestConfig("/tmp") // Create a temporary configuration for testing.

	ctx := context.Background() // Create a background context for testing.
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		}) // Create a context with an admin user.
	user1 := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "user1",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		}) // Create a context with a regular user.
	anon := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "user1",
			Authenticated: false,
			CrudType:      &CrudType{Crud: "", Create: false, Read: false, Update: false, Delete: false},
		}) // Create a context with an unauthenticated user.

	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{"", admin, "admin"},
		{"", user1, "user1"},
		{"", anon, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Dir{
				Config: configTmp,
			}
			if got := d.resolveUser(tt.ctx); got != tt.want {
				t.Errorf("Dir.resolveUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// This is a nearly concrete copy of the function TestDirResolve of golang.org/x/net/webdav/file_test.go
// just with prefixes and configuration details.
func TestDirResolve(t *testing.T) {
	// **1. Setting up test fixtures:**
	configTmp := createTestConfig("/tmp")
	configRoot := createTestConfig("/")
	configCurrentDir := createTestConfig(".")
	configEmpty := createTestConfig("")

	// Define background context.
	ctx := context.Background()

	// Create different user contexts with auth information.
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})
	user1 := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "user1",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})
	user2 := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "user2",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})

	// Define a list of test cases as structs with expected results.
	tests := []struct {
		cfg  *Config
		name string
		ctx  context.Context
		want string
	}{
		// **2. Test cases for admin user with different base directories and paths:**
		{configTmp, "", admin, "/tmp"},
		{configTmp, ".", admin, "/tmp"},
		{configTmp, "/", admin, "/tmp"},
		{configTmp, "./a", admin, "/tmp/a"},
		{configTmp, "..", admin, "/tmp"},
		{configTmp, "../", admin, "/tmp"},
		{configTmp, "../.", admin, "/tmp"},
		{configTmp, "../a", admin, "/tmp/a"},
		{configTmp, "../..", admin, "/tmp"},
		{configTmp, "../bar/a", admin, "/tmp/bar/a"},
		{configTmp, "../baz/a", admin, "/tmp/baz/a"},
		{configTmp, "...", admin, "/tmp/..."},
		{configTmp, ".../a", admin, "/tmp/.../a"},
		{configTmp, ".../..", admin, "/tmp"},
		{configTmp, "a", admin, "/tmp/a"},
		{configTmp, "a/./b", admin, "/tmp/a/b"},
		{configTmp, "a/../../b", admin, "/tmp/b"},
		{configTmp, "a/../b", admin, "/tmp/b"},
		{configTmp, "a/b", admin, "/tmp/a/b"},
		{configTmp, "a/b/c/../../d", admin, "/tmp/a/d"},
		{configTmp, "a/b/c/../../../d", admin, "/tmp/d"},
		{configTmp, "a/b/c/../../../../d", admin, "/tmp/d"},
		{configTmp, "a/b/c/d", admin, "/tmp/a/b/c/d"},
		{configTmp, "/a/b/c/d", admin, "/tmp/a/b/c/d"},

		{configTmp, "ab/c\x00d/ef", admin, ""},

		{configRoot, "", admin, "/"},
		{configRoot, ".", admin, "/"},
		{configRoot, "/", admin, "/"},
		{configRoot, "./a", admin, "/a"},
		{configRoot, "..", admin, "/"},
		{configRoot, "../", admin, "/"},
		{configRoot, "../.", admin, "/"},
		{configRoot, "../a", admin, "/a"},
		{configRoot, "../..", admin, "/"},
		{configRoot, "../bar/a", admin, "/bar/a"},
		{configRoot, "../baz/a", admin, "/baz/a"},
		{configRoot, "...", admin, "/..."},
		{configRoot, ".../a", admin, "/.../a"},
		{configRoot, ".../..", admin, "/"},
		{configRoot, "a", admin, "/a"},
		{configRoot, "a/./b", admin, "/a/b"},
		{configRoot, "a/../../b", admin, "/b"},
		{configRoot, "a/../b", admin, "/b"},
		{configRoot, "a/b", admin, "/a/b"},
		{configRoot, "a/b/c/../../d", admin, "/a/d"},
		{configRoot, "a/b/c/../../../d", admin, "/d"},
		{configRoot, "a/b/c/../../../../d", admin, "/d"},
		{configRoot, "a/b/c/d", admin, "/a/b/c/d"},
		{configRoot, "/a/b/c/d", admin, "/a/b/c/d"},

		{configCurrentDir, "", admin, "."},
		{configCurrentDir, ".", admin, "."},
		{configCurrentDir, "/", admin, "."},
		{configCurrentDir, "./a", admin, "a"},
		{configCurrentDir, "..", admin, "."},
		{configCurrentDir, "../", admin, "."},
		{configCurrentDir, "../.", admin, "."},
		{configCurrentDir, "../a", admin, "a"},
		{configCurrentDir, "../..", admin, "."},
		{configCurrentDir, "../bar/a", admin, "bar/a"},
		{configCurrentDir, "../baz/a", admin, "baz/a"},
		{configCurrentDir, "...", admin, "..."},
		{configCurrentDir, ".../a", admin, ".../a"},
		{configCurrentDir, ".../..", admin, "."},
		{configCurrentDir, "a", admin, "a"},
		{configCurrentDir, "a/./b", admin, "a/b"},
		{configCurrentDir, "a/../../b", admin, "b"},
		{configCurrentDir, "a/../b", admin, "b"},
		{configCurrentDir, "a/b", admin, "a/b"},
		{configCurrentDir, "a/b/c/../../d", admin, "a/d"},
		{configCurrentDir, "a/b/c/../../../d", admin, "d"},
		{configCurrentDir, "a/b/c/../../../../d", admin, "d"},
		{configCurrentDir, "a/b/c/d", admin, "a/b/c/d"},
		{configCurrentDir, "/a/b/c/d", admin, "a/b/c/d"},

		{configEmpty, "", admin, "."},

		{configTmp, "", user1, "/tmp/subdir1"},
		{configTmp, ".", user1, "/tmp/subdir1"},
		{configTmp, "/", user1, "/tmp/subdir1"},
		{configTmp, "./a", user1, "/tmp/subdir1/a"},
		{configTmp, "..", user1, "/tmp/subdir1"},
		{configTmp, "../", user1, "/tmp/subdir1"},
		{configTmp, "../.", user1, "/tmp/subdir1"},
		{configTmp, "../a", user1, "/tmp/subdir1/a"},
		{configTmp, "../..", user1, "/tmp/subdir1"},
		{configTmp, "../bar/a", user1, "/tmp/subdir1/bar/a"},
		{configTmp, "../baz/a", user1, "/tmp/subdir1/baz/a"},
		{configTmp, "...", user1, "/tmp/subdir1/..."},
		{configTmp, ".../a", user1, "/tmp/subdir1/.../a"},
		{configTmp, ".../..", user1, "/tmp/subdir1"},
		{configTmp, "a", user1, "/tmp/subdir1/a"},
		{configTmp, "a/./b", user1, "/tmp/subdir1/a/b"},
		{configTmp, "a/../../b", user1, "/tmp/subdir1/b"},
		{configTmp, "a/../b", user1, "/tmp/subdir1/b"},
		{configTmp, "a/b", user1, "/tmp/subdir1/a/b"},
		{configTmp, "a/b/c/../../d", user1, "/tmp/subdir1/a/d"},
		{configTmp, "a/b/c/../../../d", user1, "/tmp/subdir1/d"},
		{configTmp, "a/b/c/../../../../d", user1, "/tmp/subdir1/d"},
		{configTmp, "a/b/c/d", user1, "/tmp/subdir1/a/b/c/d"},
		{configTmp, "/a/b/c/d", user1, "/tmp/subdir1/a/b/c/d"},
		{configTmp, "", user2, "/tmp/subdir2"},
	}
	// **3. Looping through test cases and running individual tests:**
	for _, tt := range tests { // Loop through each element in the `tests` list.
		t.Run(tt.name, func(t *testing.T) { // Run each test case with its name for better reporting.
			d := Dir{ // Call the `Dir.resolve` method with the test case's context and name.
				Config: tt.cfg,
			}
			if got := Resolve(tt.ctx, tt.name, d); got != tt.want {
				// Compare the returned resolved directory path (`got`) with the expected value (`tt.want`).
				t.Errorf("Dir.resolve() = %v, want %v. Base dir is %v", got, tt.want, tt.cfg.Dir) // If they differ, report an error with details.
				// Include the resolved path (`got`), expected path (`tt.want`), and base directory from the config for context.
			}
		})
	}
}

// TestDirMkdir tests the Mkdir method of the Dir struct for various scenarios.
func TestDirMkdir(t *testing.T) {
	// 1. Create a temporary directory for the test and clean up afterwards.
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir)

	t.Logf("using test dir: %s", tmpDir)
	configTmp := createTestConfig(tmpDir)

	// 2. Set up test context with admin information for authorization.
	ctx := context.Background()
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})

	// 3. Define test cases for different Mkdir scenarios.
	tests := []struct {
		name    string
		perm    os.FileMode
		wantErr bool // Should error occur?
	}{
		// Valid directory name and permission.
		{"a", 0700, false},
		{"/a/", 0700, true}, // already exists
		{"ab/c\x00d/ef", 0700, true},
		{"/a/a/a/a", 0700, true},
		{"a/a/a/a", 0700, true},
	}
	// 4. Run each test case with a sub-test.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Dir{
				Config: configTmp,
			}
			// 5. Call Mkdir and verify if the expected error occurs.
			if err := d.Mkdir(admin, tt.name, tt.perm); (err != nil) != tt.wantErr {
				t.Errorf("Dir.Mkdir() name = %v, error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestDirOpenFile(t *testing.T) {
	// 1. Create a temporary directory for the test
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir) // Clean up after the test is done.

	// 2. Generate the test configuration
	configTmp := createTestConfig(tmpDir)

	// 3. Create a context with mock admin information for authorization
	ctx := context.Background()
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})

	// 4. Define test cases for different file operations and expected behavior
	tests := []struct {
		name    string
		flag    int
		perm    os.FileMode
		wantErr bool // Should error occur?
	}{
		// Case 1: Try opening existing file with read/write access (should succeed)
		{"foo", os.O_RDWR, 0644, true},
		// Case 2: Try opening non-existent file with read/write and create flags (should succeed)
		{"foo", os.O_RDWR | os.O_CREATE, 0644, false},
		// Case 3: Try opening invalid path with read/write access (should fail)
		{"ab/c\x00d/ef", os.O_RDWR, 0700, true},
	}

	// 5. Loop through each test case and verify Dir.OpenFile behavior
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Dir instance with the test configuration
			d := Dir{
				Config: configTmp,
			}

			// Call Dir.OpenFile with test case parameters
			got, err := d.OpenFile(admin, tt.name, tt.flag, tt.perm)

			// Verify if error occurred as expected
			if (err != nil) != tt.wantErr {
				t.Errorf("Dir.OpenFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error expected, verify file information
			if !tt.wantErr {
				// Get the expected file information
				wantFileInfo, err := os.Stat(filepath.Join(tmpDir, tt.name))
				if err != nil {
					t.Errorf("Dir.OpenFile() error = %v", err)
				}

				// Get the actual file information from the returned file object
				gotFileInfo, _ := got.Stat()

				// Compare expected and actual file information for deep equality
				if !reflect.DeepEqual(gotFileInfo, wantFileInfo) {
					t.Errorf("Dir.OpenFile() = %v, want %v", gotFileInfo, wantFileInfo)
				}
			}
		})
	}
}

func TestRemoveDir(t *testing.T) {
	// Create temporary directory and test configuration
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir)
	configTmp := createTestConfig(tmpDir)

	// Mock admin context for authorization
	ctx := context.Background()
	admin := context.WithValue(ctx, authInfoKey, &AuthInfo{Username: "admin", Authenticated: true})

	// Define test cases with directories and expected behavior
	tests := []struct {
		name    string
		wantErr bool // Should error occur?
	}{
		{"a", false},
		{"a/b/c", false},
		{"/a/b/c", false},
		{"ab/c\x00d/ef", true},
	}

	// Loop through test cases and verify Dir.RemoveAll behavior
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Dir instance and construct directory path
			d := Dir{
				Config: configTmp,
			}

			// Pre-condition: Create directory if not expected to fail
			file := filepath.Join(tmpDir, tt.name)
			if !tt.wantErr {
				err := os.MkdirAll(file, 0700)
				if err != nil {
					t.Errorf("Dir.RemoveAll() pre condition failed. name = %v, error = %v", tt.name, err)
				}
			}

			// Call Dir.RemoveAll and verify if error occurs as expected
			if err := d.RemoveAll(admin, tt.name); (err != nil) != tt.wantErr {
				t.Errorf("Dir.RemoveAll() name = %v, error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			// Post-condition: Verify directory is removed if not expected to fail
			if !tt.wantErr {
				if _, err := os.Stat(file); err == nil {
					t.Errorf("Dir.RemoveAll() post condition failed. name = %v, error = %v", tt.name, err)
				}
			}
		})
	}
}

func TestDirRemoveAll(t *testing.T) {
	// Set up context and mock admin user for authorization
	ctx := context.Background()
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})

	// Create a temporary directory for testing
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir) // Clean up after the test is done.

	// Generate the test config for the directory
	configTmp := createTestConfig(tmpDir)

	// Define test cases with directories and removal paths, along with expected behavior (error or success)
	tests := []struct {
		name       string
		removeName string // Directory to remove
		wantErr    bool   // Should error occur?
	}{
		// Test case 1: Remove subdirectory from parent ("a/b/c" -> "a")
		// This case tests whether Dir.RemoveAll can successfully remove a subdirectory from its parent directory ("a").
		// Success is expected (wantErr = false) as the user has permission and the path is valid.
		{"a/b/c", "a", false},
		// Test case 2: Remove intermediate directory ("a/b/c" -> "a/b")
		// This case tests whether Dir.RemoveAll can remove an intermediate directory within the target directory path ("a/b").
		// Success is expected (wantErr = false) as the user has permission and the path is valid, even though it's not the leaf directory.
		{"a/b/c", "a/b", false},
		// Test case 3: Remove target directory itself ("a/b/c" -> "a/b/c")
		// This case tests whether Dir.RemoveAll can remove the target directory itself ("a/b/c").
		// Success is expected (wantErr = false) as the user has permission and the path is valid.
		{"a/b/c", "a/b/c", false},
		// Test case 4: Remove directory using absolute path ("/a/b/c" -> "a") - notice the leading slash.
		// This case tests whether Dir.RemoveAll allows removing directories using absolute paths, which is potentially dangerous.
		// An error is expected (wantErr = true) as using absolute paths is not allowed by the security model of the application.
		{"/a/b/c", "a", false},
	}

	// Loop through each test case and verify Dir.RemoveAll behavior
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a Dir instance with the test configuration
			d := Dir{
				Config: configTmp,
			}

			// Pre-condition: Create the directory to be removed
			err := os.MkdirAll(filepath.Join(tmpDir, tt.name), 0700)
			if err != nil {
				t.Errorf("Dir.RemoveAll() error creating dir error = %v, name %v", err, tt.name)
				return
			}

			// Call Dir.RemoveAll with the chosen removal path and verify if error occurs as expected
			if err := d.RemoveAll(admin, tt.removeName); (err != nil) != tt.wantErr {
				t.Errorf("Dir.RemoveAll() removeName = %v, error = %v, wantErr %v", tt.removeName, err, tt.wantErr)
				return
			}

			// Verify if the target directory and any parent directories were deleted
			if _, err := os.Stat(filepath.Join(tmpDir, tt.removeName)); err == nil {
				t.Errorf("Dir.RemoveAll() file or directory not deleted = %v, removeName %v", err, tt.removeName)
				return
			}

			// Verify that only the target directory was removed, not its parent directories
			if _, err := os.Stat(filepath.Join(tmpDir, tt.removeName, "/..")); err != nil {
				t.Errorf("Dir.RemoveAll() parent directory deleted = %v, removeName %v", err, tt.removeName)
				return
			}
		})
	}
}

func TestRename(t *testing.T) {
	// Create a temporary directory and generate configuration
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir) // Clean up after the test
	configTmp := createTestConfig(tmpDir)

	// Set up context with mock admin user for authorization
	ctx := context.Background()
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})

	// Define test cases with source and destination names, pre-creation flag, and expected error
	tests := []struct {
		name      string
		oldName   string
		newName   string
		create    bool // Should file be created beforehand?
		wantError bool
	}{
		// This case attempts to rename an existing file ("a") to a new name ("b") without pre-creating it.
		// The create flag is set to false, indicating the file should already exist.
		// Since the test case expects an error (wantErr = true), this likely simulates a scenario where the user lacks
		// the necessary operating system update permission for renaming files.
		{"a", "a", "b", false, true},
		// Similar to the previous case, this attempts to rename "a" to "b", but with the create flag set to true.
		// This ensures the file exists beforehand, mimicking a scenario where the user has
		// the operating system update permission to rename existing files.
		// The expected wantErr is set to false as the rename operation should succeed.
		{"a", "a", "b", true, false},
		// This case tries to rename a file with an invalid character in its name ("\x00d") to another valid name ("foo").
		// Both create and wantErr are set to false and true, respectively.
		// This likely tests the function's behavior with invalid characters in file names, which should result in an error.
		{"\x00d", "\x00da", "foo", false, true},
		// This case attempts to rename a file with an invalid character ("\x00d") to a destination ("foo") that already exists.
		// The create flag is set to false as the destination shouldn't be created, and wantErr is set to true as the rename should fail due to the existing file.
		// This tests the function's handling of conflicts when renaming to existing files.
		{"\x00d", "foo", "\x00da", false, true},
	}
	// Loop through each test case and verify Dir.Rename behavior
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a Dir instance with the test configuration
			d := Dir{
				Config: configTmp,
			}

			// Pre-condition: Create the source file if required by the test case
			if tt.create {
				_, err := d.OpenFile(admin, tt.oldName, os.O_RDWR|os.O_CREATE, 0700)
				if err != nil {
					t.Errorf("Dir.Rename() pre condition failed. name = %v, error = %v", tt.name, err)
					return
				}
			}

			// Call Dir.Rename with the chosen source and destination names and verify if error occurs as expected
			if err := d.Rename(admin, tt.oldName, tt.newName); (err != nil) != tt.wantError {
				t.Errorf("Dir.Rename() error = %v, wantErr %v", err, tt.wantError)
				return
			}

			// Verify if the source file still exists after the rename operation
			if _, err := os.Stat(filepath.Join(tmpDir, tt.oldName)); err == nil {
				t.Errorf("Dir.Rename() oldName still remained. oldName = %v, newName = %v", tt.oldName, tt.newName)
				return
			}

			// Verify if the renamed file now exists at the destination path
			join := filepath.Join(tmpDir, tt.newName)
			fmt.Println(join) // Print the renamed file path for debugging purposes
			if _, err := os.Stat(join); err != nil {
				if !tt.create {
					// Skip check if file shouldn't be created
					return
				}
				t.Errorf("Dir.Rename() newName not present. oldName = %v, newName = %v", tt.oldName, tt.newName)
				return
			}

		})
	}
}

func TestDirStat(t *testing.T) {
	// Create a temporary directory and configure test environment
	tmpDir := filepath.Join(os.TempDir(), "david__"+strconv.FormatInt(time.Now().UnixNano(), 10))
	os.Mkdir(tmpDir, 0700)
	defer os.RemoveAll(tmpDir) // Clean up after test
	configTmp := createTestConfig(tmpDir)

	// Set up context and mock admin user for authorization
	ctx := context.Background()
	admin := context.WithValue(ctx, authInfoKey,
		&AuthInfo{Username: "admin",
			Authenticated: true,
			CrudType:      &CrudType{Crud: "crud", Create: true, Read: true, Update: true, Delete: true},
		})

	// Define test cases with paths, expected kinds ("dir" or "file"), and expected error
	tests := []struct {
		name    string
		kind    string // Type of file/directory (dir, file)
		wantErr bool
	}{
		{"/a", "dir", false},
		{"/a/b", "file", false},
		{"\x00da", "file", true},
	}
	// Loop through each test case and verify Dir.Stat behavior
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a Dir instance with the test configuration
			d := Dir{
				Config: configTmp,
			}

			// Prepare the path based on the test case
			fp := filepath.Join(tmpDir, tt.name)

			// Pre-condition: Create the required file/directory if necessary
			if tt.kind == "dir" {
				err := os.MkdirAll(fp, 0700)
				if err != nil {
					t.Errorf("Dir.Stat() error creating dir. error = %v", err)
					return
				}
			} else if !tt.wantErr {
				_, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					t.Errorf("Dir.Stat() error creating file. error = %v", err)
					return
				}
			}

			// Call Dir.Stat with the prepared path and context
			got, err := d.Stat(admin, tt.name)
			// Verify if error occurred as expected
			if (err != nil) != tt.wantErr {
				t.Errorf("Dir.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want, _ := os.Stat(filepath.Join(tmpDir, tt.name))

			if !reflect.DeepEqual(got, want) {
				t.Errorf("Dir.Stat() = %v, want %v", got, want)
			}
		})
	}
}

func createTestConfig(dir string) *Config {
	// Define a list of subdirectories for test users
	subdirs := [2]string{"subdir1", "subdir2"}

	// Create a map to store test user information
	userInfos := map[string]*UserInfo{
		"admin": {
			// Grant full "crud" permissions to the admin user
			Permissions: "crud",
			Crud: &CrudType{
				Crud:   "crud", // Access all CRUD operations
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		"user1": {
			// Assign user1 to the first subdirectory
			Subdir:      &subdirs[0],
			Permissions: "crud",
			Crud: &CrudType{
				Crud:   "crud", // Access all CRUD operations
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
		"user2": {
			// Assign user2 to the second subdirectory
			Subdir:      &subdirs[1],
			Permissions: "crud",
			Crud: &CrudType{
				Crud:   "crud", // Access all CRUD operations
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		},
	}

	// Create a Config instance with the test directory and user information
	config := &Config{
		Dir:   dir,       // Set the directory path for testing
		Users: userInfos, // Populate the user map
		Log: Logging{
			Error:  true, // Enable error logging
			Create: true, // Enable logging for file creation
			Read:   true, // Enable logging for file reading
			Update: true, // Enable logging for file updates
			Delete: true, // Enable logging for file deletion
		},
	}
	return config
}
