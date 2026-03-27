package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/users"
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  `Manage user accounts.`,
}

var UserListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all users in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := listUsers()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var UserCreateCmd = &cobra.Command{
	Use:   "create <username>",
	Short: "Create a new user",
	Long:  `Create a new user account.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := createUser(args[0])
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

var UserDeleteCmd = &cobra.Command{
	Use:   "delete <user-id>",
	Short: "Delete a user",
	Long:  `Delete a user account.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm("Are you sure you want to delete this user?") {
			fmt.Println("Operation cancelled")
			return
		}
		err := deleteUser(args[0])
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
}

func listUsers() error {
	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return fmt.Errorf("failed to get users database: %w", err)
	}

	var usersList []*storage.User
	err = usersDB.Iterate([]byte("user:"), func(key, value []byte) error {
		var u storage.User
		if err := storage.Decode(value, &u); err != nil {
			return err
		}
		usersList = append(usersList, &u)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUSERNAME\tEMAIL\tROLE")
	fmt.Fprintln(w, "--\t--------\t-----\t---")

	for _, u := range usersList {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", u.ID, u.Username, u.Email, u.Role)
	}

	w.Flush()
	return nil
}

func createUser(username string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter display name: ")
	displayName, _ := reader.ReadString('\n')
	displayName = strings.TrimSpace(displayName)

	fmt.Print("Enter role (admin, manager, user): ")
	role, _ := reader.ReadString('\n')
	role = strings.TrimSpace(role)

	if role != "admin" && role != "manager" && role != "user" {
		return fmt.Errorf("invalid role: %s (must be admin, manager, or user)", role)
	}

	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return fmt.Errorf("failed to get users database: %w", err)
	}

	userStore := users.New(usersDB)

	userID := "user-" + fmt.Sprintf("%d", time.Now().UnixNano())
	if err := userStore.Create(userID, username, password, email, displayName, role); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("User '%s' created with ID '%s'\n", username, userID)
	return nil
}

func deleteUser(userID string) error {
	reader := bufio.NewReader(os.Stdin)
	_ = reader

	dbMgr, err := openDBManager("./data/david")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer dbMgr.CloseAll()

	usersDB, err := dbMgr.Get(storage.DBUsers)
	if err != nil {
		return fmt.Errorf("failed to get users database: %w", err)
	}

	userStore := users.New(usersDB)

	if err := userStore.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	fmt.Printf("User '%s' deleted successfully\n", userID)
	return nil
}
