package autobuy

type user struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	NameKana string `toml:"name_kana"`
	Zipcode1 string `toml:"zipcode1"`
	Zipcode2 string `toml:"zipcode2"`
	Pref     string `toml:"pref"`
	City     string `toml:"city"`
	Street   string `toml:"street"`
	Building string `toml:"building"`
	Phone    string `toml:"phone"`
}
