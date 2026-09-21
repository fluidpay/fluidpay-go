package fluidpay

// Response fixtures taken from https://sandbox.fluidpay.com/docs/.

const cardSaleFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "d3bfel2cunifubaumn00",
    "user_id": "testmerchant12345678",
    "user_name": "test_merchant",
    "merchant_id": "testmerchant12345678",
    "merchant_name": "Example Merchant",
    "idempotency_key": "7df4d862-1a3d-44c4-b3df-536aadf307b0",
    "idempotency_time": 300,
    "type": "sale",
    "status": "pending_settlement",
    "response": "approved",
    "response_code": 100,
    "amount": 1299,
    "base_amount": 1299,
    "amount_authorized": 1299,
    "amount_captured": 1299,
    "amount_settled": 0,
    "amount_refunded": 0,
    "payment_adjustment": 0,
    "tip_amount": 0,
    "processor_id": "d37kgq2cuni9almp7tmg",
    "processor_type": "tsys_sierra",
    "processor_name": "TSYS default",
    "payment_method": "card",
    "payment_type": "card",
    "features": ["avs"],
    "national_tax_amount": 0,
    "duty_amount": 0,
    "ship_from_postal_code": "",
    "summary_commodity_code": "",
    "merchant_vat_registration_number": "",
    "customer_vat_registration_number": "",
    "tax_amount": 100,
    "tax_exempt": false,
    "shipping_amount": 100,
    "surcharge": 0,
    "discount_amount": 0,
    "service_fee": 0,
    "currency": "usd",
    "description": "test transaction",
    "settlement_batch_id": "",
    "order_id": "someOrderID",
    "po_number": "somePONumber",
    "ip_address": "4.2.2.2",
    "transaction_source": "api",
    "email_receipt": false,
    "email_address": "user@home.com",
    "customer_id": "",
    "customer_payment_type": "",
    "customer_payment_id": "",
    "subscription_id": "",
    "referenced_transaction_id": "",
    "response_body": {
      "card": {
        "id": "d3bffkacunifubaumn60",
        "card_type": "visa",
        "first_six": "411111",
        "last_four": "1111",
        "masked_card": "411111******1111",
        "expiration_date": "12/26",
        "response": "approved",
        "response_code": 100,
        "auth_code": "TAS000",
        "processor_response_code": "00",
        "processor_response_text": "APPROVAL TAS000 ",
        "processor_transaction_id": "000000000000000",
        "processor_type": "tsys_sierra",
        "processor_id": "d37kgq2cuni9almp7tmg",
        "bin_type": "STANDARD",
        "type": "debit",
        "avs_response_code": "Y",
        "cvv_response_code": "M",
        "processor_specific": "",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
      }
    },
    "custom_fields": {"cf_1": ["value"]},
    "line_items": null,
    "billing_address": {
      "first_name": "John",
      "last_name": "Smith",
      "company": "Test Company",
      "address_line_1": "123 Some St",
      "address_line_2": "",
      "city": "Wheaton",
      "state": "IL",
      "postal_code": "60187",
      "country": "US",
      "phone": "5555555555",
      "fax": "",
      "email": "help@website.com"
    },
    "shipping_address": {
      "first_name": "",
      "last_name": "",
      "company": "",
      "address_line_1": "",
      "address_line_2": "",
      "city": "",
      "state": "",
      "postal_code": "",
      "country": "",
      "phone": "",
      "fax": "",
      "email": ""
    },
    "created_at": "2025-09-26T20:30:09.693858Z",
    "updated_at": "2025-09-26T20:30:09.693858Z",
    "captured_at": "2025-09-26T20:30:09.702601Z",
    "settled_at": null
  }
}`

const declinedSaleFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "d3bfdeclined00000000",
    "type": "sale",
    "status": "declined",
    "response": "declined",
    "response_code": 202,
    "amount": 1299,
    "amount_authorized": 0,
    "payment_method": "card",
    "response_body": {
      "card": {
        "id": "d3bfdeclined00000001",
        "card_type": "visa",
        "masked_card": "400000******9995",
        "response": "declined",
        "response_code": 202,
        "processor_response_code": "51",
        "processor_response_text": "INSUFFICIENT FUNDS"
      }
    },
    "created_at": "2025-09-26T20:30:09.693858Z",
    "updated_at": "2025-09-26T20:30:09.693858Z",
    "captured_at": null,
    "settled_at": null
  }
}`

const terminalSaleFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "b9efr9qj8m0ge2h7tat0",
    "type": "sale",
    "status": "pending_settlement",
    "response": "approved",
    "response_code": 100,
    "amount": 500,
    "amount_authorized": 500,
    "amount_captured": 500,
    "payment_method": "terminal",
    "response_body": {
      "terminal": {
        "id": "b9efr9qj8m0ge2h7tatg",
        "card_type": "mastercard",
        "payment_type": "credit",
        "entry_type": "swiped",
        "first_four": "5424",
        "last_four": "3333",
        "masked_card": "5424********3333",
        "cardholder_name": "FDCS TEST CARD /MASTERCARD",
        "auth_code": "",
        "response_code": 100,
        "processor_response_text": "APPROVAL VTLMC1",
        "processor_specific": {"BatchNum": "8", "Tip": "0.00"},
        "signature_data": "Qk0OIQ==",
        "created_at": "2018-01-15T19:14:47.225068Z",
        "updated_at": "2018-01-15T19:15:02.335853Z"
      }
    },
    "created_at": "2018-01-15T19:14:47.108371Z",
    "updated_at": "2018-01-15T19:15:02.529558Z",
    "captured_at": "2018-01-15T19:14:47.337763Z",
    "settled_at": null
  }
}`

const achSaleFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "d7xxxxxxxxxxxxxxxx8g",
    "type": "sale",
    "status": "pending_settlement",
    "response": "approved",
    "response_code": 100,
    "amount": 8912,
    "payment_method": "ach",
    "payment_type": "ach",
    "response_body": {
      "ach": {
        "id": "d7xxxxxxxxxxxxxxxx90",
        "account_type": "checking",
        "masked_account_number": "XX**********XX",
        "routing_number": "XXXXXXXXX",
        "sec_code": "web",
        "response": "approved",
        "response_code": 100,
        "auth_code": "XXXXXXXXX",
        "processor_response_code": "000",
        "processor_response_text": "Successful",
        "processor_type": "achcom",
        "processor_id": "cixxxxxxxxxxxxxxxx40",
        "processor_specific": {},
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
      }
    },
    "features": ["cash_discount"],
    "created_at": "2026-03-25T19:16:28.386935889Z",
    "updated_at": "2026-03-25T19:16:28.386935889Z",
    "captured_at": "2026-03-25T19:16:28.401283835Z",
    "settled_at": null
  }
}`

// getTransactionFixture mirrors the documented GET /transaction/{id} example,
// which wraps the single record in an array.
const getTransactionFixture = `{
  "status": "success",
  "msg": "success",
  "data": [
    {
      "id": "b7kgflt1tlv51er0fts0",
      "type": "sale",
      "amount": 1112,
      "tax_amount": 100,
      "currency": "usd",
      "description": "test transaction",
      "order_id": "someOrderID",
      "payment_method": "card",
      "response": "approved",
      "response_code": 100,
      "status": "pending_settlement",
      "created_at": "2017-10-19T20:15:19.560708Z",
      "updated_at": "2017-10-19T20:15:20.832049Z"
    }
  ],
  "total_count": 1
}`

const searchTransactionsFixture = `{
  "status": "success",
  "msg": "",
  "total_count": 2,
  "data": [
    {
      "id": "b84vgb2j8m0jujadi4v0",
      "user_id": "aucio551tlv85l7moe60",
      "type": "sale",
      "amount": 1112,
      "amount_authorized": 1112,
      "amount_captured": 1112,
      "amount_settled": 0,
      "payment_method": "token",
      "payment_type": "card",
      "status": "pending_settlement",
      "response": "approved",
      "response_code": 100,
      "customer_id": "b81ko5qq9qq5v460r9i0",
      "transaction_source": "api",
      "response_body": {},
      "created_at": "2017-11-13T19:53:17Z",
      "updated_at": "2017-11-13T19:53:18Z",
      "captured_at": null,
      "settled_at": null
    },
    {
      "id": "b84vgb2j8m0jujadi4v1",
      "type": "refund",
      "amount": 500,
      "status": "refunded",
      "response": "approved",
      "response_code": 100,
      "transaction_source": "recurring",
      "subscription_id": "sub_1",
      "created_at": "2017-11-13T19:53:17Z",
      "updated_at": "2017-11-13T19:53:18Z"
    }
  ]
}`

const vaultCustomerFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "952f250d-fa85-40ac-a45e-3886f98032c6",
    "owner_id": "testmerchant12345678",
    "data": {
      "customer": {
        "addresses": [
          {
            "city": "Some Town",
            "company": "Some Business",
            "country": "US",
            "email": "user@somesite.com",
            "fax": "555555555",
            "first_name": "John",
            "hash": "cV9qS6UvnnRDJCkgfitHqGrKSRfTARDsLy5zfkNNTcM=",
            "id": "bl6nqk69ku6897l45fp0",
            "last_name": "Smith",
            "line_1": "123 Some St",
            "line_2": "",
            "phone": "5555555555",
            "postal_code": "60187",
            "state": "IL"
          }
        ],
        "defaults": {
          "billing_address_id": "bl6nqk69ku6897l45fp0",
          "payment_method_id": "bl6nqk69ku6897l45fq0",
          "payment_method_type": "card",
          "shipping_address_id": "bl6nqk69ku6897l45fp0"
        },
        "description": "test description",
        "notes": "vip",
        "flags": ["surcharge_exempt"],
        "payments": {
          "ach": [
            {
              "id": "bl6nqk69ku6897l45ach",
              "account_type": "checking",
              "sec_code": "web",
              "masked_account_number": "XXXXX1111",
              "routing_number": "111111111"
            }
          ],
          "cards": [
            {
              "card_type": "visa",
              "expiration_date": "2020-12-31T00:00:00Z",
              "id": "bl6nqk69ku6897l45fq0",
              "masked_number": "411111******1111",
              "flags": [],
              "processor_id": ""
            }
          ]
        }
      }
    },
    "created_at": "2019-08-09T09:04:00.786094073-05:00",
    "updated_at": "2019-08-09T09:04:00.786094116-05:00"
  }
}`

const vaultAddressCreatedFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "952f250d-fa85-40ac-a45e-3886f98032c6",
    "created_address_id": "newaddr0000000000000",
    "data": {"customer": {"addresses": [{"id": "newaddr0000000000000", "line_1": "1 Main"}], "defaults": {}, "payments": {"ach": [], "cards": []}}},
    "created_at": "2019-08-09T09:04:00Z",
    "updated_at": "2019-08-09T09:04:00Z"
  }
}`

const vaultPaymentCreatedFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "952f250d-fa85-40ac-a45e-3886f98032c6",
    "created_payment_method_id": "newcard0000000000000",
    "data": {"customer": {"addresses": [], "defaults": {}, "payments": {"ach": [], "cards": [{"id": "newcard0000000000000", "masked_number": "411111******1111"}]}}},
    "created_at": "2019-08-09T09:04:00Z",
    "updated_at": "2019-08-09T09:04:00Z"
  }
}`

const vaultSearchFixture = `{
  "status": "success",
  "msg": "success",
  "data": [
    {
      "created_at": "2019-07-17T20:20:26.255763044Z",
      "data": {
        "customer": {
          "addresses": [{"id": "bgl4vq1erttokpdi1kk0", "first_name": "TEST", "line_1": "ABC LANE"}],
          "defaults": {"billing_address_id": "bgl4vq1erttokpdi1kk0", "payment_method_id": "bmak1e1erttvtct22910", "payment_method_type": "card"},
          "description": "",
          "flags": [],
          "notes": "",
          "payments": {
            "ach": [],
            "cards": [
              {"card_type": "visa", "expiration_date": "12/19", "flags": [], "id": "bmak1e1erttvtct22910", "masked_number": "411111******1111", "processor_id": ""},
              {"card_type": "visa", "expiration_date": "12/20", "flags": [], "id": "bo572hherttokc9v84t0", "masked_number": "400551******0004", "processor_id": ""}
            ]
          }
        }
      },
      "id": "bgl4vq1erttokpdi1kj0",
      "updated_at": "2019-07-17T20:20:26.255763044Z"
    }
  ],
  "total_count": 81
}`

const legacyCustomerFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "b798ls2q9qq646ksu070",
    "description": "test description",
    "payment_method": {
      "card": {
        "id": "b798ls2q9qq646ksu080",
        "card_type": "visa",
        "first_six": "411111",
        "last_four": "1111",
        "masked_card": "411111******1111",
        "expiration_date": "12/20",
        "processor_id": "",
        "created_at": "2017-10-02T18:52:32Z",
        "updated_at": "2017-10-02T18:52:32Z"
      }
    },
    "billing_address": {
      "id": "b798ls2q9qq646ksu07g",
      "customer_id": "b798ls2q9qq646ksu070",
      "first_name": "John",
      "last_name": "Smith",
      "address_line_1": "123 Some St",
      "city": "Some Town",
      "state": "IL",
      "postal_code": "60187",
      "country": "US",
      "created_at": "2017-10-02T18:52:32Z",
      "updated_at": "2017-10-02T18:52:32Z"
    },
    "shipping_address": null,
    "created_at": "2017-10-02T18:52:32Z",
    "updated_at": "2017-10-02T18:52:32Z"
  }
}`

const legacyAddressListFixture = `{
  "status": "success",
  "msg": "success",
  "data": [
    {
      "id": "b798ls2q9qq646ksu07g",
      "customer_id": "b798ls2q9qq646ksu070",
      "first_name": "John",
      "last_name": "Smith",
      "address_line_1": "123 Some St",
      "city": "Some Town",
      "state": "IL",
      "postal_code": "60187",
      "country": "US",
      "created_at": "2017-10-02T18:52:32Z",
      "updated_at": "2017-10-02T18:52:32Z"
    }
  ],
  "count": 1
}`

const legacyCardListFixture = `{
  "status": "success",
  "msg": "success",
  "data": [
    {"id": "b798ls2q9qq646ksu080", "card_type": "visa", "last_four": "1111", "masked_card": "411111******1111", "expiration_date": "12/20"},
    {"id": "b799g1iq9qq6dk5l39i0", "card_type": "visa", "last_four": "1112", "masked_card": "411111******1112", "expiration_date": "12/20"}
  ],
  "count": 2
}`

const legacyCardCreatedFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "card": {
      "id": "b799g1iq9qq6dk5l39i0",
      "card_type": "visa",
      "first_six": "411111",
      "last_four": "1112",
      "masked_card": "411111******1112",
      "expiration_date": "12/20",
      "created_at": "2017-10-02T19:48:22.468483Z",
      "updated_at": "2017-10-02T19:48:22.468483Z"
    }
  }
}`

const addOnFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "b89ffdqj8m0o735i19i0",
    "name": "test addon",
    "description": "just a simple add-on",
    "amount": 100,
    "percentage": null,
    "duration": 0,
    "created_at": "2017-11-20T15:41:43.330315Z",
    "updated_at": "2017-11-20T15:41:43.330315Z"
  }
}`

const discountFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "b89flfqj8m0o735i19ig",
    "name": "test discount",
    "description": "this will discount the cost of the subscription",
    "amount": null,
    "percentage": 10000,
    "duration": 3,
    "created_at": "2017-11-20T15:41:43.330315Z",
    "updated_at": "2017-11-20T15:41:43.330315Z"
  }
}`

const planFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "b89g35qj8m0o735i19jg",
    "name": "test plan",
    "description": "just a simple test plan",
    "amount": 100,
    "billing_cycle_interval": 1,
    "billing_frequency": "twice_monthly",
    "billing_days": "1,15",
    "charge_on_day": false,
    "total_add_ons": 100,
    "total_discounts": 50,
    "duration": 0,
    "add_ons": [
      {
        "id": "b75cvl51tlv38t0o7o30",
        "name": "test_addon",
        "description": "this will add to the cost of the subscription",
        "amount": 100,
        "percentage": null,
        "duration": 0,
        "created_at": null,
        "updated_at": null
      }
    ],
    "discounts": [
      {
        "id": "b89flfqj8m0o735i19ig",
        "name": "test discount",
        "description": "this will discount the cost of the subscription",
        "amount": 50,
        "percentage": null,
        "duration": 0,
        "created_at": null,
        "updated_at": null
      }
    ],
    "created_at": "2017-11-20T16:23:51.990051Z",
    "updated_at": "2017-11-20T16:23:51.990051Z"
  }
}`

const subscriptionFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "b89gftaj8m0oft7upk80",
    "plan_id": "b89g35qj8m0o735i19jg",
    "status": "active",
    "description": "some description to describe the subscription",
    "customer": {
      "id": "b81ko5qq9qq5v460r9i0",
      "payment_method_type": "card",
      "payment_method_id": "b81ko5qq9qq5v460r9i0",
      "billing_address_id": "b81ko5qq9qq5v460r9i0",
      "shipping_address_id": "b81ko5qq9qq5v460r9i0"
    },
    "amount": 100,
    "total_adds": 0,
    "total_discounts": 50,
    "billing_cycle_interval": 1,
    "billing_frequency": "twice_monthly",
    "billing_days": "1,15",
    "charge_on_day": false,
    "duration": 0,
    "next_bill_date": "2017-11-22",
    "add_ons": null,
    "discounts": [
      {
        "id": "b89flfqj8m0o735i19ig",
        "description": "this will discount the cost of the subscription",
        "amount": 50,
        "percentage": null,
        "duration": 0
      }
    ],
    "created_at": "2017-11-20T16:51:01.798736Z",
    "updated_at": "2017-11-20T16:51:01.798736Z"
  }
}`

const terminalsFixture = `{
  "status": "success",
  "msg": "",
  "total_count": 1,
  "data": [
    {
      "id": "1ucio551tlv85l7moe5s",
      "merchant_id": "aucio551tlv85l7moe5g",
      "manufacturer": "dejavoo",
      "model": "z11",
      "serial_number": "1811000XXXX",
      "tpn": "1811000XXXX",
      "description": "front counter z11",
      "status": "active",
      "auth_key": "wcR1c9o1",
      "register_id": "1",
      "auto_settle": true,
      "settle_at": "00:00:00",
      "created_at": "2018-01-12T03:57:59Z",
      "updated_at": "0001-01-01T00:00:00Z"
    }
  ]
}`

const settlementSearchFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "summary": [
      {
        "merchant_id": "aaaaaaaaaaaaaaaaaaaa",
        "batch_date": "2018-12-04",
        "processor_id": "bbbbbbbbbbbbbbbbbbbb",
        "processor_name": "main mid",
        "num_transactions": 1,
        "captured": 10000,
        "credit": 0
      }
    ],
    "results": [
      {
        "id": "cccccccccccccccccccc",
        "merchant_id": "aaaaaaaaaaaaaaaaaaaa",
        "batch_date": "2018-12-04T13:31:02Z",
        "processor_id": "bbbbbbbbbbbbbbbbbbbb",
        "processor_name": "main mid",
        "processor_type": "tsys_sierra",
        "batch_number": 1,
        "num_transactions": 1,
        "amount_captured": 10000,
        "amount_credit": 0,
        "net_deposit": 10000,
        "response_code": 100,
        "response_message": " ACCEPT"
      }
    ]
  },
  "total_count": 1
}`

const binLookupFixture = `{
  "bin": "424242",
  "card_brand": "Visa",
  "issuing_bank": "Example Bank",
  "card_type": "credit",
  "card_level_generic": "standard",
  "country": "US",
  "is_surchargeable": true,
  "payment_method_type": "card"
}`

const userFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "d0ieitnsvrv16guvobeg",
    "username": "will_test",
    "name": "test user",
    "phone": "5555555555",
    "email": "user@fluidpay.com",
    "timezone": "America/Chicago",
    "status": "active",
    "role": "admin",
    "account_type": "merchant",
    "account_type_id": "testmerchant12345678",
    "flags": {"password_expired": "true"},
    "permissions": {
      "manage_users": true,
      "manage_api_keys": true,
      "process_sale": true,
      "process_credit": false
    },
    "notifications": {"merchant": {"transaction_receipts": false}},
    "defaults": {"processor_id": ""},
    "access_restrictions": {"ip": null},
    "two_factor_enabled": false,
    "api_key": "api_xxxxxxxxxxxxxxxxxxxx",
    "pub_api_key": "pub_xxxxxxxxxxxxxxxxxxx",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
  }
}`

const apiKeyFixture = `{
  "status": "success",
  "msg": "success",
  "data": {
    "id": "key00000000000000001",
    "user_id": "d0ieitnsvrv16guvobeg",
    "type": "api",
    "name": "backend",
    "api_key": "api_2EXPCBb5hEnyIa79UhE4Pa30WmH",
    "account_type": "merchant",
    "account_type_id": "testmerchant12345678",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}`

const webhookTransactionFixture = `{
  "status": "success",
  "msg": "success",
  "type": "transaction_create",
  "account_type": "merchant",
  "account_type_id": "testmerchant12345678",
  "transaction_id": "bm5s8gm9ku6ejcu15t9g",
  "action_at": "2019-09-25T19:47:14.001901427Z",
  "data": {
    "id": "bm5s8gm9ku6ejcu15t9g",
    "type": "sale",
    "status": "pending_settlement",
    "amount": 450,
    "amount_authorized": 450,
    "payment_method": "card",
    "response": "approved",
    "response_code": 100,
    "response_body": {"card": {"card_type": "visa", "masked_card": "411111******1111", "auth_code": "TAS000"}},
    "transaction_source": "api",
    "created_at": "2019-09-25T19:47:14.001901427Z",
    "captured_at": "2019-09-25T19:47:14.031268117Z",
    "settled_at": null,
    "updated_at": "2019-09-25T19:47:14.031268348Z"
  }
}
`
