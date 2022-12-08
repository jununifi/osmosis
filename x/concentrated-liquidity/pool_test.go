package concentrated_liquidity_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v13/x/concentrated-liquidity/internal/math"
	"github.com/osmosis-labs/osmosis/v13/x/concentrated-liquidity/types"
)

func (s *KeeperTestSuite) TestCalcAndSwapOutAmtGivenIn() {
	tests := map[string]struct {
		positionAmount0  sdk.Int
		positionAmount1  sdk.Int
		addPositions     func(ctx sdk.Context, poolId uint64)
		tokenIn          sdk.Coin
		tokenOutDenom    string
		priceLimit       sdk.Dec
		expectedTokenIn  sdk.Coin
		expectedTokenOut sdk.Coin
		expectedTick     sdk.Int
		newLowerPrice    sdk.Dec
		newUpperPrice    sdk.Dec
		poolLiqAmount0   sdk.Int
		poolLiqAmount1   sdk.Int
		expectErr        bool
	}{
		//  One price range
		//
		//          5000
		//  4545 -----|----- 5500
		"single position within one tick: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(5004),
			// params
			// liquidity: 		 1517818840.967515822610790519
			// sqrtPriceNext:    70.738349405152439867 which is 5003.914076565430543175 https://www.wolframalpha.com/input?i=70.710678118654752440+%2B+42000000+%2F+1517818840.967515822610790519
			// sqrtPriceCurrent: 70.710678118654752440 which is 5000
			// expectedTokenIn:  42000000.0000 rounded up https://www.wolframalpha.com/input?i=1517818840.967515822610790519+*+%2870.738349405152439867+-+70.710678118654752440%29
			// expectedTokenOut: 8396.714105 rounded down https://www.wolframalpha.com/input?i=%281517818840.967515822610790519+*+%2870.738349405152439867+-+70.710678118654752440+%29%29+%2F+%2870.710678118654752440+*+70.738349405152439867%29
			// expectedTick: 	 85184.0 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C5003.914076565430543175%5D
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(8396)),
			expectedTick:     sdk.NewInt(85184),
		},
		"single position within one tick: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(13370)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4993),
			// params
			// liquidity: 		 1517818840.967515822610790519
			// sqrtPriceNext:    70.666662070529219856 which is 4993.777128190373086350 https://www.wolframalpha.com/input?i=%28%281517818840.967515822610790519%29%29+%2F+%28%28%281517818840.967515822610790519%29+%2F+%2870.710678118654752440%29%29+%2B+%2813370%29%29
			// expectedTokenIn:  13369.9999 rounded up https://www.wolframalpha.com/input?i=%281517818840.967515822610790519+*+%2870.710678118654752440+-+70.666662070529219856+%29%29+%2F+%2870.666662070529219856+*+70.710678118654752440%29
			// expectedTokenOut: 66808387.149 rounded down https://www.wolframalpha.com/input?i=1517818840.967515822610790519+*+%2870.710678118654752440+-+70.666662070529219856%29
			// expectedTick: 	 85163.7 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C4993.777128190373086350%5D
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(13370)),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(66808387)),
			expectedTick:     sdk.NewInt(85163),
		},
		//  Two equal price ranges
		//
		//          5000
		//  4545 -----|----- 5500
		//  4545 -----|----- 5500
		"two positions within one tick: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// add second position
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(5002),
			// params
			// liquidity: 		 3035637681.935031645221581038
			// sqrtPriceNext:    70.724513761903596153 which is 5001.956846857691162236 https://www.wolframalpha.com/input?i=70.710678118654752440%2B%2842000000+%2F+3035637681.935031645221581038%29
			// sqrtPriceCurrent: 70.710678118654752440 which is 5000
			// expectedTokenIn:  41999999.999 rounded up https://www.wolframalpha.com/input?i=3035637681.935031645221581038+*+%2870.724513761903596153+-+70.710678118654752440%29
			// expectedTokenOut: 8398.3567 rounded down https://www.wolframalpha.com/input?i=%283035637681.935031645221581038+*+%2870.724513761903596153+-+70.710678118654752440+%29%29+%2F+%2870.710678118654752440+*+70.724513761903596153%29
			// expectedTick:     85180.1 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C5003.914076565430543175%5D
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(8398)),
			expectedTick:     sdk.NewInt(85180),
			// two positions with same liquidity entered
			poolLiqAmount0: sdk.NewInt(1000000).MulRaw(2),
			poolLiqAmount1: sdk.NewInt(5000000000).MulRaw(2),
		},
		"two positions within one tick: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// add second position
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(13370)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4996),
			// params
			// liquidity: 		 3035637681.935031645221581038
			// sqrtPriceNext:    70.688663242671855280 which is 4996.887111035867053835 https://www.wolframalpha.com/input?i=%28%283035637681.935031645221581038%29%29+%2F+%28%28%283035637681.935031645221581038%29+%2F+%2870.710678118654752440%29%29+%2B+%2813370%29%29
			// sqrtPriceCurrent: 70.710678118654752440 which is 5000
			// expectedTokenIn:  13369.9999 rounded up https://www.wolframalpha.com/input?i=%283035637681.935031645221581038+*+%2870.710678118654752440+-+70.688663242671855280+%29%29+%2F+%2870.688663242671855280+*+70.710678118654752440%29
			// expectedTokenOut: 66829187.096 rounded down https://www.wolframalpha.com/input?i=3035637681.935031645221581038+*+%2870.710678118654752440+-+70.688663242671855280%29
			// expectedTick: 	 85170.00 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C4996.887111035867053835%5D
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(13370)),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(66829187)),
			expectedTick:     sdk.NewInt(85169), // TODO: should be 85170, is 85169 due to log precision
			// two positions with same liquidity entered
			poolLiqAmount0: sdk.NewInt(1000000).MulRaw(2),
			poolLiqAmount1: sdk.NewInt(5000000000).MulRaw(2),
		},
		//  Consecutive price ranges
		//
		//          5000
		//  4545 -----|----- 5500
		//             5500 ----------- 6250
		//
		"two positions with consecutive price ranges: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    74.160724590951092256 which is 5499.813071854898049815 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?i2d=true&i=Power%5B1.0001%2CDivide%5B86129%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  5236545537.865 rounded up https://www.wolframalpha.com/input?i=1517818840.967515822610790519+*+%2874.160724590951092256+-+70.710678118654752440%29
				// expectedTokenOut: 998587.023 rounded down https://www.wolframalpha.com/input?i=%281517818840.967515822610790519+*+%2874.161984870956629487+-+70.710678118654752440+%29%29+%2F+%2870.710678118654752440+*+74.161984870956629487%29
				// expectedTick:     86129.0 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C5499.813071854898049815%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(5500)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 86129
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// liquidity (2nd):  1198107969.043944887658592210
				// sqrtPriceNext:    78.136538612066568296 which is 6105.473934424522538231 https://www.wolframalpha.com/input?i=74.160724590951092256+%2B+4763454462.135+%2F+1198107969.043944887658592210
				// sqrtPriceCurrent: 74.160724590951092256 which is 5499.813071854898049815
				// expectedTokenIn:  4763454462.135 rounded up https://www.wolframalpha.com/input?i=1198107969.043944887658592210+*+%2878.136538612066568296+-+74.160724590951092256%29
				// expectedTokenOut: 822041.769 rounded down https://www.wolframalpha.com/input?i=%281198107969.043944887658592210+*+%2878.136538612066568296+-+74.160724590951092256+%29%29+%2F+%2874.160724590951092256+*+78.136538612066568296%29
				// expectedTick:     87173.8 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C6105.473934424522538231%5D
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(6106),
			// expectedTokenIn:  5236545537.865 + 4763454462.135 = 1000000000 usdc
			// expectedTokenOut: 998587.023 + 822041.769 = 1820628.792 round down = 1.820628 eth
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(1820628)),
			expectedTick:     sdk.NewInt(87173),
			newLowerPrice:    sdk.NewDec(5500),
			newUpperPrice:    sdk.NewDec(6250),
		},
		//  Consecutive price ranges
		//
		//                     5000
		//             4545 -----|----- 5500
		//  4000 ----------- 4545
		//
		"two positions with consecutive price ranges: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    67.416477345120534991 which is 4545 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?key=&i2d=true&i=Power%5B1.0001%2CDivide%5B84222%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  1048863.4367 rounded up https://www.wolframalpha.com/input?key=&i=%281517818840.967515822610790519+*+%2870.710678118654752440+-+67.416477345120534991%29%29+%2F+%2867.416477345120534991+*+70.710678118654752440%29
				// expectedTokenOut: 5000000000.000 rounded down https://www.wolframalpha.com/input?key=&i=1517818840.967515822610790519+*+%2870.710678118654752440-+67.416477345120534991%29
				// expectedTick:     84222.0 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4545%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(4000)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 82944
				newUpperPrice := sdk.NewDec(4545)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 84222

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// liquidity (2nd):  1198190689.904476625405593381
				// sqrtPriceNext:    63.991892380355787289 which is 4094.962290419 https://www.wolframalpha.com/input?key=&i=%28%281198190689.904476625405593381%29%29+%2F+%28%28%281198190689.904476625405593381%29+%2F+%2867.416477345120534991%29%29+%2B+%28951136.5633%29%29
				// sqrtPriceCurrent: 67.416477345120534991 which is 4545
				// expectedTokenIn:  951136.563300 rounded up https://www.wolframalpha.com/input?key=&i=%281198190689.904476625405593381+*+%2867.416477345120534991+-+63.991892380355787289%29%29+%2F+%2863.991892380355787289+*+67.416477345120534991%29
				// expectedTokenOut: 4103305821.5679708 rounded down https://www.wolframalpha.com/input?key=&i=1198190689.904476625405593381+*+%2867.416477345120534991-+63.991892380355787289%29
				// expectedTick:     83179.3 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4094.962290419%5D
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(2000000)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4094),
			// expectedTokenIn:  1048863.4367 + 951136.563300 = 2000000 eth
			// expectedTokenOut: 5000000000.000 + 4103305821.5679708 = 9103305821.5679708 round down = 9103.305821 usdc
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(2000000)),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(9103305821)),
			expectedTick:     sdk.NewInt(83179),
			newLowerPrice:    sdk.NewDec(4000),
			newUpperPrice:    sdk.NewDec(4545),
		},
		//  Partially overlapping price ranges
		//
		//          5000
		//  4545 -----|----- 5500
		//        5001 ----------- 6250
		//
		"two positions with partially overlapping price ranges: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    74.160724590951092256 which is 5499.813071854898049815 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?i2d=true&i=Power%5B1.0001%2CDivide%5B86129%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  5236545537.864897 rounded up https://www.wolframalpha.com/input?i=1517818840.967515822610790519+*+%2874.160724590951092256+-+70.710678118654752440%29
				// expectedTokenOut: 998934.824728 rounded down https://www.wolframalpha.com/input?i=%281517818840.967515822610790519+*+%2874.161984870956629487+-+70.710678118654752440+%29%29+%2F+%2870.710678118654752440+*+74.161984870956629487%29
				// expectedTick:     86129.0 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C5499.813071854898049815%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(5001)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 85178
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// liquidity (2nd):  670565280.937711183763565748
				// sqrtPriceNext:    77.820305833374877582 which is 6056.0000000000000000 we hit the price limit here, so we just use the user defined max (6056)
				// sqrtPriceCurrent: 70.717075849691272487 which is 5000.9048167309886086
				// expectedTokenIn:  4763179409.57397 rounded up https://www.wolframalpha.com/input?i=670565280.937711183763565748+*+%2877.820305833374877582+-+70.717075849691272487%29
				// expectedTokenOut: 865525.190 rounded down https://www.wolframalpha.com/input?i=%28670565280.937711183763565748+*+%2877.820305833374877582+-+70.717075849691272487+%29%29+%2F+%2870.717075849691272487+*+77.820305833374877582%29
				// expectedTick:     87092.4 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C6056.0000000000000000%5D
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(6056),
			// expectedTokenIn:  5236545537.865 + 4763179409.57397 = 9999724947.43897 = 999972.49 usdc
			// expectedTokenOut: 998587.023 + 865525.190 = 1864112.213 round down = 1.864112 eth
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(9999724947)),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(1864112)),
			expectedTick:     sdk.NewInt(87092),
			newLowerPrice:    sdk.NewDec(5001),
			newUpperPrice:    sdk.NewDec(6250),
		},
		"two positions with partially overlapping price ranges, not utilizing full liquidity of second position: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    74.160724590951092256 which is 5499.813071854898049815 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?i2d=true&i=Power%5B1.0001%2CDivide%5B86129%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  5236545537.864897 rounded up https://www.wolframalpha.com/input?i=1517818840.967515822610790519+*+%2874.160724590951092256+-+70.710678118654752440%29
				// expectedTokenOut: 998934.824728 rounded down https://www.wolframalpha.com/input?i=%281517818840.967515822610790519+*+%2874.161984870956629487+-+70.710678118654752440+%29%29+%2F+%2870.710678118654752440+*+74.161984870956629487%29
				// expectedTick:     86129.0 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C5499.813071854898049815%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(5001)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 85178
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// liquidity (2nd):  670565280.937711183763565748
				// sqrtPriceNext:    75.583797338121934796 which is 5712.9104200502884761 https://www.wolframalpha.com/input?key=&i=70.717075849691272487+%2B+3263454462.135103+%2F+670565280.937711183763565748
				// sqrtPriceCurrent: 70.717075849691272487 which is 5000.9048167309886086
				// expectedTokenIn:  3263454462.13510 rounded up https://www.wolframalpha.com/input?key=&i=670565280.937711183763565748+*+%2875.583797338121934796+-+70.717075849691272487%29
				// expectedTokenOut: 610554.667 rounded down https://www.wolframalpha.com/input?key=&i=%28670565280.937711183763565748+*+%2875.583797338121934796+-+70.717075849691272487+%29%29+%2F+%2870.717075849691272487+*+75.583797338121934796%29
				// expectedTick:     86509.2 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C5712.9104200502884761%5D
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(8500000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(6056),
			// expectedTokenIn:  5236545537.865 + 3263454462.13510 = 8500000000.000 = 8500.00 usdc
			// expectedTokenOut: 998587.023 + 610554.667 = 1609141.69 round down = 1.609141 eth
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(8500000000)),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(1609141)),
			expectedTick:     sdk.NewInt(86509),
			newLowerPrice:    sdk.NewDec(5001),
			newUpperPrice:    sdk.NewDec(6250),
		},
		//  Partially overlapping price ranges
		//
		//                5000
		//        4545 -----|----- 5500
		//  4000 ----------- 4999
		//
		"two positions with partially overlapping price ranges: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    67.416477345120534991 which is 4545 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?key=&i2d=true&i=Power%5B1.0001%2CDivide%5B84222%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  1048863.4367 rounded up https://www.wolframalpha.com/input?key=&i=%281517818840.967515822610790519+*+%2870.710678118654752440+-+67.416477345120534991%29%29+%2F+%2867.416477345120534991+*+70.710678118654752440%29
				// expectedTokenOut: 5000000000.000 rounded down https://www.wolframalpha.com/input?key=&i=1517818840.967515822610790519+*+%2870.710678118654752440-+67.416477345120534991%29
				// expectedTick:     84222.0 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4545%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(4000)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 82944
				newUpperPrice := sdk.NewDec(4999)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 85174

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// liquidity (2nd):  670293788.068821610959382388
				// sqrtPriceNext:    64.256329884039174301 which is 4128.8759302 https://www.wolframalpha.com/input?key=&i=%28%28670293788.068821610959382388%29%29+%2F+%28%28%28670293788.068821610959382388%29+%2F+%2870.702934555750545592%29%29+%2B+%28951136.5633%29%29
				// sqrtPriceCurrent: 70.702934555750545592 which is 4998.9049548
				// expectedTokenIn:  951136.5633 rounded up https://www.wolframalpha.com/input?key=&i=%28670293788.068821610959382388+*+%2870.702934555750545592+-+64.256329884039174301%29%29+%2F+%2864.256329884039174301+*+70.702934555750545592%29
				// expectedTokenOut: 4321119065.5835772240 rounded down https://www.wolframalpha.com/input?key=&i=670293788.068821610959382388+*+%2870.702934555750545592-+64.256329884039174301%29
				// expectedTick:     83261.9 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4128.9472754%5D
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(2000000)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4128),
			// expectedTokenIn:  1048863.4367 + 951136.5633 = 2000000 eth
			// expectedTokenOut: 5000000000.000 + 4321119065.5835772240 = 9321119065.583577224 round down = 9321.119065 usdc
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(2000000)),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(9321119065)),
			expectedTick:     sdk.NewInt(83261),
			newLowerPrice:    sdk.NewDec(4000),
			newUpperPrice:    sdk.NewDec(4999),
		},
		"two positions with partially overlapping price ranges, not utilizing full liquidity of second position: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    67.416477345120534991 which is 4545 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?key=&i2d=true&i=Power%5B1.0001%2CDivide%5B84222%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  1048863.4367 rounded up https://www.wolframalpha.com/input?key=&i=%281517818840.967515822610790519+*+%2870.710678118654752440+-+67.416477345120534991%29%29+%2F+%2867.416477345120534991+*+70.710678118654752440%29
				// expectedTokenOut: 5000000000.000 rounded down https://www.wolframalpha.com/input?key=&i=1517818840.967515822610790519+*+%2870.710678118654752440-+67.416477345120534991%29
				// expectedTick:     84222.0 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4545%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(4000)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 82944
				newUpperPrice := sdk.NewDec(4999)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 85174

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// 1800000 - 1048863.4367 = 751136.5633
				// liquidity (2nd):  670293788.068821610959382388
				// sqrtPriceNext:    65.512371527657899703 which is 4291.8708232 https://www.wolframalpha.com/input?key=&i=%28%28670293788.068821610959382388%29%29+%2F+%28%28%28670293788.068821610959382388%29+%2F+%2870.702934555750545592%29%29+%2B+%28751136.5633%29%29
				// sqrtPriceCurrent: 70.702934555750545592 which is 4998.9049548
				// expectedTokenIn:  751136.5633 rounded up https://www.wolframalpha.com/input?key=&i=%28670293788.068821610959382388+*+%2870.702934555750545592+-+65.512371527657899703%29%29+%2F+%2865.512371527657899703+*+70.702934555750545592%29
				// expectedTokenOut: 3479202154.310192937 rounded down https://www.wolframalpha.com/input?key=&i=670293788.068821610959382388+*+%2870.702934555750545592-+65.512371527657899703%29
				// expectedTick:     83649.0 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4291.8708232%5D
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(1800000)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4128),
			// expectedTokenIn:  1048863.4367 + 751136.5633 = 1.800000 eth
			// expectedTokenOut: 5000000000.000 + 3479202154.310192937 = 8479202154.310192937 round down = 8479.202154 usdc
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(1800000)),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(8479202154)),
			expectedTick:     sdk.NewInt(83648),
			newLowerPrice:    sdk.NewDec(4000),
			newUpperPrice:    sdk.NewDec(4999),
		},
		//  Sequential price ranges with a gap
		//
		//          5000
		//  4545 -----|----- 5500
		//              5501 ----------- 6250
		//
		"two sequential positions with a gap": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
				// params
				// liquidity (1st):  1517818840.967515822610790519
				// sqrtPriceNext:    74.160724590951092256 which is 5499.813071854898049815 (this is calculated by finding the closest tick LTE the upper range of the first range) https://www.wolframalpha.com/input?i2d=true&i=Power%5B1.0001%2CDivide%5B86129%2C2%5D%5D
				// sqrtPriceCurrent: 70.710678118654752440 which is 5000
				// expectedTokenIn:  5236545537.864897 rounded up https://www.wolframalpha.com/input?i=1517818840.967515822610790519+*+%2874.160724590951092256+-+70.710678118654752440%29
				// expectedTokenOut: 998934.824728 rounded down https://www.wolframalpha.com/input?i=%281517818840.967515822610790519+*+%2874.161984870956629487+-+70.710678118654752440+%29%29+%2F+%2870.710678118654752440+*+74.161984870956629487%29
				// expectedTick:     86129.0 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C5499.813071854898049815%5D

				// create second position parameters
				newLowerPrice := sdk.NewDec(5501)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 86131
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
				// params
				// liquidity (2nd):  1200046517.432645168443803695
				// sqrtPriceNext:    78.137532176937376749 which is 6105.473934701923906716 https://www.wolframalpha.com/input?i=74.168140663410187419++%2B++4763454462.135+%2F+1200046517.432645168443803695
				// sqrtPriceCurrent: 74.168140663410187419 which is 5500.913089467399755950
				// expectedTokenIn:  4763454462.135 rounded up https://www.wolframalpha.com/input?i=1200046517.432645168443803695+*+%2878.137532176937376749+-+74.168140663410187419%29
				// expectedTokenOut: 821949.120898 rounded down https://www.wolframalpha.com/input?i=%281200046517.432645168443803695+*+%2878.137532176937376749+-+74.168140663410187419+%29%29+%2F+%2874.168140663410187419+*+78.137532176937376749%29
				// expectedTick:     87173.8 rounded down https://www.wolframalpha.com/input?i2d=true&i=Log%5B1.0001%2C6105.473934424522538231%5D
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(6106),
			// expectedTokenIn:  5236545537.865 + 4763454462.135 = 1000000000 usdc
			// expectedTokenOut: 998587.023 + 821949.120898 = 1820536.143 round down = 1.820536 eth
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(1820536)),
			expectedTick:     sdk.NewInt(87173),
			newLowerPrice:    sdk.NewDec(5501),
			newUpperPrice:    sdk.NewDec(6250),
		},
		// Slippage protection doesn't cause a failure but interrupts early.
		"single position within one tick, trade completes but slippage protection interrupts trade early: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(13370)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4994),
			// params
			// liquidity: 		 1517818840.967515822610790519
			// sqrtPriceNext:    70.668238976219012614 which is 4994 https://www.wolframalpha.com/input?i=70.710678118654752440+%2B+42000000+%2F+1517818840.967515822610790519
			// sqrtPriceCurrent: 70.710678118654752440 which is 5000
			// expectedTokenIn:  12890.72275 rounded up https://www.wolframalpha.com/input?key=&i=%281517818840.967515822610790519+*+%2870.710678118654752440+-+70.668238976219012614+%29%29+%2F+%2870.710678118654752440+*+70.668238976219012614%29
			// expectedTokenOut: 64414929.9834 rounded down https://www.wolframalpha.com/input?key=&i=1517818840.967515822610790519+*+%2870.710678118654752440+-+70.668238976219012614%29
			// expectedTick: 	 85164.2 rounded down https://www.wolframalpha.com/input?key=&i2d=true&i=Log%5B1.0001%2C4994%5D
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(12891)),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(64414929)),
			expectedTick:     sdk.NewInt(85164),
		},
		"single position within one tick, trade does not complete due to lack of liquidity: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("usdc", sdk.NewInt(5300000000)),
			tokenOutDenom: "eth",
			priceLimit:    sdk.NewDec(6000),
			expectErr:     true,
		},
		"single position within one tick, trade does not complete due to lack of liquidity: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenIn:       sdk.NewCoin("eth", sdk.NewInt(1100000)),
			tokenOutDenom: "usdc",
			priceLimit:    sdk.NewDec(4000),
			expectErr:     true,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			s.Setup()
			s.FundAcc(s.TestAccs[0], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))
			s.FundAcc(s.TestAccs[1], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))

			// Create default CL pool
			pool := s.PrepareDefaultPool(s.Ctx)

			// add positions
			test.addPositions(s.Ctx, pool.GetId())

			// perform calc
			tokenIn, tokenOut, updatedTick, updatedLiquidity, _, err := s.App.ConcentratedLiquidityKeeper.CalcOutAmtGivenInInternal(
				s.Ctx,
				test.tokenIn, test.tokenOutDenom,
				DefaultZeroSwapFee, test.priceLimit, pool.GetId())
			if test.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				s.Require().Equal(test.expectedTokenIn.String(), tokenIn.String())
				s.Require().Equal(test.expectedTokenOut.String(), tokenOut.String())
				s.Require().Equal(test.expectedTick.String(), updatedTick.String())

				if test.newLowerPrice.IsNil() && test.newUpperPrice.IsNil() {
					test.newLowerPrice = DefaultLowerPrice
					test.newUpperPrice = DefaultUpperPrice
				}

				newLowerTick := math.PriceToTick(test.newLowerPrice)
				newUpperTick := math.PriceToTick(test.newUpperPrice)

				lowerSqrtPrice, err := math.TickToSqrtPrice(newLowerTick)
				s.Require().NoError(err)
				upperSqrtPrice, err := math.TickToSqrtPrice(newUpperTick)
				s.Require().NoError(err)

				if test.poolLiqAmount0.IsNil() && test.poolLiqAmount1.IsNil() {
					test.poolLiqAmount0 = DefaultAmt0
					test.poolLiqAmount1 = DefaultAmt1
				}

				expectedLiquidity := math.GetLiquidityFromAmounts(DefaultCurrSqrtPrice, lowerSqrtPrice, upperSqrtPrice, test.poolLiqAmount0, test.poolLiqAmount1)
				s.Require().Equal(expectedLiquidity.String(), updatedLiquidity.String())
			}

			// perform swap
			_, tokenOut, _, _, _, err = s.App.ConcentratedLiquidityKeeper.SwapOutAmtGivenIn(
				s.Ctx,
				test.tokenIn, test.tokenOutDenom,
				DefaultZeroSwapFee, test.priceLimit, pool.GetId())
			if test.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				pool, err = s.App.ConcentratedLiquidityKeeper.GetPoolById(s.Ctx, pool.GetId())
				s.Require().NoError(err)

				// check that the pool's current tick was updated correctly
				s.Require().Equal(test.expectedTick.String(), pool.GetCurrentTick().String())
				// check that we produced the same token out from the swap function that we expected
				s.Require().Equal(test.expectedTokenOut.String(), tokenOut.String())

				// the following is needed to get the expected liquidity to later compare to what the pool was updated to
				if test.newLowerPrice.IsNil() && test.newUpperPrice.IsNil() {
					test.newLowerPrice = DefaultLowerPrice
					test.newUpperPrice = DefaultUpperPrice
				}

				newLowerTick := math.PriceToTick(test.newLowerPrice)
				newUpperTick := math.PriceToTick(test.newUpperPrice)

				lowerSqrtPrice, err := math.TickToSqrtPrice(newLowerTick)
				s.Require().NoError(err)
				upperSqrtPrice, err := math.TickToSqrtPrice(newUpperTick)
				s.Require().NoError(err)

				if test.poolLiqAmount0.IsNil() && test.poolLiqAmount1.IsNil() {
					test.poolLiqAmount0 = DefaultAmt0
					test.poolLiqAmount1 = DefaultAmt1
				}

				expectedLiquidity := math.GetLiquidityFromAmounts(DefaultCurrSqrtPrice, lowerSqrtPrice, upperSqrtPrice, test.poolLiqAmount0, test.poolLiqAmount1)
				// check that the pools liquidity was updated correctly
				s.Require().Equal(expectedLiquidity.String(), pool.GetLiquidity().String())

				// TODO: need to figure out a good way to test that the currentSqrtPrice that the pool is set to makes sense
				// right now we calculate this value through iterations, so unsure how to do this here / if its needed
			}
		})

	}
}

func (s *KeeperTestSuite) TestCalcAndSwapInAmtGivenOut() {
	tests := map[string]struct {
		positionAmount0  sdk.Int
		positionAmount1  sdk.Int
		addPositions     func(ctx sdk.Context, poolId uint64)
		tokenOut         sdk.Coin
		tokenInDenom     string
		priceLimit       sdk.Dec
		expectedTokenIn  sdk.Coin
		expectedTokenOut sdk.Coin
		expectedTick     sdk.Int
		newLowerPrice    sdk.Dec
		newUpperPrice    sdk.Dec
		poolLiqAmount0   sdk.Int
		poolLiqAmount1   sdk.Int
		expectErr        bool
	}{
		//  One price range
		//
		//          5000
		//  4545 -----|----- 5500
		"single position within one tick: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			tokenInDenom:     "eth",
			priceLimit:       sdk.NewDec(5004),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(8396)),
			expectedTick:     sdk.NewInt(85184),
		},
		"single position within one tick: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("eth", sdk.NewInt(13370)),
			tokenInDenom:     "usdc",
			priceLimit:       sdk.NewDec(4993),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(13370)),
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(66808387)),
			expectedTick:     sdk.NewInt(85163),
		},
		//  Two equal price ranges
		//
		//          5000
		//  4545 -----|----- 5500
		//  4545 -----|----- 5500
		"two positions within one tick: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// add second position
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			tokenInDenom:     "eth",
			priceLimit:       sdk.NewDec(5002),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(42000000)),
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(8398)),
			expectedTick:     sdk.NewInt(85180),
			// two positions with same liquidity entered
			poolLiqAmount0: sdk.NewInt(1000000).MulRaw(2),
			poolLiqAmount1: sdk.NewInt(5000000000).MulRaw(2),
		},
		"two positions within one tick: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// add second position
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("eth", sdk.NewInt(13370)),
			tokenInDenom:     "usdc",
			priceLimit:       sdk.NewDec(4996),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(13370)),
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(66829187)),
			expectedTick:     sdk.NewInt(85169), // TODO: should be 85170, is 85169 due to log precision
			// two positions with same liquidity entered
			poolLiqAmount0: sdk.NewInt(1000000).MulRaw(2),
			poolLiqAmount1: sdk.NewInt(5000000000).MulRaw(2),
		},
		//  Consecutive price ranges
		//
		//          5000
		//  4545 -----|----- 5500
		//             5500 ----------- 6250
		//
		"two positions with consecutive price ranges: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(5500)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 86129
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			tokenInDenom:     "eth",
			priceLimit:       sdk.NewDec(6106),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(1820628)),
			expectedTick:     sdk.NewInt(87173),
			newLowerPrice:    sdk.NewDec(5500),
			newUpperPrice:    sdk.NewDec(6250),
		},
		//  Consecutive price ranges
		//
		//                     5000
		//             4545 -----|----- 5500
		//  4000 ----------- 4545
		//
		"two positions with consecutive price ranges: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(4000)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 82944
				newUpperPrice := sdk.NewDec(4545)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 84222

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("eth", sdk.NewInt(2000000)),
			tokenInDenom:     "usdc",
			priceLimit:       sdk.NewDec(4094),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(2000000)),
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(9103305821)),
			expectedTick:     sdk.NewInt(83179),
			newLowerPrice:    sdk.NewDec(4000),
			newUpperPrice:    sdk.NewDec(4545),
		},
		//  Partially overlapping price ranges
		//
		//          5000
		//  4545 -----|----- 5500
		//        5001 ----------- 6250
		//
		"two positions with partially overlapping price ranges: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(5001)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 85178
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			tokenInDenom:     "eth",
			priceLimit:       sdk.NewDec(6056),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(9999724947)),
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(1864112)),
			expectedTick:     sdk.NewInt(87092),
			newLowerPrice:    sdk.NewDec(5001),
			newUpperPrice:    sdk.NewDec(6250),
		},
		"two positions with partially overlapping price ranges, not utilizing full liquidity of second position: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(5001)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 85178
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("usdc", sdk.NewInt(8500000000)),
			tokenInDenom:     "eth",
			priceLimit:       sdk.NewDec(6056),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(8500000000)),
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(1609141)),
			expectedTick:     sdk.NewInt(86509),
			newLowerPrice:    sdk.NewDec(5001),
			newUpperPrice:    sdk.NewDec(6250),
		},
		//  Partially overlapping price ranges
		//
		//                5000
		//        4545 -----|----- 5500
		//  4000 ----------- 4999
		//
		"two positions with partially overlapping price ranges: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(4000)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 82944
				newUpperPrice := sdk.NewDec(4999)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 85174

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("eth", sdk.NewInt(2000000)),
			tokenInDenom:     "usdc",
			priceLimit:       sdk.NewDec(4128),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(2000000)),
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(9321119065)),
			expectedTick:     sdk.NewInt(83261),
			newLowerPrice:    sdk.NewDec(4000),
			newUpperPrice:    sdk.NewDec(4999),
		},
		"two positions with partially overlapping price ranges, not utilizing full liquidity of second position: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(4000)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 82944
				newUpperPrice := sdk.NewDec(4999)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 85174

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("eth", sdk.NewInt(1800000)),
			tokenInDenom:     "usdc",
			priceLimit:       sdk.NewDec(4128),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(1800000)),
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(8479202154)),
			expectedTick:     sdk.NewInt(83648),
			newLowerPrice:    sdk.NewDec(4000),
			newUpperPrice:    sdk.NewDec(4999),
		},
		//  Sequential price ranges with a gap
		//
		//          5000
		//  4545 -----|----- 5500
		//              5501 ----------- 6250
		//
		"two sequential positions with a gap": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)

				// create second position parameters
				newLowerPrice := sdk.NewDec(5501)
				s.Require().NoError(err)
				newLowerTick := math.PriceToTick(newLowerPrice) // 86131
				newUpperPrice := sdk.NewDec(6250)
				s.Require().NoError(err)
				newUpperTick := math.PriceToTick(newUpperPrice) // 87407

				// add position two with the new price range above
				_, _, _, err = s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[1], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), newLowerTick.Int64(), newUpperTick.Int64())
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			tokenInDenom:     "eth",
			priceLimit:       sdk.NewDec(6106),
			expectedTokenOut: sdk.NewCoin("usdc", sdk.NewInt(10000000000)),
			expectedTokenIn:  sdk.NewCoin("eth", sdk.NewInt(1820536)),
			expectedTick:     sdk.NewInt(87173),
			newLowerPrice:    sdk.NewDec(5501),
			newUpperPrice:    sdk.NewDec(6250),
		},
		// Slippage protection doesn't cause a failure but interrupts early.
		"single position within one tick, trade completes but slippage protection interrupts trade early: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:         sdk.NewCoin("eth", sdk.NewInt(13370)),
			tokenInDenom:     "usdc",
			priceLimit:       sdk.NewDec(4994),
			expectedTokenOut: sdk.NewCoin("eth", sdk.NewInt(12891)),
			expectedTokenIn:  sdk.NewCoin("usdc", sdk.NewInt(64414929)),
			expectedTick:     sdk.NewInt(85164),
		},
		"single position within one tick, trade does not complete due to lack of liquidity: usdc -> eth": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:     sdk.NewCoin("usdc", sdk.NewInt(5300000000)),
			tokenInDenom: "eth",
			priceLimit:   sdk.NewDec(6000),
			expectErr:    true,
		},
		"single position within one tick, trade does not complete due to lack of liquidity: eth -> usdc": {
			addPositions: func(ctx sdk.Context, poolId uint64) {
				// add first position
				_, _, _, err := s.App.ConcentratedLiquidityKeeper.CreatePosition(ctx, poolId, s.TestAccs[0], DefaultAmt0, DefaultAmt1, sdk.ZeroInt(), sdk.ZeroInt(), DefaultLowerTick, DefaultUpperTick)
				s.Require().NoError(err)
			},
			tokenOut:     sdk.NewCoin("eth", sdk.NewInt(1100000)),
			tokenInDenom: "usdc",
			priceLimit:   sdk.NewDec(4000),
			expectErr:    true,
		},
	}

	for name, test := range tests {
		s.Run(name, func() {
			s.Setup()
			s.FundAcc(s.TestAccs[0], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))
			s.FundAcc(s.TestAccs[1], sdk.NewCoins(sdk.NewCoin("eth", sdk.NewInt(10000000000000)), sdk.NewCoin("usdc", sdk.NewInt(1000000000000))))

			// Create default CL pool
			pool := s.PrepareDefaultPool(s.Ctx)

			// add positions
			test.addPositions(s.Ctx, pool.GetId())

			// perform calc
			tokenIn, tokenOut, updatedTick, updatedLiquidity, _, err := s.App.ConcentratedLiquidityKeeper.CalcInAmtGivenOutInternal(
				s.Ctx,
				test.tokenOut, test.tokenInDenom,
				DefaultZeroSwapFee, test.priceLimit, pool.GetId())
			if test.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				s.Require().Equal(test.expectedTokenOut.String(), tokenOut.String())
				s.Require().Equal(test.expectedTokenIn.String(), tokenIn.String())
				s.Require().Equal(test.expectedTick.String(), updatedTick.String())

				if test.newLowerPrice.IsNil() && test.newUpperPrice.IsNil() {
					test.newLowerPrice = DefaultLowerPrice
					test.newUpperPrice = DefaultUpperPrice
				}

				newLowerTick := math.PriceToTick(test.newLowerPrice)
				newUpperTick := math.PriceToTick(test.newUpperPrice)

				lowerSqrtPrice, err := math.TickToSqrtPrice(newLowerTick)
				s.Require().NoError(err)
				upperSqrtPrice, err := math.TickToSqrtPrice(newUpperTick)
				s.Require().NoError(err)

				if test.poolLiqAmount0.IsNil() && test.poolLiqAmount1.IsNil() {
					test.poolLiqAmount0 = DefaultAmt0
					test.poolLiqAmount1 = DefaultAmt1
				}

				expectedLiquidity := math.GetLiquidityFromAmounts(DefaultCurrSqrtPrice, lowerSqrtPrice, upperSqrtPrice, test.poolLiqAmount0, test.poolLiqAmount1)
				s.Require().Equal(expectedLiquidity.String(), updatedLiquidity.String())
			}

			// perform swap
			tokenIn, _, _, _, _, err = s.App.ConcentratedLiquidityKeeper.SwapInAmtGivenOut(
				s.Ctx,
				test.tokenOut, test.tokenInDenom,
				DefaultZeroSwapFee, test.priceLimit, pool.GetId())
			if test.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				pool, err = s.App.ConcentratedLiquidityKeeper.GetPoolById(s.Ctx, pool.GetId())
				s.Require().NoError(err)

				// check that the pool's current tick was updated correctly
				s.Require().Equal(test.expectedTick.String(), pool.GetCurrentTick().String())
				// check that we produced the same token out from the swap function that we expected
				s.Require().Equal(test.expectedTokenIn.String(), tokenIn.String())

				// the following is needed to get the expected liquidity to later compare to what the pool was updated to
				if test.newLowerPrice.IsNil() && test.newUpperPrice.IsNil() {
					test.newLowerPrice = DefaultLowerPrice
					test.newUpperPrice = DefaultUpperPrice
				}

				newLowerTick := math.PriceToTick(test.newLowerPrice)
				newUpperTick := math.PriceToTick(test.newUpperPrice)

				lowerSqrtPrice, err := math.TickToSqrtPrice(newLowerTick)
				s.Require().NoError(err)
				upperSqrtPrice, err := math.TickToSqrtPrice(newUpperTick)
				s.Require().NoError(err)

				if test.poolLiqAmount0.IsNil() && test.poolLiqAmount1.IsNil() {
					test.poolLiqAmount0 = DefaultAmt0
					test.poolLiqAmount1 = DefaultAmt1
				}

				expectedLiquidity := math.GetLiquidityFromAmounts(DefaultCurrSqrtPrice, lowerSqrtPrice, upperSqrtPrice, test.poolLiqAmount0, test.poolLiqAmount1)
				// check that the pools liquidity was updated correctly
				s.Require().Equal(expectedLiquidity.String(), pool.GetLiquidity().String())

				// TODO: need to figure out a good way to test that the currentSqrtPrice that the pool is set to makes sense
				// right now we calculate this value through iterations, so unsure how to do this here / if its needed
			}
		})

	}
}

func (s *KeeperTestSuite) TestOrderInitialPoolDenoms() {
	denom0, denom1, err := types.OrderInitialPoolDenoms("axel", "osmo")
	s.Require().NoError(err)
	s.Require().Equal(denom0, "axel")
	s.Require().Equal(denom1, "osmo")

	denom0, denom1, err = types.OrderInitialPoolDenoms("usdc", "eth")
	s.Require().NoError(err)
	s.Require().Equal(denom0, "eth")
	s.Require().Equal(denom1, "usdc")

	denom0, denom1, err = types.OrderInitialPoolDenoms("usdc", "usdc")
	s.Require().Error(err)

}

func (s *KeeperTestSuite) TestGetPoolById() {
	const validPoolId = 1

	tests := []struct {
		name        string
		poolId      uint64
		expectedErr error
	}{
		{
			name:   "Get existing pool",
			poolId: validPoolId,
		},
		{
			name:        "Get non-existing pool",
			poolId:      2,
			expectedErr: types.PoolNotFoundError{PoolId: 2},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()

			// Create default CL pool
			pool := s.PrepareDefaultPool(s.Ctx)

			// Get pool defined in test case
			getPool, err := s.App.ConcentratedLiquidityKeeper.GetPoolById(s.Ctx, test.poolId)

			if test.expectedErr == nil {
				// Ensure no error is returned
				s.Require().NoError(err)

				// Ensure that pool returned matches the default pool attributes
				s.Require().Equal(pool.GetId(), getPool.GetId())
				s.Require().Equal(pool.GetAddress(), getPool.GetAddress())
				s.Require().Equal(pool.GetCurrentSqrtPrice(), getPool.GetCurrentSqrtPrice())
				s.Require().Equal(pool.GetCurrentTick(), getPool.GetCurrentTick())
				s.Require().Equal(pool.GetLiquidity(), getPool.GetLiquidity())
			} else {
				// Ensure specified error is returned
				s.Require().Error(err)
				s.Require().ErrorIs(err, test.expectedErr)

				// Check that GetPoolById returns a nil pool object due to error
				s.Require().Nil(getPool)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPoolExists() {
	s.SetupTest()

	// Create default CL pool
	pool := s.PrepareDefaultPool(s.Ctx)

	// Check that the pool exists
	poolExists := s.App.ConcentratedLiquidityKeeper.PoolExists(s.Ctx, pool.GetId())
	s.Require().True(poolExists)

	// try checking for a non-existent pool
	poolExists = s.App.ConcentratedLiquidityKeeper.PoolExists(s.Ctx, 2)

	// ensure that this returns false
	s.Require().False(poolExists)
}
