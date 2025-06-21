package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func ClearConsole() {
	// Очистка консоли (Windows)
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func search_data() {
	if _, err := os.Stat("data.json"); err == nil {
		return
	} else if os.IsNotExist(err) {
		fmt.Println("database not found. Creating...")
		file, err := os.Create("data.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		fmt.Println("database created.")
		return
	} else {
		fmt.Println("Error checking file:", err)
		return
	}
}
func GeneratePassword() string {
	const (
		lettersLower = "abcdefghijklmnopqrstuvwxyz"
		lettersUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits       = "0123456789"
	)
	fmt.Println("Generating password...")
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(7) + 9 // Random length between 9 and 15
	password := []byte{
		lettersUpper[rand.Intn(len(lettersUpper))],
		digits[rand.Intn(len(digits))],
	}
	// Остальные символы — случайные буквы (англ) и цифры
	allChars := lettersLower + lettersUpper + digits
	for len(password) < length {
		password = append(password, allChars[rand.Intn(len(allChars))])
	}

	// Перемешиваем символы
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})
	fmt.Println("Password generated successfully!")
	return string(password)
}
func SavePassword(name, password string) {
	// Читаем существующие данные
	data := make(map[string]string)
	fileBytes, err := ioutil.ReadFile("data.json")
	if err == nil && len(fileBytes) > 0 {
		json.Unmarshal(fileBytes, &data)
	}

	// Добавляем новую запись
	data[name] = password

	// Сохраняем обратно в файл
	newBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	err = ioutil.WriteFile("data.json", newBytes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Printf("Password '%s' for '%s' saved successfully!\n", password, name)
}
func RemovePassword(delname string) {
	fileBytes, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	data := make(map[string]string)
	if len(fileBytes) > 0 {
		err = json.Unmarshal(fileBytes, &data)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	}

	if _, exists := data[delname]; exists {
		delete(data, delname)
		newBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		err = ioutil.WriteFile("data.json", newBytes, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		fmt.Printf("Password for '%s' removed successfully!\n", delname)
	} else {
		fmt.Printf("No password found for '%s'.\n", delname)
	}
}
func Settings() {
	for {
		ClearConsole()
		fmt.Println("Settings Menu:")
		fmt.Println("1. Remove password")
		fmt.Println("2. Change password settings (not implemented yet.)")
		fmt.Println("3. Back to main menu")

		var choice int
		fmt.Print("Your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			ClearConsole()
			fmt.Println("Enter the name of the password to remove:")
			var delname string
			fmt.Scanln(&delname)
			RemovePassword(delname)

		case 2:
			ClearConsole()
			fmt.Println("Change password settings is not implemented yet.")
		case 3:
			return // Возврат в главное меню
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 3.")
		}
	}
}
func main() {
	search_data()
	for {
		ClearConsole()
		fmt.Println("Hi, I'll come up with a password for you.")
		fmt.Println("Please choose one of the following options:")
		fmt.Println("1. Generate a password")
		fmt.Println("2. View passwords")
		fmt.Println("3. Settings")
		fmt.Println("4. Exit")

		var chise int
		fmt.Print("Your choice: ")
		fmt.Scanln(&chise)

		switch chise {
		case 1:
			ClearConsole()
			fmt.Print("Enter name of the password: ")
			var name string
			fmt.Scanln(&name)
			password := GeneratePassword()
			fmt.Println("Your password:", password)
			SavePassword(name, password)
			fmt.Println("\nPress Enter to return to menu...")
			fmt.Scanln() // просто ждёт ввода

		case 2:
			fileBytes, err := ioutil.ReadFile("data.json")
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			data := make(map[string]string)
			if len(fileBytes) > 0 {
				err = json.Unmarshal(fileBytes, &data)
				if err != nil {
					fmt.Println("Error decoding JSON:", err)
					continue
				}
				if len(data) == 0 {
					fmt.Println("No passwords saved yet.")
				} else {
					fmt.Println("Saved passwords:")
					for name, password := range data {
						fmt.Printf("%s: %s\n", name, password)
					}
				}
			} else {
				fmt.Println("No passwords saved yet.")
			}
			fmt.Println("\nPress Enter to return to menu...")
			fmt.Scanln() // просто ждёт ввода

		case 3:
			Settings()
		case 4:
			fmt.Println("Goodbye!")
			return // завершение программы

		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 4.")
		}
	}
}
