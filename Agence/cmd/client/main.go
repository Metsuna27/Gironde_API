package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"example.com/Gironde_API_G5/pb"

	"golang.org/x/term"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("\n===== GIRONDE API - AUTH =====")
		fmt.Println("1. Inscription")
		fmt.Println("2. Connexion")
		fmt.Println("3. Quitter")

		choice := readInput(reader, "Votre choix : ")

		switch choice {

		case "1":
			register(client, reader)

		case "2":
			login(client, reader)

		case "3":
			fmt.Println("Au revoir !")
			return

		default:
			fmt.Println("Choix invalide")
		}
	}
}

func readInput(reader *bufio.Reader, message string) string {

	fmt.Print(message)

	value, _ := reader.ReadString('\n')

	return strings.TrimSpace(value)
}

func readPassword() string {

	fmt.Print("Mot de passe : ")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))

	fmt.Println()

	if err != nil {
		log.Fatal(err)
	}

	return string(password)
}

func register(
	client pb.AuthServiceClient,
	reader *bufio.Reader,
) {

	firstname := readInput(reader, "Prénom : ")
	lastname := readInput(reader, "Nom : ")
	email := readInput(reader, "Email : ")
	password := readPassword()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	response, err := client.Register(
		ctx,
		&pb.RegisterRequest{
			Firstname: firstname,
			Lastname:  lastname,
			Email:     email,
			Password:  password,
		},
	)

	if err != nil {
		fmt.Println("Erreur gRPC :", err)
		return
	}

	fmt.Println(response.GetMessage())

	if response.GetSuccess() {
		fmt.Println("Utilisateur ID :", response.GetUserId())
	}
}

func login(
	client pb.AuthServiceClient,
	reader *bufio.Reader,
) {

	email := readInput(reader, "Email : ")
	password := readPassword()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	response, err := client.Login(
		ctx,
		&pb.LoginRequest{
			Email:    email,
			Password: password,
		},
	)

	if err != nil {
		fmt.Println("Erreur gRPC :", err)
		return
	}

	fmt.Println(response.GetMessage())

	if response.GetSuccess() {
		fmt.Println("Utilisateur ID :", response.GetUserId())
	}
}
