// verification.go

package main


import (
    "database/sql"
    "fmt"
    "github.com/gin-gonic/gin"
    _ "github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
    "gopkg.in/gomail.v2"
	"net/http"
    "math/rand"
    "time"
)

var (
    dbFile                     = "./verify.db" // SQLite database file
    smtpUsername               = "kaminoverify@gmail.com"
    smtpPassword               = "#ad6pYEHX8q@!b3A"
    smtpHost                   = "smtp.gmail.com"
    smtpPort                   = 587
    verificationCodeExpireDuration = 15 * time.Minute // Adjust as needed
)

// User struct to represent a registered user
type User struct {
	Username          string
	Email             string
	VerificationCode  string
	VerificationExpiry time.Time
	IsVerified        bool
}

func confirmation(c *gin.Context) {
	// Parse form data
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Generate a random verification code
	verificationCode := generateVerificationCode()

	// Save user data in the database
	user := User{
		Username:          username,
		Email:             email,
		VerificationCode:  verificationCode,
		VerificationExpiry: time.Now().Add(verificationCodeExpireDuration),
		IsVerified:        false,
		// Set other fields as needed
	}

	// Save the user in the database
	err = saveUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Send verification email
	sendVerificationEmail(user)

	// You can redirect the user to a confirmation page or return a success message
	c.JSON(200, gin.H{"message": "Registration successful. Please check your email for verification."})
}

func sendVerificationEmail(user User) {
    // Set up the email message
    m := gomail.NewMessage()
    m.SetHeader("From", "noreply@example.com")
    m.SetHeader("To", user.Email)
    m.SetHeader("Subject", "Email Verification")
    m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", user.VerificationCode))

    // Set up the SMTP client
    d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

    // Send the email
    if err := d.DialAndSend(m); err != nil {
        fmt.Println("Error sending email:", err)
        return
    }

    fmt.Println("Email sent successfully to:", user.Email)
}

func verify(c *gin.Context) {
	// Parse form data
	email := c.PostForm("email")
	verificationCode := c.PostForm("verificationCode")

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
