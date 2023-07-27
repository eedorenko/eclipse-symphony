package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCampaignMatch(t *testing.T) {
	campaign1 := CampaignSpec{
		Name: "name",
	}
	targetSelector := TargetSelector{
		Name: "name",
	}
	equal, err := campaign1.DeepEquals(targetSelector)
	assert.Nil(t, err)
	assert.True(t, equal)
}

func TestCampaignMatchOneEmpty(t *testing.T) {
	campaign1 := CampaignSpec{
		Name: "name",
	}
	res, err := campaign1.DeepEquals(nil)
	assert.Errorf(t, err, "parameter is not a CampaignSpec type")
	assert.False(t, res)
}

func TestCampaignRoleNotMatch(t *testing.T) {
	campaign1 := CampaignSpec{
		Name: "name",
	}
	targetSelector := TargetSelector{
		Name: "name1",
	}
	equal, err := campaign1.DeepEquals(targetSelector)
	assert.Nil(t, err)
	assert.False(t, equal)
}
