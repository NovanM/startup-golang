package user

type UserFormater struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageUrl   string `json:"image_url"`
}

func FormaterUser(user User, token string) UserFormater {
	formater := UserFormater{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
		ImageUrl:   user.AvatarFileName,
	}
	return formater
}
