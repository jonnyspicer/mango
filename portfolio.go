package mango

// LivePortfolioMetrics represents a user's current portfolio state.
type LivePortfolioMetrics struct {
	InvestmentValue     float64 `json:"investmentValue"`
	CashInvestmentValue float64 `json:"cashInvestmentValue"`
	Balance             float64 `json:"balance"`
	CashBalance         float64 `json:"cashBalance"`
	SpiceBalance        float64 `json:"spiceBalance"`
	TotalDeposits       float64 `json:"totalDeposits"`
	TotalCashDeposits   float64 `json:"totalCashDeposits"`
	LoanTotal           float64 `json:"loanTotal"`
	Timestamp           int64   `json:"timestamp"`
	Profit              float64 `json:"profit,omitempty"`
	UserId              string  `json:"userId"`
	DailyProfit         float64 `json:"dailyProfit"`
}

// PortfolioMetrics represents a snapshot of portfolio metrics at a point in time.
type PortfolioMetrics struct {
	InvestmentValue float64 `json:"investmentValue"`
	Balance         float64 `json:"balance"`
	TotalDeposits   float64 `json:"totalDeposits"`
	LoanTotal       float64 `json:"loanTotal"`
	Timestamp       int64   `json:"timestamp"`
	Profit          float64 `json:"profit,omitempty"`
}

// PortfolioPeriod represents a time period for portfolio history queries.
type PortfolioPeriod string

const (
	PeriodDaily   PortfolioPeriod = "daily"
	PeriodWeekly  PortfolioPeriod = "weekly"
	PeriodMonthly PortfolioPeriod = "monthly"
	PeriodAllTime PortfolioPeriod = "allTime"
)
