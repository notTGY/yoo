# yoo
YooMoney Client

## API

```go
c, _ := yoo.NewClient(YOOKASSA_SHOP_ID, YOOKASSA_SECRET_KEY)

id, crt, _ := c.CreatePayment(9900, "example@example.com", "тестовый платеж", "hello-world")

log.Printf("ID: %s\nCRT: %s\n", id, crt)

r, _ := c.GetPayment(id)
log.Println(r)
```


## Проверить создание платежа

чтобы создать тестовый платеж выполни:
```bash
YOOKASSA_SHOP_ID=xxx YOOKASSA_SECRET_KEY=xxx go run cmd/main.go
```

получить информацию о платеже (ID, это ID из предыдущего шага):
```bash
ID=xxx YOOKASSA_SHOP_ID=xxx YOOKASSA_SECRET_KEY=xxx go run cmd/main.go
```

