type User struct {
	Username string
	Password string
}

type Item struct {
	Name   string
	Price  float64
	Rating float64
}

type Rating struct {
	UserID int
	ItemID int
	Score  float64
}

func (u *User) Register(username, password string) {
	u.Username = username
	u.Password = password
}

func (u *User) Authorize(username, password string) bool {
	return u.Username == username && u.Password == password
}

func (i *Item) Search(name string) []Item {
	var items []Item
	// implementation to search items based on name and add to items slice
	return items
}

func (i *Item) Filter(price, rating float64) []Item {
	var items []Item
	// implementation to filter items based on price and rating and add to items slice
	return items
}

func (r *Rating) GiveRating(userID, itemID int, score float64) {
	r.UserID = userID
	r.ItemID = itemID
	r.Score = score
}