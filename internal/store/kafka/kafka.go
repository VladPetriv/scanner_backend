package kafka

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

func connectAsConsumer(addr string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer([]string{addr}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return conn, nil
}

func SaveChannelsFromQueueToDB(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) {
	worker, err := connectAsConsumer(cfg.KafkaAddr)
	if err != nil {
		log.Error(err)
	}

	consumer, err := worker.ConsumePartition("channels.get", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error(err)
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error(err)

				return

			case messages := <-consumer.Messages():
				channel := model.DBChannel{}

				json.Unmarshal(messages.Value, &channel)

				candidate, err := srvManager.Channel.GetChannelByName(channel.Name)
				if !errors.Is(err, service.ErrChannelNotFound) && err != nil {
					log.Error(err)
				}

				if candidate == nil {
					err := srvManager.Channel.CreateChannel(&channel)
					if err != nil {
						log.Error(err)
					}
				} else {
					log.Infof("channel with name %s is exist", channel.Name)
				}
			}
		}
	}()
}

func SaveDataFromQueueToDB(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) {
	worker, err := connectAsConsumer(cfg.KafkaAddr)
	if err != nil {
		log.Error(err)
	}

	consumer, err := worker.ConsumePartition("messages.get", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error(err)
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error(err)

				return
			case messages := <-consumer.Messages():
				telegramMessage := model.TgMessage{}

				json.Unmarshal(messages.Value, &telegramMessage)

				channel, err := srvManager.Channel.GetChannelByName(telegramMessage.PeerID.Username)
				if err != nil {
					log.Error(err)
				}

				userID, err := srvManager.User.CreateUser(&model.User{
					Username: telegramMessage.FromID.Username,
					FullName: telegramMessage.FromID.Fullname,
					ImageURL: telegramMessage.FromID.ImageURL,
				})
				if err != nil {
					log.Error(err)
				}

				messageID, err := srvManager.Message.CreateMessage(&model.DBMessage{
					ChannelID:  channel.ID,
					UserID:     userID,
					Title:      telegramMessage.Message,
					MessageURL: telegramMessage.MessageURL,
					ImageURL:   telegramMessage.ImageURL,
				})
				if err != nil {
					log.Error(err)
				}

				for _, replie := range telegramMessage.Replies.Messages {
					userID, err := srvManager.User.CreateUser(&model.User{
						Username: replie.FromID.Username,
						FullName: replie.FromID.Fullname,
						ImageURL: replie.FromID.ImageURL,
					})
					if err != nil {
						log.Error(err)
					}

					err = srvManager.Replie.CreateReplie(&model.DBReplie{
						MessageID: messageID,
						UserID:    userID,
						Title:     replie.Message,
						ImageURL:  replie.ImageURL,
					})
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}()
}
