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
	"github.com/rs/zerolog/log"
)

type kafka struct {
	SrvManager *service.Manager
	Cfg        *config.Config
	Log        *logger.Logger
}

func New(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) Queue {
	return kafka{
		SrvManager: srvManager,
		Cfg:        cfg,
		Log:        log,
	}
}

func (k kafka) SaveChannelsData() {
	consumer, err := connectAsConsumer(k.Cfg.KafkaAddr, "groups")
	if err != nil {
		k.Log.Error().Err(err).Msg("connect to queue as consumer")
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				k.Log.Error().Err(err).Msg("get data from queue")

				return
			case messages := <-consumer.Messages():
				var channel model.DBChannel

				err := json.Unmarshal(messages.Value, &channel)
				if err != nil {
					k.Log.Error().Err(err).Msg("unmarshal channel data")

					continue
				}

				err = k.SrvManager.Channel.CreateChannel(&channel)
				if err != nil {
					if errors.Is(err, service.ErrChannelExists) {
						k.Log.Warn().Err(err).Msgf("channel with name %s already exists", channel.Name)

						continue
					}

					k.Log.Error().Err(err).Msg("create channel")
				}
			}
		}
	}()
}

func (k kafka) SaveMessagesData() {
	consumer, err := connectAsConsumer(k.Cfg.KafkaAddr, "messages")
	if err != nil {
		k.Log.Error().Err(err).Msg("connect to queue as consumer")
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				k.Log.Error().Err(err).Msg("get data from queue")

				return
			case data := <-consumer.Messages():
				var telegramMessage model.TgMessage

				err = json.Unmarshal(data.Value, &telegramMessage)
				if err != nil {
					k.Log.Error().Err(err).Msg("unmarshal message data")

					continue
				}

				channel, err := k.SrvManager.Channel.GetChannelByName(telegramMessage.PeerID.Username)
				if err != nil {
					if errors.Is(err, service.ErrChannelNotFound) {
						k.Log.Warn().Err(err)

						continue
					}

					log.Error().Err(err).Msg("get channel by name")

					continue
				}

				userID, err := k.SrvManager.User.CreateUser(&model.User{
					Username: telegramMessage.FromID.Username,
					FullName: telegramMessage.FromID.Fullname,
					ImageURL: telegramMessage.FromID.ImageURL,
				})
				if err != nil {
					log.Error().Err(err).Msg("create user")

					continue
				}

				messageID, err := k.SrvManager.Message.CreateMessage(&model.DBMessage{
					ChannelID:  channel.ID,
					UserID:     userID,
					Title:      telegramMessage.Message,
					MessageURL: telegramMessage.MessageURL,
					ImageURL:   telegramMessage.ImageURL,
				})
				if err != nil {
					log.Error().Err(err).Msg("create message")

					continue
				}

				k.processReplyData(messageID, &telegramMessage)
			}
		}
	}()
}

func (k kafka) processReplyData(messageID int, telegramMessage *model.TgMessage) {
	for _, reply := range telegramMessage.Replies.Messages {
		userID, err := k.SrvManager.User.CreateUser(&model.User{
			Username: reply.FromID.Username,
			FullName: reply.FromID.Fullname,
			ImageURL: reply.FromID.ImageURL,
		})
		if err != nil {
			k.Log.Error().Err(err).Msg("create user for reply")

			continue
		}

		err = k.SrvManager.Reply.CreateReply(&model.DBReply{
			MessageID: messageID,
			UserID:    userID,
			Title:     reply.Message,
			ImageURL:  reply.ImageURL,
		})
		if err != nil {
			k.Log.Error().Err(err).Msg("create reply")
		}
	}
}

func connectAsConsumer(addr string, topic string) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	worker, err := sarama.NewConsumer([]string{addr}, config)
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, fmt.Errorf("create partition consumer: %w", err)
	}

	return consumer, nil
}
