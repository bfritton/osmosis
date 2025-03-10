package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/x/lockup/types"
)

func (suite *KeeperTestSuite) LockTokens(addr sdk.AccAddress, coins sdk.Coins, duration time.Duration) {
	suite.app.BankKeeper.SetBalances(suite.ctx, addr, coins)
	_, err := suite.app.LockupKeeper.LockTokens(suite.ctx, addr, coins, duration)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) BeginUnlocking(addr sdk.AccAddress) {
	_, _, err := suite.app.LockupKeeper.BeginUnlockAllNotUnlockings(suite.ctx, addr)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestModuleBalance() {
	suite.SetupTest()

	// initial check
	res, err := suite.app.LockupKeeper.ModuleBalance(sdk.WrapSDKContext(suite.ctx), &types.ModuleBalanceRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// final check
	res, err = suite.app.LockupKeeper.ModuleBalance(sdk.WrapSDKContext(suite.ctx), &types.ModuleBalanceRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)
}

func (suite *KeeperTestSuite) TestModuleLockedAmount() {
	// test for module locked balance check
	suite.SetupTest()

	// initial check
	res, err := suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))

	// lock coins
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// current module locked balance check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)

	// module locked balance after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))

	// module locked balance after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))
}

func (suite *KeeperTestSuite) TestAccountUnlockableCoins() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address unlockable coins check
	res, err := suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: ""})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// initial check
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// check before start unlocking
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	suite.BeginUnlocking(addr1)

	// check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// check after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.AccountUnlockableCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)

	// check after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.AccountUnlockableCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)
}

func (suite *KeeperTestSuite) TestAccountUnlockingCoins() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address unlockable coins check
	res, err := suite.app.LockupKeeper.AccountUnlockingCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockingCoinsRequest{Owner: ""})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})
	// initial check
	res, err = suite.app.LockupKeeper.AccountUnlockingCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockingCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// check before start unlocking
	res, err = suite.app.LockupKeeper.AccountUnlockingCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockingCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	suite.BeginUnlocking(addr1)

	// check at unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountUnlockingCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockingCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{sdk.NewInt64Coin("stake", 10)})

	// check after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.AccountUnlockingCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.AccountUnlockingCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// check after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountUnlockingCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.AccountUnlockingCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})
}

func (suite *KeeperTestSuite) TestAccountLockedCoins() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locked coins check
	res, err := suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedCoinsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))
	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)

	// check after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.AccountLockedCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))

	// check after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.AccountLockedCoinsRequest{Owner: addr1.String()})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins(nil))
}

func (suite *KeeperTestSuite) TestAccountLockedPastTime() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address locks check
	res, err := suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: "", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// check after 1 second = unlockTime
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1.String(), Timestamp: now.Add(time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// check after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1.String(), Timestamp: now.Add(2 * time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestAccountLockedPastTimeNotUnlockingOnly() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address locks check
	res, err := suite.app.LockupKeeper.AccountLockedPastTimeNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeNotUnlockingOnlyRequest{Owner: "", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeNotUnlockingOnlyRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// check when not start unlocking
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeNotUnlockingOnlyRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// begin unlocking
	suite.BeginUnlocking(addr1)

	// check after start unlocking
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeNotUnlockingOnlyRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestAccountUnlockedBeforeTime() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address unlockables check
	res, err := suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: "", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial check
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1.String(), Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// check after 1 second = unlockTime
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1.String(), Timestamp: now.Add(time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// check after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1.String(), Timestamp: now.Add(2 * time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)
}

func (suite *KeeperTestSuite) TestAccountLockedPastTimeDenom() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address locks by denom check
	res, err := suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: "", Denom: "stake", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1.String(), Denom: "stake", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1.String(), Denom: "stake", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks by not available denom
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1.String(), Denom: "stake2", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks by denom after 1 second = unlockTime
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1.String(), Denom: "stake", Timestamp: now.Add(time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks by denom after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1.String(), Denom: "stake", Timestamp: now.Add(2 * time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// try querying with prefix coins like "stak" for potential attack
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1.String(), Denom: "stak", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestLockedByID() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// lock by not avaialble id check
	res, err := suite.app.LockupKeeper.LockedByID(sdk.WrapSDKContext(suite.ctx), &types.LockedRequest{LockId: 0})
	suite.Require().Error(err)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// lock by available availble id check
	res, err = suite.app.LockupKeeper.LockedByID(sdk.WrapSDKContext(suite.ctx), &types.LockedRequest{LockId: 1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Lock.ID, uint64(1))
	suite.Require().Equal(res.Lock.Owner, addr1.String())
	suite.Require().Equal(res.Lock.Coins, coins)
	suite.Require().Equal(res.Lock.Duration, time.Second)
	suite.Require().Equal(res.Lock.EndTime, time.Time{})
	suite.Require().Equal(res.Lock.IsUnlocking(), false)
}

func (suite *KeeperTestSuite) TestAccountLockedLongerDuration() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locks longer than duration check
	res, err := suite.app.LockupKeeper.AccountLockedLongerDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: "", Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedLongerDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1.String(), Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// account locks longer than duration check, duration = 0s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1.String(), Duration: 0})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1.String(), Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check, duration = 2s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1.String(), Duration: 2 * time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestAccountLockedLongerDurationNotUnlockingOnly() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locks longer than duration check
	res, err := suite.app.LockupKeeper.AccountLockedLongerDurationNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationNotUnlockingOnlyRequest{Owner: "", Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationNotUnlockingOnlyRequest{Owner: addr1.String(), Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account locks longer than duration check before start unlocking, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationNotUnlockingOnlyRequest{Owner: addr1.String(), Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	suite.BeginUnlocking(addr1)

	// account locks longer than duration check after start unlocking, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationNotUnlockingOnly(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationNotUnlockingOnlyRequest{Owner: addr1.String(), Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestAccountLockedLongerDurationDenom() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locks longer than duration by denom check
	res, err := suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: "", Duration: time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial check
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1.String(), Duration: time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)
	suite.BeginUnlocking(addr1)

	// account locks longer than duration check by denom, duration = 0s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1.String(), Duration: 0, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check by denom, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1.String(), Duration: time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check by not available denom, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1.String(), Duration: time.Second, Denom: "stake2"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks longer than duration check by denom, duration = 2s
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1.String(), Duration: 2 * time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// try querying with prefix coins like "stak" for potential attack
	res, err = suite.app.LockupKeeper.AccountLockedLongerDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1.String(), Duration: 0, Denom: "sta"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestLockedDenom() {
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Hour)

	// test with single lockup
	res, err := suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: 0})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(10))

	duration, _ := time.ParseDuration("30m")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(10))

	duration, _ = time.ParseDuration("1h")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(10))

	duration, _ = time.ParseDuration("2h")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(0))

	// adds different account and lockup for testing
	addr2 := sdk.AccAddress([]byte("addr2---------------"))

	coins = sdk.Coins{sdk.NewInt64Coin("stake", 20)}
	suite.LockTokens(addr2, coins, time.Hour*2)

	duration, _ = time.ParseDuration("30m")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(30))

	duration, _ = time.ParseDuration("1h")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(30))

	duration, _ = time.ParseDuration("2h")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(20))

	// test unlocking
	suite.BeginUnlocking(addr2)
	duration, _ = time.ParseDuration("2h")
	res, err = suite.app.LockupKeeper.LockedDenom(sdk.WrapSDKContext(suite.ctx), &types.LockedDenomRequest{Denom: "stake", Duration: duration})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Amount, sdk.NewInt(0))
}
