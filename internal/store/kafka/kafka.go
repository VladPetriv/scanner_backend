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
