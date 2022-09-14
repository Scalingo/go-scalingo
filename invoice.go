package scalingo

import (
	"context"
	"time"

	"gopkg.in/errgo.v1"
)

type InvoicesService interface {
}

var _ InvoicesService = (*Client)(nil)

type Invoice struct {
	ID                string        `json:"id"`
	TotalPrice        int           `json:"total_price"`
	TotalPriceWithVat int           `json:"total_price_with_vat"`
	BillingMonth      time.Time     `json:"billing_month"`
	PdfUrl            string        `json:"pdf_url"`
	InvoiceNumber     string        `json:"invoice_number"`
	State             string        `json:"state"`
	VatRate           int           `json:"vat_rate"`
	Items             []interface{} `json:"items"`
	DetailedItems     []interface{} `json:"detailed_items"`
}

type Invoices []Invoice

type InvoicesRes struct {
	Invoices Invoices `json:"invoices"`
	Meta     struct {
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

func (c *Client) ListInvoices(ctx context.Context, opts PaginationOpts) (Invoices, PaginationMeta, error) {
	var invoicesRes InvoicesRes
	err := c.ScalingoAPI().ResourceList(ctx, "account/invoices", opts.ToMap(), &invoicesRes)
	if err != nil {
		return nil, PaginationMeta{}, errgo.Mask(err)
	}
	return invoicesRes.Invoices, invoicesRes.Meta.PaginationMeta, nil
}
