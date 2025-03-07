package auth

import (
	"crypto/rand"
	"errors"
	"math/big"
	"maxwarden/config"
	"maxwarden/entries"
	"maxwarden/security"
	"maxwarden/users"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

const (
	UNAUTHORIZED_MESSAGE string = "Unauthorized access"
)

type Identity struct {
	UserID        int32
	SecurityStamp string
	MasterKey     string
	Permissions   users.Permissions
	Authenticated bool
	Expiration    time.Time
}

func NewIdentity(userid int32, securityStamp string, masterKey string, rememberMe bool) *Identity {
	expirationDuration := time.Duration(time.Hour * 24 * time.Duration(config.IDENTITY_COOKIE_EXPIRY_DAYS))
	expiration := time.Now().Add(expirationDuration)

	// hash password input
	mk := security.SHA512_58(masterKey)

	return &Identity{
		UserID: userid,
		SecurityStamp: securityStamp,
		Authenticated: true,
		MasterKey:     mk,
		Expiration:    expiration,
	}
}

func Authenticate(username string, password string) (int32, string, bool) {
	// time attack partial mitigation
	// adds up to 0.5 seconds to the response time

	// this technically does not prevent a time attack, since there is still time variance without the randomness added.
	// you could theoretically take an average of a 'valid user; incorrect password' vs 'invalid user' response times
	// to figure out if a user exists, but you would need a lot of data to do that.
	// this should make it *extremely* unlikely to do when paired with 'n login attempt per ip/minute/fingerprint'
	// since you would need way more than `n` login attempts to collect an accurate average

	// https://security.stackexchange.com/questions/96489/can-i-prevent-timing-attacks-with-random-delays/96493#96493
	// https://www.reddit.com/r/PHP/comments/kn6ezp/have_you_secured_your_signup_process_against_a/

	randomSeconds, _ := rand.Int(rand.Reader, big.NewInt(500))
	randomDuration := time.Duration(randomSeconds.Int64()) * time.Millisecond

	time.Sleep(randomDuration)

	user, userErr := users.FetchByUsername(username)

	if userErr != nil || user.FailedAttempts > int32(config.MAX_LOGIN_ATTEMPTS) {
		// set user password to dummy password to keep timing consistent when validating password
		user.Password = "$2a$14$KW5OO1wZqGGq3SrpBFj0Oema5DG8Ph7lZJvq0ECkkYBpNFom6b9vO"
		security.ComparePasswords(password, user.Password)
		return 0, "", false
	}

	result := security.ComparePasswords(password, user.Password)

	if !result {
		user.FailedAttempts += 1
		users.Update(user)
	} else {
		user.FailedAttempts = 0

		// seed data
		if user.Data == nil || len(user.Data) == 0 {
			secrets := []entries.Secret{}

			// THE FUNNY THING
			// for range 2500000 {
			// 	secrets = append(secrets, entries.Secret{
			// 		ID: uuid.New().String()
			// 		Description: "some website",
			// 		URL: "https://example.com",
			// 		Notes: "test notes something here i like writing notes lalalalala test test",
			// 		Username: "username2345",
			// 		Password: "laksjdflkjasdlfkj2934829384sldkfj",
			// 	})
			// }

			mk := security.SHA512_58(password)
			user.Data, _ = security.EncryptDataWithKey(&secrets, mk)
		}

		users.Update(user)
	}

	return user.ID, user.SecurityStamp, result
}

func CheckPasswordCriteria(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("Password cannot be blank.")
	}

	if len(password) < config.PASSWORD_MIN_LENGTH {
		return errors.New("Password must be at least " + strconv.Itoa(config.PASSWORD_MIN_LENGTH) + " characters long.")
	}

	uppercaseCount := 0
	lowercaseCount := 0
	numberCount := 0
	symbolCount := 0

	for _, r := range password {
		if unicode.IsUpper(r) {
			uppercaseCount += 1
		}

		if unicode.IsLower(r) {
			lowercaseCount += 1
		}

		if unicode.IsNumber(r) {
			numberCount += 1
		}

		if !unicode.IsNumber(r) && !unicode.IsLower(r) && !unicode.IsUpper(r) {
			symbolCount += 1
		}
	}

	if uppercaseCount < config.PASSWORD_REQUIRED_UPPERCASE {
		return errors.New("Password must contain at least " + strconv.Itoa(config.PASSWORD_REQUIRED_UPPERCASE) + " uppercase character(s).")
	}

	if lowercaseCount < config.PASSWORD_REQUIRED_LOWERCASE {
		return errors.New("Password must contain at least " + strconv.Itoa(config.PASSWORD_REQUIRED_LOWERCASE) + " lowercase character(s).")
	}

	if numberCount < config.PASSWORD_REQUIRED_NUMBERS {
		return errors.New("Password must contain at least " + strconv.Itoa(config.PASSWORD_REQUIRED_NUMBERS) + " number(s).")
	}

	if symbolCount < config.PASSWORD_REQUIRED_SYMBOLS {
		return errors.New("Password must contain at least " + strconv.Itoa(config.PASSWORD_REQUIRED_SYMBOLS) + " symbol(s).")
	}

	return nil
}

// func ChangePassword(user models.User, oldPassword string, newPassword string, passwordConfirm string, noCheck bool) (models.User, error) {
// 	// noCheck skips criteria validation, and confirmation validation
// 	if noCheck {
// 		if newPassword != passwordConfirm {
// 			log.Println("Passwords do not match")
// 			return user, errors.New("passwords do not match")
// 		}

// 		critErr := CheckPasswordCriteria(newPassword)

// 		if critErr != nil {
// 			return user, critErr
// 		}

// 		if !security.ComparePasswords(oldPassword, user.Password) {
// 			return user, errors.New("old password incorrect")
// 		}
// 	}

// 	newHash, hashErr := security.HashPassword(newPassword)
// 	if hashErr != nil {
// 		return user, errors.New("could not hash password")
// 	}

// 	user.Password = newHash

// 	updateErr := database.UpdateUser(user)

// 	if updateErr != nil {
// 		return user, errors.New("could not update user")
// 	}

// 	return user, nil
// }

// // wrapper with less args for skipping validation, confirmation
// func ChangePasswordNoCheck(user models.User, newPassword string) (models.User, error) {
// 	return ChangePassword(user, "", newPassword, "", true)
// }

// // hard reset user password without confirmation or record.
// // should only be used for developer purposes
// func ResetPasswordNoConfirm(userid int) (models.User, error) {
// 	user, err := database.FetchUserById(userid)
// 	if err != nil {
// 		return user, err
// 	}

// 	hash, hashErr := security.HashPassword(config.GetConfig().IdentityDefaultPassword)
// 	if hashErr != nil {
// 		return user, hashErr
// 	}

// 	user.Password = hash

// 	updateErr := database.UpdateUser(user)
// 	if updateErr != nil {
// 		return user, updateErr
// 	}

// 	return user, nil
// }
