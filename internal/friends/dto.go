package friends

type UserProfile struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type AddFriendRequest struct {
	Name string `json:"name"`
}

type User struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type FriendRequest struct {
	ID         uint64 `json:"id"`
	SenderID   uint64 `json:"sender_id"`
	ReceiverID uint64 `json:"receiver_id"`
	Status     string `json:"status"`
}

type FriendDto struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type FriendsResponse struct {
	Friends []FriendDto `json:"friends"`
}
