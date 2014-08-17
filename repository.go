package main

//TODO create type for users slice

type Repository struct {
	readItems []ReadItem
	users     []User
}

func (repo *Repository) Init() {
	item1 := ReadItem{
		Name:  "Effective Go",
		Link:  "http://golang.org/doc/effective_go.html",
		Descr: "GO lang details",
		Tag:   "Go"}

	item2 := ReadItem{
		Name:  "Go: Best Practices for Production Environments",
		Link:  "http://peter.bourgon.org/go-in-production/",
		Descr: "Go best practices",
		Tag:   "Go"}
	repo.readItems = []ReadItem{item1, item2}

	user1 := User{
		UserName: "turgay",
		Password: "heybe",
		Email:    "tk@heybe.com"}
	repo.users = []User{user1}
}

func (repo Repository) LoadItems() ([]ReadItem, error) {
	return repo.readItems, nil
}

func (repo *Repository) AddUser(user User) {
	repo.users = append(repo.users, user)
}

func (repo *Repository) AddItem(item ReadItem) {
	repo.readItems = append(repo.readItems, item)
}

func (repo *Repository) FindUser(userName string) *User {
	for _, user := range repo.users {
		if user.UserName == userName {
			return &user
		}
	}
	return nil
}
