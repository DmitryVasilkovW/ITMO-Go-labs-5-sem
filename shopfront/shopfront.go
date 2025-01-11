//go:build !solution

package shopfront

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type ShopCounters struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) Counters {
	return ShopCounters{redisClient: redisClient}
}

func (c ShopCounters) RecordView(ctx context.Context, itemID ItemID, userID UserID) error {
	itemKey := c.generateItemKey(itemID)
	viewCountKey := c.generateViewCountKey(itemID)
	userIdentifier := int64ToString(int64(userID))

	pipe := c.redisClient.TxPipeline()

	c.queueIncrementViewCount(ctx, pipe, viewCountKey)
	c.queueAddUserToItemViews(ctx, pipe, itemKey, userIdentifier)

	return c.execute(ctx, pipe)
}

func (c ShopCounters) queueIncrementViewCount(ctx context.Context, pipe redis.Pipeliner, viewCountKey string) {
	pipe.Incr(ctx, viewCountKey)
}

func (c ShopCounters) queueAddUserToItemViews(ctx context.Context, pipe redis.Pipeliner, itemKey string, userIdentifier string) {
	pipe.SAdd(ctx, itemKey, userIdentifier)
}

func (c ShopCounters) execute(ctx context.Context, pipe redis.Pipeliner) error {
	_, err := pipe.Exec(ctx)
	return err
}

func (c ShopCounters) GetItems(ctx context.Context, itemIDs []ItemID, userID UserID) ([]Item, error) {
	itemKeys := c.generateItemKeys(itemIDs)
	viewCountKeys := c.generateViewCountKeys(itemIDs)
	userIdentifier := int64ToString(int64(userID))

	pipe := c.redisClient.Pipeline()

	viewCountCommands := c.queueViewCountCommands(ctx, pipe, viewCountKeys)
	viewedCommands := c.queueViewedCommands(ctx, pipe, itemKeys, userIdentifier)

	if err := c.executePipeline(ctx, pipe); err != nil {
		return nil, err
	}

	return c.buildItemsFromRedisResults(viewCountCommands, viewedCommands)
}

func (c ShopCounters) queueViewCountCommands(ctx context.Context, pipe redis.Pipeliner, viewCountKeys []string) []*redis.StringCmd {
	viewCountCommands := make([]*redis.StringCmd, len(viewCountKeys))
	for i, key := range viewCountKeys {
		viewCountCommands[i] = pipe.Get(ctx, key)
	}
	return viewCountCommands
}

func (c ShopCounters) queueViewedCommands(ctx context.Context, pipe redis.Pipeliner, itemKeys []string, userIdentifier string) []*redis.BoolCmd {
	viewedCommands := make([]*redis.BoolCmd, len(itemKeys))
	for i, key := range itemKeys {
		viewedCommands[i] = pipe.SIsMember(ctx, key, userIdentifier)
	}
	return viewedCommands
}

func (c ShopCounters) executePipeline(ctx context.Context, pipe redis.Pipeliner) error {
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (c ShopCounters) buildItemsFromRedisResults(viewCountCommands []*redis.StringCmd, viewedCommands []*redis.BoolCmd) ([]Item, error) {
	items := make([]Item, len(viewCountCommands))

	for i := range items {
		viewCount, err := c.getViewCount(viewCountCommands[i])
		if err != nil {
			return nil, err
		}
		isViewed := c.getViewedStatus(viewedCommands[i])
		items[i] = c.createItem(viewCount, isViewed)
	}
	return items, nil
}

func (c ShopCounters) getViewCount(viewCountCmd *redis.StringCmd) (int, error) {
	viewCount, err := viewCountCmd.Int()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return viewCount, nil
}

func (c ShopCounters) getViewedStatus(viewedCmd *redis.BoolCmd) bool {
	return viewedCmd.Val()
}

func (c ShopCounters) createItem(viewCount int, isViewed bool) Item {
	return Item{
		ViewCount: viewCount,
		Viewed:    isViewed,
	}
}

func (c ShopCounters) generateItemKey(itemID ItemID) string {
	return "item_" + int64ToString(int64(itemID))
}

func (c ShopCounters) generateItemKeys(itemIDs []ItemID) []string {
	keys := make([]string, len(itemIDs))
	for i, itemID := range itemIDs {
		keys[i] = c.generateItemKey(itemID)
	}
	return keys
}

func (c ShopCounters) generateViewCountKey(itemID ItemID) string {
	return "item_" + int64ToString(int64(itemID)) + "_count"
}

func (c ShopCounters) generateViewCountKeys(itemIDs []ItemID) []string {
	keys := make([]string, len(itemIDs))
	for i, itemID := range itemIDs {
		keys[i] = c.generateViewCountKey(itemID)
	}
	return keys
}

func int64ToString(number int64) string {
	return strconv.FormatInt(number, 10)
}
