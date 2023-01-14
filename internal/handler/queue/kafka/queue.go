package kafka

type Queue interface {
	SaveChannelsData()
	SaveMessagesData()
}
