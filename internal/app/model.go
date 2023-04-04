package app

import (
	"fmt"
	"github.com/johnnyzhao/retail-ai-api/internal/domain"
	"strings"
)

type GetUserResponse struct {
	Message string      `json:"message"`
	User    domain.User `json:"user"`
}

type ItemResponse struct {
	Message string      `json:"message"`
	Recipe  domain.User `json:"recipe"`
}

type CreateUserResponse struct {
	Message string      `json:"message"`
	User    domain.User `json:"user"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

type CreateUserPayload struct {
	UserID   string `json:"user_id,omitempty"`
	Password string `json:"password,omitempty"`
}

//func (p *CreateUserPayload) ToUpdateValues() map[string]interface{} {
//	result := make(map[string]interface{})
//	if p.Title != nil {
//		result["title"] = *p.Title
//	}
//	if p.MakingTime != nil {
//		result["making_time"] = *p.MakingTime
//	}
//	if p.Serves != nil {
//		result["serves"] = *p.Serves
//	}
//	if p.Ingredients != nil {
//		result["ingredients"] = *p.Ingredients
//	}
//	if p.Cost != nil {
//		result["cost"] = *p.Cost
//	}
//	return result
//}

func (p *CreateUserPayload) ValidateRequired() error {
	result := make([]string, 0)
	if len(p.UserID) == 0 {
		result = append(result, "user_id")
	}
	if len(p.Password) == 0 {
		result = append(result, "password")
	}

	if len(result) > 0 {
		return fmt.Errorf("required %s", strings.Join(result, " and "))
	}
	if len(p.UserID) < 6 || len(p.UserID) > 20 {
		return fmt.Errorf("user_id should be 6 to 20 characters")
	}
	if len(p.Password) < 8 || len(p.Password) > 20 {
		return fmt.Errorf("password should be 8 to 20 characters")
	}

	if !isAlphanumeric(p.UserID) {
		return fmt.Errorf("invalid user_id")
	}

	if !isASCIIString(p.Password) {
		return fmt.Errorf("invalid password")
	}

	return nil
}

type PatchUserPayload struct {
	Nickname *string `json:"nickname,omitempty"`
	Comment  *string `json:"comment,omitempty"`
}

func (p *PatchUserPayload) ValidateRequired() error {
	if p.Nickname != nil && len(*p.Nickname) == 0 && p.Comment != nil && len(*p.Comment) == 0 {
		return fmt.Errorf("required nickname or comment")
	}

	if p.Nickname != nil && len(*p.Nickname) >= 30 {
		return fmt.Errorf("invalid nickname")
	}
	if p.Comment != nil && len(*p.Comment) >= 100 {
		return fmt.Errorf("invalid comment")
	}

	if p.Nickname != nil && !isValidChar(*p.Nickname) {
		return fmt.Errorf("invalid nickname")
	}
	if p.Comment != nil && !isValidChar(*p.Comment) {
		return fmt.Errorf("invalid comment")
	}

	return nil
}
