package wibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	costGetCoins sdk.Gas = 10
	costHasCoins sdk.Gas = 10
)

type Keeper interface {
	SendKeeper
	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
}

// BaseKeeper manages transfers between accounts. It implements the Keeper
// interface.
type BaseKeeper struct {
	am auth.AccountKeeper
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(am auth.AccountKeeper) BaseKeeper {
	return BaseKeeper{am: am}
}

//______________________________________________________________________________________________

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	InputOutputCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins, forTime uint) (sdk.Tags, sdk.Error)
}

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	am auth.AccountKeeper
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(am auth.AccountKeeper) BaseViewKeeper {
	return BaseViewKeeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

//______________________________________________________________________________________________

func getCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress) sdk.Coins {
	ctx.GasMeter().ConsumeGas(costGetCoins, "getCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) bool {
	ctx.GasMeter().ConsumeGas(costHasCoins, "hasCoins")
	return getCoins(ctx, am, addr).IsAllGTE(amt)
}
