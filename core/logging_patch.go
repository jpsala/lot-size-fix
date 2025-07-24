package core

import (
	"fmt"
	"regexp"
)

// LotSizeLogging is a patch that adds detailed logging to the lot size calculation.
var LotSizeLogging = Patch{
	Name:        "LotSizeLogging",
	Description: "Adds detailed logging to the sqMMFixedAmount function to trace lot size calculation.",
	Apply: func(content string) (string, error) {
		originalContent := content

		// Regex to find the function signature
		reFunc := regexp.MustCompile(`(double\s+sqMMFixedAmount\s*\([^)]*\)\s*\{)`)

		// 1. Log entry and parameters
		logEntry := `
	PrintFormat("JP: >> Entering sqMMFixedAmount. Symbol: %s, OrderType: %s, Price: %.5f, SL: %.5f, RiskedMoney: %.2f, Decimals: %d, LotsIfNoMM: %.2f, MaxLots: %.2f, Multiplier: %.2f", symbol, EnumToString(orderType), price, sl, RiskedMoney, decimals, LotsIfNoMM, MaximumLots, multiplier);`

		if reFunc.MatchString(content) {
			content = reFunc.ReplaceAllString(content, fmt.Sprintf("$1%s", logEntry))
		}

		// 2. Log invalid RiskedMoney
		reInvalidRisk := regexp.MustCompile(`Verbose\("Computing Money Management - Incorrect RiskedMoney value, it must be above 0"\);`)
		logInvalidRisk := `PrintFormat("JP: !! Invalid RiskedMoney: %.2f. Must be > 0. Returning 0 lots.", RiskedMoney);`
		if reInvalidRisk.MatchString(content) {
			content = reInvalidRisk.ReplaceAllString(content, logInvalidRisk)
		}

		// 3. Log when LotSize is <= 0
		reLotSizeZero := regexp.MustCompile(`Verbose\("Calculated LotSize is <= 0\. Using LotsIfNoMM value: ", DoubleToString\(LotsIfNoMM\), "\)"\);`)
		logLotSizeZero := `PrintFormat("JP: ## Calculated LotSize was <= 0. Using default LotsIfNoMM: %.2f", LotsIfNoMM);`
		if reLotSizeZero.MatchString(content) {
			content = reLotSizeZero.ReplaceAllString(content, logLotSizeZero)
		}

		// 4. Log when LotSize is too small
		reLotSizeSmall := regexp.MustCompile(`Verbose\("Calculated LotSize is too small\. Minimal allowed lot size from the broker is: ", DoubleToString\(Smallest_Lot\), "\. Please, increase your risk or set fixed LotSize\."\);`)
		logLotSizeSmall := `PrintFormat("JP: !! Calculated LotSize %.2f is smaller than broker's minimum %.2f. Returning 0 lots.", LotSize, Smallest_Lot);`
		if reLotSizeSmall.MatchString(content) {
			content = reLotSizeSmall.ReplaceAllString(content, logLotSizeSmall)
		}

		// 5. Log the final return value
		reReturn := regexp.MustCompile(`(return\s*\(\s*LotSize\s*\);)`)
		logReturn := `
	PrintFormat("JP: << Exiting sqMMFixedAmount. Final LotSize: %.2f", LotSize);
	$1`
		if reReturn.MatchString(content) {
			content = reReturn.ReplaceAllString(content, logReturn)
		}

		// 6. Prefix all Verbose and VerboseLog calls
		reVerbose := regexp.MustCompile(`Verbose\("`)
		content = reVerbose.ReplaceAllString(content, `Verbose("JP: `)

		reVerboseLog := regexp.MustCompile(`VerboseLog\("`)
		content = reVerboseLog.ReplaceAllString(content, `VerboseLog("JP: `)

		// 7. Add detailed logging for broker-specific values
		reLotStep := regexp.MustCompile(`(double\s+LotStep\s*=\s*SymbolInfoDouble\([^)]*\);)`)
		logBrokerVals := `$1
	PrintFormat("JP: Broker Values - Smallest Lot: %.3f, Largest Lot: %.2f, Lot Step: %.3f", Smallest_Lot, Largest_Lot, LotStep);`
		if reLotStep.MatchString(content) {
			content = reLotStep.ReplaceAllString(content, logBrokerVals)
		}

		// 8. Log the exact values used for profit calculation
		reCalcProfit := regexp.MustCompile(`(if\(!OrderCalcProfit\([^)]*\))\s*\{`)
		logCalcParams := `
	PrintFormat("JP: Calculating profit with OpenPrice: %.5f, SL: %.5f", openPrice, sl);
	$1`
		if reCalcProfit.MatchString(content) {
			content = reCalcProfit.ReplaceAllString(content, logCalcParams)
		}

		if originalContent != content {
			return content, nil
		}

		return content, fmt.Errorf("no logging changes applied")
	},
}
