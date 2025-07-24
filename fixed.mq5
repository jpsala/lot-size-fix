// This is the relevant part of an .mql script patched by the patcher.go script
// for you to understand the changes made to the original script.
double sqMMFixedAmount(string symbol, ENUM_ORDER_TYPE orderType, double price, double sl, double RiskedMoney, int decimals, double LotsIfNoMM, double MaximumLots, double multiplier) {
   Verbose("Computing Money Management for order - Fixed amount");
   
   if(UseMoneyManagement == false) {
      Verbose("Use Money Management = false, MM not used");
      return (mmLotsIfNoMM);
   }
      
   string correctedSymbol = correctSymbol(symbol);
   sl = NormalizeDouble(sl, (int) SymbolInfoInteger(correctedSymbol, SYMBOL_DIGITS));
   
   double openPrice = price > 0 ? price : SymbolInfoDouble(correctedSymbol, isLongOrder(orderType) ? SYMBOL_ASK : SYMBOL_BID);
   double LotSize=0;

   if(RiskedMoney <= 0 ) {
      Verbose("Computing Money Management - Incorrect RiskedMoney value, it must be above 0");
      return(0);
   }
   
    
   double Smallest_Lot = SymbolInfoDouble(correctedSymbol, SYMBOL_VOLUME_MIN);
   double Largest_Lot = SymbolInfoDouble(correctedSymbol, SYMBOL_VOLUME_MAX);    
   double LotStep = SymbolInfoDouble(correctedSymbol, SYMBOL_VOLUME_STEP);
		
   
	Verbose(StringFormat("JP: >> Entering sqMMFixedAmount. Symbol: %s, OrderType: %s, Price: %.5f, SL: %.5f, RiskedMoney: %.2f, Decimals: %d, LotsIfNoMM: %.2f, MaxLots: %.2f, Multiplier: %.2f", symbol, EnumToString(orderType), price, sl, RiskedMoney, decimals, LotsIfNoMM, MaximumLots, multiplier));

	// Calculate profit/loss for a 1-lot trade to determine the exact drawdown
	double oneLotSLDrawdown;
	Verbose(StringFormat("JP: Calculating profit with OpenPrice: %.5f, SL: %.5f", openPrice, sl));
	if(!OrderCalcProfit(isLongOrder(orderType) ? ORDER_TYPE_BUY : ORDER_TYPE_SELL, correctedSymbol, 1.0, openPrice, sl, oneLotSLDrawdown)) {
		Verbose("JP: OrderCalcProfit failed. Error: ", GetLastError());
		return 0;
	}
	oneLotSLDrawdown = MathAbs(oneLotSLDrawdown);
	Verbose(StringFormat("JP: Money to risk: %.2f, One Lot SL Drawdown: %.2f, Open Price: %.5f, SL: %.5f, Distance: %.5f", RiskedMoney, oneLotSLDrawdown, openPrice, sl, MathAbs(openPrice - sl)));
	// --- FIX END ---

		
   if(oneLotSLDrawdown > 0) {
	  LotSize = roundDown(RiskedMoney / oneLotSLDrawdown, decimals);
   }
   else {
	  LotSize = 0;
   }
   LotSize = LotSize * multiplier;
   //--- MAXLOT and MINLOT management

   Verbose("Computing Money Management - Smallest_Lot: ", DoubleToString(Smallest_Lot), ", Largest_Lot: ", DoubleToString(Largest_Lot), ", Computed LotSize: ", DoubleToString(LotSize));
   // Verbose("Money to risk: ", DoubleToString(RiskedMoney), ", Max 1 lot trade drawdown: ", DoubleToString(oneLotSLDrawdown), ", Point value: ", DoubleToString(PointValue));

   if(LotSize <= 0) {
      Verbose("Calculated LotSize is <= 0. Using LotsIfNoMM value: ", DoubleToString(LotsIfNoMM), ")");
			LotSize = LotsIfNoMM;
	 }                              

   if (LotSize < Smallest_Lot) {
      Verbose("Calculated LotSize is too small. Minimal allowed lot size from the broker is: ", DoubleToString(Smallest_Lot), ". Please, increase your risk or set fixed LotSize.");
      LotSize = 0;
   }
   else if (LotSize > Largest_Lot) {
      Verbose("LotSize is too big. LotSize set to maximal allowed market value: ", DoubleToString(Largest_Lot));
      LotSize = Largest_Lot;
   }

   if(LotSize > MaximumLots) {
      Verbose("LotSize is too big. LotSize set to maximal allowed value (MaximumLots): ", DoubleToString(MaximumLots));
      LotSize = MaximumLots;
   }

   //--------------------------------------------

   
	// --- JPS - Margin Check START ---
	double margin_required;
	if(!OrderCalcMargin(isLongOrder(orderType) ? ORDER_TYPE_BUY : ORDER_TYPE_SELL, correctedSymbol, LotSize, openPrice, margin_required)) {
		Verbose("JP: OrderCalcMargin failed for initial LotSize. Error: ", GetLastError());
	} else {
		double free_margin = AccountInfoDouble(ACCOUNT_MARGIN_FREE);
		Verbose(StringFormat("JP: Margin Check - Initial LotSize: %.2f, Required Margin: %.2f, Free Margin: %.2f", LotSize, margin_required, free_margin));

		if(margin_required > free_margin) {
			Verbose("JP: Not enough free margin. Adjusting LotSize down...");
			while(margin_required > free_margin && LotSize > Smallest_Lot) {
				LotSize -= LotStep;
				if(LotSize < Smallest_Lot) {
					LotSize = 0;
					break;
				}
				if(!OrderCalcMargin(isLongOrder(orderType) ? ORDER_TYPE_BUY : ORDER_TYPE_SELL, correctedSymbol, LotSize, openPrice, margin_required)){
					Verbose("JP: OrderCalcMargin failed during adjustment. Error: ", GetLastError());
					LotSize = 0; // Fail safe
					break;
				}
				Verbose(StringFormat("JP: Margin Check (Adjusting) - New LotSize: %.2f, Required Margin: %.2f", LotSize, margin_required));
			}

			if(LotSize > 0) {
				Verbose(StringFormat("JP: Final Adjusted LotSize to fit margin: %.2f", LotSize));
			} else {
				Verbose("JP: Could not adjust LotSize to fit margin. LotSize set to 0.");
			}
		}
	}
	// --- JPS - Margin Check END ---

	Verbose(StringFormat("JP: << Exiting sqMMFixedAmount. Final LotSize: %.2f", LotSize));
	return (LotSize);
}
