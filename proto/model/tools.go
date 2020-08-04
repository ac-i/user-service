package model

import (
	"encoding/json"
)

func (x *User) DoMarshal() ([]byte, error) {
	return json.Marshal(&x)
}

func (x *Users) DoMarshal() ([]byte, error) {
	return json.Marshal(&x)
}

func (x *UserSelect) DoMarshal() ([]byte, error) {
	return json.Marshal(&x)
}

func DoStringUserErr(x *User, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	jb, err := x.DoMarshal()
	if err != nil {
		return "error: " + err.Error()
	}
	return string(jb)
}
func DoStringUsersErr(x *Users, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	jb, err := x.DoMarshal()
	if err != nil {
		return "error: " + err.Error()
	}
	return string(jb)
}
func DoStringUserSelectErr(x *UserSelect, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	jb, err := x.DoMarshal()
	if err != nil {
		return "error: " + err.Error()
	}
	return string(jb)
}

// Validate
// if in.Name == "" {
// 	return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument: user name is empty")
// }
