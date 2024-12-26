package request

type User struct {
    ID              int64  `json:"id"`
    FirstName       string `json:"first_name"`
    LastName        string `json:"last_name"`
    Username        string `json:"username"`
    LanguageCode    string `json:"language_code"`
    AllowsWriteToPM bool   `json:"allows_write_to_pm"`
    PhotoURL        string `json:"photo_url"`
}
