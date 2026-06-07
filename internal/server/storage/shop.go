package storage

import (
	"context"
	"fmt"
)

// ShopItemDef stores an item definition.
type ShopItemDef struct {
	ID          string
	Name        string
	Description string
	Category    string // card_back, avatar, title
	Price       int
}

// AllShopItems returns the full shop catalog.
func AllShopItems() []ShopItemDef {
	return []ShopItemDef{
		{ID: "cardback_classic", Name: "经典卡背", Description: "经典黄色卡背", Category: "card_back", Price: 0},
		{ID: "cardback_sanguo", Name: "三国卡背", Description: "三国风格卡背", Category: "card_back", Price: 500},
		{ID: "cardback_dark", Name: "暗黑卡背", Description: "黑金卡背", Category: "card_back", Price: 1000},
		{ID: "cardback_golden", Name: "金光卡背", Description: "闪耀金光卡背", Category: "card_back", Price: 3000},
		{ID: "avatar_frame_silver", Name: "银色头像框", Description: "简单银色头像框", Category: "avatar", Price: 200},
		{ID: "avatar_frame_gold", Name: "金色头像框", Description: "华丽金色头像框", Category: "avatar", Price: 800},
		{ID: "avatar_frame_dragon", Name: "龙纹头像框", Description: "龙纹镶金头像框", Category: "avatar", Price: 2000},
		{ID: "title_warrior", Name: "战士称号", Description: "“战士”称号", Category: "title", Price: 0},
		{ID: "title_hero", Name: "英雄称号", Description: "“英雄”称号", Category: "title", Price: 1500},
		{ID: "title_legend", Name: "传奇称号", Description: "“传奇”称号", Category: "title", Price: 5000},
	}
}

// GetShopItemByID finds a shop item by ID.
func GetShopItemByID(id string) *ShopItemDef {
	for _, item := range AllShopItems() {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

// PurchaseResult holds the result of a purchase attempt.
type PurchaseResult struct {
	Success bool
	ItemID  string
	Coins   int
	Error   string
}

// PurchaseItem attempts to purchase an item for a player.
func (lm *LeaderboardManager) PurchaseItem(ctx context.Context, playerID, playerName, itemID string) (*PurchaseResult, error) {
	item := GetShopItemByID(itemID)
	if item == nil {
		return &PurchaseResult{Success: false, Error: "物品不存在"}, nil
	}

	stats, err := lm.getOrCreateStats(ctx, playerID, playerName)
	if err != nil {
		return nil, err
	}

	// Check if already owned
	for _, owned := range stats.Inventory {
		if owned == itemID {
			return &PurchaseResult{Success: false, Error: "已拥有此物品"}, nil
		}
	}

	// Check coins
	if stats.Coins < item.Price {
		return &PurchaseResult{Success: false, Error: fmt.Sprintf("元子不足，需要 %d", item.Price)}, nil
	}

	// Deduct coins and add to inventory
	stats.Coins -= item.Price
	stats.LastCoinChange = -item.Price
	stats.Inventory = append(stats.Inventory, itemID)

	if err := lm.SavePlayerStats(ctx, stats); err != nil {
		return nil, err
	}

	return &PurchaseResult{
		Success: true,
		ItemID:  itemID,
		Coins:   stats.Coins,
	}, nil
}

// GetPlayerInventory returns the player's owned item IDs.
func (lm *LeaderboardManager) GetPlayerInventory(ctx context.Context, playerID, playerName string) ([]string, error) {
	stats, err := lm.getOrCreateStats(ctx, playerID, playerName)
	if err != nil {
		return nil, err
	}
	return stats.Inventory, nil
}
