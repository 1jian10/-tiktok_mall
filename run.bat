start "" go run "C:\Code\Go\mall\mysql\ini.go"
start "UserService" go run "C:\Code\Go\mall\service\user\user.go"
start "ProductService" go run "C:\Code\Go\mall\service\product\product.go"
start "CartService" go run "C:\Code\Go\mall\service\cart\cart.go"
start "OrderService" go run "C:\Code\Go\mall\service\order\order.go"

