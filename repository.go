package main

type UserList []User

func (list UserList) FindBy(userName string) *User {
	for _, user := range list {
		if user.UserName == userName {
			return &user
		}
	}
	return nil
}

type Repository struct {
	readItems []HeybeItem
	userList  UserList
}

func (repo *Repository) Init() {
	item1 := HeybeItem{
		Name:  "Effective Go",
		Link:  "http://golang.org/doc/effective_go.html",
		Descr: "GO lang details",
		Tags:  []string{"Go", "Books"}}

	item2 := HeybeItem{
		Name:  "Go: Best Practices for Production Environments",
		Link:  "http://peter.bourgon.org/go-in-production/",
		Descr: "Go best practices",
		Tags:  []string{"Go", "Best Practices"}}
	item3 := HeybeItem{
		Name:  "Java Memory Model Pragmatics",
		Link:  "http://shipilev.net/blog/2014/jmm-pragmatics/",
		Descr: "Java memory model",
		Tags:  []string{"Java", "Performance"}}

	repo.readItems = []HeybeItem{item1, item2, item3}

	user1 := User{
		UserName: "turgay",
		Password: "heybe",
		Email:    "tk@heybe.com"}
	repo.userList = []User{user1}
}

func (repo Repository) LoadItems() ([]HeybeItem, error) {
	return repo.readItems, nil
}

func (repo *Repository) AddUser(user *User) {
	repo.userList = append(repo.userList, *user)
}

func (repo *Repository) AddItem(item HeybeItem) {
	repo.readItems = append(repo.readItems, item)
}

func (repo *Repository) FindUser(userName string) *User {
	return repo.userList.FindBy(userName)
}
