package lib

import (
	"math"
	"socious/src/apps/models"
)

// Fee multipliers
const IMPACT_ORG_FEE = 0.02
const IMPACT_USER_FEE = 0.05
const ORG_FEE = 0.03
const USER_FEE = 0.1
const STRIPE_FEE = 0.036
const REFERRED_ORG_FEE_DISCOUNT = 0.5
const REFERRED_USER_FEE_DISCOUNT = 0.5

type AmountsOptions struct {
	Amount             float64
	Round              *float64
	IsVerified         bool
	Service            models.PaymentService
	OrgReferredWallet  *string
	UserReferredWallet *string
	OrgFeeDiscount     bool
	UserFeeDiscount    bool
}

func AmountsOptionsFromContract(contract models.Contract, orgReferrer *models.Referring, userReferrer *models.Referring) AmountsOptions {
	round := 1.0
	service := models.PaymentServiceStripe

	if contract.Currency != nil && *contract.Currency != models.JPY {
		round = 100.0
	}
	if *contract.PaymentType == models.PaymentModeTypeCrypto {
		round = 100000.0
		service = models.PaymentServiceCrypto
	}

	isVerified := contract.Provider.MetaMap["verified_impact"]
	if isVerified != nil {
		isVerified = isVerified.(bool)
	} else {
		isVerified = false
	}

	//Referrerings
	var orgReferrerWallet, userReferrerWallet *string = nil, nil
	orgReferrerFeeDiscount, userReferrerFeeDiscount := false, false
	if orgReferrer != nil {
		orgReferrerWallet = orgReferrer.WalletAddress
		orgReferrerFeeDiscount = orgReferrer.FeeDiscount
	}
	if userReferrer != nil {
		userReferrerWallet = userReferrer.WalletAddress
		userReferrerFeeDiscount = userReferrer.FeeDiscount
	}

	return AmountsOptions{
		Amount:             contract.TotalAmount,
		Round:              &round,
		IsVerified:         isVerified.(bool),
		OrgReferredWallet:  orgReferrerWallet,
		UserReferredWallet: userReferrerWallet,
		OrgFeeDiscount:     orgReferrerFeeDiscount,
		UserFeeDiscount:    userReferrerFeeDiscount,
		Service:            service,
	}
}

func CalculateAmounts(options AmountsOptions) map[string]any {
	orgFeeRate, userFeeRate := ORG_FEE, USER_FEE
	if options.IsVerified {
		orgFeeRate, userFeeRate = IMPACT_ORG_FEE, IMPACT_USER_FEE
	}

	if options.OrgReferredWallet != nil && options.OrgFeeDiscount {
		orgFeeRate *= REFERRED_ORG_FEE_DISCOUNT
	}

	if options.UserReferredWallet != nil && options.UserFeeDiscount {
		userFeeRate *= REFERRED_USER_FEE_DISCOUNT
	}

	amount := options.Amount
	fee := amount * orgFeeRate

	//rounding
	round := *options.Round
	if options.Round == nil || *options.Round < 1 {
		*options.Round = 1
	}

	stripeFee := 0.0
	if models.PaymentService(options.Service) == models.PaymentServiceStripe {
		stripeFee = (fee + options.Amount) * STRIPE_FEE
	}

	total := math.Ceil((amount+fee+stripeFee)*round) / round
	payoutFee := amount * userFeeRate
	payout := amount - payoutFee

	//Referrings
	orgReferredWallet, userReferredWallet := "", ""
	if options.OrgReferredWallet != nil {
		orgReferredWallet = *options.OrgReferredWallet
	}
	if options.UserReferredWallet != nil {
		userReferredWallet = *options.UserReferredWallet
	}

	return map[string]any{
		"amount":               amount,
		"fee":                  fee,
		"stripe_fee":           stripeFee,
		"total":                total,
		"payout":               payout,
		"app_fee":              fee + payoutFee,
		"org_referrer_wallet":  orgReferredWallet,
		"user_referrer_wallet": userReferredWallet,
		"org_fee_discount":     options.OrgFeeDiscount,
		"user_fee_discount":    options.UserFeeDiscount,
	}
}
