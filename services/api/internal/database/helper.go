package database

import "errors"

func ListingAccessModeFromInt(i int) (ListingAccessMode, error) {
	mode := ListingAccessMode(i)
	switch mode {
	case Private, Public:
		return mode, nil
	default:
		return 0, errors.New("invalid ListingAccessMode")
	}
}
