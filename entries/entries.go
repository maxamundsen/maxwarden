package entries

import (
	"errors"
	"maxwarden/database"
	"maxwarden/security"
	"maxwarden/users"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID string
	Description string
	URL string
	Notes string
	Username string
	Password string
	Created time.Time
	Modified time.Time
}

type EntryFilter struct {
	Filter database.Filter
	MasterKey string
	UserId int32
}

func OrderByDescription(secret []Secret, desc bool) []Secret {
    sort.Slice(secret, func(i, j int) bool {
		if desc {
			return strings.ToLower(secret[i].Description) > strings.ToLower(secret[j].Description)
		} else {
			return strings.ToLower(secret[i].Description) < strings.ToLower(secret[j].Description)
		}
    })
    return secret
}

func Filter(f EntryFilter) ([]Secret, error) {
	user, _ := users.FetchById(f.UserId)

	// we need to do the rest in memory because the data is encrypted, so we need to decrypt the data
	secrets, decErr := security.DecryptDataWithKey[[]Secret](user.Data, f.MasterKey)
	if decErr != nil {
		return nil, decErr
	}

	if secrets == nil {
		return nil, errors.New("secrets list is null")
	}

	output := []Secret{}

	// decrypt data inside each entry retrieved
	for _, v := range *secrets {
		search := f.Filter.Search["description"]
		if search != "" {
			search = strings.ToLower(search)
			searchable := strings.ToLower(v.Description)

			if strings.Contains(searchable, search) {
				output = append(output, v)
			}
		} else {
			output = append(output, v)
		}
	}

	//order by
	output = OrderByDescription(output, f.Filter.OrderDescending)

	// pagination
	output = database.PaginateSlice(output, f.Filter)

	return output, nil
}

// linear search because that's the only option
func FetchSecretFromID(userId int32, masterKey string, secretId string) (Secret, error) {
	user, _ := users.FetchById(userId)

	secrets, decErr := security.DecryptDataWithKey[[]Secret](user.Data, masterKey)
	if decErr != nil {
		return Secret{}, decErr
	}

	if secrets == nil {
		return Secret{}, errors.New("secrets are nil")
	}

	for _, v := range *secrets {
		if v.ID == secretId {
			return v, nil
		}
	}

	return Secret{}, errors.New("no secret found")
}

func DeleteSecret(userId int32, masterKey string, secretId string) error {
	user, _ := users.FetchById(userId)

	secrets, decErr := security.DecryptDataWithKey[[]Secret](user.Data, masterKey)
	if decErr != nil {
		return decErr
	}

	if secrets == nil {
		return errors.New("attempt to range over null secret array")
	}

	output := []Secret{}

	for _, v := range *secrets {
		if v.ID != secretId {
			output = append(output, v)
		}
	}

	enc, _ := security.EncryptDataWithKey(&output, masterKey)
	user.Data = enc

	_, userErr := users.Update(user)

	return userErr
}

func Update(userId int32, masterKey string, secret Secret) error {
	user, _ := users.FetchById(userId)

	secrets, _ := security.DecryptDataWithKey[[]Secret](user.Data, masterKey)
	if secrets == nil {
		return errors.New("user secrets are null")
	}

	secret.Modified = time.Now()

	// linear search and replace
	for i, v := range *secrets {
		if v.ID == secret.ID {
			created := (*secrets)[i].Created
			secret.Created = created

			(*secrets)[i] = secret
		}
	}

	enc, _ := security.EncryptDataWithKey(secrets, masterKey)

	user.Data = enc
	_, updateErr := users.Update(user)

	return updateErr
}

func Add(userId int32, masterKey string, secret Secret) error {
	user, _ := users.FetchById(userId)

	secrets, _ := security.DecryptDataWithKey[[]Secret](user.Data, masterKey)
	if secrets == nil {
		return errors.New("user secrets are null")
	}

	secret.ID = uuid.New().String()
	secret.Modified = time.Now()
	secret.Created = time.Now()

	*secrets = append(*secrets, secret)

	enc, _ := security.EncryptDataWithKey(secrets, masterKey)

	user.Data = enc
	_, updateErr := users.Update(user)

	return updateErr
}