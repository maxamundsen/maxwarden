package users

import "encoding/json"

type Permissions struct {
	Admin bool

	Orders struct {
		View bool
		Edit bool
	}

	Users struct {
		View bool
		Edit bool
	}
}

func GetPermissionsFromJSON(jsonInput string) Permissions {
	var perms Permissions

	json.Unmarshal([]byte(jsonInput), &perms)

	return perms
}

func GetPermissionsJSON(p Permissions) string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}

func SetUserPermissions(p Permissions) error {
	// implement
	return nil
}
