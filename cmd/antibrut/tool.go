package main

import (
	"errors"
	"fmt"

	"github.com/romsar/antibrut/config"
	proto "github.com/romsar/antibrut/proto/antibrut/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	toolCmd = &cobra.Command{
		Use:   "tool",
		Short: "Полезные команды для работы с приложением.",
	}

	resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Выполнить сброс данных о попытках сделать запрос.",
		RunE:  reset,
	}

	whiteListCmd = &cobra.Command{
		Use:   "wl",
		Short: "Управление белым списком.",
	}

	whiteListAddIPCmd = &cobra.Command{
		Use:   "add",
		Short: "Добавить подсеть в белый список.",
		RunE:  addIPToWhiteList,
	}

	whiteListDeleteIPCmd = &cobra.Command{
		Use:   "delete",
		Short: "Удалить подсеть из белого списка.",
		RunE:  deleteIPFromWhiteList,
	}

	blackListCmd = &cobra.Command{
		Use:   "bl",
		Short: "Управление черным списком.",
	}

	blackListAddIPCmd = &cobra.Command{
		Use:   "add",
		Short: "Добавить подсеть в черным список.",
		RunE:  addIPToBlackList,
	}

	blackListDeleteIPCmd = &cobra.Command{
		Use:   "delete",
		Short: "Удалить подсеть из черного списка.",
		RunE:  deleteIPFromBlackList,
	}
)

func init() {
	rootCmd.AddCommand(toolCmd)

	toolCmd.AddCommand(resetCmd)
	resetCmd.Flags().String("login", "", "Логин.")
	resetCmd.Flags().String("ip", "", "IP-адрес.")

	toolCmd.AddCommand(whiteListCmd)
	whiteListCmd.AddCommand(whiteListAddIPCmd)
	whiteListCmd.AddCommand(whiteListDeleteIPCmd)

	toolCmd.AddCommand(blackListCmd)
	blackListCmd.AddCommand(blackListAddIPCmd)
	blackListCmd.AddCommand(blackListDeleteIPCmd)
}

func reset(cmd *cobra.Command, args []string) error {
	// flags
	login, err := cmd.Flags().GetString("login")
	if err != nil {
		return err
	}

	ip, err := cmd.Flags().GetString("ip")
	if err != nil {
		return err
	}

	// validation
	if login == "" && ip == "" {
		return errors.New("укажите логин или IP")
	}

	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// grpc
	conn, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer conn.Close()

	// client
	service := proto.NewAntiBrutServiceClient(conn)

	// do request
	_, err = service.Reset(cmd.Context(), &proto.ResetRequest{
		Login: login,
		Ip:    ip,
	})
	if err != nil {
		return fmt.Errorf("произошла ошибка в процессе работы: %w", err)
	}

	cmd.Println("успешно!")

	return nil
}

func addIPToWhiteList(cmd *cobra.Command, args []string) error {
	// validation
	if len(args) < 1 {
		return errors.New("укажите подсеть")
	}

	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// grpc
	conn, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer conn.Close()

	// client
	service := proto.NewAntiBrutServiceClient(conn)

	// do request
	_, err = service.AddIPToWhiteList(cmd.Context(), &proto.AddIPToWhiteListRequest{
		Subnet: args[0],
	})
	if err != nil {
		return fmt.Errorf("произошла ошибка в процессе работы: %w", err)
	}

	cmd.Println("успешно!")

	return nil
}

func deleteIPFromWhiteList(cmd *cobra.Command, args []string) error {
	// validation
	if len(args) < 1 {
		return errors.New("укажите подсеть")
	}

	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// grpc
	conn, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer conn.Close()

	// client
	service := proto.NewAntiBrutServiceClient(conn)

	// do request
	_, err = service.DeleteIPFromWhiteList(cmd.Context(), &proto.DeleteIPFromWhiteListRequest{
		Subnet: args[0],
	})
	if err != nil {
		return fmt.Errorf("произошла ошибка в процессе работы: %w", err)
	}

	cmd.Println("успешно!")

	return nil
}

func addIPToBlackList(cmd *cobra.Command, args []string) error {
	// validation
	if len(args) < 1 {
		return errors.New("укажите подсеть")
	}

	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// grpc
	conn, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer conn.Close()

	// client
	service := proto.NewAntiBrutServiceClient(conn)

	// do request
	_, err = service.AddIPToBlackList(cmd.Context(), &proto.AddIPToBlackListRequest{
		Subnet: args[0],
	})
	if err != nil {
		return fmt.Errorf("произошла ошибка в процессе работы: %w", err)
	}

	cmd.Println("успешно!")

	return nil
}

func deleteIPFromBlackList(cmd *cobra.Command, args []string) error {
	// validation
	if len(args) < 1 {
		return errors.New("укажите подсеть")
	}

	// config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// grpc
	conn, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу: %w", err)
	}
	defer conn.Close()

	// client
	service := proto.NewAntiBrutServiceClient(conn)

	// do request
	_, err = service.DeleteIPFromBlackList(cmd.Context(), &proto.DeleteIPFromBlackListRequest{
		Subnet: args[0],
	})
	if err != nil {
		return fmt.Errorf("произошла ошибка в процессе работы: %w", err)
	}

	cmd.Println("успешно!")

	return nil
}
