package structs

type Orders struct {
	ID                string     `json:"order_uid" db:"uid"`
	TrackNumber       string     `json:"track_number,omitempty" db:"track_number"`
	Entry             string     `json:"entry,omitempty" db:"entry"`
	Delivery          Deliveries `json:"delivery"`
	Payments          Payments   `json:"payments"`
	Items             []Item     `json:"items,omitempty" fakesize:"2,3"`
	Locale            string     `json:"locale,omitempty" db:"locale"`
	InternalSignature string     `json:"internal_signature,omitempty" db:"internal_signature"`
	CustomerID        string     `json:"customer_id,omitempty" db:"customer_id"`
	DeliveryService   string     `json:"delivery_service,omitempty" db:"delivery_service"`
	ShardKey          string     `json:"shardkey,omitempty" db:"shardkey"`
	SmID              int        `json:"sm_id,omitempty" db:"sm_id" fake:"{number: 1, 10000}"`
	DateCreated       string     `json:"date_created" db:"date_created"`
	OofShard          string     `json:"oof_shard,omitempty" db:"oof_shard"`
}

type Ord struct {
	ID                string     `json:"order_uid" db:"uid"`
	TrackNumber       string     `json:"track_number,omitempty" db:"track_number"`
	Entry             string     `json:"entry,omitempty" db:"entry"`
	Delivery          Deliveries `json:"delivery"`
	Payment           Payments   `json:"payment"`
	Locale            string     `json:"locale,omitempty" db:"locale"`
	InternalSignature string     `json:"internal_signature,omitempty" db:"internal_signature"`
	CustomerID        string     `json:"customer_id,omitempty" db:"customer_id"`
	DeliveryService   string     `json:"delivery_service,omitempty" db:"delivery_service"`
	ShardKey          string     `json:"shardkey,omitempty" db:"shardkey"`
	SmID              int        `json:"sm_id,omitempty" db:"sm_id" fake:"{number: 1, 10000}"`
	DateCreated       string     `json:"date_created" db:"date_created"`
	OofShard          string     `json:"oof_shard,omitempty" db:"oof_shard"`
}

type Deliveries struct {
	Name        string `json:"name,omitempty" db:"name"`
	PhoneNumber string `json:"phone,omitempty" db:"phone"`
	Zip         string `json:"zip,omitempty" db:"zip"`
	City        string `json:"city,omitempty" db:"city"`
	Address     string `json:"address,omitempty" db:"address"`
	Region      string `json:"region,omitempty" db:"region"`
	Email       string `json:"email,omitempty" db:"email"`
}

type Payments struct {
	Transaction string `json:"transaction,omitempty" db:"transaction"`
	RequestID   string `json:"request_id,omitempty" db:"request_id"`
	Currency    string `json:"currency,omitempty" db:"currency"`
	Provider    string `json:"provider,omitempty" db:"provider"`
	Amount      int    `json:"amount,omitempty" db:"amount" fake:"{number: 1, 10000}"`
	PaymentDT   int    `json:"payment_dt,omitempty" db:"payment_dt" fake:"{number: 1, 10000}"`
	BankName    string `json:"bank,omitempty" db:"bank"`
	Cost        int    `json:"delivery_cost,omitempty" db:"delivery_cost" fake:"{number: 1, 10000}"`
	TotalGoods  int    `json:"goods_total,omitempty" db:"goods_total" fake:"{number: 1, 10000}"`
	CustomFee   int    `json:"custom_fee,omitempty" db:"custom_fee" fake:"{number: 1, 10000}"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id,omitempty" db:"chrt_id" fake:"{number: 1, 10000}"`
	TrackNumber string `json:"track_number,omitempty" db:"track_number"`
	Price       int    `json:"price,omitempty" db:"price" fake:"{number: 1, 10000}"`
	RID         string `json:"rid,omitempty" db:"rid"`
	Name        string `json:"name,omitempty" db:"name"`
	Sale        int    `json:"sale,omitempty" db:"sale" fake:"{number: 1, 10000}"`
	Size        string `json:"size,omitempty" db:"size"`
	TotalPrice  int    `json:"total_price,omitempty" db:"total_price" fake:"{number: 1, 10000}"`
	NmID        int    `json:"nm_id,omitempty" db:"nm_id" fake:"{number: 1, 10000}"`
	Brand       string `json:"brand,omitempty" db:"brand"`
	Status      int    `json:"status,omitempty" db:"status" fake:"{number: 1, 10000}"`
}
