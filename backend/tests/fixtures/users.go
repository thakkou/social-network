package fixtures

// UserInfo holds key profile details from the seeder data.
// Tests can reference these constants instead of hardcoding IDs.
type UserInfo struct {
	ID        int
	Firstname string
	Lastname  string
	Email     string
	Nickname  string
	IsPrivate bool
}

// Named user constants matching the seeder_test.sql data.
var (
	Alice = UserInfo{ID: 1, Firstname: "Alice", Lastname: "Martin", Email: "alice@example.com", Nickname: "ali_m", IsPrivate: false}
	Bob   = UserInfo{ID: 2, Firstname: "Bob", Lastname: "Nguyen", Email: "bob@example.com", Nickname: "", IsPrivate: false}
	Chloe = UserInfo{ID: 3, Firstname: "Chloe", Lastname: "Dubois", Email: "chloe@example.com", Nickname: "Chloe", IsPrivate: true}
	David = UserInfo{ID: 4, Firstname: "David", Lastname: "Smith", Email: "david@example.com", Nickname: "david", IsPrivate: false}
	Emma  = UserInfo{ID: 5, Firstname: "Emma", Lastname: "Wilson", Email: "emma@example.com", Nickname: "", IsPrivate: true}
	Farid = UserInfo{ID: 6, Firstname: "Farid", Lastname: "El Amrani", Email: "farid@example.com", Nickname: "", IsPrivate: false}
	Grace = UserInfo{ID: 7, Firstname: "Grace", Lastname: "Lee", Email: "grace@example.com", Nickname: "", IsPrivate: false}
	Hugo  = UserInfo{ID: 8, Firstname: "Hugo", Lastname: "Costa", Email: "hugo@example.com", Nickname: "costa77", IsPrivate: false}
	Bella = UserInfo{ID: 9, Firstname: "Isabella", Lastname: "Rossi", Email: "isabella@example.com", Nickname: "bella", IsPrivate: false}
	Jack  = UserInfo{ID: 10, Firstname: "Jack", Lastname: "Thompson", Email: "jack@example.com", Nickname: "", IsPrivate: false}
)

// AllUsers returns all seeded users for use in table-driven tests.
func AllUsers() []UserInfo {
	return []UserInfo{Alice, Bob, Chloe, David, Emma, Farid, Grace, Hugo, Bella, Jack}
}

// PublicUsers returns only users with public profiles.
func PublicUsers() []UserInfo {
	return []UserInfo{Alice, Bob, David, Farid, Grace, Hugo, Bella, Jack}
}

// PrivateUsers returns only users with private profiles.
func PrivateUsers() []UserInfo {
	return []UserInfo{Chloe, Emma}
}

// UserByID returns the UserInfo for the given ID, or nil if not found.
func UserByID(id int) *UserInfo {
	for _, u := range AllUsers() {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

// CommonPassword is the password hash used for all seeded users.
const CommonPassword = "password123"
