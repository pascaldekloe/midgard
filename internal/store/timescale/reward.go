package timescale

import (
	"fmt"

	"github.com/pkg/errors"

	"gitlab.com/thorchain/midgard/internal/models"
)

func (s *Client) CreateRewardRecord(record models.EventReward) error {
	err := s.CreateEventRecord(record.Event)
	if err != nil {
		return errors.Wrap(err, "Failed to create event record")
	}

	query := fmt.Sprintf(`
		INSERT INTO %v (
			time,
			event_id,
			pool,
			units,
      from_address
		)  VALUES ( $1, $2, $3, $4, $5 ) RETURNING event_id`, models.ModelStakesTable)

	for _, reward := range record.PoolRewards {
		_, err := s.db.Exec(query,
			record.Event.Time,
			record.Event.ID,
			reward.Pool.String(),
			reward.Amount,
			"BLOCK_REWARD",
		)

		if err != nil {
			s.logger.Error().Err(err).Msg("failed to prepareNamed query for EventRecord")
		}
	}
	return nil
}
