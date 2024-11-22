package server

import (
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"log"
	"tender-backend/gen_proto"
	"tender-backend/model"
	request_model "tender-backend/model/request"
	"tender-backend/internal/usecase/rabbit_mq"
	"tender-backend/internal/usecase/web_socket"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

func (s *NotificationService) CreateNotification(notification *request_model.CreateNotificationReq) (*model.Notification, error) {
	newNotification := model.Notification{
		UserID:      notification.UserID,
		Message:     notification.Message,
		IsDelivered: false,
		DeliveredAt: nil,
	}

	if err := s.db.Create(&newNotification).Error; err != nil {
		return nil, err
	}

	return &newNotification, nil
}

func (s *NotificationService) ConsumeNotifications() {
	messages, err := rabbit_mq.Consume("notifications")
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	for msg := range messages {
		var notification gen_proto.Notification
		if err := proto.Unmarshal(msg.Body, &notification); err != nil {
			log.Printf("Failed to unmarshal notification: %v", err)
			continue
		}

		// check if notification delivery is already handled
		notificationFromDb, err := s.getNotificationByID(notification.Id)
		if err != nil {
			log.Printf("Failed to get notification from DB: %v", err)
			continue
		}

		if notificationFromDb.IsDelivered {
			continue
		}

		err = s.db.Transaction(func(tx *gorm.DB) error {
			// Mark notification as delivered
			err := s.markNotificationAsDelivered(notification.Id)

			if err != nil {
				return err
			}

			// Check WebSocket connection
			return web_socket.SendNotification(notification.UserId, msg.Body)
		})

		if err != nil {
			log.Printf("Failed to handle notification: %v", err)
		}
	}
}

func (s *NotificationService) getNotificationByID(id int64) (*model.Notification, error) {
	var notification model.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationService) markNotificationAsDelivered(id int64) error {
	notification, err := s.getNotificationByID(id)
	if err != nil {
		return err
	}

	notification.IsDelivered = true
	if err := s.db.Save(notification).Error; err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) PublishNotDeliveredNotificationsForUser(userID int64) error {
	var notifications []model.Notification

	if err := s.db.Where("user_id = ? AND is_delivered = ?", userID, false).Find(&notifications).Error; err != nil {
		return err
	}

	for _, notification := range notifications {
		notificationProto := &gen_proto.Notification{
			Id:      notification.ID,
			UserId:  notification.UserID,
			Message: notification.Message,
		}

		notificationBytes, err := proto.Marshal(notificationProto)
		if err != nil {
			return err
		}

		if err := rabbit_mq.Publish("notifications", notificationBytes); err != nil {
			return err
		}
	}

	return nil
}
