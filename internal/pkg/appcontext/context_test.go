package appcontext

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AppContextTestSuite struct {
	suite.Suite
}

func TestAppContextTestSuite(t *testing.T) {
	suite.Run(t, new(AppContextTestSuite))
}

func (s *AppContextTestSuite) TestSSOId() {
	s.Run("Returns SSO ID from context", func() {
		// Assign
		ctx := context.Background()
		ctx = WithSSOId(ctx, "test-sso-id")

		// Act
		result := SSOIdFromContext(ctx)

		// Assert
		s.Equal("test-sso-id", result)
	})

	s.Run("Returns empty string when SSO ID not in context", func() {
		// Assign
		ctx := context.Background()

		// Act
		result := SSOIdFromContext(ctx)

		// Assert
		s.Empty(result)
	})
}

func (s *AppContextTestSuite) TestContextTrackingId() {
	s.Run("Returns tracking id from context", func() {
		ctx := context.Background()
		ctx = WithTrackingID(ctx, "test-tracking-id")

		// Act
		result := TrackingIDFromContext(ctx)

		// Assert
		s.Equal("test-tracking-id", result)
	})

	s.Run("Returns empty string when tracking id not in context", func() {
		// Assign
		ctx := context.Background()

		// Act
		result := TrackingIDFromContext(ctx)

		// Assert
		s.Empty(result)
	})
}

func (s *AppContextTestSuite) TestContextTimes() {
	s.Run("Returns start time from context", func() {
		// Assign
		ctx := context.Background()
		startTime := time.Now()
		ctx = WithStartTime(ctx, startTime)

		// Act
		result := StartTimeFromContext(ctx)

		// Assert
		s.Equal(startTime, result)
	})

	s.Run("Returns zero time when start time not in context", func() {
		// Assign
		ctx := context.Background()

		// Act
		result := StartTimeFromContext(ctx)

		// Assert
		s.Equal(time.Time{}, result)
	})

	s.Run("Returns end time from context", func() {
		// Assign
		ctx := context.Background()
		endTime := time.Now()
		ctx = WithEndTime(ctx, endTime)

		// Act
		result := EndTimeFromContext(ctx)

		// Assert
		s.Equal(endTime, result)
	})

	s.Run("Returns zero time when end time not in context", func() {
		// Assign
		ctx := context.Background()

		// Act
		result := EndTimeFromContext(ctx)

		// Assert
		s.Equal(time.Time{}, result)
	})
}
