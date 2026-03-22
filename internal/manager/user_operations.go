package manager

import (
	"context"
	"fmt"
)

func (s *Manager) PrintUsers(ctx context.Context) error {
	if err := s.Repository.TxnManager.WithDB(ctx, func(ctx context.Context) error {
		users, err := s.UserService.AllUsers(ctx)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("--------------------------------------------")
		fmt.Println("Printing all user names: ")

		for _, u := range users {
			fmt.Println(u.Username)
		}
		fmt.Println("--------------------------------------------")

		return nil
	}); err != nil {
		return fmt.Errorf("error printing users: %w", err)
	}
	return nil
}

func (s *Manager) ResetUserPassword(ctx context.Context, username string) error {
	if err := s.Repository.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		pw, err := s.UserService.ResetUserPassword(ctx, username)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("--------------------------------------------")
		fmt.Printf("Password for user '%s' has been reset.\n", username)
		fmt.Printf("New password: %s\n", pw)
		fmt.Println("--------------------------------------------")

		return nil
	}); err != nil {
		return fmt.Errorf("error resetting user password: %w", err)
	}
	return nil
}
