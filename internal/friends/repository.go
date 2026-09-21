package friends

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type Repository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) GetUserByUsernameOrEmail(usernameOrEmail string, ctx fiber.Ctx) (*User, error) {
	r.logger.Debug(
		"querying user",
		"identifier", usernameOrEmail,
	)

	var user User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, display_name
		 FROM users
		 WHERE username = ? OR email = ?`,
		usernameOrEmail,
		usernameOrEmail,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
	)

	if err != nil {
		r.logger.Warn(
			"user query failed",
			"identifier", usernameOrEmail,
			"error", err,
		)

		return nil, err
	}

	r.logger.Debug(
		"user found",
		"user_id", user.ID,
		"username", user.Username,
	)

	return &user, nil
}

func (r *Repository) CheckFriendRequest(sender uint64, receiver uint64, ctx fiber.Ctx) (uint64, error) {
	var requestID uint64

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id
			 FROM friend_requests
			 WHERE sender_id = ?
			   AND receiver_id = ?
			   AND status = 'pending'
			 LIMIT 1`,
		receiver,
		sender,
	).Scan(&requestID)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug(
				"no reverse friend request found",
				"sender_id", sender,
				"receiver_id", receiver,
			)
		} else {
			r.logger.Warn(
				"failed to check reverse friend request",
				"sender_id", sender,
				"receiver_id", receiver,
				"error", err,
			)
		}

		return 0, err
	}

	r.logger.Debug(
		"reverse friend request found",
		"request_id", requestID,
		"sender_id", sender,
		"receiver_id", receiver,
	)

	return requestID, nil
}

func (r *Repository) AddFriend(sender uint64, receiver uint64, ctx fiber.Ctx) error {
	var count int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
		 FROM friend_requests
		 WHERE sender_id = ?
		   AND receiver_id = ?`,
		sender,
		receiver,
	).Scan(&count)

	if err != nil {
		r.logger.Warn(
			"failed to check existing friend request",
			"sender_id", sender,
			"receiver_id", receiver,
			"error", err,
		)

		return err
	}

	if count > 0 {
		r.logger.Debug(
			"friend request already exists",
			"sender_id", sender,
			"receiver_id", receiver,
		)

		return nil
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO friend_requests (sender_id, receiver_id)
		 VALUES (?, ?)`,
		sender,
		receiver,
	)

	if err != nil {
		r.logger.Warn(
			"failed to insert friend request",
			"sender_id", sender,
			"receiver_id", receiver,
			"error", err,
		)

		return err
	}

	r.logger.Debug(
		"friend request inserted",
		"sender_id", sender,
		"receiver_id", receiver,
	)

	return nil
}

func (r *Repository) AcceptFriendRequest(requestID uint64, userID uint64, ctx fiber.Ctx) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var senderID, receiverID uint64

	err = tx.QueryRowContext(
		ctx,
		`SELECT sender_id, receiver_id
		 FROM friend_requests
		 WHERE id = ?
		   AND receiver_id = ?
		   AND status = 'pending'
		 FOR UPDATE`,
		requestID,
		userID,
	).Scan(&senderID, &receiverID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn(
				"friend request not found or already processed",
				"request_id", requestID,
				"user_id", userID,
			)

			return fiber.NewError(
				fiber.StatusNotFound,
				"friend request not found",
			)
		}

		r.logger.Warn(
			"failed to get friend request",
			"request_id", requestID,
			"error", err,
		)

		return err
	}

	result, err := tx.ExecContext(
		ctx,
		`UPDATE friend_requests
		 SET status = 'accepted'
		 WHERE id = ?
		   AND receiver_id = ?
		   AND status = 'pending'`,
		requestID,
		userID,
	)
	if err != nil {
		r.logger.Warn(
			"failed to accept friend request",
			"request_id", requestID,
			"error", err,
		)

		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fiber.NewError(
			fiber.StatusNotFound,
			"friend request not found",
		)
	}

	userOneID := senderID
	userTwoID := receiverID

	if userOneID > userTwoID {
		userOneID, userTwoID = userTwoID, userOneID
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO conversations (
			user_one_id,
			user_two_id
		) VALUES (?, ?)`,
		userOneID,
		userTwoID,
	)
	if err != nil {
		r.logger.Warn(
			"failed to create conversation",
			"request_id", requestID,
			"user_one_id", userOneID,
			"user_two_id", userTwoID,
			"error", err,
		)

		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Warn(
			"failed to commit friend request transaction",
			"request_id", requestID,
			"error", err,
		)

		return err
	}

	r.logger.Debug(
		"friend request accepted and conversation created",
		"request_id", requestID,
		"user_one_id", userOneID,
		"user_two_id", userTwoID,
	)

	return nil
}

func (r *Repository) GetFriends(userID uint64, ctx fiber.Ctx) (FriendsResponse, error) {
	r.logger.Debug("Hit Get Friends Handler")
	var friends FriendsResponse

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			c.id,
			u.id,
			u.username,
			u.email,
			u.display_name
		FROM conversations c
		JOIN users u
			ON u.id = CASE
				WHEN c.user_one_id = ? THEN c.user_two_id
				ELSE c.user_one_id
			END
		WHERE c.user_one_id = ?
		   OR c.user_two_id = ?`,
		userID,
		userID,
		userID,
	)
	if err != nil {
		return friends, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend FriendDto

		err := rows.Scan(
			&friend.ConversationID,
			&friend.ID,
			&friend.Username,
			&friend.Email,
			&friend.DisplayName,
		)
		if err != nil {
			return friends, err
		}

		friends.Friends = append(friends.Friends, friend)
	}

	if err := rows.Err(); err != nil {
		return friends, err
	}

	return friends, nil
}

func (r *Repository) GetRequests(userID uint64, ctx fiber.Ctx) ([]FriendRequestDto, error) {
	r.logger.Debug("Hit Get Requests Handler")

	var requests []FriendRequestDto

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			id,
			sender_id,
			receiver_id,
			status
		FROM friend_requests
		WHERE receiver_id = ? OR sender_id = ?`,
		userID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var request FriendRequestDto

		err := rows.Scan(
			&request.ID,
			&request.SenderID,
			&request.ReceiverID,
			&request.Status,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) RejectFriendRequest(requestID uint64, userID uint64, ctx fiber.Ctx) error {
	r.logger.Debug("Hit Reject Friend Request Handler")

	_, err := r.db.ExecContext(
		ctx,
		`UPDATE friend_requests
		SET status = 'rejected'
		WHERE id = ? AND receiver_id = ?`,
		requestID,
		userID,
	)

	return err
}
