package usecase_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/domain"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/router/usecase"
	"github.com/osmosis-labs/osmosis/v20/x/gamm/pool-models/balancer"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v20/x/poolmanager/types"
)

// TestPrepareResult prepares the result of the quote for output to the client.
// First, it strips away unnecessary fields from each pool in the route.
// Additionally, it computes the effective spread factor from all routes.
func (s *RouterTestSuite) TestPrepareResult() {
	var (
		takerFeeOne = osmomath.NewDecWithPrec(2, 2)
		takerFeeTwo = osmomath.NewDecWithPrec(3, 3)

		defaultAmount = sdk.NewInt(100_000_00)

		totalInAmount  = defaultAmount
		totalOutAmount = defaultAmount.MulRaw(4)

		poolOneBalances = sdk.NewCoins(
			sdk.NewCoin(USDC, defaultAmount.MulRaw(5)),
			sdk.NewCoin(ETH, defaultAmount),
		)

		poolTwoBalances = sdk.NewCoins(
			sdk.NewCoin(ETH, defaultAmount),
			sdk.NewCoin(USDC, defaultAmount.MulRaw(4)),
		)
	)

	// Prepare 2 pools

	// Pool USDC / ETH -> 0.01 spread factor & 5 USDC for 1 ETH
	poolIDOne := s.PrepareCustomBalancerPool([]balancer.PoolAsset{
		{
			Token:  sdk.NewCoin(USDC, defaultAmount.MulRaw(5)),
			Weight: sdk.NewInt(100),
		},
		{
			Token:  sdk.NewCoin(ETH, defaultAmount),
			Weight: sdk.NewInt(100),
		},
	}, balancer.PoolParams{
		SwapFee: sdk.NewDecWithPrec(1, 2),
		ExitFee: osmomath.ZeroDec(),
	})

	poolOne, err := s.App.PoolManagerKeeper.GetPool(s.Ctx, poolIDOne)
	s.Require().NoError(err)

	// Pool ETH / USDC -> 0.005 spread factor & 4 USDC for 1 ETH
	poolIDTwo := s.PrepareCustomBalancerPool([]balancer.PoolAsset{
		{
			Token:  sdk.NewCoin(ETH, defaultAmount),
			Weight: sdk.NewInt(100),
		},
		{
			Token:  sdk.NewCoin(USDC, defaultAmount.MulRaw(4)),
			Weight: sdk.NewInt(100),
		},
	}, balancer.PoolParams{
		SwapFee: sdk.NewDecWithPrec(5, 3),
		ExitFee: osmomath.ZeroDec(),
	})

	poolTwo, err := s.App.PoolManagerKeeper.GetPool(s.Ctx, poolIDTwo)
	s.Require().NoError(err)

	testQuote := &usecase.QuoteImpl{
		AmountIn:  sdk.NewCoin(ETH, totalInAmount),
		AmountOut: totalOutAmount,

		// 2 routes with 50-50 split, each single hop
		Route: []domain.SplitRoute{

			// Route 1
			&usecase.RouteWithOutAmount{
				Route: &usecase.RouteImpl{
					Pools: []domain.RoutablePool{
						&usecase.RoutableCFMMPoolImpl{
							PoolI:         domain.NewPool(poolOne, poolOne.GetSpreadFactor(sdk.Context{}), poolOneBalances),
							TokenOutDenom: USDC,
							TakerFee:      takerFeeOne,
						},
					},
				},

				InAmount:  totalInAmount.QuoRaw(2),
				OutAmount: totalOutAmount.QuoRaw(2),
			},

			// Route 2
			&usecase.RouteWithOutAmount{
				Route: &usecase.RouteImpl{
					Pools: []domain.RoutablePool{
						&usecase.RoutableCFMMPoolImpl{
							PoolI:         domain.NewPool(poolTwo, poolTwo.GetSpreadFactor(sdk.Context{}), poolTwoBalances),
							TokenOutDenom: USDC,
							TakerFee:      takerFeeTwo,
						},
					},
				},

				InAmount:  totalInAmount.QuoRaw(2),
				OutAmount: totalOutAmount.QuoRaw(2),
			},
		},
		EffectiveSpreadFactor: osmomath.OneDec(),
	}

	expectedRoutes := []domain.SplitRoute{

		// Route 1
		&usecase.RouteWithOutAmount{
			Route: &usecase.RouteImpl{
				Pools: []domain.RoutablePool{
					&usecase.RoutableResultPoolImpl{
						ID:            poolIDOne,
						Type:          poolmanagertypes.Balancer,
						Balances:      poolOneBalances,
						SpreadFactor:  poolOne.GetSpreadFactor(sdk.Context{}),
						TokenOutDenom: USDC,
						TakerFee:      takerFeeOne,
					},
				},
			},

			InAmount:  totalInAmount.QuoRaw(2),
			OutAmount: totalOutAmount.QuoRaw(2),
		},

		// Route 2
		&usecase.RouteWithOutAmount{
			Route: &usecase.RouteImpl{
				Pools: []domain.RoutablePool{
					&usecase.RoutableResultPoolImpl{
						ID:            poolIDTwo,
						Type:          poolmanagertypes.Balancer,
						Balances:      poolTwoBalances,
						SpreadFactor:  poolTwo.GetSpreadFactor(sdk.Context{}),
						TokenOutDenom: USDC,
						TakerFee:      takerFeeTwo,
					},
				},
			},

			InAmount:  totalInAmount.QuoRaw(2),
			OutAmount: totalOutAmount.QuoRaw(2),
		},
	}

	// 0.01 * 0.5  + 0.005 * 0.5
	expectedEffectiveSpreadFactor := osmomath.MustNewDecFromStr("0.0075")

	// System under test
	routes, effectiveSpreadFactor := testQuote.PrepareResult()

	// Validate routes.
	s.validateRoutes(expectedRoutes, routes)

	// Validate effective spread factor.
	s.Require().Equal(expectedEffectiveSpreadFactor.String(), effectiveSpreadFactor.String())
}

// validateRoutes validates that the given routes are equal.
// Specifically, validates:
// - Pools
// - In amount
// - Out amount
func (s *RouterTestSuite) validateRoutes(expectedRoutes []domain.SplitRoute, actualRoutes []domain.SplitRoute) {
	s.Require().Equal(len(expectedRoutes), len(actualRoutes))
	for i, expectedRoute := range expectedRoutes {
		actualRoute := actualRoutes[i]

		// Validate pools
		s.validateRoutePools(expectedRoute.GetPools(), actualRoute.GetPools())

		// Validate in amount
		s.Require().Equal(expectedRoute.GetAmountIn().String(), actualRoute.GetAmountIn().String())

		// Validate out amount
		s.Require().Equal(expectedRoute.GetAmountOut().String(), actualRoute.GetAmountOut().String())
	}
}
