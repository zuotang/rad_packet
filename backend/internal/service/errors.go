package service

import "errors"

var (
	ErrAlreadyBound      = errors.New("already bound")
	ErrBindSelf          = errors.New("cannot bind self")
	ErrReferralCode      = errors.New("invalid referral code")
	ErrAlreadyClaimed    = errors.New("task already claimed")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrWithdrawBelowMin  = errors.New("withdraw amount below minimum")
	ErrRiskCheckFailed   = errors.New("risk check failed")
	ErrWithdrawState     = errors.New("invalid withdraw state transition")
	ErrWithdrawNotFound  = errors.New("withdraw request not found")
	ErrNoSpinChance      = errors.New("no spin chance")
)
