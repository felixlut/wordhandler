/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)
var (
	HOST = os.Getenv("TARGET_HOST")
	PORT = os.Getenv("TARGET_PORT")
	TYPE = "tcp"
)

func readFromConnection(connection net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return "", err
	}
	message := string(buffer[:mLen])
	return message, nil
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get stats for a specific word",
	Long:  `Get stats for a specific word`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]
		connection, err := net.Dial(TYPE, HOST+":"+PORT)
		if err != nil {
			log.Println("Unable to send 'get word' request:", err)
			os.Exit(1)
		}
		defer connection.Close()

		_, err = connection.Write([]byte(word))
		if err != nil {
			fmt.Println("Failed to write command")
			panic(err)
		}

		response, err := readFromConnection(connection)
		if err != nil {
			fmt.Println("Failed to read response")
			panic(err)
		}

		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
