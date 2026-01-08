package customers

type CustomerRegistered struct {
	*Customer
}

func (e CustomerRegistered) Evolve(c *Customer) {
	*c = *e.Customer
}

type CustomerOrderClosed struct {
	CustomerID CustomerID
	OrderID    OrderID
}

func (CustomerOrderClosed) Evolve(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	CustomerID CustomerID
	OrderID    OrderID
}

func (OrderAccepted) Evolve(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	OrderID OrderID
	Reason  string
}

func (OrderRejected) Evolve(c *Customer) {

}
