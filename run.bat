@echo off

start cmd /c "cd C:\Code\Go\mall\service\user && go run user.go"
start cmd /c "cd C:\Code\Go\mall\service\product && go run product.go"
start cmd /c "cd C:\Code\Go\mall\service\Order && go run Order.go"
start cmd /c "cd C:\Code\Go\mall\service\cart && go run cart.go"


