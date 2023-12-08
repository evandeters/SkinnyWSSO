// verification.go

package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbFile                         = "./verify.db" // SQLite database file
	smtpUsername                   = "kaminoverify@gmail.com"
	smtpPassword                   = "#ad6pYEHX8q@!b3A"
	smtpHost                       = "smtp.gmail.com"
	smtpPort                       = 587
	verificationCodeExpireDuration = 15 * time.Minute // Adjust as needed
)

func verify(c *gin.Context) {
	fmt.Println("Got to verify")

	// Parse form data
	email := c.PostForm("email")
	verificationCode := c.PostForm("verificationCode")

	fmt.Println("Received verification request:")
	fmt.Println("Email:", email)
	fmt.Println("Verification Code:", verificationCode)

	// Retrieve the user from the database based on email
	user, err := getUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Check if the verification code matches and has not expired
	if isValidVerificationCode(user, verificationCode) {
		// Mark the user as verified in the database
		user.IsVerified = true
		err = updateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// You can redirect the user to a success page or return a success JSON response
		c.JSON(200, gin.H{"message": "Email verification successful."})
	} else {
		// You can redirect the user to a failure page or return a failure JSON response
		c.JSON(401, gin.H{"error": "Invalid or expired verification code."})
	}
}

func initializeDatabase() error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute a query to create the table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			email TEXT,
			verification_code TEXT,
			verification_expiry DATETIME,
			is_verified BOOLEAN
		);
	`)
	if err != nil {
		return err
	}

	// Check if the table is empty (you might want to refine this check based on your needs)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	// If the table is empty, insert dummy data
	if count == 0 {
		// Insert dummy data
		dummyUser := User{
			Username:           "test",
			Email:              "test@test.com",
			VerificationCode:   "123456",
			VerificationExpiry: time.Now().Add(verificationCodeExpireDuration),
			IsVerified:         true,
		}

		_, err := db.Exec(`
			INSERT INTO users (username, email, verification_code, verification_expiry, is_verified)
			VALUES (?, ?, ?, ?, ?)
		`, dummyUser.Username, dummyUser.Email, dummyUser.VerificationCode, dummyUser.VerificationExpiry, dummyUser.IsVerified)

		if err != nil {
			return err
		}
	}

	return nil
}

// User struct to represent a registered user
type User struct {
	Username           string
	Email              string
	VerificationCode   string
	VerificationExpiry time.Time
	IsVerified         bool
}

func confirmation(c *gin.Context) {
	// Parse form data
	username := c.PostForm("username")
	email := c.PostForm("email")

	// Generate a random verification code
	verificationCode := generateVerificationCode()

	// Save user data in the database
	user := User{
		Username:           username,
		Email:              email,
		VerificationCode:   verificationCode,
		VerificationExpiry: time.Now().Add(verificationCodeExpireDuration),
		IsVerified:         false,
		// Set other fields as needed
	}

	// Save the user in the database
	err := saveUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Send verification email
	sendVerificationEmail(user)

	// You can redirect the user to a confirmation page or return a success message
	c.JSON(200, gin.H{"message": "Registration successful. Please check your email for verification.", "user": user})
}

func sendVerificationEmail(user User) {
	from := "skinnywsso@gmail.com"
	pass := "bzei uxxz ecef sdmi"
	to := user.Email

	body := fmt.Sprintf("Your verification code is: %s", user.VerificationCode)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Email Verification\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully to:", user.Email)
}

func generateVerificationCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func isValidVerificationCode(user User, enteredCode string) bool {
	return user.VerificationCode == enteredCode && time.Now().Before(user.VerificationExpiry)
}

func saveUser(user User) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO users (username, email, verification_code, verification_expiry, is_verified)
		VALUES (?, ?, ?, ?, ?)
	`, user.Username, user.Email, user.VerificationCode, user.VerificationExpiry, false)

	return err
}

func getUserByEmail(email string) (User, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return User{}, err
	}
	defer db.Close()

	user := User{}
	err = db.QueryRow(`
		SELECT id, username, email, verification_code, verification_expiry, is_verified
		FROM users
		WHERE email = ?
	`, email).Scan(
		&user.Username, &user.Email, &user.VerificationCode, &user.VerificationExpiry, &user.IsVerified,
	)

	return user, err
}

func updateUser(user User) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE users
		SET is_verified = ?
		WHERE id = ?
	`, user.IsVerified)

	return err
}
