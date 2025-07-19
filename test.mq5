input bool UseMoneyManagement = true;
input double mmLotsIfNoMM = 0.01;

double sqMMFixedAmount(string symbol, ENUM_ORDER_TYPE orderType, double price, double sl, double RiskedMoney, int decimals, double LotsIfNoMM, double MaximumLots, double multiplier) {
    Verbose("Computing Money Management for order: Fixed amount");
    Verbose("price: ", DoubleToString(price));
    Verbose("sl: ", DoubleToString(sl));
    Verbose("Difference: ", DoubleToString(sl - price));
    Verbose("RiskedMoney: ", DoubleToString(RiskedMoney));
    Verbose("decimals: ", IntegerToString(decimals));
    Verbose("LotsIfNoMM: ", DoubleToString(LotsIfNoMM));
    Verbose("MaximumLots: ", DoubleToString(MaximumLots));
    Verbose("multiplier: ", DoubleToString(multiplier));
    if(UseMoneyManagement == false) {
        Verbose("Use Money Management = false, MM not used");
        return(mmLotsIfNoMM);
    }
      
    string correctedSymbol = correctSymbol(symbol);
    sl = NormalizeDouble(sl, (int) SymbolInfoInteger(correctedSymbol, SYMBOL_DIGITS));
   
    double openPrice = price > 0 ? price : SymbolInfoDouble(correctedSymbol, isLongOrder(orderType) ? SYMBOL_ASK : SYMBOL_BID);
    double LotSize = 0;
    Verbose("openPrice: ", DoubleToString(openPrice));
    if(RiskedMoney <= 0 ) {
        Verbose("Computing Money Management - Incorrect RiskedMoney value, it must be above 0");
        return(0);
    }
   
    double PointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE);
    double Smallest_Lot = SymbolInfoDouble(correctedSymbol, SYMBOL_VOLUME_MIN);
    double Largest_Lot = SymbolInfoDouble(correctedSymbol, SYMBOL_VOLUME_MAX);
    double LotStep = SymbolInfoDouble(correctedSymbol, SYMBOL_VOLUME_STEP);
    Verbose("PointValue: ", DoubleToString(PointValue));
    Verbose("--- Symbol Info for Margin ---");
    Verbose("SYMBOL_TRADE_TICK_VALUE: ", DoubleToString(SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_VALUE)));
    Verbose("SYMBOL_TRADE_TICK_SIZE: ", DoubleToString(SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_SIZE)));
    Verbose("SYMBOL_TRADE_CONTRACT_SIZE: ", DoubleToString(SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE)));
    Verbose("--- End Symbol Info ---");
  
   //Maximum drawdown of this order if we buy 1 lot
    double oneLotSLDrawdown = PointValue * MathAbs(openPrice - sl);
  
    if(oneLotSLDrawdown > 0) {
        LotSize = roundDown(RiskedMoney / oneLotSLDrawdown, decimals);
    }
    else {
        LotSize = 0;
    }
    LotSize = LotSize * multiplier;
   //--- MAXLOT and MINLOT management
    Verbose("LotSize(before checks): ", DoubleToString(LotSize));

    Verbose("Smallest_Lot: ", DoubleToString(Smallest_Lot));
    Verbose("Largest_Lot: ", DoubleToString(Largest_Lot));
    Verbose("Money to risk: ", DoubleToString(RiskedMoney));
    Verbose("Max 1 lot trade drawdown: ", DoubleToString(oneLotSLDrawdown));

    if(LotSize <= 0) {
        Verbose("LotSize <= 0, using LotsIfNoMM: ", DoubleToString(LotsIfNoMM));
        LotSize = LotsIfNoMM;
    }

    if(LotSize < Smallest_Lot) {
        Verbose("Calculated LotSize is too small. Minimal allowed lot size from the broker is: ", DoubleToString(Smallest_Lot), ". Please, increase your risk or set fixed LotSize.");
        LotSize = 0;
    }
    else if(LotSize > Largest_Lot) {
        Verbose("LotSize > Largest_Lot, set to: ", DoubleToString(Largest_Lot));
        LotSize = Largest_Lot;
    }

    if(LotSize > MaximumLots) {
        Verbose("LotSize > MaximumLots, set to: ", DoubleToString(MaximumLots));
        LotSize = MaximumLots;
    }


   //--------------------------------------------
    //--- Enhanced Margin Debugging
    double margin_required_debug;
    double free_margin_debug = AccountInfoDouble(ACCOUNT_MARGIN_FREE);
    
    if(OrderCalcMargin(isLongOrder(orderType) ? ORDER_TYPE_BUY : ORDER_TYPE_SELL, correctedSymbol, LotSize, openPrice, margin_required_debug)) {
        Verbose("--- Margin Debug ---");
        Verbose("Account Free Margin: ", DoubleToString(free_margin_debug));
        Verbose("Margin required for ", DoubleToString(LotSize), " lots: ", DoubleToString(margin_required_debug));
        if(margin_required_debug > free_margin_debug) {
            Verbose("Margin required is GREATER than free margin. This will cause an error.");
        } else {
            Verbose("Margin required is LESS than or EQUAL to free margin. This should be OK.");
        }
        Verbose("--- End Margin Debug ---");
    } else {
        Verbose("--- Margin Debug ---");
        Verbose("OrderCalcMargin failed for debugging. Error: ", IntegerToString(GetLastError()));
        Verbose("--- End Margin Debug ---");
    }

    Verbose("LotSize(final): ", DoubleToString(LotSize));
    return(LotSize);
}

bool isLongOrder(ENUM_ORDER_TYPE orderType){
   return orderType == ORDER_TYPE_BUY || orderType == ORDER_TYPE_BUY_LIMIT || orderType == ORDER_TYPE_BUY_STOP;
}
string correctSymbol(string symbol){
    if(symbol == NULL || symbol == "NULL" || symbol == "Current" || symbol == "0" || symbol == "Same as main chart") {
        return Symbol();
    }
        else return symbol;
}

double roundDown(double value, int decimals) {
  	double p = 0;
  	
  	switch(decimals) {
  		case 0: return (int) value; 
  		case 1: p = 10; break;
  		case 2: p = 100; break;
  		case 3: p = 1000; break;
  		case 4: p = 10000; break;
  		case 5: p = 100000; break;
  		case 6: p = 1000000; break;
  		default: p = MathPow(10, decimals);
  	}

  	value = value * p;
  	double tmp = MathFloor(value + 0.00000001);
  	return NormalizeDouble(tmp/p, decimals);
}