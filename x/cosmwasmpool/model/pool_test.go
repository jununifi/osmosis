package model_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/osmoutils"
	"github.com/osmosis-labs/osmosis/v20/app/apptesting"
)

type CosmWasmPoolSuite struct {
	apptesting.KeeperTestHelper
}

const (
	denomA = "axlusdc"
	denomB = "gravusdc"
)

func TestPoolModuleSuite(t *testing.T) {
	suite.Run(t, new(CosmWasmPoolSuite))
}

func (s *CosmWasmPoolSuite) SetupTest() {
	s.Setup()
}

// TestGetSpreadFactor validates that spread factor is set to zero.
func (s *CosmWasmPoolSuite) TestGetSpreadFactor() {
	var (
		expectedSwapFee = osmomath.ZeroDec()
	)

	pool := s.PrepareCosmWasmPool()

	actualSwapFee := pool.GetSpreadFactor(s.Ctx)

	s.Require().Equal(expectedSwapFee, actualSwapFee)
}

// TestSpotPrice validates that spot price is returned as one.
func (s *CosmWasmPoolSuite) TestSpotPrice() {
	var (
		expectedSpotPrice = osmomath.OneBigDec()
	)

	pool := s.PrepareCosmWasmPool()

	actualSpotPrice, err := pool.SpotPrice(s.Ctx, denomA, denomB)
	s.Require().NoError(err)

	s.Require().Equal(expectedSpotPrice, actualSpotPrice)

	actualSpotPrice, err = pool.SpotPrice(s.Ctx, denomB, denomA)
	s.Require().NoError(err)

	s.Require().Equal(expectedSpotPrice, actualSpotPrice)
}

func (s *CosmWasmPoolSuite) TestGetPoolDenoms() {
	// var (
	// 	expectedSwapFee = osmomath.ZeroDec()
	// )

	pool := s.PrepareCosmWasmPool()

	poolLiquidity := pool.GetTotalPoolLiquidity(s.Ctx)
	fmt.Println("poolLiquidity", poolLiquidity)

	fmt.Println(osmoutils.CoinsDenoms(poolLiquidity))

	// s.Require().Equal(expectedSwapFee, actualSwapFee)
}
