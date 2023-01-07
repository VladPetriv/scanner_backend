package kafka

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

func SaveChannelsFromQueueToDB(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) {
	worker, err := connectAsConsumer(cfg.KafkaAddr)
	if err != nil {
		log.Error().Err(err).Msg("create queue consumer")
	}

	consumer, err := worker.ConsumePartition("groups", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error().Err(err).Msg("consume partition")
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error().Err(err).Msg("get data from queuer")

				return

			case messages := <-consumer.Messages():
				channel := model.DBChannel{}

				err = json.Unmarshal(messages.Value, &channel)
				if err != nil {
					log.Error().Err(err).Msg("unmarshal channel data")
				}

				candidate, err := srvManager.Channel.GetChannelByName(channel.Name)
				if err != nil && !errors.Is(err, pg.ErrChannelNotFound) {
					log.Error().Err(err).Msg("get channel by name")
				}

				if candidate == nil {
					err = srvManager.Channel.CreateChannel(&channel)
					if err != nil {
						log.Error().Err(err).Msg("create channel")
					}
				} else {
					log.Info().Msgf("channel with name %s is exist", channel.Name)
				}
			}
		}
	}()
}

func SaveDataFromQueueToDB(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) {
	worker, err := connectAsConsumer(cfg.KafkaAddr)
	if err != nil {
		log.Error().Err(err).Msg("create consumer")
	}

	consumer, err := worker.ConsumePartition("messages", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error().Err(err).Msg("consume partition")
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error().Err(err).Msg("get data from queue")

				return
			case messages := <-consumer.Messages():
				telegramMessage := model.TgMessage{}

				err = json.Unmarshal(messages.Value, &telegramMessage)
				if err != nil {
					log.Error().Err(err).Msg("unmarshal message")
				}

				channel, err := srvManager.Channel.GetChannelByName(telegramMessage.PeerID.Username)
				if err != nil {
					log.Error().Err(err).Msg("get channel by name")

					return
				}

				userID, err := srvManager.User.CreateUser(&model.User{
					Username: telegramMessage.FromID.Username,
					FullName: telegramMessage.FromID.Fullname,
					ImageURL: telegramMessage.FromID.ImageURL,
				})
				if err != nil {
					log.Error().Err(err).Msg("create user")
				}

				messageID, err := srvManager.Message.CreateMessage(&model.DBMessage{
					ChannelID:  channel.ID,
					UserID:     userID,
					Title:      telegramMessage.Message,
					MessageURL: telegramMessage.MessageURL,
					ImageURL:   telegramMessage.ImageURL,
				})
				if err != nil {
					log.Error().Err(err).Msg("create message")
				}

				for _, replie := range telegramMessage.Replies.Messages {
					userID, err = srvManager.User.CreateUser(&model.User{
						Username: replie.FromID.Username,
						FullName: replie.FromID.Fullname,
						ImageURL: replie.FromID.ImageURL,
					})
					if err != nil {
						log.Error().Err(err).Msg("create user for reply")
					}

					err = srvManager.Reply.CreateReply(&model.DBReply{
						MessageID: messageID,
						UserID:    userID,
						Title:     replie.Message,
						ImageURL:  replie.ImageURL,
					})
					if err != nil {
						log.Error().Err(err).Msg("create reply")
					}
				}
			}
		}
	}()
}

func connectAsConsumer(addr string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer([]string{addr}, config)
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	return conn, nil
}
