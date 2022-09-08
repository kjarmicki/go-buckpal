package account_domain

type Money struct {
	amount int64
}

func NewMoney(amount int64) Money {
	return Money{
		amount: amount,
	}
}

func MoneyAdd(a, b Money) Money {
	return NewMoney(a.amount + b.amount)
}

func MoneySubtract(a, b Money) Money {
	return NewMoney(a.amount - b.amount)
}

func (m Money) IsPositiveOrZero() bool {
	return m.amount >= 0
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

func (m Money) IsPositive() bool {
	return m.amount > 0
}

func (m Money) IsGreaterThanOrEqualTo(money Money) bool {
	return m.amount >= 0
}

func (m Money) IsGreaterThan(money Money) bool {
	return m.amount >= 1
}

func (m Money) Minus(money Money) Money {
	return NewMoney(m.amount - money.amount)
}

func (m Money) Plus(money Money) Money {
	return NewMoney(m.amount + money.amount)
}

func (m Money) Negate() Money {
	return NewMoney(-m.amount)
}

func (m Money) GetAmount() int64 {
	return m.amount
}
